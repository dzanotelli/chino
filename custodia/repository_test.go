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



func TestRepositoryCRUDL(t *testing.T) {
    // ResponseInnerRepository will be included in responses
    type ResponseInnerRepository struct {
        RepositoryId string `json:"repository_id"`
        Description string `json:"description"`
        InsertDate string `json:"insert_date"`
        LastUpdate string `json:"last_update"`
        IsActive bool `json:"is_active"`
    }

    // RepoResponse will be marshalled to create an API-like reponse
    type RepoResponse struct {
        Repository ResponseInnerRepository `json:"repository"`
    }

    // ReposResponse will be marshalled to create an API-like reponse
    type ReposResponse struct {
        Count int `json:"count"`
        TotalCount int `json:"total_count"`
        Limit int `json:"limit"`
        Offset int `json:"offset"`
        Repositories []ResponseInnerRepository
    }

    // init stuff
    dummyRepository := ResponseInnerRepository{
        RepositoryId: uuid.New().String(),
        Description: "unittest",
        InsertDate: "2015-02-24T21:48:16.332",
        LastUpdate: "2015-02-24T21:48:16.332",
        IsActive: false,
    }

    writeRepoResponse := func(w http.ResponseWriter) {
        data, _ := json.Marshal(RepoResponse{dummyRepository})
        envelope := CustodiaEnvelope{
            Result: "success",
            ResultCode: 200,
            Message: nil,
            Data: data,
        }
        out, _ := json.Marshal(envelope)

        w.WriteHeader(http.StatusOK)
        w.Write(out)
    }

    // mock calls
    mockHandler := func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/api/v1/repositories" && r.Method == "POST" {
            // mock CREATE response
            writeRepoResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/repositories/%s",
            dummyRepository.RepositoryId) && r.Method == "GET" {
            // mock READ response
            writeRepoResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/repositories/%s",
            dummyRepository.RepositoryId) && r.Method == "PUT" {
            // mock UPDATE response
            dummyRepository.Description = "changed"
            dummyRepository.IsActive = false
            writeRepoResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/repositories/%s",
            dummyRepository.RepositoryId) && r.Method == "DELETE" {
            // mock DELETE response
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == "/api/v1/repositories" && r.Method == "GET" {
            // mock LIST response
            repositoriesResp := ReposResponse {
                Count: 1,
                TotalCount: 1,
                Limit: 100,
                Offset: 0,
                Repositories: []ResponseInnerRepository{dummyRepository},
            }
            data, _ := json.Marshal(repositoriesResp)
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            envelope.Data = data
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

    client := common.NewClient(server.URL, common.GetFakeAuth())
    custodia := NewCustodiaAPIv1(client)

    // test CREATE
    repo, err := custodia.CreateRepository("unittest", false)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if repo != nil {
        if (*repo).Id != dummyRepository.RepositoryId {
            t.Errorf("bad RepositoryId, got: %v want: %v",
                     repo.Id, dummyRepository.RepositoryId)
        }
        if (*repo).Description != dummyRepository.Description {
            t.Errorf("bad Description, got: %v want: %s",
                     repo.Description,
                     dummyRepository.Description)
        }
        if (*repo).InsertDate.Year() != 2015 {
            t.Errorf("bad insert_date year, got: %v want: 2015",
                (*repo).InsertDate.Year())
        }
        if (*repo).LastUpdate.Year() != 2015 {
            t.Errorf("bad last_update year, got: %v want: 2015",
                (*repo).InsertDate.Year())
        }
        if (*repo).IsActive != false {
            t.Errorf("bad isActive, got: %v want: false", (*repo).IsActive)
        }
    } else {
        t.Errorf("unexpected: both repository and error are nil!")
    }

    // test READ
    repo, err = custodia.ReadRepository(dummyRepository.RepositoryId)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if repo != nil {
        if repo.Id != dummyRepository.RepositoryId {
            t.Errorf("bad RepositoryId, got: %v want: %v",
                     repo.Id, dummyRepository.RepositoryId)
        }
        if repo.Description != dummyRepository.Description {
            t.Errorf("bad Description, got: %v want: %s",
                     repo.Description,
                    dummyRepository.Description)
        }
        if repo.InsertDate.Year() != 2015 {
            t.Errorf("bad insert_date year, got: %v want: 2015",
                repo.InsertDate.Year())
        }
        if repo.LastUpdate.Year() != 2015 {
            t.Errorf("bad last_update year, got: %v want: 2015",
                repo.InsertDate.Year())
        }
        if repo.IsActive != false {
            t.Errorf("bad isActive, got: %v want: false", (*repo).IsActive)
        }
    } else {
        t.Errorf("unexpected: both repository and error are nil!")
    }

    // test UPDATE
    repo, err = custodia.UpdateRepository(repo.Id, "changed", false)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if repo != nil {
        if repo.Id != dummyRepository.RepositoryId {
            t.Errorf("bad RepositoryId, got: %v want: %v",
                     repo.Id, dummyRepository.RepositoryId)
        }
        if repo.Description != dummyRepository.Description {
            t.Errorf("bad Description, got: %v want: %s",
                     repo.Description,
                    dummyRepository.Description)
        }
        if repo.InsertDate.Year() != 2015 {
            t.Errorf("bad insert_date year, got: %v want: 2015",
                repo.InsertDate.Year())
        }
        if repo.LastUpdate.Year() != 2015 {
            t.Errorf("bad last_update year, got: %v want: 2015",
                repo.InsertDate.Year())
        }
        if repo.IsActive != false {
            t.Errorf("bad isActive, got: %v want: false", repo.IsActive)
        }
    } else {
        t.Errorf("unexpected: both repository and error are nil!")
    }

    // test DELETE
    err = custodia.DeleteRepository(repo.Id, true)
    if err != nil {
        t.Errorf("error while deleting repository. Details: %v", err)
    }

    // test LIST
    repos, err := custodia.ListRepositories()
    if err != nil {
        t.Errorf("error while listing repositories. Details: %v", err)
    }
    if len(repos) != 1 {
        t.Errorf("bad repositories lenght, got: %v want: 1", len(repos))
    }
    if repos[0].Id != dummyRepository.RepositoryId {
        t.Errorf("bad repository id, got: %v want: %v",
            dummyRepository.RepositoryId, repos[0].Id)
    }
}
