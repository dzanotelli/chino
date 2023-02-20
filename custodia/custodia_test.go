package storage

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dzanotelli/chino/common"
	"github.com/google/uuid"
)


func TestRepositoryCRUD(t *testing.T) {
	// mock calls
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/repositories/" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"description": "unittest", "is_active": true}`)) 
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
	storage := NewStorageAPIv1(client)

	// test READ
	repositoryId := uuid.New().String()
	repository, err := storage.GetRepository(repositoryId)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if repository != nil {
		if repository.RepositoryId != repositoryId {
			t.Errorf("bad repositoryId, got: %v want: %v", 
					 repository.RepositoryId, repositoryId)
		}
		if repository.IsActive != true {
			t.Errorf("bad isActive, got: %v want: true", repository.IsActive)
		}
	} else {
		t.Errorf("unexpected: both repository and error are nil!")
	}
}