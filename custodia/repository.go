package custodia

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dzanotelli/chino/common"
)

// Repository represent a repository stored in Custodia
type Repository struct {
	RepositoryId string `json:"repository_id,omitempty"`
	Description string
	IsActive bool `json:"is_active"`
}

// [C]reate a new repository
func (ca *CustodiaAPIv1) CreateRepository(description string, isActive bool) (
	*Repository, error) {
	repository := Repository{Description: description, IsActive: isActive}
	url := "/repositories"
	data, err := json.Marshal(repository)
	if err != nil {
		return nil, err
	}
	resp, err := ca.Post(url, string(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// JSON: unmarshal resp content
	if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
		return nil, err
	}
	return &repository, nil
}

// [R]ead an existent repository
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

// [U]pdate an existent repository
func (ca *CustodiaAPIv1) UpdateRepository(repository *Repository, 
	description string, isActive bool) (*Repository, error) {
	// FIXME
	return repository, nil
}

// [D]elete an existent repository
func (ca *CustodiaAPIv1) DeleteRepository(repository *Repository) (error) {
	// FIXME
	return nil
}

// [L]ist all the repositories
func (ca *CustodiaAPIv1) ListRepositories() ([]*Repository, error) {
	// FIXME
	return nil, nil
}

