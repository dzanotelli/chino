package custodia

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/dzanotelli/chino/common"
	"github.com/google/uuid"
)



func TestRepositoryCRUDL(t *testing.T) {
    envelope := CustodiaEnvelope{
        Result: "success",
        ResultCode: 200,
        Message: nil,
    }
    dummyUUID := uuid.New()

    repoCreateResp := map[string]interface{}{
        "repository_id": dummyUUID.String(),
        "description": "unittest",
        "insert_date": "2015-04-14T05:09:54.915Z",
        "last_update": "2015-04-14T05:09:54.915Z",
        "is_active": true,
    }
    repoUpdateResp := map[string]interface{}{
        "repository_id": dummyUUID.String(),
        "description": "changed",
        "insert_date": "2025-04-14T05:09:54.915Z",
        "last_update": "2025-04-14T05:09:54.915Z",
        "is_active": false,
    }
    // repoListResp := map[string]interface{}{
    //     "repositories": []map[string]interface{}{
    //         repoCreateResp,
    //         repoUpdateResp,
    //     },
    // }

    // // ResponseInnerRepository will be included in responses
    // type ResponseInnerRepository struct {
    //     RepositoryId string `json:"repository_id"`
    //     Description string `json:"description"`
    //     InsertDate string `json:"insert_date"`
    //     LastUpdate string `json:"last_update"`
    //     IsActive bool `json:"is_active"`
    // }

    // // RepoResponse will be marshalled to create an API-like reponse
    // type RepoResponse struct {
    //     Repository ResponseInnerRepository `json:"repository"`
    // }

    // // ReposResponse will be marshalled to create an API-like reponse
    // type ReposResponse struct {
    //     Count int `json:"count"`
    //     TotalCount int `json:"total_count"`
    //     Limit int `json:"limit"`
    //     Offset int `json:"offset"`
    //     Repositories []ResponseInnerRepository
    // }

    // // init stuff
    // dummyRepository := ResponseInnerRepository{
    //     RepositoryId: uuid.New().String(),
    //     Description: "unittest",
    //     InsertDate: "2015-02-24T21:48:16.332",
    //     LastUpdate: "2015-02-24T21:48:16.332",
    //     IsActive: false,
    // }

    // writeRepoResponse := func(w http.ResponseWriter) {
    //     data, _ := json.Marshal(RepoResponse{dummyRepository})
    //     envelope := CustodiaEnvelope{
    //         Result: "success",
    //         ResultCode: 200,
    //         Message: nil,
    //         Data: data,
    //     }
    //     out, _ := json.Marshal(envelope)

    //     w.WriteHeader(http.StatusOK)
    //     w.Write(out)
    // }

    // mock calls
    mockHandler := func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/api/v1/repositories" && r.Method == "POST" {
            // mock CREATE repository
            w.WriteHeader(http.StatusOK)
            data := map[string]interface{}{
                "repository": repoCreateResp,
            }
			envelope.Data, _ = json.Marshal(data)
			out, _ := json.Marshal(envelope)
			w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/repositories/%s",
            dummyUUID) && r.Method == "GET" {
            // mock READ repository
            w.WriteHeader(http.StatusOK)
            data := map[string]interface{}{
                "repository": repoCreateResp,
            }
            envelope.Data, _ = json.Marshal(data)
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/repositories/%s",
            dummyUUID) && r.Method == "PUT" {
            // mock UPDATE repository
            w.WriteHeader(http.StatusOK)
            data := map[string]interface{}{
                "repository": repoUpdateResp,
            }
            envelope.Data, _ = json.Marshal(data)
            out, _ := json.Marshal(envelope)
            w.Write(out)
            } else if r.URL.Path == fmt.Sprintf("/api/v1/repositories/%s",
            dummyUUID) && r.Method == "DELETE" {
            // mock DELETE response
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            envelope.Data = nil
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == "/api/v1/repositories" && r.Method == "GET" {
            // mock LIST response
            data := map[string]interface{}{
                "count": 1,
                "total_count": 1,
                "limit": 100,
                "offset": 0,
                "repositories": []map[string]interface{}{
                    repoCreateResp,
                    repoUpdateResp,
                },
            }
            envelope.Data, _ = json.Marshal(data)
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
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {repo.Id.String(), dummyUUID.String()},
            {repo.Description, "unittest"},
            {repo.InsertDate.Year(), 2015},
            {repo.LastUpdate.Year(), 2015},
            {repo.IsActive, true},
        }

        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("CreateRepository %d: bad value, got: %v want: %v",
                    i, test.got, test.want)
            }
        }
    } else {
        t.Errorf("unexpected: both repository and error are nil!")
    }

    // test READ
    repo, err = custodia.ReadRepository(dummyUUID)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if repo != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {repo.Id.String(), dummyUUID.String()},
            {repo.Description, "unittest"},
            {repo.InsertDate.Year(), 2015},
            {repo.LastUpdate.Year(), 2015},
            {repo.IsActive, true},
        }

        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("ReadRepository %d: bad value, got: %v want: %v",
                    i, test.got, test.want)
            }
        }
    } else {
        t.Errorf("unexpected: both repository and error are nil!")
    }

    // test UPDATE
    repo, err = custodia.UpdateRepository(dummyUUID, "changed", false)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if repo != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {repo.Id.String(), dummyUUID.String()},
            {repo.Description, "changed"},
            {repo.InsertDate.Year(), 2025},
            {repo.LastUpdate.Year(), 2025},
            {repo.IsActive, false},
        }
        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("UpdateRepository %d: bad value, got: %v want: %v",
                    i, test.got, test.want)
            }
        }
    } else {
        t.Errorf("unexpected: both repository and error are nil!")
    }

    // test DELETE
    err = custodia.DeleteRepository(dummyUUID, true)
    if err != nil {
        t.Errorf("error while deleting repository. Details: %v", err)
    }

    // test LIST
    repos, err := custodia.ListRepositories()
    if err != nil {
        t.Errorf("error while listing repositories. Details: %v", err)
    }
    if len(repos) != 2 {
        t.Errorf("bad repositories lenght, got: %v want: 1", len(repos))
    }
    if repos[0].Id.String() != dummyUUID.String() {
        t.Errorf("bad repository id, got: %v want: %v",
            dummyUUID.String(), repos[0].Id.String())
    }
}
