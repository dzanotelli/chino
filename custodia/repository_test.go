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

    repoCreateResp := map[string]any{
        "repository_id": dummyUUID.String(),
        "description": "unittest",
        "insert_date": "2015-04-14T05:09:54.915Z",
        "last_update": "2015-04-14T05:09:54.915Z",
        "is_active": true,
    }
    repoUpdateResp := map[string]any{
        "repository_id": dummyUUID.String(),
        "description": "changed",
        "insert_date": "2025-04-14T05:09:54.915Z",
        "last_update": "2025-04-14T05:09:54.915Z",
        "is_active": false,
    }

    // mock calls
    mockHandler := func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/api/v1/repositories" && r.Method == "POST" {
            // mock CREATE repository
            data := map[string]any{
                "repository": repoCreateResp,
            }
			envelope.Data, _ = json.Marshal(data)
			out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusCreated)
			w.Write(out)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/repositories/%s", dummyUUID,
        ) && r.Method == "GET" {
            // mock READ repository
            data := map[string]any{
                "repository": repoCreateResp,
            }
            envelope.Data, _ = json.Marshal(data)
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/repositories/%s", dummyUUID,
        ) && r.Method == "PUT" {
            // mock UPDATE repository
            data := map[string]any{
                "repository": repoUpdateResp,
            }
            envelope.Data, _ = json.Marshal(data)
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
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
            data := map[string]any{
                "count": 1,
                "total_count": 1,
                "limit": 100,
                "offset": 0,
                "repositories": []map[string]any{
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
            want any
            got any
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
            want any
            got any
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
            want any
            got any
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
    queryParams := map[string]string{
        "offset": "0",
        "limit": "100",
    }
    repos, err := custodia.ListRepositories(queryParams)
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
