// OAuth is performed against custodia system. Then, it can be used in other
// packages such as Consenta (that is, you use this API to get the access and
// refresh tokens, and then you use the `ClientAuth.AuthType=Bearer` with any
// client).

package custodia

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dzanotelli/chino/common"
)

// Define `grant_type` enum
type GrantType int

const (
	GrantAuthorizationCode GrantType = iota + 1
	GrantPassword
	GrantRefreshToken
)

func (gt GrantType) Choices() []string {
	return []string{"authorization-code", "password", "refresh_token"}
}

func (gt GrantType) String() string {
	return gt.Choices()[gt-1]
}

func (gt GrantType) MarshalJSON() ([]byte, error) {
    return json.Marshal(gt.String())
}

func (gt *GrantType) UnmarshalJSON(data []byte) (err error) {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	intValue := indexOf(value, gt.Choices()) + 1   // enum starts from 1
	if intValue < 1 {
		return fmt.Errorf("GrantType: received unknown value '%v'", value)
	}

	*gt = GrantType(intValue)
	return nil
}

// Define `client_type` enum
type ClientType int

const (
	ClientPublic ClientType = iota +1
	ClientConfidential
)

func (ct ClientType) Choices() ([]string) {
	return []string{"public", "confidential"}
}

func (ct ClientType) String() string {
	return ct.Choices()[ct-1]
}

func (ct ClientType) EnumIndex() int {
	return int(ct)
}

func (ct ClientType) MarshalJSON() ([]byte, error) {
    return json.Marshal(ct.String())
}

func (ct *ClientType) UnmarshalJSON(data []byte) (err error) {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	intValue := indexOf(value, ct.Choices()) + 1   // enum starts from 1
	if intValue < 1 {
		return fmt.Errorf("ClientType: received unknown value '%v'", value)
	}

	*ct = ClientType(intValue)
	return nil
}

// helper struct to performs OAuth calls
type oauthResponseData struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int `json:"expires_in"`
	Scope string `json:"scope"`
}

type TokenInfo struct {
	Active bool
	Scope string
	Expiration int `json:"exp"`
	ApplicationId string `json:"client_id"`
	Username string
}


// Application represent an application stored in Custodia
type Application struct {
	Secret string `json:"app_secret,omitempty"`
	GrantType GrantType `json:"grant_type"`
	ClientType ClientType `json:"client_type"`
	Name string `json:"app_name"`
	Id string `json:"app_id,omitempty"`
	RedirectUrl string `json:"redirect_url,omitempty"`
}

// ApplicationEnvelope: used to unmarshal the CRU responses
type ApplicationEnvelope struct {
	Application *Application `json:"application"`
}

// ApplicationsEnvelope: used to unmarshal the L response
type ApplicationsEnvelope struct {
	Applications []Application `json:"applications"`
}

// [C]reate a new application
func (ca *CustodiaAPIv1) CreateApplication(name string, grantType GrantType,
	clientType ClientType, redirectUrl string) (*Application, error) {

	if grantType == GrantPassword && redirectUrl != "" {
		err := fmt.Errorf("redirectUrl must be empty when grantType is '%s'",
			grantType)
		return nil, err
	}

	application := Application{
		Name: name,
		GrantType: grantType,
		ClientType: clientType,
		RedirectUrl: redirectUrl,
	}
	url := "/auth/applications"
	params := map[string]interface{}{"data": application}
	resp, err := ca.Call("POST", url, params)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	appEnvelope := ApplicationEnvelope{}
	if err := json.Unmarshal([]byte(resp), &appEnvelope); err != nil {
		return nil, err
	}

	return appEnvelope.Application, nil
}

// [R]ead an existent application
func (ca *CustodiaAPIv1) ReadApplication(id string) (*Application, error) {
	url := fmt.Sprintf("/auth/applications/%s", id)
	resp, err := ca.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	appEnvelope := ApplicationEnvelope{}
	if err := json.Unmarshal([]byte(resp), &appEnvelope); err != nil {
		return nil, err
	}
	return appEnvelope.Application, nil
}

// [U]pdate an existent application
func (ca *CustodiaAPIv1) UpdateApplication(id string, name string,
	grantType GrantType, clientType ClientType, redirectUrl string) (
		*Application, error) {

	if grantType == GrantPassword && redirectUrl != "" {
		err := fmt.Errorf("redirectUrl must be empty when grantType is '%s'",
			grantType)
		return nil, err
	}

	application := Application{
		Name: name,
		GrantType: grantType,
		ClientType: clientType,
		RedirectUrl: redirectUrl,
	}
	url := fmt.Sprintf("/auth/applications/%s", id)
	params := map[string]interface{}{"data": application}
	resp, err := ca.Call("PUT", url, params)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content and return a fresh application instance
	appEnvelope := ApplicationEnvelope{}
	if err := json.Unmarshal([]byte(resp), &appEnvelope); err != nil {
		return nil, err
	}

	return appEnvelope.Application, nil
}

// [D]elete an existent application
func (ca *CustodiaAPIv1) DeleteApplication(id string) (error) {
	url := fmt.Sprintf("/auth/applications/%s", id)
	_, err := ca.Call("DELETE", url, nil)
	if err != nil {
		return err
	}
	return nil
}

// [L]ist all the applications
func (ca *CustodiaAPIv1) ListApplications() ([]*Application, error) {
	url := "/auth/applications"
	resp, err := ca.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	appsEnvelope := ApplicationsEnvelope{}
	if err := json.Unmarshal([]byte(resp), &appsEnvelope); err != nil {
		return nil, err
	}

	result := []*Application{}
	for _, app := range appsEnvelope.Applications {
		result = append(result, &app)
	}

	return result, nil
}


// Login a user
// This is the oauth flow of type "password" where via api call a client
// immediately get access_token and refresh token
func (ca *CustodiaAPIv1) LoginUser(username string, password string,
	application Application) (*common.ClientAuth, error) {
	url := "/auth/token"
	auth := *common.NewClientAuth()    // defaults to no auth

	data := map[string]string{
		"grant_type": GrantPassword.String(),
		"username": username,
		"password": password,
		"client_id": application.Id,
	}

	// when client is not public, we need to set the application secret as well
	if application.ClientType == ClientConfidential {
		data["client_secret"] = application.Secret
	}

	params := map[string]interface{}{
		"data": data,
		"contentType": "multipart/form-data",
	}
	resp, err := ca.Call("POST", url, params)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	respData := oauthResponseData{}
	if err := json.Unmarshal([]byte(resp), &respData); err != nil {
		return nil, err
	}

	// compute the unixtime of expiration
	expiration := int(time.Now().Unix()) + respData.ExpiresIn

	auth.SetOAuth(respData.AccessToken, expiration, respData.RefreshToken)

	return &auth, nil
}

// LoginCode
// Get the access_token and refresh_token using the authorization code
// retrieved in the previous steps of the authorization code flow.
func (ca *CustodiaAPIv1) LoginAuthCode(code string,
	application Application) (*common.ClientAuth, error) {
	url := "/auth/token"
	auth := *common.NewClientAuth()    // defaults to no auth

	data := map[string]string{
		"grant_type": GrantAuthorizationCode.String(),
		"code": code,
		"redirect_uri": application.RedirectUrl,
		"client_id": application.Id,
		"scope": "read write",
	}
	if application.ClientType == ClientConfidential {
		data["client_secret"] = application.Secret
	}

	params := map[string]interface{}{
		"data": data,
		"contentType": "multipart/form-data",
	}
	resp, err := ca.Call("POST", url, params)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	respData := oauthResponseData{}
	if err := json.Unmarshal([]byte(resp), &respData); err != nil {
		return nil, err
	}

	// compute the unixtime of expiration
	expiration := int(time.Now().Unix()) + respData.ExpiresIn

	auth.SetOAuth(respData.AccessToken, expiration, respData.RefreshToken)

	return &auth, nil
}

// Refresh the access token
func (ca *CustodiaAPIv1) RefreshToken(auth common.ClientAuth,
	application Application) (*common.ClientAuth, error) {
	url := "/auth/refresh"

	data := map[string]string{
		"refresh_token": auth.GetRefreshToken(),
		"grant_type": GrantRefreshToken.String(),
		"client_id": application.Id,
		"scope": "read write",
	}

	// when client is not public, we need to set the application secret as well
	if application.ClientType == ClientConfidential {
		data["client_secret"] = application.Secret
	}

	params := map[string]interface{}{
		"data": data,
		"contentType": "multipart/form-data",
	}

	resp, err := ca.Call("POST", url, params)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	respData := oauthResponseData{}
	if err := json.Unmarshal([]byte(resp), &respData); err != nil {
		return nil, err
	}

	// compute the unixtime of expiration
	expiration := int(time.Now().Unix()) + respData.ExpiresIn

	auth.SetAccessToken(respData.AccessToken)
	auth.SetAccessTokenExpire(expiration)
	auth.SetRefreshToken(respData.RefreshToken)

	return &auth, nil
}

// Revoke an existing token
func (ca *CustodiaAPIv1) RevokeToken(auth common.ClientAuth,
	application Application) error {
	url := "/auth/revoke_token"

	data := map[string]string{
		"token": auth.GetAccessToken(),
		"client_id": application.Id,
		"client_secret": application.Secret,
	}

	params := map[string]interface{}{
		"data": data,
		"contentType": "multipart/form-data",
	}
	_, err := ca.Call("POST", url, params)
	if err != nil {
		return err
	}

	return nil
}

// Introspect token
// wants Basic auth using application id/secret
func (ca *CustodiaAPIv1) IntrospectToken(token string) (*TokenInfo, error) {
	url := "/auth/revoke_token"

	data := map[string]string{
		"token": token,
	}

	params := map[string]interface{}{
		"data": data,
		"contentType": "application/x-www-form-urlencoded",
	}
	resp, err := ca.Call("POST", url, params)
	if err != nil {
		return nil, err
	}
	tokenInfo := TokenInfo{}
	if err := json.Unmarshal([]byte(resp), &tokenInfo); err != nil {
		return nil, err
	}

	return &tokenInfo, nil
}

// User info
func (ca *CustodiaAPIv1) UserInfo(schema *UserSchema) (*User, error) {
	url := "/users/me"

	resp, err := ca.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// JSON: unmarshal resp content
	userEnvelope := UserEnvelope{}
	if err := json.Unmarshal([]byte(resp), &userEnvelope); err != nil {
		return nil, err
	}

	if (schema == nil) {
		// get the user schema
		userSchema, err := ca.ReadUserSchema(userEnvelope.User.UserSchemaId)
		if err != nil {
			return nil, err
		}
		schema = userSchema
	}

	// convert values to concrete types
	converted, ee := convertData(userEnvelope.User.Attributes, schema)
	if len(ee) > 0 {
		err := fmt.Errorf("conversion errors: %w", errors.Join(ee...))
		return userEnvelope.User, err
	}

	// all good, assign the new content to doc and return it
	userEnvelope.User.Attributes = converted
	return userEnvelope.User, nil
}
