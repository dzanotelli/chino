package client

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

// Client holds the configuration (url, auth) and wraps http Requests
type Client struct {
	rootUrl *url.URL
	username string
	password string
}

// NewClients configures and returns a new Client
func NewClient(serverUrl, username, password string) Client {
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
		username: username,
		password: password,
	}
}

// call performs a HTTP call using Client configuration
func (c *Client) call(method, path string, data ...string) *http.Response {
	url := c.rootUrl.String() + path  //FIXME join strings
	var req *http.Request
	var err error
	
	// FIXME auth missing

	switch method {
	case "GET", "DELETE":
		req, err = http.NewRequest(method, url, nil)
	case "POST", "PUT", "PATCH":
		if len(data) > 0 {
			var jsonStr = []byte(data[0])
		} else {
			var jsonStr = []byte("")
		}
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
	default:
		err := fmt.Errorf("unsupported HTTP method %q", method)
		panic(err)
	}

	// FIXME: dunno what to do, just panic for the time being
	if err != nil {
		panic(err)
	}

	// perform the call
	client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

	return resp
}

// Get wraps call to perform a HTTP GET call
func (c *Client) Get(path string) *http.Response {
	return c.call("GET", path)
}

// Post wraps call to perform a HTTP POST call
func (c *Client) Post(path, payload string) *http.Response {
	return c.call("POST", path, payload)
}

// Put wraps call to perform a HTTP PUT call
func (c *Client) Put(path, payload string) *http.Response {
	return c.call("PUT", path, payload)
}

// Patch wraps call to perform a HTTP PATCH call
func (c *Client) Patch(path, payload string) *http.Response {
	return c.call("PATCH", path, payload)
}

// Delete wraps call to perform a HTTP DELETE call
func (c *Client) Delete(path string) *http.Response {
	return c.call("DELETE", path)
}
