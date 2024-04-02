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
	ClientTypePublic ClientType = iota +1
	ClientTypeConfidential
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
	AppSecret string `json:"app_secret,omitempty"`
	GrantType GrantType `json:"grant_type,omitempty"`
	AppName string `json:"app_name"`
	AppId string `json:"app_id"`
}

// ApplicationEnvelope: used to unmarshal the CRU responses
type ApplicationEnvelope struct {
	Application *Application `json:"application"`
}

// ApplicationsEnvelope: used to unmarshal the L response
type ApplicationsEnvelope struct {
	Applications []Application `json:"applications"`
}

// FIXME: missing funcs to marshal/unmarshal


// [C]reate a new application
func (ca *CustodiaAPIv1) CreateApplication(name string, grantType GrantType,
	clientType ClientType) (*Application, error) {

	// FIXME

	return nil, nil
}
