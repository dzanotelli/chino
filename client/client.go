package client

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const NoAuth = "No Auth"
const BasicAuth = "Basic"
const OAuth = "Bearer"

// Client holds the configuration (url, auth) and wraps http Requests
type Client struct {
	rootUrl *url.URL
	auth *ClientAuth
}

// ClientAuth keeps the authentication details - Basic vs Bearer (OAuth)
type ClientAuth struct {
	authType string      // basic or bearer
	username string  
	password string  
	token string  		 // only for OAuth
	refreshToken string  // only for OAuth
}

// NewClientAuth sets up the authentication details to pass to Client
// Args:
// - authType
// - username
// - password
// - token (OAuth only)
// - refreshToken (OAuth only)
func NewClientAuth(conf ...string) *ClientAuth {
	if len(conf) == 0 {
		panic("first argument needed: authType")
	}

	authType := conf[0]
	clientAuth := ClientAuth{}
	clientAuth.SetAuthType(authType)
	
	switch authType {
	case NoAuth:
	case BasicAuth:
		if len(conf) < 3 {
			panic(fmt.Sprintf("too few arguments for auth %s", authType))
		}
		clientAuth.SetUsername(conf[1])
		clientAuth.SetPassword(conf[2])
	case OAuth:
		if len(conf) < 5 {
			panic(fmt.Sprintf("too few arguments for auth %s", authType))
		}
		clientAuth.SetUsername(conf[1])
		clientAuth.SetPassword(conf[2])
		clientAuth.SetToken(conf[3])
		clientAuth.SetRefreshToken(conf[4])
	default:
		panic(fmt.Sprintf("unknown auth %s", authType))
	}
	return &clientAuth
}

// NewClient configures and returns a new Client
func NewClient(serverUrl string, auth *ClientAuth) Client {
	parsedUrl, err := url.Parse(serverUrl)
	if err != nil {
		panic(err)
	}

	if parsedUrl.Scheme == "" {
		panic("serverUrl has no schema")
	} 
	// FIXME: anything else to handle?

	return Client{
		rootUrl: parsedUrl,
		auth: auth,
	}
}

func (a *ClientAuth) SetAuthType(authType string) {
	authType = strings.Title(authType)
	if authType != NoAuth && authType != BasicAuth && authType != OAuth {
		err := fmt.Sprintf("authType: bad value, expected %q, %q or %q", 
			NoAuth, BasicAuth, OAuth)
		panic(err)
	}
	a.authType = authType
}

func (a *ClientAuth) GetAuthType() string {
	return a.authType
}

func (a *ClientAuth) SetUsername(username string) {
	username = strings.Trim(username, " ")
	if len(username) == 0 {
		panic("username is empty")
	}
	a.username = username
}

func (a *ClientAuth) GetUsername() string {
	return a.username
}

func (a *ClientAuth) SetPassword(password string) {
	password = strings.Trim(password, " ")
	if len(password) == 0 {
		panic("password is empty")
	}
	a.password = password
}

func (a *ClientAuth) SetToken(token string) {
	token = strings.Trim(token, " ")
	if len(token) == 0 {
		panic("token is empty")
	}
	a.token = token
}

func (a *ClientAuth) GetToken() string {
	return a.token
}

func (a *ClientAuth) SetRefreshToken(token string) {
	token = strings.Trim(token, " ")
	if len(token) == 0 {
		panic("refreshToken is empty")
	}
	a.refreshToken = token
}

func (a *ClientAuth) GetRefreshToken() string {
	return a.refreshToken
}

// call performs a HTTP call using Client configuration
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
