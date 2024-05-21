// Package common provides a basic HTTP client to perform HTTP calls.
// It supports Basic and OAuth authentication methods.
package common

import (
	"bytes"
	"errors"
	"fmt"
	"json"
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
	authType string         // basic or bearer
	username string
	password string
	accessToken string      // only for OAuth
	accessTokenExpire int   // only for OAuth
	refreshToken string     // only for OAuth
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
	ca.accessToken = ""
	ca.accessTokenExpire = 0
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
	ca.accessToken = ""
	ca.accessTokenExpire = 0
	ca.refreshToken = ""

	return nil
}

// Set ClientAuth authType to OAuth
func (ca *ClientAuth) SetOAuth(username, password, accessToken string,
	accessTokenExpire int, refreshToken string) error {
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
	ca.accessToken = accessToken
	ca.accessTokenExpire = accessTokenExpire
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

func (ca *ClientAuth) SetAccessToken(token string) error {
	if (ca.authType != OAuth) {
		return errors.New("token can be set only to OAuth client")
	}

	token = strings.Trim(token, " ")
	if (len(token) == 0) {
		return errors.New("token cannot be empty")
	}

	ca.accessToken = token
	return nil
}

func (ca *ClientAuth) GetAccessToken() string {
	return ca.accessToken
}

func (ca *ClientAuth) SetAccessTokenExpire(expire int) error {
	if (ca.authType != OAuth) {
		return errors.New("refreshToken can be set only to OAuth client")
	}

	ca.accessTokenExpire = expire
	return nil
}

func (ca *ClientAuth) GetAccessTokenExpire() int {
	return ca.accessTokenExpire
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

// Performs a HTTP Call using Client configuration
// `params` map may contain two keys:
// 	- "content_type" is the content type of the request, it defaults to
// 	  "application/json"
//  - "data" is the data to be sent in the request body.
// 		- For "application/json" content type, data is serialized as JSON
// 		- For "application/x-www-form-urlencoded" content type, data is
//  	  serialized as a form
func (c *Client) Call(method, path string, params map[string]interface{}) (
	*http.Response,	error) {
	fullPath := c.rootUrl.String() + path  //FIXME join strings
	var req *http.Request
	var err error

	switch method {
	case "GET", "DELETE":
		req, err = http.NewRequest(method, fullPath, nil)
	case "POST", "PUT", "PATCH":
		contentType, ok := params["content_type"]
		if !ok {
			contentType = "application/json"
		}

		if contentType == "appliaction/json" {
			data := params["data"].(map[string]interface{})
			var dataJson []byte
			if len(data) > 0 {
				// Replaces all occurrences of "\u003c" with "<" in dataJson
				// This is necessary because the json.Marshal function escapes
				// '<' character as "\u003c", but some backend services (like
				// Custodia) don't expect this.
				dataJson = bytes.Replace(dataJson, []byte("\\u003c"),
					[]byte("<"), -1)
			} else {
				dataJson = []byte("")
			}

			dataJson, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}
			req, err = http.NewRequest(method, fullPath, bytes.NewBuffer(dataJson))
			req.Header.Set("Content-Type", "application/json")
		} else if contentType == "application/x-www-form-urlencoded" {
			values := url.Values{}
			formData := params["data"].(map[string]string)
			for key, value := range formData {
				values.Add(key, value)
			}
			req, err = http.NewRequest("POST", fullPath,
				strings.NewReader(values.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else {
			panic(fmt.Sprintf("unsupported content type %q", contentType))
		}
	default:
		err = fmt.Errorf("unsupported HTTP method %q", method)
	}

	if err != nil {
		return nil, err
	}

	// set the headers
	req.Header.Set("User-Agent", userAgent)
	req.Header.Add("Accept", "application/json")

	// handle auth
	switch c.auth.authType {
	case NoAuth:
		// do nothing
	case BasicAuth:
		req.SetBasicAuth(c.auth.username, c.auth.password)
	case OAuth:
		bearer := "Bearer: " + c.auth.accessToken
		req.Header.Add("Authorization", bearer)
	default:
		panic(fmt.Sprintf("Unsupported auth type %q", c.auth.authType))
	}

	// perform the call
	client := &http.Client{}
    resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// Get wraps call to perform a HTTP GET call
func (c *Client) Get(path string) (*http.Response, error) {
	return c.Call("GET", path)
}

// Post wraps call to perform a HTTP POST call
func (c *Client) Post(path, payload string) (*http.Response, error) {
	return c.Call("POST", path, payload)
}

// Put wraps call to perform a HTTP PUT call
func (c *Client) Put(path, payload string) (*http.Response, error) {
	return c.Call("PUT", path, payload)
}

// Patch wraps call to perform a HTTP PATCH call
func (c *Client) Patch(path, payload string) (*http.Response, error) {
	return c.Call("PATCH", path, payload)
}

// Delete wraps call to perform a HTTP DELETE call
func (c *Client) Delete(path string) (*http.Response, error) {
	return c.Call("DELETE", path)
}
