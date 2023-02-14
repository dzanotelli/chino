package custodia

import (
	"encoding/json"
	"errors"

	"github.com/dzanotelli/chino/common"
)

type CustodiaAPI struct {
	client common.Client
}

// NewCustodiaAPI returns a new CustodiaAPI object to interact
// with the Custodia APIs
func NewCustodiaAPI(client common.Client) *CustodiaAPI {
	capi := &CustodiaAPI{}
	capi.client = client
	return capi
}

func (*CustodiaAPI) CreateRepository(description string) (Repository, error) {
	// FIXME
	return Repository{}, nil
}

func (capi *CustodiaAPI) GetRepository(id string) (Repository, error) {
	if common.IsValidUUID(id) == false {
		return Repository{}, errors.New("id is not a valid UUID: " + id)
	}

	resp, err := capi.client.Get("/api/v1/repositories")
	if err != nil {
		return Repository{}, err
	}
	
	// JSON: unmarshal resp content
	repository := Repository{}
	if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
		return Repository{}, err
	}
	return repository, nil
}