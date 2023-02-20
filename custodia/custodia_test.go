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
	repoId := uuid.New().String()

	// mock calls
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == fmt.Sprintf("/api/v1/repositories/%s", repoId) {
			// test READ
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
	custodia := NewCustodiaAPIv1(client)

	// test READ
	repoPtr, err := custodia.GetRepository(repoId)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if repoPtr != nil {
		if (*repoPtr).RepositoryId != repoId {
			t.Errorf("bad RepositoryId, got: %v want: %v", 
					 repoPtr.RepositoryId, repoId)
		}
		if (*repoPtr).Description != "unittest" {
			t.Errorf("bad Description, got: %v want: unittest", 
					 repoPtr.Description)
		}

		if (*repoPtr).IsActive != true {
			t.Errorf("bad isActive, got: %v want: true", (*repoPtr).IsActive)
		}
	} else {
		t.Errorf("unexpected: both repository and error are nil!")
	}
}