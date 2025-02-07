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

func TestSchemaCRUDL(t *testing.T) {
    // ResponseInnerSchema will be included in responses
    type ResponseInnerSchema struct {
        SchemaId string `json:"schema_id"`
        RepositoryId string `json:"repository_id"`
        Description string `json:"description"`
        InsertDate string `json:"insert_date"`
        LastUpdate string `json:"last_update"`
        IsActive bool `json:"is_active"`
        Structure []SchemaField `json:"structure"`
    }

    // SchemaResponse will be marshalled to create an API-like response
    type SchemaResponse struct {
        Schema ResponseInnerSchema `json:"schema"`
    }

    // SchemasResponse will be marshalled to create an API-like response
    type SchemasResponse struct {
        Count int `json:"count"`
        TotalCount int `json:"total_count"`
        Limit int `json:"limit"`
        Offset int `json:"offset"`
        Schemas []ResponseInnerSchema `json:"schemas"`
    }

    // init stuff
    dummySchema := ResponseInnerSchema{
        SchemaId: uuid.New().String(),
        RepositoryId: uuid.New().String(),
        Description: "unittest",
        InsertDate: "2015-02-24T21:48:16.332",
        LastUpdate: "2015-02-24T21:48:16.332",
        IsActive: false,
        // Structure: json.RawMessage{},
        Structure: []SchemaField{
            {Name: "IntField", Type: "integer", Indexed: true, Default: 42},
            {Name: "StrField", Type: "string", Indexed: true, Default: "asd"},
            {Name: "FloatField", Type: "float", Indexed: false, Default: 3.14},
            {Name: "BoolField", Type: "bool", Indexed: false},
            {Name: "DateField", Type: "date", Default: "2023-03-15"},
            {Name: "TimeField", Type: "time", Default: "11:43:04.058"},
            {Name: "DateTimeField", Type: "datetime",
                Default: "2023-03-15T11:43:04.058"},
        },
    }

    // shortcut
    repoId := dummySchema.RepositoryId

    writeSchemaResponse := func(w http.ResponseWriter) {
        data, _ := json.Marshal(SchemaResponse{dummySchema})
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
        if r.URL.Path == fmt.Sprintf("/api/v1/repositories/%s/schemas",
            repoId) && r.Method == "POST" {
            // mock CREATE response
            writeSchemaResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/schemas/%s",
            dummySchema.SchemaId) && r.Method == "GET" {
            // mock READ response
            writeSchemaResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/schemas/%s",
            dummySchema.SchemaId) && r.Method == "PUT" {
            // mock UPDATE response
            dummySchema.Description = "changed"
            // dummySchema.Structure[0].Default = 21
            writeSchemaResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/schemas/%s",
            dummySchema.SchemaId) && r.Method == "DELETE" {
            // mock DELETE response
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/repositories/%s/schemas",
            repoId) &&  r.Method == "GET" {
            // mock LIST response
            schemasResp := SchemasResponse{
                Count: 1,
                TotalCount: 1,
                Limit: 100,
                Offset: 0,
                Schemas: []ResponseInnerSchema{dummySchema},
            }
            data, _ := json.Marshal(schemasResp)
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

    // test CREATE: we submit an empty field list, since the response is mocked
    // and we will still get a working structure. The purpose here is to test
    // that the received data are correctly populating the objects
    structure := []SchemaField{}
    // we init a Repository with just the right id, don't need other data
    repo := Repository{Id: dummySchema.RepositoryId}
    schema, err := custodia.CreateSchema(&repo, "unittest", true, structure)

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if schema != nil {
        if schema.RepositoryId != repoId {
            t.Errorf("bad RepositoryId, got: %v want: %v",
                schema.RepositoryId, repoId)
        }
        if schema.Id != dummySchema.SchemaId {
            t.Errorf("bad SchemaId, got: %v want: %v",
                schema.Id, dummySchema.SchemaId)
        }
        if schema.Description != dummySchema.Description {
            t.Errorf("bad Description, got: %v want: %s",
                     schema.Description,
                     dummySchema.Description)
        }
        if schema.InsertDate.Year() != 2015 {
            t.Errorf("bad insert_date year, got: %v want: 2015",
                schema.InsertDate.Year())
        }
        if schema.LastUpdate.Year() != 2015 {
            t.Errorf("bad last_update year, got: %v want: 2015",
                schema.InsertDate.Year())
        }
        if schema.IsActive != false {
            t.Errorf("bad isActive, got: %v want: false", schema.IsActive)
        }

        expectedFields := dummySchema.Structure
        for i, want := range expectedFields {
            got := schema.Structure[i]
            if want != got {
                t.Errorf("bad field received, got: %v want: %v", got, want)
            }
        }
    } else {
        t.Errorf("unexpected: both schema and error are nil!")
    }

    // test READ
    schema, err = custodia.ReadSchema(dummySchema.SchemaId)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if schema != nil {
        if schema.RepositoryId != repoId {
            t.Errorf("bad RepositoryId, got: %v want: %v",
                schema.RepositoryId, repoId)
        }
        if schema.Description != dummySchema.Description {
            t.Errorf("bad Description, got: %v want: %s",
                     schema.Description,
                     dummySchema.Description)
        }
        if schema.InsertDate.Year() != 2015 {
            t.Errorf("bad insert_date year, got: %v want: 2015",
                schema.InsertDate.Year())
        }
        if schema.LastUpdate.Year() != 2015 {
            t.Errorf("bad last_update year, got: %v want: 2015",
                schema.InsertDate.Year())
        }
        if schema.IsActive != false {
            t.Errorf("bad isActive, got: %v want: false", schema.IsActive)
        }

        expectedFields := dummySchema.Structure
        for i, want := range expectedFields {
            got := schema.Structure[i]
            if want != got {
                t.Errorf("bad field received, got: %v want: %v", got, want)
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
        if schema.RepositoryId != repoId {
            t.Errorf("bad RepositoryId, got: %v want: %v",
                schema.RepositoryId, repoId)
        }
        if schema.Description != dummySchema.Description {
            t.Errorf("bad Description, got: %v want: %s",
                     schema.Description,
                     dummySchema.Description)
        }
        if schema.InsertDate.Year() != 2015 {
            t.Errorf("bad insert_date year, got: %v want: 2015",
                schema.InsertDate.Year())
        }
        if schema.LastUpdate.Year() != 2015 {
            t.Errorf("bad last_update year, got: %v want: 2015",
                schema.InsertDate.Year())
        }
        if schema.IsActive != false {
            t.Errorf("bad isActive, got: %v want: false", schema.IsActive)
        }

        expectedFields := dummySchema.Structure
        for i, want := range expectedFields {
            got := schema.Structure[i]
            if want != got {
                t.Errorf("bad field received, got: %v want: %v", got, want)
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
    schemas, err := custodia.ListSchemas(repoId)
    if err != nil {
        t.Errorf("error while listing schemas. Details: %v", err)
    }
    if len(schemas) != 1 {
        t.Errorf("bad schemas lenght, got: %v want: 1", len(schemas))
    }
    if schemas[0].Id != dummySchema.SchemaId {
        t.Errorf("bad schema id, got: %v want: %v",
        dummySchema.SchemaId, schemas[0].Id)
    }
}
