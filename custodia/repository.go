package custodia

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dzanotelli/chino/common"
	"github.com/simplereach/timeutils"
)

// Repository represent a repository stored in Custodia
type Repository struct {
	RepositoryId string `json:"repository_id,omitempty"`
	Description string `json:"description"`
	InsertDate timeutils.Time `json:"insert_date"`
	LastUpdate timeutils.Time `json:"last_update"`
	IsActive bool `json:"is_active"`
}
// RepositoryEnvelope: used the unmarshal the CRU responses
type RepositoryEnvelope struct {
	Repository *Repository `json:"repository"`
}

// RepositoriesEnvelope: used the unmarshal the L response
type RepositoriesEnvelope struct {
	Repositories []Repository `json:"repositories"`
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
func (ca *CustodiaAPIv1) ReadRepository(id string) (*Repository, error) {
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

// [D]elete an existent repository
// if force=true recursively deletes all the repository content, else the
// repository is just deactivated
func (ca *CustodiaAPIv1) DeleteRepository(repository *Repository, force bool) (
	error) {
	url := fmt.Sprintf("/repositories/%s", repository.RepositoryId)
	url += fmt.Sprintf("?force=%v", force)

	_, err := ca.Call("DELETE", url)
	if err != nil {
		return err
	}
	return nil
}

// [L]ist all the repositories
func (ca *CustodiaAPIv1) ListRepositories() ([]*Repository, error) {
	resp, err := ca.Call("GET", "/repositories")
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	reposEnvelope := RepositoriesEnvelope{}
	if err := json.Unmarshal([]byte(resp), &reposEnvelope); err != nil {
		return nil, err
	}

	result := []*Repository{}
	for _, repo := range reposEnvelope.Repositories {
		result = append(result, &repo)
	}
	return result, nil
}

