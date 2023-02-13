// Package common provides a basic HTTP client to perform HTTP calls.
// It supports Basic and OAuth authentication methods.
package common

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const NoAuth = "No Auth"
const BasicAuth = "Basic"
const OAuth = "Bearer"
const userAgent = "golang/chino-" + Version

// ClientAuth keeps the authentication details - Basic vs Bearer (OAuth)
type ClientAuth struct {
	authType string      // basic or bearer
	username string  
	password string  
	token string  		 // only for OAuth
	refreshToken string  // only for OAuth
}

// Client holds the configuration (url, auth) and wraps http Requests
type Client struct {
	rootUrl *url.URL
	auth *ClientAuth
	userAgent string
}

// NewClientAuth returns a new ClientAuth with auth set to NoAuth
func NewClientAuth() *ClientAuth {
	ca := &ClientAuth{}
	ca.SetNoAuth()
	return ca
}

// Set ClientAuth authType to NoAuth removing other attributes
func (ca *ClientAuth) SetNoAuth() {
	ca.authType = NoAuth
	ca.username = ""
	ca.password = ""
	ca.token = ""
	ca.refreshToken = ""
}

// Set ClientAuth authType to BasicAuth removing tokens
func (ca *ClientAuth) SetBasicAuth(username, password string) error {
	if !IsValidUUID(username) {
		return errors.New("username must be a valid UUID")
	} else if !IsValidUUID(password) {
		return errors.New("password must be a valid UUID")
	}

	ca.authType = BasicAuth
	ca.username = username
	ca.password = password
	ca.token = ""
	ca.refreshToken = ""

	return nil
}

// Set ClientAuth authType to OAuth
func (ca *ClientAuth) SetOAuth(username, password, token, refreshToken string) error {
	username = strings.Trim(username, " ")
	password = strings.Trim(password, " ")
	
	if (len(username) == 0) {
		return errors.New("username cannot be empty")
	} else if (len(password) == 0) {
		return errors.New("password cannot be empty")
	}

	ca.authType = OAuth
	ca.username = username
	ca.password = password
	ca.token = token
	ca.refreshToken = refreshToken

	return nil
}

func (ca *ClientAuth) GetAuthType() string {
	return ca.authType
}

func (ca *ClientAuth) SetUsername(username string) error {
	switch ca.authType {
	case NoAuth:
		return errors.New("cannot set username to NoAuth client")
	case BasicAuth:
		if !IsValidUUID(username) {
			return errors.New("username must be a valid UUID")
		}
	case OAuth:
		username = strings.Trim(username, " ")
		if (len(username) == 0) {
			return errors.New("username cannot be empty")
		}
	}

	ca.username = username
	return nil
}

func (ca *ClientAuth) GetUsername() string {
	return ca.username
}

func (ca *ClientAuth) SetPassword(password string) error {
	switch ca.authType {
	case NoAuth:
		return errors.New("cannot set password to NoAuth client")
	case BasicAuth:
		if !IsValidUUID(password) {
			return errors.New("password must be a valid UUID")
		}
	case OAuth:
		password = strings.Trim(password, " ")
		if (len(password) == 0) {
			return errors.New("password cannot be empty")
		}
	}

	ca.password = password
	return nil
}

func (ca *ClientAuth) SetToken(token string) error {
	if (ca.authType != OAuth) {
		return errors.New("token can be set only to OAuth client")
	}

	token = strings.Trim(token, " ")
	if (len(token) == 0) {
		return errors.New("token cannot be empty")
	}

	ca.token = token
	return nil
}

func (ca *ClientAuth) GetToken() string {
	return ca.token
}

func (ca *ClientAuth) SetRefreshToken(refreshToken string) error {
	if (ca.authType != OAuth) {
		return errors.New("refreshToken can be set only to OAuth client")
	}

	refreshToken = strings.Trim(refreshToken, " ")
	if (len(refreshToken) == 0) {
		return errors.New("refreshToken cannot be empty")
	}
	ca.refreshToken = refreshToken
	return nil
}

func (ca *ClientAuth) GetRefreshToken() string {
	return ca.refreshToken
}

// NewClient configures and returns a new Client
func NewClient(serverUrl string, auth *ClientAuth) *Client {
	parsedUrl, err := url.Parse(serverUrl)
	if err != nil {
		panic(err)
	}

	if parsedUrl.Scheme == "" {
		panic("serverUrl has no schema")
	} 
	// FIXME: anything else to handle?

	return &Client{
		rootUrl: parsedUrl,
		auth: auth,
	}
}

// Performs a HTTP call using Client configuration
func (c *Client) call(method, path string, data ...string) (*http.Response, 
	error) {
	url := c.rootUrl.String() + path  //FIXME join strings
	var req *http.Request
	var err error

	switch method {
	case "GET", "DELETE":
		req, err = http.NewRequest(method, url, nil)
	case "POST", "PUT", "PATCH":
		var jsonStr []byte
		if len(data) > 0 {
			jsonStr = []byte(data[0])
		} else {
			jsonStr = []byte("")
		}
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
	default:
		err = fmt.Errorf("unsupported HTTP method %q", method)
	}

	if err != nil {
		return nil, err
	}

	// set the User-Agent	
	req.Header.Set("User-Agent", userAgent)

	// handle auth
	switch c.auth.authType {
	case NoAuth:
		// do nothing
	case BasicAuth:
		req.SetBasicAuth(c.auth.username, c.auth.password)
	case OAuth:
		bearer := "Bearer: " + c.auth.token
		req.Header.Add("Authorization", bearer)
	default:
		panic(fmt.Sprintf("Unsupported auth type %q", c.auth.authType))
	}

	// perform the call
	client := &http.Client{}
    resp, err := client.Do(req)
    defer resp.Body.Close()

	return resp, err
}

// Get wraps call to perform a HTTP GET call
func (c *Client) Get(path string) (*http.Response, error) {
	return c.call("GET", path)
}

// Post wraps call to perform a HTTP POST call
func (c *Client) Post(path, payload string) (*http.Response, error) {
	return c.call("POST", path, payload)
}

// Put wraps call to perform a HTTP PUT call
func (c *Client) Put(path, payload string) (*http.Response, error) {
	return c.call("PUT", path, payload)
}

// Patch wraps call to perform a HTTP PATCH call
func (c *Client) Patch(path, payload string) (*http.Response, error) {
	return c.call("PATCH", path, payload)
}

// Delete wraps call to perform a HTTP DELETE call
func (c *Client) Delete(path string) (*http.Response, error) {
	return c.call("DELETE", path)
}
