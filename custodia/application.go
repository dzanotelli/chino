// OAuth is performed against custodia system. Then, it can be used in other
// packages such as Consenta (that is, you use this API to get the access and
// refresh tokens, and then you use the `ClientAuth.AuthType=Bearer` with any
// client).

package custodia

import (
	"encoding/json"
	"fmt"

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

// helper structs to performs OAuth calls (marshal, unmarshal)
type oauthRequestData struct {
	GrantType GrantType `json:"grant_type"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	ClientId string `json:"client_id"`
	ClientSecret string `json:"client_secret,omitempty"`
	Scope string `json:"scope"`
}
type oauthResponseData struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int `json:"expires_in"`
	Scope string `json:"scope"`
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
	data, err := json.Marshal(application)
	if err != nil {
		return nil, err
	}
	resp, err := ca.Call("POST", url, string(data))
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
	resp, err := ca.Call("GET", url)
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
	data, err := json.Marshal(application)
	if err != nil {
		return nil, err
	}
	resp, err := ca.Call("PUT", url, string(data))
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
	_, err := ca.Call("DELETE", url)
	if err != nil {
		return err
	}
	return nil
}

// [L]ist all the applications
func (ca *CustodiaAPIv1) ListApplications() ([]*Application, error) {
	url := "/auth/applications"
	resp, err := ca.Call("GET", url)
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
func (ca *CustodiaAPIv1) LoginUser(username string, password string,
	application Application) (common.ClientAuth, error) {
	url := "/auth/token"
	auth := *common.NewClientAuth()    // defaults to no auth

	data := oauthRequestData {
		GrantType: GrantPassword,
		Username: username,
		Password: password,
		ClientId: application.Id,
	}
	// when client is not public, we need to set the application secret as well
	if application.ClientType == ClientConfidential {
		data.ClientSecret = application.Secret
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return auth, err
	}

	resp, err := ca.Call("POST", url, string(payload))
	if err != nil {
		return auth, err
	}

	// JSON: unmarshal resp content
	respData := oauthResponseData{}
	if err := json.Unmarshal([]byte(resp), &respData); err != nil {
		return auth, err
	}

	auth.SetOAuth(username, password, respData.AccessToken, respData.ExpiresIn,
		respData.RefreshToken)

	return auth, nil
}

// Refresh the access token
func (ca *CustodiaAPIv1) RefreshToken(auth common.ClientAuth,
	application Application) (common.ClientAuth, error) {
	url := "/auth/token"

	data := oauthRequestData {
		GrantType: GrantRefreshToken,
		ClientId: application.Id,
		Scope: "read write",
	}
	// when client is not public, we need to set the application secret as well
	if application.ClientType == ClientConfidential {
		data.ClientSecret = application.Secret
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return auth, err
	}

	resp, err := ca.Call("POST", url, string(payload))
	if err != nil {
		return auth, err
	}

	// JSON: unmarshal resp content
	respData := oauthResponseData{}
	if err := json.Unmarshal([]byte(resp), &respData); err != nil {
		return auth, err
	}

	auth.SetAccessToken(respData.AccessToken)
	auth.SetAccessTokenExpire(respData.ExpiresIn)
	auth.SetRefreshToken(respData.RefreshToken)

	return auth, nil
}