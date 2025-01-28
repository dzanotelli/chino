// Package common provides a basic HTTP client to perform HTTP calls.
// It supports Basic and OAuth authentication methods.
package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

const userAgent = "golang/chino-" + Version

type AuthType int
const (
	NoAuth AuthType = iota + 1
	CustomerAuth
	UserAuth
	ApplicationAuth
)
func (at AuthType) Choices() []string {
	return []string{"Customer", "User", "Application"}
}

func (at AuthType) String() string {
	return at.Choices()[at-1]
}

// ClientAuth keeps the authentication details - Basic vs Bearer (OAuth)
type ClientAuth struct {
	currentAuthType AuthType
	prevAuthType AuthType
	customerId string          // only for Customer auth
	customerKey string		   // only for Customer auth
	accessToken string         // only for OAuth
	accessTokenExpire int      // only for OAuth
	refreshToken string        // only for OAuth
	applicationId string       // only for Application auth
	applicationSecret string   // only for Application auth
}

// Client holds the configuration (url, auth) and wraps http Requests
type Client struct {
	rootUrl *url.URL
	auth *ClientAuth
}

// NewClientAuth returns a new ClientAuth with auth set to NoAuth
func NewClientAuth(data map[string]interface{}) *ClientAuth {
	ca := &ClientAuth{}
	ca.SetNoAuth()
	ca.Update(data)
	return ca
}

// Set ClientAuth authType to NoAuth removing other attributes
func (ca *ClientAuth) SetNoAuth() {
	ca.currentAuthType = NoAuth
	ca.prevAuthType = NoAuth
	ca.customerId = ""
	ca.customerKey = ""
	ca.accessToken = ""
	ca.accessTokenExpire = 0
	ca.refreshToken = ""
	ca.applicationId = ""
	ca.applicationSecret = ""
}

func (ca *ClientAuth) SwitchTo(authType AuthType) {
	ca.prevAuthType = ca.currentAuthType
	ca.currentAuthType = authType
}

func (ca *ClientAuth) SwitchBack() {
	ca.currentAuthType, ca.prevAuthType = ca.prevAuthType, ca.currentAuthType
}


func (ca *ClientAuth) Update(data map[string]interface{}) error {
	for key, value := range data {
		switch key {
		case "customerId":
			if customerId, ok := value.(string); ok { //FIXME: use UUID type
				if !IsValidUUID(customerId) {
					return errors.New("customerId must be a valid UUID")
				}
				ca.customerId = customerId
			}
		case "customerKey":
			if customerKey, ok := value.(string); ok { //FIXME: use UUID type
				if !IsValidUUID(customerKey) {
					return errors.New("customerKey must be a valid UUID")
				}
				ca.customerKey = customerKey
			}
		case "accessToken":
			if accessToken, ok := value.(string); ok {
				ca.accessToken = accessToken
			}
		case "accessTokenExpire":
			if accessTokenExpire, ok := value.(int); ok {
				ca.accessTokenExpire = accessTokenExpire
			}
		case "refreshToken":
			if refreshToken, ok := value.(string); ok {
				ca.refreshToken = refreshToken
			}
		case "applicationId":
			if applicationId, ok := value.(string); ok {
				ca.applicationId = applicationId
			}
		case "applicationSecret":
			if applicationSecret, ok := value.(string); ok {
				ca.applicationSecret = applicationSecret
			}
		default:
			return errors.New("unknown attribute: " + key)
		}
	}
	return nil
}

// Set ClientAuth authType to BasicAuth removing tokens
func (ca *ClientAuth) SetCustomerAuth(id, key string) error {
	if !IsValidUUID(id) {
		return errors.New("customerId must be a valid UUID")
	} else if !IsValidUUID(key) {
		return errors.New("customerKey must be a valid UUID")
	}

	ca.customerId = id
	ca.customerKey = key
	return nil
}

func (ca *ClientAuth) SetUserAuth(accessToken string, tokenExpire int,
	refreshToken string) error {
	ca.accessToken = accessToken
	ca.accessTokenExpire = tokenExpire
	ca.refreshToken = refreshToken
	return nil
}

func (ca *ClientAuth) SetApplicationAuth(id, secret string) error {
	ca.applicationId = id
	ca.applicationSecret = secret
	return nil
}


func (ca *ClientAuth) GetAuthType() AuthType {
	return ca.currentAuthType
}

func (ca *ClientAuth) GetCustomerId() string {
	return ca.customerId
}

func (ca *ClientAuth) GetAccessToken() string {
	return ca.accessToken
}

func (ca *ClientAuth) GetAccessTokenExpire() int {
	return ca.accessTokenExpire
}

func (ca *ClientAuth) GetRefreshToken() string {
	return ca.refreshToken
}

func (ca *ClientAuth) GetApplicationId() string {
	return ca.applicationId
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

	if auth == nil {
		auth = &ClientAuth{}
	}

	return &Client{
		rootUrl: parsedUrl,
		auth: auth,
	}
}

func (c *Client) GetAuth() *ClientAuth {
	return c.auth
}

// Performs a HTTP Call using Client configuration
// `params` map may contain the following keys:
//  - "_data" is the data to be sent in the request body.
// 		- For "application/json" content type, data is serialized as JSON
// 		- For "application/x-www-form-urlencoded" content type, data is
//  	  serialized as a form
//  - "_rawResponse" is a boolean that indicates if the response should be
//    returned as is
// 	- "Content-Type" is the content type of the request, it defaults to
// 	  "application/json"
//  - any other key-value pairs which doesn't start wiht a '_' are added
//   as request headers
func (c *Client) Call(method, path string, params map[string]interface{}) (
	*http.Response,	error) {
	fullPath := strings.TrimRight(c.rootUrl.String(), "/")
	fullPath += "/" + strings.TrimLeft(path, "/")

	var req *http.Request
	var err error

	switch method {
	case "GET", "DELETE":
		req, err = http.NewRequest(method, fullPath, nil)
	case "POST", "PUT", "PATCH":
		contentType, ok := params["Content-Type"].(string)
		if !ok {
			contentType = "application/json"
		}

		switch contentType {
		case "application/json":
			data := params["_data"]
			jsonData, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}
			req, err = http.NewRequest(method, fullPath,
				bytes.NewBuffer(jsonData))
			if err != nil {
				return nil, err
			}
		case "application/x-www-form-urlencoded":
			values := url.Values{}
			formData, ok := params["_data"].(map[string]string)
			if !ok {
				return nil, errors.New("_data must be a map[string]string")
			}
			for key, value := range formData {
				values.Add(key, value)
			}
			req, err = http.NewRequest("POST", fullPath,
				strings.NewReader(values.Encode()))
		case "multipart/form-data":
			data, ok := params["_data"].(map[string]string)
			if !ok {
				return nil, errors.New("_data must be a map[string]string")
			}
			body := &bytes.Buffer{}
			w := multipart.NewWriter(body)
			for key, value := range data {
				fw, err := w.CreateFormField(key)
				if err != nil {
					return nil, err
				}
				_, err = io.WriteString(fw, value)
				if err != nil {
					return nil, err
				}
			}
			w.Close()
			req, err = http.NewRequest(method, fullPath, body)
			if err != nil {
				return nil, err
			}
		case "application/octet-stream":
			data, ok := params["_data"].([]byte)
			if !ok {
				return nil, errors.New("_data must be []byte")
			}
			req, err = http.NewRequest(method, fullPath,
				bytes.NewBuffer(data))
		default:
			panic(fmt.Sprintf("unsupported content type %q", contentType))
		}

		// this may have changed before. reassign it to params since later
		// we will add all the keys but body as headers to the request
		if params == nil {
			params = map[string]interface{}{}
		}
		params["Content-Type"] = contentType
	default:
		err = fmt.Errorf("unsupported HTTP method %q", method)
	}

	if err != nil {
		return nil, err
	}

	// default headers (we use `Set`, not `Add` cos we want a single value)
	req.Header.Set("User-Agent", userAgent)

	isRawResponse, ok := params["_rawResponse"].(bool)
	if !ok {
		isRawResponse = false
	}
	if isRawResponse {
		req.Header.Set("Accept", "*/*")
	} else {
		req.Header.Set("Accept", "application/json")
	}

	// add all the params as headers which doesn't start with "_"
	for key, value := range params {
		if strings.HasPrefix(key, "_") {
			continue
		}
		req.Header.Set(key, fmt.Sprint(value))
	}

	// handle auth
	switch c.auth.currentAuthType {
	case NoAuth:
		// do nothing
	case CustomerAuth:
		req.SetBasicAuth(c.auth.customerId, c.auth.customerKey)
	case UserAuth:
		bearer := "Bearer: " + c.auth.accessToken
		req.Header.Add("Authorization", bearer)
	case ApplicationAuth:
		req.SetBasicAuth(c.auth.applicationId, c.auth.applicationSecret)
	default:
		panic(fmt.Sprintf("Unsupported auth type %q", c.auth.currentAuthType))
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
	return c.Call("GET", path, nil)
}

// Post wraps call to perform a HTTP POST call
func (c *Client) Post(path string, params map[string]interface{}) (
	*http.Response, error) {
	return c.Call("POST", path, params)
}

// Put wraps call to perform a HTTP PUT call
func (c *Client) Put(path string, params map[string]interface{}) (
	*http.Response, error) {
	return c.Call("PUT", path, params)
}

// Patch wraps call to perform a HTTP PATCH call
func (c *Client) Patch(path string, params map[string]interface{}) (
	*http.Response, error) {
	return c.Call("PATCH", path, params)
}

// Delete wraps call to perform a HTTP DELETE call
func (c *Client) Delete(path string) (*http.Response, error) {
	return c.Call("DELETE", path, nil)
}
