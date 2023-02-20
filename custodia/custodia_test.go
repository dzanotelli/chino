package custodia

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dzanotelli/chino/common"
	"github.com/google/uuid"
)


func TestRepositoryCRUD(t *testing.T) {
	// init stuff
	dummyRepository := Repository{
		RepositoryId: uuid.New().String(),
		Description: "unittest",
		IsActive: false,
	}

	// mock calls
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/repositories" && r.Method == "POST" {
			// test CREATE
			w.WriteHeader(http.StatusOK)
			out := fmt.Sprintf(`{"repository_id": "%s", "description": ` +
							   `"%s", "is_active": %v}`, 
							   dummyRepository.RepositoryId,
							   dummyRepository.Description,
							   dummyRepository.IsActive)
			w.Write([]byte(out))
		} else if r.URL.Path == fmt.Sprintf("/api/v1/repositories/%s", 
			dummyRepository.RepositoryId) {
			// test READ
			dummyRepository.IsActive = true  // turn to true
			w.WriteHeader(http.StatusOK)
			out := fmt.Sprintf(`{"description": "%v", "is_active": %v}`, 
							  dummyRepository.Description,
							  dummyRepository.IsActive)
			w.Write([]byte(out))
		} else {
			err := `{"result": "error", "result_code": 404, "data": null, `
			err += `"message": "Resource not found (you may have a '/' at `
			err += `the end)"}`
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err))
		}
	}
	server := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer server.Close()

	auth := common.NewClientAuth()  // auth is tested elsewhere
	client := common.NewClient(server.URL, auth)
	custodia := NewCustodiaAPIv1(client)

	// test CREATE
	repo, err := custodia.CreateRepository("unittest", false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if repo != nil {
		if (*repo).RepositoryId != dummyRepository.RepositoryId {
			t.Errorf("bad RepositoryId, got: %v want: %v", 
					 repo.RepositoryId, dummyRepository.RepositoryId)
		}
		if (*repo).Description != dummyRepository.Description {
			t.Errorf("bad Description, got: %v want: %s", 
					 repo.Description,
					dummyRepository.Description)
		}

		if (*repo).IsActive != false {
			t.Errorf("bad isActive, got: %v want: false", (*repo).IsActive)
		}
	} else {
		t.Errorf("unexpected: both repository and error are nil!")
	}

	// test READ
	repo, err = custodia.GetRepository(dummyRepository.RepositoryId)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if repo != nil {
		if (*repo).RepositoryId != dummyRepository.RepositoryId {
			t.Errorf("bad RepositoryId, got: %v want: %v", 
					 repo.RepositoryId, dummyRepository.RepositoryId)
		}
		if (*repo).Description != dummyRepository.Description {
			t.Errorf("bad Description, got: %v want: %s", 
					 repo.Description,
					dummyRepository.Description)
		}

		if (*repo).IsActive != true {
			t.Errorf("bad isActive, got: %v want: true", (*repo).IsActive)
		}
	} else {
		t.Errorf("unexpected: both repository and error are nil!")
	}
}