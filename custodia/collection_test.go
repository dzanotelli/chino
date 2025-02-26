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


func TestCollectionCRUDL(t *testing.T) {
    envelope := CustodiaEnvelope{
        Result: "success",
        ResultCode: 200,
        Message: nil,
    }

    dummyUUID := uuid.New()

    createResponse := map[string]interface{}{
        "collection_id": dummyUUID.String(),
        "name": "unittest",
        "insert_date": "2015-04-14T05:09:54.915Z",
        "last_update": "2015-04-14T05:09:54.915Z",
        "is_active": true,
    }
    updateResponse := map[string]interface{}{
        "collection_id": dummyUUID.String(),
        "name": "changed",
        "insert_date": "2025-04-14T05:09:54.915Z",
        "last_update": "2025-04-14T05:09:54.915Z",
        "is_active": false,
    }
    listResponse := map[string]interface{}{
        "collections": []map[string]interface{}{
            createResponse,
            updateResponse,
        },
    }

    // mock calls
    mockHandler := func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/api/v1/collections" && r.Method == "POST" {
            // mock CREATE
            w.WriteHeader(http.StatusOK)
            envelope.Data, _ = json.Marshal(createResponse)
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/collections/%s",
            dummyUUID) && r.Method == "GET" {
            // mock READ
            w.WriteHeader(http.StatusOK)
            envelope.Data, _ = json.Marshal(createResponse) // same as CREATE
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/collections/%s",
            dummyUUID) && r.Method == "PUT" {
            // mock UPDATE
            w.WriteHeader(http.StatusOK)
            envelope.Data, _ = json.Marshal(updateResponse)
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/collections/%s",
            dummyUUID) && r.Method == "DELETE" {
            // mock DELETE
            w.WriteHeader(http.StatusOK)
            envelope.Data = nil
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else if r.URL.Path == "/api/v1/collections" && r.Method == "GET" {
            // mock LIST
            w.WriteHeader(http.StatusOK)
            envelope.Data, _ = json.Marshal(listResponse)
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else {
            err := `{"result": "error", "result_code": 404, "data": null, `
            err += `"message": "Resource not found (you may have a '/' at `
            err += `the end)"}`
            fmt.Print(err)
            w.WriteHeader(http.StatusNotFound)
            w.Write([]byte(err))
        }
    }

    server := httptest.NewServer(http.HandlerFunc(mockHandler))
    defer server.Close()

    client := common.NewClient(server.URL, common.GetFakeAuth())
    custodia := NewCustodiaAPIv1(client)

    // == test CRUDL ==
    // test Create
    collection, err := custodia.CreateCollection("unittest")
    if err != nil {
        t.Errorf("error while processing request: %s", err)
        return // stop execution here
    } else {
        var tests = []struct {
            want interface{}
            got  interface{}
        }{
            {dummyUUID.String(), collection.Id.String()},
            {"unittest", collection.Name},
            {2015, collection.InsertDate.Year()},
            {4, int(collection.InsertDate.Month())},
            {14, collection.InsertDate.Day()},
            {5, collection.InsertDate.Hour()},
            {9, collection.InsertDate.Minute()},
            {54, collection.InsertDate.Second()},
            {2015, collection.LastUpdate.Year()},
            {4, int(collection.LastUpdate.Month())},
            {14, collection.LastUpdate.Day()},
            {5, collection.LastUpdate.Hour()},
            {9, collection.LastUpdate.Minute()},
            {54, collection.LastUpdate.Second()},
            {true, collection.IsActive},
        }

        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("CreateCollection %d: expected %v, got %v", i, test.want,
                test.got)
            }
        }
    }

    // test Read
    collection, err = custodia.ReadCollection(dummyUUID)
    if err != nil {
        t.Errorf("error while processing request: %s", err)
    } else {
        var tests = []struct {
            want interface{}
            got  interface{}
        }{
            {dummyUUID.String(), collection.Id.String()},
            {"unittest", collection.Name},
            {2015, collection.InsertDate.Year()},
            {4, int(collection.InsertDate.Month())},
            {14, collection.InsertDate.Day()},
            {5, collection.InsertDate.Hour()},
            {9, collection.InsertDate.Minute()},
            {54, collection.InsertDate.Second()},
            {2015, collection.LastUpdate.Year()},
            {4, int(collection.LastUpdate.Month())},
            {14, collection.LastUpdate.Day()},
            {5, collection.LastUpdate.Hour()},
            {9, collection.LastUpdate.Minute()},
            {54, collection.LastUpdate.Second()},
            {true, collection.IsActive},
        }

        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("ReadCollection %d: expected %v, got %v", i,
                test.want, test.got)
            }
        }
    }

    // test Update
    collection, err = custodia.UpdateCollection(dummyUUID, "unittest")
    if err != nil {
        t.Errorf("error while processing request: %s", err)
    } else {
        var tests = []struct {
            want interface{}
            got  interface{}
        }{
            {dummyUUID.String(), collection.Id.String()},
            {"changed", collection.Name},
            {2025, collection.InsertDate.Year()},
            {4, int(collection.InsertDate.Month())},
            {14, collection.InsertDate.Day()},
            {5, collection.InsertDate.Hour()},
            {9, collection.InsertDate.Minute()},
            {54, collection.InsertDate.Second()},
            {2025, collection.LastUpdate.Year()},
            {4, int(collection.LastUpdate.Month())},
            {14, collection.LastUpdate.Day()},
            {5, collection.LastUpdate.Hour()},
            {9, collection.LastUpdate.Minute()},
            {54, collection.LastUpdate.Second()},
            {false, collection.IsActive},
        }

        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("UpdateCollection %d: expected %v, got %v", i,
                test.want, test.got)
            }
        }
    }

    // test Delete
    err = custodia.DeleteCollection(dummyUUID, true)
    if err != nil {
        t.Errorf("error while processing request: %s", err)
    }

    // test List
    collections, err := custodia.ListCollections()
    if err != nil {
        t.Errorf("error while processing request: %s", err)
    } else {
        if len(collections) != 2 {
            t.Errorf("ListCollections: want %v, got %v", 2,
                len(collections))
        }
        // we don't check every field, just some here and there
        var tests = []struct {
            want interface{}
            got  interface{}
        }{
            {dummyUUID.String(), collections[0].Id.String()},
            {"unittest", collections[0].Name},
            {2015, collections[0].InsertDate.Year()},
            {dummyUUID.String(), collections[1].Id.String()},
            {"changed", collections[1].Name},
            {2025, collections[1].InsertDate.Year()},
        }

        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("ListCollections %d: expected %v, got %v", i,
                test.want, test.got)
            }
        }
    }
}

func TestCollectionAndDocuments(t *testing.T) {
    envelope := CustodiaEnvelope{
        Result: "success",
        ResultCode: 200,
        Message: nil,
    }

    dummyUUID := uuid.New()

    collectionData := map[string]interface{}{
        "collection_id": dummyUUID.String(),
        "name": "unittest",
        "insert_date": "2015-04-14T05:09:54.915Z",
        "last_update": "2015-04-14T05:09:54.915Z",
        "is_active": true,
    }

    documentResponse := map[string]interface{}{
        "document_id": dummyUUID.String(),
        "repository_id": dummyUUID.String(),
        "schema_id": dummyUUID.String(),
        "insert_date": "2015-04-14T05:09:54.915Z",
        "last_update": "2015-04-14T05:09:54.915Z",
        "is_active": true,
        "content": map[string]interface{}{
            "field": 42,
        },
    }
    searchCollectionResponse := map[string]interface{}{
        "collections": []map[string]interface{}{
            {
                "collection_id": dummyUUID.String(),
                "name": "unittest1",
                "insert_date": "2015-04-14T05:09:54.915Z",
                "last_update": "2015-04-14T05:09:54.915Z",
                "is_active": true,
            },
            {
                "collection_id": dummyUUID.String(),
                "name": "unittest2",
                "insert_date": "2035-04-14T05:09:54.915Z",
                "last_update": "2035-04-14T05:09:54.915Z",
                "is_active": true,

            },
        },
    }

    // mock calls
    mockHandler := func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == fmt.Sprintf("/api/v1/collections/documents/%s",
            dummyUUID) && r.Method == "GET" {
            w.WriteHeader(http.StatusOK)
            data := map[string]interface{}{
                "collections": []map[string]interface{}{
                    collectionData,
                },
            }
            envelope.Data, _ = json.Marshal(data)
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/collections/%s/documents",
            dummyUUID) && r.Method == "GET" {
            w.WriteHeader(http.StatusOK)
            data := map[string]interface{}{
                "documents": []map[string]interface{}{
                    documentResponse,
                },
            }
            envelope.Data, _ = json.Marshal(data)
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/collections/%s/documents/%s", dummyUUID, dummyUUID) &&
            r.Method == "POST" {
            w.WriteHeader(http.StatusOK)
            envelope.Data = nil
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/collections/%s/documents/%s", dummyUUID, dummyUUID) &&
            r.Method == "DELETE" {
            w.WriteHeader(http.StatusOK)
            envelope.Data = nil
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else if r.URL.Path == "/api/v1/collections/search" &&
            r.Method == "POST" {
            w.WriteHeader(http.StatusOK)
            envelope.Data, _ = json.Marshal(searchCollectionResponse)
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else {
            err := `{"result": "error", "result_code": 404, "data": null, `
            err += `"message": "Resource not found (you may have a '/' at `
            err += `the end)"}`
            fmt.Print(err)
            w.WriteHeader(http.StatusNotFound)
            w.Write([]byte(err))
        }
    }

    server := httptest.NewServer(http.HandlerFunc(mockHandler))
    defer server.Close()

    client := common.NewClient(server.URL, common.GetFakeAuth())
    custodia := NewCustodiaAPIv1(client)

    // test ListDocumentCollections
    collections, err := custodia.ListDocumentCollections(dummyUUID)
    if err != nil {
        t.Errorf("error while processing request: %s", err)
    } else {
        if len(collections) != 1 {
            t.Errorf("ListDocumentCollections: want %v, got %v", 1,
                len(collections))
        }
        var tests = []struct {
            want interface{}
            got  interface{}
        }{
            {dummyUUID.String(), collections[0].Id.String()},
            {"unittest", collections[0].Name},
            {2015, collections[0].InsertDate.Year()},
        }

        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("ListDocumentCollections %d: expected %v, got %v", i,
                test.want, test.got)
            }
        }
    }

    // test ListCollectionDocuments
    documents, err := custodia.ListCollectionDocuments(dummyUUID, true)
    if err != nil {
        t.Errorf("error while processing request: %s", err)
    } else {
        if len(documents) != 1 {
            t.Errorf("ListCollectionDocuments: want %v, got %v", 1,
                len(documents))
        }
        var tests = []struct {
            want interface{}
            got  interface{}
        }{
            {dummyUUID.String(), documents[0].Id.String()},
            {dummyUUID.String(), documents[0].RepositoryId.String()},
            {dummyUUID.String(), documents[0].SchemaId.String()},
            {2015, documents[0].InsertDate.Year()},
            {true, documents[0].IsActive},
            // FIXME: need to fix how Document.Content is handled:
            //  in ReadDocument we convert the underlying type to the
            //  types defined in Schema, here we don't (for now)
            {float64(42), documents[0].Content["field"].(float64)},
        }

        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("ListCollectionDocuments %d: expected %v, got %v", i,
                test.want, test.got)
            }
        }
    }

    // test AddDocumentToCollection
    err = custodia.AddDocumentToCollection(dummyUUID, dummyUUID)
    if err != nil {
        t.Errorf("error while processing request: %s", err)
    }

    // test RemoveDocumentFromCollection
    err = custodia.RemoveDocumentFromCollection(dummyUUID, dummyUUID)
    if err != nil {
        t.Errorf("error while processing request: %s", err)
    }

    // test SearchCollection
    collections, err = custodia.SearchCollection("unittest", true)
    if err != nil {
        t.Errorf("error while processing request: %s", err)
    } else {
        if len(collections) != 2 {
            t.Errorf("SearchCollection: want %v, got %v", 2,
                len(collections))
        }
        var tests = []struct {
            want interface{}
            got  interface{}
        }{
            {dummyUUID.String(), collections[0].Id.String()},
            {"unittest1", collections[0].Name},
            {2015, collections[0].InsertDate.Year()},
            {dummyUUID.String(), collections[1].Id.String()},
            {"unittest2", collections[1].Name},
            {2035, collections[1].InsertDate.Year()},
        }

        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("SearchCollection %d: expected %v, got %v", i,
                test.want, test.got)
            }
        }
    }
}