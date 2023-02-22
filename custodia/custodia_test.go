package custodia

import (
	"encoding/json"
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
		InsertDate: "2015-02-24T21:48:16.332",
		LastUpdate: "2015-02-24T21:48:16.332",
		IsActive: false,
	}

	writeRepoResponse := func(w http.ResponseWriter) {
		data := fmt.Sprintf(`{"repository": {"repository_id": "%s", ` +
		`"description": "%s", "insert_date": "%s", `+ 
		`"last_update": "%s", "is_active": %v}}`, 
		dummyRepository.RepositoryId,
		dummyRepository.Description,
		dummyRepository.InsertDate,
		dummyRepository.LastUpdate,
		dummyRepository.IsActive)

		envelope := CustodiaEnvelope{
			Result: "success",
			ResultCode: 200,
			Message: nil,
			Data: []byte(data),
		}
		out, _ := json.Marshal(envelope)

		w.WriteHeader(http.StatusOK)
		w.Write(out)
	}

	// mock calls
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/repositories" && r.Method == "POST" {
			// test CREATE
			writeRepoResponse(w)
		} else if r.URL.Path == fmt.Sprintf("/api/v1/repositories/%s", 
			dummyRepository.RepositoryId) && r.Method == "GET" {
			// test READ
			writeRepoResponse(w)
		} else if r.URL.Path == fmt.Sprintf("/api/v1/repositories/%s", 
			dummyRepository.RepositoryId) && r.Method == "PUT" {
			// test UPDATE
			dummyRepository.Description = "changed"
			dummyRepository.IsActive = false
			writeRepoResponse(w)
		} else if r.URL.Path == fmt.Sprintf("/api/v1/repositories/%s", 
			dummyRepository.RepositoryId) && r.Method == "DELETE" {
			// test DELETE
			w.WriteHeader(http.StatusOK)
			envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
			out, _ := json.Marshal(envelope)
			w.WriteHeader(http.StatusOK)			
			w.Write(out)
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
		if (*repo).InsertDate == "" {
			t.Error("insert_date is empty")
		}
		if (*repo).LastUpdate == "" {
			t.Error("last_update is empty")
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
		if (*repo).InsertDate == "" {
			t.Error("insert_date is empty")
		}
		if (*repo).LastUpdate == "" {
			t.Error("last_update is empty")
		}
		if (*repo).IsActive != false {
			t.Errorf("bad isActive, got: %v want: false", (*repo).IsActive)
		}
	} else {
		t.Errorf("unexpected: both repository and error are nil!")
	}

	// test UPDATE
	repo, err = custodia.UpdateRepository(&dummyRepository, "changed", false)
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
		if (*repo).InsertDate == "" {
			t.Error("insert_date is empty")
		}
		if (*repo).LastUpdate == "" {
			t.Error("last_update is empty")
		}
		if (*repo).IsActive != false {
			t.Errorf("bad isActive, got: %v want: false", (*repo).IsActive)
		}
	} else {
		t.Errorf("unexpected: both repository and error are nil!")
	}
	
	// // test DELETE
	// err = custodia.DeleteRepository(&dummyRepository)
	// if err != nil {
	// 	t.Errorf("error while deleting repository. Details: %v", err)
	// }
}