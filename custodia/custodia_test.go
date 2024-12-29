package custodia

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

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

func TestDocumentCRUDL(t *testing.T) {
    // ResponseInnerDocument will be included in responses
    type ResponseInnerDocument struct {
        DocumentId string `json:"document_id"`
        SchemaId string `json:"schema_id"`
        RepositoryId string `json:"repository_id"`
        InsertDate string `json:"insert_date"`
        LastUpdate string `json:"last_update"`
        IsActive bool `json:"is_active"`
        Content map[string]interface{} `json:"content,omitempty"`
    }

    // DocumentResponse will be marshalled to create and API-like response
    type DocumentResponse struct {
        Document ResponseInnerDocument `json:"document"`
    }

    // DocumentsResponse will be marshalled to crete an API-like response
    type DocumentsResponse struct {
        Count int `json:"count"`
        TotalCount int `json:"total_count"`
        Limit int `json:"limit"`
        Offset int `json:"offset"`
        Documents []ResponseInnerDocument `json:"documents"`
    }

    // init stuff
    dummyDoc := ResponseInnerDocument{
        DocumentId: uuid.New().String(),
        SchemaId: uuid.New().String(),
        RepositoryId: uuid.New().String(),
        InsertDate: "2015-02-24T21:48:16.332",
        LastUpdate: "2015-02-24T21:48:16.332",
        IsActive: false,
    }
    dummyContent := map[string]interface{}{
        "integerField": 42,
        "flaotField": 3.14,
        "stringField": "antani",
        "textField": "this is not a very long string, but should be",
        "boolField": true,
        "dateField": "1970-01-01",
        "timeField": "00:01:30",
        "datetimeField": "2001-03-08T23:31:42",
        "base64Field": "VGhpcyBpcyBhIGJhc2UtNjQgZW5jb2RlZCBzdHJpbmcu",
        "jsonField": `{"success": true}`,
        "blobField": uuid.New().String(),
        "arrayIntegerField": `[0, 1, 1, 2, 3, 5]`,
        "arrayFloatField": `[1.1, 2.2, 3.3, 4.4]`,
        "arrayStringField": `["Hello", "world", "!"]`,
    }
    dummyDoc.Content = dummyContent

    // shortcuts
    schemaId := dummyDoc.SchemaId
    docId := dummyDoc.DocumentId

    writeDocResponse := func(w http.ResponseWriter) {
        data, _ := json.Marshal(DocumentResponse{dummyDoc})
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
        if r.URL.Path == fmt.Sprintf("/api/v1/schemas/%s/documents",
            schemaId) && r.Method == "POST" {
            // mock CREATE response
            writeDocResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/documents/%s",
            docId) && r.Method == "GET" {
            // mock READ response
            writeDocResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/documents/%s",
            docId) && r.Method == "PUT" {
            // mock UPDATE response
            dummyDoc.IsActive = true
            dummyContent["stringField"] = "brematurata"
            writeDocResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/documents/%s",
            docId) && r.Method == "DELETE" {
            // mock DELETE response
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/schemas/%s/documents",
            schemaId) && r.Method == "GET" {
            // mock LIST response
            documentsResp := DocumentsResponse{
                Count: 1,
                TotalCount: 1,
                Limit: 100,
                Offset: 0,
                Documents: []ResponseInnerDocument{dummyDoc},
            }
            data, _ := json.Marshal(documentsResp)
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

    // test CREATE: we submit no content, since the response is mocked
    // we init instead a Schema with just the right ids
    schema := Schema{
        RepositoryId: dummyDoc.RepositoryId,
        Id: dummyDoc.SchemaId,
        Description: "unittest",
        Structure: []SchemaField{
            {Name: "integerField", Type: "integer"},
            {Name: "flaotField", Type: "float"},
            {Name: "stringField", Type: "string"},
            {Name: "textField", Type: "text"},
            {Name: "boolField", Type: "boolean"},
            {Name: "dateField", Type: "date"},
            {Name: "timeField", Type: "time"},
            {Name: "datetimeField", Type: "datetime"},
            {Name: "base64Field", Type: "base64"},
            {Name: "jsonField", Type: "json"},
            {Name: "blobField", Type: "blob"},
            {Name: "arrayIntegerField", Type: "array[integer]"},
            {Name: "arrayFloatField", Type: "array[float]"},
            {Name: "arrayStringField", Type: "array[string]"},
        },
    }
    content := map[string]interface{}{}
    document, err := custodia.CreateDocument(&schema, false, content)

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if document != nil {
        if (*document).RepositoryId != dummyDoc.RepositoryId {
            t.Errorf("bad RepositoryId, got: %v want: %v",
            document.RepositoryId, dummyDoc.RepositoryId)
        }
        if (*document).SchemaId != dummyDoc.SchemaId {
            t.Errorf("bad SchemaId, got: %v want: %v",
            document.SchemaId, dummyDoc.SchemaId)
        }
        if (*document).Id != dummyDoc.DocumentId {
            t.Errorf("bad DocumentId, got: %v want: %v",
            document.Id, dummyDoc.DocumentId)
        }
        if (*document).InsertDate.Year() != 2015 {
            t.Errorf("bad insert_date year, got: %v want: 2015",
                (*document).InsertDate.Year())
        }
        if (*document).LastUpdate.Year() != 2015 {
            t.Errorf("bad last_update year, got: %v want: 2015",
                (*document).InsertDate.Year())
        }
        if (*document).IsActive != false {
            t.Errorf("bad isActive, got: %v want: false", (*document).IsActive)
        }
    } else {
        t.Errorf("unexpected: both document and error are nil!")
    }

    // test READ
    doc, err := custodia.ReadDocument(schema, dummyDoc.DocumentId)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if doc != nil {
        if doc.Id != dummyDoc.DocumentId {
            t.Errorf("bad DocumentId, got: %v want: %v", doc.Id,
                dummyDoc.DocumentId)
        }
        if doc.SchemaId != dummyDoc.SchemaId {
            t.Errorf("bad SchemaId, got: %v want: %v", doc.SchemaId,
                dummyDoc.SchemaId)
        }
        if doc.RepositoryId != dummyDoc.RepositoryId {
            t.Errorf("bad RepositoryId, got: %v want: %v", doc.RepositoryId,
                dummyDoc.RepositoryId)
        }
        if doc.InsertDate.Year() != 2015 {
            t.Errorf("bad insert_date year, got: %v want: 2015",
                doc.InsertDate.Year())
        }
        if doc.LastUpdate.Year() != 2015 {
            t.Errorf("bad last_update year, got: %v want: 2015",
                doc.InsertDate.Year())
        }
        if doc.IsActive != false {
            t.Errorf("bad isActive, got: %v want: false", doc.IsActive)
        }

        // test the content
        if doc.Content["integerField"] != int64(42) {
            t.Errorf("content: bad integerField, got: %v want: %v",
                doc.Content["integerField"], int64(42))
        }
        if doc.Content["flaotField"] != 3.14 {
            t.Errorf("content: bad flaotField, got: %v want: %v",
                doc.Content["flaotField"], 3.14 )
        }
        if doc.Content["stringField"] != "antani" {
            t.Errorf("content: bad stringField, got: %v want: %v",
                doc.Content["stringField"], "antani")
        }
        if doc.Content["textField"] !=
            "this is not a very long string, but should be" {
            t.Errorf("content: bad textField, got: %v want: %v",
                doc.Content["textField"],
                "this is not a very long string, but should be")
        }
        if doc.Content["boolField"] != true {
            t.Errorf("content: bad boolField, got: %v want: %v",
                doc.Content["boolField"], true)
        }

        // for date we check just the yyyy-mm-dd part (ignoring time)
        dateField, _ := doc.Content["dateField"].(time.Time)
        if dateField.Year() != 1970 {
            t.Errorf("content: bad dateField year, got: %v want: 1970",
                dateField.Year())
        }
        if dateField.Month() != 1 {
            t.Errorf("content: bad dateField month, got: %v want: 1",
                dateField.Month())
        }
        if dateField.Day() != 1 {
            t.Errorf("content: bad dateField day, got: %v want: 1",
                dateField.Day())
        }

        // for time we check just the HH:MM:SS part (ignoring date)
        timeField, _ := doc.Content["timeField"].(time.Time)
        if timeField.Hour() != 0 {
            t.Errorf("content: bad timeField hour, got: %v want: 0",
                timeField.Hour())
        }
        if timeField.Minute() != 1 {
            t.Errorf("content: bad timeField minute, got: %v want: 1",
                timeField.Minute())
        }
        if timeField.Second() != 30 {
            t.Errorf("content: bad timeField second, got: %v want: 30",
                timeField.Second())
        }

        // for datetime we check all
        dateTimeField, _ := doc.Content["datetimeField"].(time.Time)
        if dateTimeField.Year() != 2001 {
            t.Errorf("content: bad dateTimeField year, got: %v want: 2001",
                dateTimeField.Year())
        }
        if dateTimeField.Month() != 3 {
            t.Errorf("content: bad dateTimeField month, got: %v want: 3",
                dateTimeField.Month())
        }
        if dateTimeField.Day() != 8 {
            t.Errorf("content: bad dateTimeField day, got: %v want: 8",
                dateTimeField.Day())
        }
        if dateTimeField.Hour() != 23 {
            t.Errorf("content: bad dateTimeField hour, got: %v want: 23",
                dateTimeField.Hour())
        }
        if dateTimeField.Minute() != 31 {
            t.Errorf("content: bad dateTimeField minute, got: %v want: 31",
                dateTimeField.Minute())
        }
        if dateTimeField.Second() != 42 {
            t.Errorf("content: bad dateTimeField second, got: %v want: 42",
                dateTimeField.Second())
        }

        // base64, json, and blob are just string fields. We don't care if the
        // string in it actually converts to real data. This is a problem of
        // the user, his duty to encode/decode data correctly
        if doc.Content["base64Field"] !=
            "VGhpcyBpcyBhIGJhc2UtNjQgZW5jb2RlZCBzdHJpbmcu" {
            t.Errorf("content: bad base64Field, got: %v want: %v",
            doc.Content["base64Field"],
            "VGhpcyBpcyBhIGJhc2UtNjQgZW5jb2RlZCBzdHJpbmcu")
        }
        if doc.Content["jsonField"] != `{"success": true}` {
            t.Errorf("content: bad jsonField, got: %v want: {\"success\": true}",
                doc.Content["jsonField"])
        }
        wantBlobField, _ := dummyContent["blobField"].(string)
        if doc.Content["blobField"] != wantBlobField {
            t.Errorf("content: bad blobField, got: %v want: %v",
                doc.Content["blobField"], wantBlobField)
        }

        // test array fields
        if reflect.DeepEqual(doc.Content["arrayIntegerField"],
            []int{0, 1, 2, 3, 4, 5}) {
            t.Errorf("content: bad arrayIntegerField, got: %v want: %v",
                doc.Content["arrayIntegerField"], []int{0, 1, 2, 3, 4, 5})
        }
        if reflect.DeepEqual(doc.Content["arrayFloatField"],
            []float64{1.1, 2.2, 3.3, 4.4}) {
            t.Errorf("content: bad arrayIntegerField, got: %v want: %v",
                doc.Content["arrayFloatField"], []float64{1.1, 2.2, 3.3, 4.4})
        }
        if reflect.DeepEqual(doc.Content["arrayStringField"],
            []string{"Hello", "world", "!"}) {
            t.Errorf("content: bad arrayIntegerField, got: %v want: %v",
                doc.Content["arrayStringField"],
                []string{"Hello", "world", "!"})
        }
    } else {
        t.Errorf("unexpected: both document and error are nil!")
    }

    // test UPDATE
    // we still pass empty content, non influential, response is mocked
    doc, err = custodia.UpdateDocument(doc.Id, true, content)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if doc != nil {
        if doc.IsActive != true {
            t.Errorf("bad isActive, got: %v want: true", doc.IsActive)
        }
        if doc.Content["stringField"] != "brematurata" {
            t.Errorf("bad stringField, got: %v want: brematurata",
                doc.Content["stringField"])
        }
    } else {
        t.Errorf("unexpected: both document and error are nil!")
    }

    // test DELETE
    err = custodia.DeleteDocument(doc.Id, true, true)
    if err != nil {
        t.Errorf("error while deleting document. Details: %v", err)
    }

    // test LIST
    // test we gave a wrong argument
    params := map[string]interface{}{"antani": 42}
    _, err = custodia.ListDocuments(schema.Id, params)
    if err == nil {
        t.Errorf("ListDocuments is not giving error with wrong param %v",
            params)
    }
    // test that all the other params are accepted instead
    goodParams := map[string]interface{}{
		"full_document": true,
		"is_active": true,
		"insert_date__gt": time.Time{},
		"insert_date__lt": time.Time{},
		"last_update__gt": time.Time{},
		"last_update__lt": time.Time{},
	}
    documents, err := custodia.ListDocuments(schema.Id, goodParams)
    if err != nil {
        t.Errorf("error while listing documents. Details: %v", err)
    } else if reflect.TypeOf(documents) != reflect.TypeOf([]*Document{}) {
        t.Errorf("documents is not list of Documents, got: %T want: %T",
            documents, []*Document{})
    }
}

func TestApplicationCRULD(t *testing.T) {
    // ResponseInnerApp will be included in responses
    type ResponseInnerApp struct {
        AppSecret string `json:"app_secret"`
        ClientType string `json:"client_type"`
        GrantType string `json:"grant_type"`
        AppName string `json:"app_name"`
        RedirectUrl string `json:"redirect_url"`
        AppId string `json:"app_id"`
    }

    type ApplicationResponse struct {
        Application ResponseInnerApp `json:"application"`
    }

    type ApplicationsResponse struct {
        Count int `json:"count"`
        TotalCount int `json:"total_count"`
        Limit int `json:"limit"`
        Offset int `json:"offset"`
        Applications []ResponseInnerApp `json:"applications"`
    }

    // init stuff
    aid := "MyAppId42"
    dummyApp := ResponseInnerApp{
        AppId: aid,
        AppSecret: "123456",
        ClientType: "public",
        GrantType: "password",
        AppName: "antani",
        RedirectUrl: "",
    }

    writeAppResponse := func(w http.ResponseWriter) {
        data, _ := json.Marshal(ApplicationResponse{dummyApp})
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
        if r.URL.Path == "/api/v1/auth/applications" && r.Method == "POST" {
            // mock CREATE response
            writeAppResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/auth/applications/%s",
            aid) && r.Method == "GET" {
            // mock READ response
            writeAppResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/auth/applications/%s",
            aid) && r.Method == "PUT" {
            // mock UPDATE response
            dummyApp.GrantType = GrantAuthorizationCode.String()
            dummyApp.ClientType = ClientConfidential.String()
            dummyApp.RedirectUrl = "http://antani.org"
            writeAppResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/auth/applications/%s",
            aid) && r.Method == "DELETE" {
            // mock DELETE response
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == "/api/v1/auth/applications" &&
            r.Method == "GET" {
            // mock LIST response
            appsResp := ApplicationsResponse{
                Count: 1,
                TotalCount: 1,
                Limit: 100,
                Offset: 0,
                Applications: []ResponseInnerApp{dummyApp},
            }
            data, _ := json.Marshal(appsResp)
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
    app, err := custodia.CreateApplication("antani", GrantPassword,
        ClientConfidential, "")

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if app != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyApp.AppId, app.Id},
            {GrantPassword, app.GrantType},
            {dummyApp.AppName, app.Name},
            {dummyApp.AppSecret, app.Secret},

        }
        for _, test := range tests {
            if test.want != test.got {
                t.Errorf("Application CREATE: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    }

    // test READ
    app, err = custodia.ReadApplication(dummyApp.AppId)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if app != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyApp.AppId, app.Id},
            {GrantPassword, app.GrantType},
            {dummyApp.AppName, app.Name},
            {dummyApp.AppSecret, app.Secret},

        }
        for _, test := range tests {
            if test.want != test.got {
                t.Errorf("Application GET: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    }

    // test UPDATE
    app, err = custodia.UpdateApplication(dummyApp.AppId, "antani",
        GrantAuthorizationCode, ClientConfidential, "http://antani.org")
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if app != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyApp.AppId, app.Id},
            {ClientConfidential, app.ClientType},
            {GrantAuthorizationCode, app.GrantType},
            {"antani", app.Name},
            {dummyApp.AppSecret, app.Secret},

        }
        for _, test := range tests {
            if test.want != test.got {
                t.Errorf("Application GET: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    }

    // test DELETE
    err = custodia.DeleteApplication(dummyApp.AppId)
    if err != nil {
        t.Errorf("error while deleting application. Details: %v", err)
    }

    // test LIST
    apps, err := custodia.ListApplications()
    if err != nil {
        t.Errorf("error while listing applications. Details: %v", err)
    } else if reflect.TypeOf(apps) != reflect.TypeOf([]*Application{}) {
        t.Errorf("apps is not list of Applications, got: %T want: %T",
            apps, []*Application{})
    }
}

func TestUserSchemaCRUDL(t *testing.T) {
    // ResponseInnerUserSchema will be included in responses
    type ResponseInnerUserSchema struct {
        UserSchemaId string `json:"user_schema_id"`
        Description string `json:"description"`
        Groups []string `json:"groups"`
        InsertDate string `json:"insert_date"`
        LastUpdate string `json:"last_update"`
        IsActive bool `json:"is_active"`
        Structure []SchemaField `json:"structure"`
    }

    // SchemaResponse will be marshalled to create an API-like response
    type UserSchemaResponse struct {
        Schema ResponseInnerUserSchema `json:"user_schema"`
    }

    // SchemasResponse will be marshalled to create an API-like response
    type UserSchemasResponse struct {
        Count int `json:"count"`
        TotalCount int `json:"total_count"`
        Limit int `json:"limit"`
        Offset int `json:"offset"`
        UserSchemas []ResponseInnerUserSchema `json:"user_schemas"`
    }

    dummyUserSchema := ResponseInnerUserSchema{
        UserSchemaId: uuid.New().String(),
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

    writeUserSchemaResponse := func(w http.ResponseWriter) {
        data, _ := json.Marshal(UserSchemaResponse{dummyUserSchema})
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
        if r.URL.Path == "/api/v1/user_schemas" && r.Method == "POST" {
            // mock CREATE response
            writeUserSchemaResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/user_schemas/%s",
            dummyUserSchema.UserSchemaId) && r.Method == "GET" {
            // mock READ response
            writeUserSchemaResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/user_schemas/%s",
            dummyUserSchema.UserSchemaId) && r.Method == "PUT" {
            // mock UPDATE response
            dummyUserSchema.Description = "changed"
            dummyUserSchema.IsActive = true
            // dummyUserSchema.Structure[0].Default = 21
            writeUserSchemaResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/user_schemas/%s",
            dummyUserSchema.UserSchemaId) && r.Method == "DELETE" {
            // mock DELETE response
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == "/api/v1/user_schemas" &&  r.Method == "GET" {
            // mock LIST response
            schemasResp := UserSchemasResponse{
                Count: 1,
                TotalCount: 1,
                Limit: 100,
                Offset: 0,
                UserSchemas: []ResponseInnerUserSchema{dummyUserSchema},
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

    userSchema, err := custodia.CreateUserSchema("unittest", true, structure);

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if userSchema != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyUserSchema.UserSchemaId, userSchema.Id},
            {dummyUserSchema.Description, userSchema.Description},
            {dummyUserSchema.Groups, userSchema.Groups},
            {2015, userSchema.InsertDate.Year()},
            {2015, userSchema.LastUpdate.Year()},
            {false, userSchema.IsActive},
        }
        for _, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("UserSchema CREATE: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    } else {
        t.Errorf("unexpected: both userSchema and error are nil!")
    }

    // test READ
    userSchema, err = custodia.ReadUserSchema(dummyUserSchema.UserSchemaId)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if userSchema != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyUserSchema.UserSchemaId, userSchema.Id},
            {dummyUserSchema.Description, userSchema.Description},
            {dummyUserSchema.Groups, userSchema.Groups},
            {2015, userSchema.InsertDate.Year()},
            {2015, userSchema.LastUpdate.Year()},
            {false, userSchema.IsActive},
        }
        for _, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("UserSchema CREATE: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    } else {
        t.Errorf("unexpected: both userSchema and error are nil!")
    }

    // test UPDATE
    userSchema, err = custodia.UpdateUserSchema(dummyUserSchema.UserSchemaId,
        "antani2", true, dummyUserSchema.Structure)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if userSchema != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyUserSchema.UserSchemaId, userSchema.Id},
            {dummyUserSchema.Description, userSchema.Description},
            {dummyUserSchema.Groups, userSchema.Groups},
            {2015, userSchema.InsertDate.Year()},
            {2015, userSchema.LastUpdate.Year()},
            {true, userSchema.IsActive},
        }
        for _, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("UserSchema CREATE: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    } else {
        t.Errorf("unexpected: both userSchema and error are nil!")
    }

    // test DELETE
    err = custodia.DeleteUserSchema(dummyUserSchema.UserSchemaId, false)
    if err != nil {
        t.Errorf("error while deleting user schema. Details: %v", err)
    }

    // test LIST
    uss, err := custodia.ListUserSchemas()
    if err != nil {
        t.Errorf("error while listing user schemas. Details: %v", err)
    } else if reflect.TypeOf(uss) != reflect.TypeOf([]*UserSchema{}) {
        t.Errorf("uss is not list of UserSchemas, got: %T want: %T",
            uss, []*UserSchema{})
    }
}

func TestUserCRUDL(t *testing.T) {
    // ResponseInnerUser will be included in responses
    type ResponseInnerUser struct {
        UserId string `json:"user_id"`
        UserSchemaId string `json:"schema_id"`
        Username string `json:"username"`
        InsertDate string `json:"insert_date"`
        LastUpdate string `json:"last_update"`
        IsActive bool `json:"is_active"`
        Attributes map[string]interface{} `json:"content,omitempty"`
        Groups []string `json:"groups"`
    }

    // UserResponse will be marshalled to create and API-like response
    type UserResponse struct {
        User ResponseInnerUser `json:"user"`
    }

    // UsersResponse will be marshalled to crete an API-like response
    type UsersResponse struct {
        Count int `json:"count"`
        TotalCount int `json:"total_count"`
        Limit int `json:"limit"`
        Offset int `json:"offset"`
        Users []ResponseInnerUser `json:"users"`
    }

    // init stuff
    dummyUser := ResponseInnerUser{
        UserId: uuid.New().String(),
        UserSchemaId: uuid.New().String(),
        Username: "unittest",
        InsertDate: "2015-02-24T21:48:16.332",
        LastUpdate: "2015-02-24T21:48:16.332",
        IsActive: false,
        Groups: []string{},
    }
    dummyAttrs := map[string]interface{}{
        "integerField": 42,
        "flaotField": 3.14,
        "stringField": "antani",
        "textField": "this is not a very long string, but should be",
        "boolField": true,
        "dateField": "1970-01-01",
        "timeField": "00:01:30",
        "datetimeField": "2001-03-08T23:31:42",
        "base64Field": "VGhpcyBpcyBhIGJhc2UtNjQgZW5jb2RlZCBzdHJpbmcu",
        "jsonField": `{"success": true}`,
        "blobField": uuid.New().String(),
        "arrayIntegerField": `[0, 1, 1, 2, 3, 5]`,
        "arrayFloatField": `[1.1, 2.2, 3.3, 4.4]`,
        "arrayStringField": `["Hello", "world", "!"]`,
    }
    dummyUser.Attributes = dummyAttrs

    // shortcuts
    userSchemaId := dummyUser.UserSchemaId
    userId := dummyUser.UserId

    writeDocResponse := func(w http.ResponseWriter) {
        data, _ := json.Marshal(UserResponse{dummyUser})
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
        if r.URL.Path == fmt.Sprintf("/api/v1/user_schemas/%s/users",
            userSchemaId) && r.Method == "POST" {
            // mock CREATE response
            writeDocResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/users/%s",
            userId) && r.Method == "GET" {
            // mock READ response
            writeDocResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/users/%s",
            userId) && r.Method == "PUT" {
            // mock UPDATE response
            dummyUser.IsActive = true
            dummyAttrs["stringField"] = "brematurata"
            writeDocResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/users/%s",
            userId) && r.Method == "DELETE" {
            // mock DELETE response
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/user_schemas/%s/users",
            userSchemaId) && r.Method == "GET" {
            // mock LIST response
            usersResp := UsersResponse{
                Count: 1,
                TotalCount: 1,
                Limit: 100,
                Offset: 0,
                Users: []ResponseInnerUser{dummyUser},
            }
            data, _ := json.Marshal(usersResp)
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            envelope.Data = data
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
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

    // test CREATE: we submit no content, since the response is mocked
    // we init instead a UserSchema with just the right ids
    userSchema := UserSchema{
        Id: dummyUser.UserSchemaId,
        Description: "unittest",
        IsActive: true,
        Structure: []SchemaField{},
    }
    attributes := map[string]interface{}{}
    user, err := custodia.CreateUser(&userSchema, false, attributes)

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if user != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyUser.UserId, user.Id},
            {dummyUser.UserSchemaId, userSchema.Id},
            {dummyUser.Username, user.Username},
            {2015, user.InsertDate.Year()},
            {2, int(user.LastUpdate.Month())},
            {false, user.IsActive},
            {reflect.TypeOf(map[string]interface{}{}), reflect.TypeOf(user.Attributes)},
            {reflect.TypeOf([]string{}), reflect.TypeOf(user.Groups)},
        }
        for _, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("Users CREATE: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    } else {
        t.Errorf("unexpected: both user and error are nil!")
    }

    // test UPDATE
    user, err = custodia.UpdateUser(dummyUser.UserId, true, attributes)

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if user != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyUser.UserId, user.Id},
            {dummyUser.UserSchemaId, userSchema.Id},
            {dummyUser.Username, user.Username},
            {2015, user.InsertDate.Year()},
            {2, int(user.LastUpdate.Month())},
            {true, user.IsActive},
            {reflect.TypeOf(map[string]interface{}{}), reflect.TypeOf(user.Attributes)},
            {reflect.TypeOf([]string{}), reflect.TypeOf(user.Groups)},
        }
        for _, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("Users CREATE: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    } else {
        t.Errorf("unexpected: both user and error are nil!")
    }

    // test DELETE
    err = custodia.DeleteUser(dummyUser.UserId, false, false)
    if err != nil {
        t.Errorf("error while deleting user. Details: %v", err)
    }

    // test LIST
    // test we gave a wrong argument
    params := map[string]interface{}{"antani": 42}
    _, err = custodia.ListUsers(dummyUser.UserSchemaId, params)
    if err == nil {
        t.Errorf("ListUsers is not giving error with wrong param %v",
            params)
    }

    // test that all the other params are accepted instead
    goodParams := map[string]interface{}{
        "full_user": true,
        "is_active": true,
        "insert_date__gt": time.Time{},
        "insert_date__lt": time.Time{},
        "last_update__gt": time.Time{},
        "last_update__lt": time.Time{},
    }
    users, err := custodia.ListUsers(dummyUser.UserSchemaId, goodParams)

    if err != nil {
        t.Errorf("error while listing users: %v", err)
    } else if reflect.TypeOf(users) != reflect.TypeOf([]*User{}) {
        t.Errorf("users is not list of Users, got: %T want: %T",
            users, []*User{})
    }
}

func TestOAuth(t *testing.T) {
    envelope := CustodiaEnvelope{
        Result: "success",
        ResultCode: 200,
        Message: nil,
    }
    responseLogin := map[string]interface{}{
        "access_token": "ans2fN08sliGpIOLMGg3fv4BpPhWRq",
        "token_type": "Bearer",
        "expires_in": 36000,
        "refresh_token": "vL0durAhdhNNYFI27F3zGGHXeNLwcO",
        "scope": "read write",
    }
    responseRefresh := map[string]interface{}{
        "access_token": "Qg3fv4BpPhWRqXeNLwcOa2fN08sliGpIOLMg3",
        "token_type": "Bearer",
        "expires_in": 36000,
        "refresh_token": "vL0durAhdhNNYFI27F3zGGHXeNLwcO",
        "scope": "read write",
    }

    mockHandler := func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/api/v1/auth/token" && r.Method == "POST" {
            data, _ := json.Marshal(responseLogin)
            envelope.Data = data
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == "/api/v1/auth/refresh" && r.Method == "POST" {
            data, _ := json.Marshal(responseRefresh)
            envelope.Data = data
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
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
    app := Application{
        Id: "test",
        Secret: "test",
        ClientType: ClientConfidential,
    }

    // test LOGIN
    err := custodia.LoginUser("test", "test", app)
    auth := custodia.client.GetAuth()
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {common.UserAuth, auth.GetAuthType()},
            {"ans2fN08sliGpIOLMGg3fv4BpPhWRq", auth.GetAccessToken()},
            {"vL0durAhdhNNYFI27F3zGGHXeNLwcO", auth.GetRefreshToken()},
            // Go is super quick, so this should be true
            {36000, auth.GetAccessTokenExpire() - int(time.Now().Unix())},
        }

        for _, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("User Login: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    }

    // test REFRESH Token
    err = custodia.RefreshToken(app)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if auth != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {"Qg3fv4BpPhWRqXeNLwcOa2fN08sliGpIOLMg3", auth.GetAccessToken()},
            {"vL0durAhdhNNYFI27F3zGGHXeNLwcO", auth.GetRefreshToken()},
            // Go is super quick, so this should be true
            {36000, auth.GetAccessTokenExpire() - int(time.Now().Unix())},
        }

        for _, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("User RefreshToken: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    }
}

func TestGroupCRUDL(t *testing.T) {
    envelope := CustodiaEnvelope{
        Result: "success",
        ResultCode: 200,
        Message: nil,
    }

    dummyGroup := map[string]interface{}{
        "group_id": uuid.New().String(),
        "group_name": "unittest",
        "attributes": map[string]interface{}{"antani": 3.14},
        "is_active": true,
        "insert_date": "2015-02-07T12:14:46.754",
        "last_update": "2015-03-13T18:06:21.242",
    }
    gid, _ := dummyGroup["group_id"].(string)

    userSchemaId := uuid.New().String()
    dummyUser := map[string]interface{}{
        "username": "unittest",
        "schemas_id": userSchemaId,
        "user_id": uuid.New().String(),
        "insert_date": "2015-02-07T12:14:46.754",
        "last_update": "2015-03-13T18:06:21.242",
        "is_active": true,
        "attributes": map[string]interface{}{"antani": 3.14},
        "groups": []string{gid},
    }
    uid, _ := dummyUser["user_id"].(string)

    responseGroup := map[string]interface{}{
        "group": dummyGroup,
    }

    mockHandler := func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/api/v1/groups" && r.Method == "POST" {
            // mock CREATE response
            envelope.Data, _ = json.Marshal(responseGroup)
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/groups/%s", gid) &&
            // mock READ response
            r.Method == "GET" {
            data, _ := json.Marshal(responseGroup)
            envelope.Data = data
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/groups/%s", gid ) &&
            r.Method == "PUT" {
            // mock UPDATE response
            dummyGroup["group_name"] = "changed"
            dummyGroup["is_active"] = false
            dummyGroup["attributes"] = map[string]interface{}{
                "antani": 3.14, "something": "else"}
            envelope.Data, _ = json.Marshal(responseGroup)
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/groups/%s",
            gid) && r.Method == "DELETE" {
            // mock DELETE response
            data, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(data)
        } else if r.URL.Path == "/api/v1/groups" && r.Method == "GET" {
            // mock LIST response
            groupsResp := map[string]interface{}{
                "count": 1,
                "total_count": 1,
                "limit": 1,
                "offset": 0,
                "groups": []interface{}{dummyGroup},
            }
            data, _ := json.Marshal(groupsResp)
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            envelope.Data = data
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/groups/%s/users", gid) &&
            r.Method == "GET" {
            listResp := map[string]interface{}{
                "count": 1,
                "total_count": 1,
                "limit": 1,
                "offset": 0,
                "users": []map[string]interface{}{dummyUser},
            }
            data, _ := json.Marshal(listResp)
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            envelope.Data = data
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/groups/%s/users/%s", gid,
            uid) {
            // we don't care the method. GET, POST or DELETE, they always
            // return just 200 success with empty message
            data, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(data)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/groups/%s/user_schemas/%s", gid, userSchemaId) {
            // we don't care the method. GET, POST or DELETE, they always
            // return just 200 success with empty message
            data, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(data)
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

    // init stuff
    client := common.NewClient(server.URL, common.GetFakeAuth())
    custodia := NewCustodiaAPIv1(client)

    // test CREATE
    group, err := custodia.CreateGroup("unittest", true,
        map[string]interface{}{})
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {gid, group.Id},
            {"unittest", group.Name},
            {map[string]interface{}{"antani": 3.14}, group.Attributes},
            {true, group.IsActive},
            {2015, group.InsertDate.Year()},
            {2, int(group.InsertDate.Month())},
            {7, int(group.InsertDate.Day())},
            {12, int(group.InsertDate.Hour())},
            {14, int(group.InsertDate.Minute())},
            {46, int(group.InsertDate.Second())},
            {2015, group.LastUpdate.Year()},
            {3, int(group.LastUpdate.Month())},
            {13, int(group.LastUpdate.Day())},
            {18, int(group.LastUpdate.Hour())},
            {6, int(group.LastUpdate.Minute())},
            {21, int(group.LastUpdate.Second())},
        }

        for _, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("Group Create: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    }

    // test UPDATE
    // response is mocked, so we don't need to pass the right data
    group, err = custodia.UpdateGroup(gid, "changed", false,
        map[string]interface{}{})
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {gid, group.Id},
            {"changed", group.Name},
            {false, group.IsActive},
            {map[string]interface{}{"antani": 3.14, "something": "else"},
                group.Attributes},
            {2015, group.InsertDate.Year()},
            {2, int(group.InsertDate.Month())},
            {7, int(group.InsertDate.Day())},
            {12, int(group.InsertDate.Hour())},
            {14, int(group.InsertDate.Minute())},
            {46, int(group.InsertDate.Second())},
            {2015, group.LastUpdate.Year()},
            {3, int(group.LastUpdate.Month())},
            {13, int(group.LastUpdate.Day())},
            {18, int(group.LastUpdate.Hour())},
            {6, int(group.LastUpdate.Minute())},
            {21, int(group.LastUpdate.Second())},
        }

        for _, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("Group Update: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    }

    // test DELETE
    err = custodia.DeleteGroup(gid, false)
    if err != nil {
        t.Errorf("error while deleting group: %v", err)
    }
    err = custodia.DeleteGroup(gid, true)
    if err != nil {
        t.Errorf("error while deleting group: %v", err)
    }

    // test LIST
    groups, err := custodia.ListGroups()
    if err != nil {
        t.Errorf("error while listing groups: %v", err)
    } else if reflect.TypeOf(groups) != reflect.TypeOf([]Group{}) {
        t.Errorf("groups is not list of Groups, got: %T want: %T",
            groups, []*Group{})
    }

    // Group members tests
    users, err := custodia.ListGroupUsers(gid)
    if err!= nil {
        t.Errorf("error while listing users in group: %v", err)
    } else {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {uid, users[0].Id},
            {"unittest", users[0].Username},
            {2015, int(users[0].InsertDate.Year())},
            {2, int(users[0].InsertDate.Month())},
            {7, int(users[0].InsertDate.Day())},
            {12, int(users[0].InsertDate.Hour())},
            {14, int(users[0].InsertDate.Minute())},
            {46, int(users[0].InsertDate.Second())},
            {2015, users[0].LastUpdate.Year()},
            {3, int(users[0].LastUpdate.Month())},
            {13, int(users[0].LastUpdate.Day())},
        }

        for _, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("Group Create: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    }

    // Add user to group
    err = custodia.AddUserToGroup(uid, gid)
    if err!= nil {
        t.Errorf("error while adding user to group: %v", err)
    }

    // Remove user from group
    err = custodia.RemoveUserFromGroup(uid, gid)
    if err!= nil {
        t.Errorf("error while removing user from group: %v", err)
    }

    // Add all users of schema to group
    err = custodia.AddUsersFromUserSchemaToGroup(userSchemaId, gid)
    if err!= nil {
        t.Errorf("error while adding users to group from schema: %v", err)
    }

    // Remove all users of schema from group
    err = custodia.RemoveUsersFromUserSchemaFromGroup(userSchemaId, gid)
    if err!= nil {
        t.Errorf("error while removing users from group from schema: %v", err)
    }
}

func TestPermissions(t *testing.T) {
    envelope := CustodiaEnvelope{
        Result: "success",
        ResultCode: 200,
        Message: nil,
    }

    dummyUUID := uuid.New().String()

    // dummy data to return from ReadAllPermissions
    allPermissions := []map[string]interface{}{
        {
            "access": "Structure",
            "parent_id": nil,
            "resource_type": "Repository",
            "owner_id": dummyUUID,
            "owner_type": "users",
            "permission": map[string][]string{
                "Manage": []string{
                  "R",
                },
            },
        },
        {
            "access": "Data",
            "resource_id": dummyUUID,
            "resource_type": "Schema",
            "owner_id": dummyUUID,
            "owner_type": "users",
            "permission": map[string][]string{
                "Authorize": []string{
                  "A",
                },
                "Manage": []string{
                  "R",
                  "U",
                  "L",
                },
            },
        },
    }

    mockHandler := func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == fmt.Sprintf(
            "/api/v1/perms/grant/repositories/users/%s", dummyUUID) &&
            r.Method == "POST" {
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/perms/grant/repositories/%s/groups/%s", dummyUUID,
            dummyUUID) && r.Method == "POST" {
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/perms/revoke/repositories/%s/schemas/groups/%s",
            dummyUUID, dummyUUID) && r.Method == "POST" {
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/apit/v1/perms") &&
            r.Method == "GET" {

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

    // Test Permission on Resources (multiple)
    perms := map[PermissionScope][]PermissionType{
        PermissionScopeManage: {PermissionTypeCreate, PermissionTypeList,
            PermissionTypeRead,},
        PermissionScopeAuthorize: {},
    }
    err := custodia.PermissionOnResources(PermissionActionGrant,
        ResourceRepository, ResourceUser, dummyUUID, perms)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    // Test Permission on Resource (single)
    err = custodia.PermissionOnResource(PermissionActionGrant,
        ResourceRepository, dummyUUID, ResourceGroup, dummyUUID, perms)
    if err!= nil {
        t.Errorf("unexpected error: %v", err)
    }

    // Test Permission on Resource children
    err = custodia.PermissionOnResourceChildren(PermissionActionRevoke,
        ResourceRepository, dummyUUID, ResourceSchema, ResourceGroup,
        dummyUUID, perms)
    if err!= nil {
        t.Errorf("unexpected error: %v", err)
    }

    allPerms, err := custodia.ReadAllPermissions()
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else {
        var tests = []struct {
            want interface{}
            got interface{}
        }{

        }
        for _, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("Group Create: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    }


}
