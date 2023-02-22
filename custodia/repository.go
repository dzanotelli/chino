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
	InsertDate string `json:"insert_date"`
	LastUpdate string `json:"last_update"`
	IsActive bool `json:"is_active"`
}

type RepositoryEnvelope struct {
	Repository *Repository
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
	resp, err := ca.Call("POST", url, string(data))
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	repoEnvelope := RepositoryEnvelope{}
	if err := json.Unmarshal([]byte(resp), &repoEnvelope); err != nil {
		return nil, err
	}

	return repoEnvelope.Repository, nil
}

// [R]ead an existent repository
func (ca *CustodiaAPIv1) GetRepository(id string) (*Repository, error) {
	if !common.IsValidUUID(id) {
		return nil, errors.New("id is not a valid UUID: " + id)
	}

	url := fmt.Sprintf("/repositories/%s", id)
	resp, err := ca.Call("GET", url)
	if err != nil {
		return nil, err
	}
	
	// JSON: unmarshal resp content
	repoEnvelope := RepositoryEnvelope{}
	if err := json.Unmarshal([]byte(resp), &repoEnvelope); err != nil {
		return nil, err
	}
	return repoEnvelope.Repository, nil
}

// [U]pdate an existent repository
func (ca *CustodiaAPIv1) UpdateRepository(repository *Repository, 
	description string, isActive bool) (*Repository, error) {	
	url := fmt.Sprintf("/repositories/%s", (*repository).RepositoryId)

	// get a copy and update the values, so we can easily marshal it
	copy := *repository	
	copy.RepositoryId = ""  // we must not send this
	copy.Description = description
	copy.IsActive = isActive
	data, err := json.Marshal(copy)
	if err != nil {
		return nil, err
	}
	resp, err := ca.Call("PUT", url, string(data))
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content overwriting the old repository
	repoEnvelope := RepositoryEnvelope{}
	if err := json.Unmarshal([]byte(resp), &repoEnvelope); err != nil {
		return nil, err
	}
	return repoEnvelope.Repository, nil
}

// // [D]elete an existent repository
// func (ca *CustodiaAPIv1) DeleteRepository(repository *Repository) (error) {
// 	url := fmt.Sprintf("/repositories/%s", repository.RepositoryId)
// 	resp, err := ca.Delete(url)
// 	if err != nil {
// 		return err
// 	}
// 	resp.Body.Close()	
// 	return nil
// }

// // [L]ist all the repositories
// func (ca *CustodiaAPIv1) ListRepositories() ([]*Repository, error) {
// 	// FIXME
// 	return nil, nil
// }

