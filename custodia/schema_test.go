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

func TestSchemaCRUDL(t *testing.T) {
    envelope := CustodiaEnvelope{
        Result: "success",
        ResultCode: 200,
        Message: nil,
    }
    dummyUUID := uuid.New()

    createResp := map[string]interface{}{
        "schema_id": dummyUUID.String(),
        "repository_id": dummyUUID.String(),
        "description": "unittest",
        "insert_date": "2015-04-14T05:09:54.915Z",
        "last_update": "2015-04-14T05:09:54.915Z",
        "is_active": true,
        "structure": []map[string]interface{}{
            {
                "name": "IntField",
                "type": "integer",
                "default": 42,
                "indexed": true,
            },
            {
                "name": "StrField",
                "type": "string",
                "default": "asd",
            },
            {
                "name": "FloatField",
                "type": "number",
                "default": 3.14,
            },
            {
                "name": "BoolField",
                "type": "boolean",
                "default": false,
            },
            {
                "name": "DateField",
                "type": "date",
                "default": "2023-03-15",
            },
            {
                "name": "TimeField",
                "type": "time",
                "default": "11:43:04.058",
            },
            {
                "name": "DateTimeField",
                "type": "datetime",
                "default": "2023-03-15T11:43:04.058",
            },
        },
    }
    updateResp := map[string]interface{}{
        "schema_id": dummyUUID.String(),
        "repository_id": dummyUUID.String(),
        "description": "changed",
        "insert_date": "2025-04-14T05:09:54.915Z",
        "last_update": "2025-04-14T05:09:54.915Z",
        "is_active": false,
        "structure": createResp["structure"],
    }




    // // ResponseInnerSchema will be included in responses
    // type ResponseInnerSchema struct {
    //     SchemaId string `json:"schema_id"`
    //     RepositoryId string `json:"repository_id"`
    //     Description string `json:"description"`
    //     InsertDate string `json:"insert_date"`
    //     LastUpdate string `json:"last_update"`
    //     IsActive bool `json:"is_active"`
    //     Structure []SchemaField `json:"structure"`
    // }

    // // SchemaResponse will be marshalled to create an API-like response
    // type SchemaResponse struct {
    //     Schema ResponseInnerSchema `json:"schema"`
    // }

    // // SchemasResponse will be marshalled to create an API-like response
    // type SchemasResponse struct {
    //     Count int `json:"count"`
    //     TotalCount int `json:"total_count"`
    //     Limit int `json:"limit"`
    //     Offset int `json:"offset"`
    //     Schemas []ResponseInnerSchema `json:"schemas"`
    // }

    // init stuff
    // dummySchema := ResponseInnerSchema{
    //     SchemaId: uuid.New().String(),
    //     RepositoryId: uuid.New().String(),
    //     Description: "unittest",
    //     InsertDate: "2015-02-24T21:48:16.332",
    //     LastUpdate: "2015-02-24T21:48:16.332",
    //     IsActive: false,
    //     // Structure: json.RawMessage{},
    //     Structure: []SchemaField{
    //         {Name: "IntField", Type: "integer", Indexed: true, Default: 42},
    //         {Name: "StrField", Type: "string", Indexed: true, Default: "asd"},
    //         {Name: "FloatField", Type: "float", Indexed: false, Default: 3.14},
    //         {Name: "BoolField", Type: "bool", Indexed: false},
    //         {Name: "DateField", Type: "date", Default: "2023-03-15"},
    //         {Name: "TimeField", Type: "time", Default: "11:43:04.058"},
    //         {Name: "DateTimeField", Type: "datetime",
    //             Default: "2023-03-15T11:43:04.058"},
    //     },
    // }

    // shortcut
    // repoId := dummyUUID.String()

    // writeSchemaResponse := func(w http.ResponseWriter) {
    //     data, _ := json.Marshal(SchemaResponse{dummySchema})
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
        if r.URL.Path == fmt.Sprintf(
            "/api/v1/repositories/%s/schemas", dummyUUID,
        ) && r.Method == "POST" {
            // mock CREATE response
            w.WriteHeader(http.StatusCreated)
            envelope.Data, _ = json.Marshal(createResp)
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/schemas/%s", dummyUUID,
        ) && r.Method == "GET" {
            // mock READ response
            w.WriteHeader(http.StatusOK)
            envelope.Data, _ = json.Marshal(createResp)
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/schemas/%s",
            dummyUUID,
        ) && r.Method == "PUT" {
            // mock UPDATE response
            w.WriteHeader(http.StatusOK)
            envelope.Data, _ = json.Marshal(updateResp)
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/schemas/%s", dummyUUID,
        ) && r.Method == "DELETE" {
            // mock DELETE response
            w.WriteHeader(http.StatusOK)
            envelope.Data = nil
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/repositories/%s/schemas", dummyUUID,
        ) &&  r.Method == "GET" {
            // mock LIST response
            w.WriteHeader(http.StatusOK)
            data := map[string]interface{}{
                "count": 1,
                "total_count": 1,
                "limit": 100,
                "offset": 0,
                "schemas": []map[string]interface{}{
                    createResp,
                    updateResp,
                },
            }
            envelope.Data, _ = json.Marshal(data)
            out, _ := json.Marshal(envelope)
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

    // test CREATE: we submit an empty field list, since the response is mocked
    // and we will still get a working structure. The purpose here is to test
    // that the received data are correctly populating the objects
    structure := []SchemaField{}
    // we init a Repository with just the right id, don't need other data
    schema, err := custodia.CreateSchema(dummyUUID, "unittest", true,
        structure)

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if schema != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyUUID.String(), schema.RepositoryId.String()},
            {dummyUUID.String(), schema.Id.String()},
            {"unittest", schema.Description},
            {2015, schema.InsertDate.Year()},
            {04, schema.InsertDate.Month()},
            {14, schema.InsertDate.Day()},
            {5, schema.InsertDate.Hour()},
            {54, schema.InsertDate.Minute()},
            {91, schema.InsertDate.Second()},
            {2015, schema.LastUpdate.Year()},
            {04, schema.LastUpdate.Month()},
            {14, schema.LastUpdate.Day()},
            {5, schema.LastUpdate.Hour()},
            {54, schema.LastUpdate.Minute()},
            {91, schema.LastUpdate.Second()},
            {true, schema.IsActive},
            // test some content
            {"IntField", schema.Structure[0].Name},
            {"integer", schema.Structure[0].Type},
            {true, schema.Structure[0].Indexed},
            {42, schema.Structure[0].Default},
            {"StrField", schema.Structure[1].Name},
            {"string", schema.Structure[1].Type},
            {"asd", schema.Structure[1].Default},
            {false, schema.Structure[1].Indexed},
            {"BoolField", schema.Structure[2].Name},
            {"boolean", schema.Structure[2].Type},
            {false, schema.Structure[2].Indexed},
            {false, schema.Structure[2].Default},
        }
        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("CreateSchema %d: expected %v, got %v", i, test.want,
                test.got)
            }
        }
    } else {
        t.Errorf("unexpected: both schema and error are nil!")
    }

    // test READ
    schema, err = custodia.ReadSchema(dummyUUID)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if schema != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyUUID.String(), schema.RepositoryId.String()},
            {dummyUUID.String(), schema.Id.String()},
            {"unittest", schema.Description},
            {2015, schema.InsertDate.Year()},
            {04, schema.InsertDate.Month()},
            {14, schema.InsertDate.Day()},
            {5, schema.InsertDate.Hour()},
            {54, schema.InsertDate.Minute()},
            {91, schema.InsertDate.Second()},
            {2015, schema.LastUpdate.Year()},
            {04, schema.LastUpdate.Month()},
            {14, schema.LastUpdate.Day()},
            {5, schema.LastUpdate.Hour()},
            {54, schema.LastUpdate.Minute()},
            {91, schema.LastUpdate.Second()},
            {true, schema.IsActive},
            // test some content
            {"IntField", schema.Structure[0].Name},
            {"integer", schema.Structure[0].Type},
            {true, schema.Structure[0].Indexed},
            {42, schema.Structure[0].Default},
            {"StrField", schema.Structure[1].Name},
            {"string", schema.Structure[1].Type},
            {"asd", schema.Structure[1].Default},
            {false, schema.Structure[1].Indexed},
            {"BoolField", schema.Structure[2].Name},
            {"boolean", schema.Structure[2].Type},
            {false, schema.Structure[2].Indexed},
            {false, schema.Structure[2].Default},
        }
        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("CreateSchema %d: expected %v, got %v", i, test.want,
                test.got)
            }
        }
    } else {
        t.Errorf("unexpected: both schema and error are nil!")
    }

    // test UPDATE
    schema, err = custodia.UpdateSchema(schema.Id, "changed", true,
        structure)

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if schema != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyUUID.String(), schema.RepositoryId.String()},
            {dummyUUID.String(), schema.Id.String()},
            {"changed", schema.Description},  // changed to 'changed'
            {2015, schema.InsertDate.Year()},
            {04, schema.InsertDate.Month()},
            {14, schema.InsertDate.Day()},
            {5, schema.InsertDate.Hour()},
            {54, schema.InsertDate.Minute()},
            {91, schema.InsertDate.Second()},
            {2015, schema.LastUpdate.Year()},
            {04, schema.LastUpdate.Month()},
            {14, schema.LastUpdate.Day()},
            {5, schema.LastUpdate.Hour()},
            {54, schema.LastUpdate.Minute()},
            {91, schema.LastUpdate.Second()},
            {false, schema.IsActive},  // changed to 'false'
            // test some content
            {"IntField", schema.Structure[0].Name},
            {"integer", schema.Structure[0].Type},
            {true, schema.Structure[0].Indexed},
            {42, schema.Structure[0].Default},
            {"StrField", schema.Structure[1].Name},
            {"string", schema.Structure[1].Type},
            {"asd", schema.Structure[1].Default},
            {false, schema.Structure[1].Indexed},
            {"BoolField", schema.Structure[2].Name},
            {"boolean", schema.Structure[2].Type},
            {false, schema.Structure[2].Indexed},
            {false, schema.Structure[2].Default},
        }
        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("CreateSchema %d: expected %v, got %v", i, test.want,
                test.got)
            }
        }
    } else {
        t.Errorf("unexpected: both schema and error are nil!")
    }

    // test DELETE
    err = custodia.DeleteSchema(schema.Id, true, true)
    if err != nil {
        t.Errorf("error while deleting schema. Details: %v", err)
    }

    // test LIST
    schemas, err := custodia.ListSchemas(dummyUUID)
    if err != nil {
        t.Errorf("error while listing schemas. Details: %v", err)
    }
    var tests = []struct {
        want interface{}
        got interface{}
    }{
        {2, len(schemas)},
        {dummyUUID.String(), schemas[0].RepositoryId.String()},
        {dummyUUID.String(), schemas[0].Id.String()},
        {"unittest", schemas[1].Description},
        {true, schemas[1].IsActive},
        {dummyUUID.String(), schemas[1].RepositoryId.String()},
        {dummyUUID.String(), schemas[1].Id.String()},
        {"changed", schemas[1].Description},
        {false, schemas[1].IsActive},
    }
    for i, test := range tests {
        if !reflect.DeepEqual(test.want, test.got) {
            t.Errorf("CreateSchema %d: expected %v, got %v", i, test.want,
            test.got)
        }
    }
}
