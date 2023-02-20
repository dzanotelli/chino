package custodia

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dzanotelli/chino/common"
)

type Repository struct {
	RepositoryId string
	Description string
	IsActive bool `json:"is_active"`
}

// API methods to handle Repository

func (ca *CustodiaAPIv1) CreateRepository(description string) (Repository, 
	error) {
	// FIXME
	return Repository{}, nil
}

func (ca *CustodiaAPIv1) GetRepository(id string) (*Repository, error) {
	if !common.IsValidUUID(id) {
		return nil, errors.New("id is not a valid UUID: " + id)
	}

	url := fmt.Sprintf("/repositories/%s", id)
	resp, err := ca.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// JSON: unmarshal resp content
	repository := Repository{RepositoryId: id}
	if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
		return nil, err
	}
	return &repository, nil
}