package custodia

import (
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

func (*CustodiaAPI) CreateRepository(description string) (*Repository, error) {
	// FIXME
	return &Repository{}, nil
}

func (capi *CustodiaAPI) GetRepository(id string) (*Repository, error) {
	if IsValidUUID(id) == false {
		return &Repository{}, errors.New("id is not a valid UUID: " + id)
	}

	resp, err := capi.client.Get("/api/v1/repositories")
	if err != nil {
		return &Repository{}, err
	}

	// FIXME: unmarshal resp content

	// return the new Repository object
	repository := Repository{
		repository_id: "FIXME",
		description: "FIXME",
	}
	return &repository, nil
}