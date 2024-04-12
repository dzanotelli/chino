package custodia

import (
	"encoding/json"
	"fmt"
)

// Define `grant_type` enum
type GrantType int

const (
	GrantAuthorizationCode GrantType = iota + 1
	GrantPassword
)

func (gt GrantType) Choices() []string {
	return []string{"authorization-code", "password"}
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