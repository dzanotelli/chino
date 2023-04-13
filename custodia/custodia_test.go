package custodia

import (
	"encoding/json"
	"errors"
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
        if repo.RepositoryId != dummyRepository.RepositoryId {
            t.Errorf("bad RepositoryId, got: %v want: %v", 
                     repo.RepositoryId, dummyRepository.RepositoryId)
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
    repo, err = custodia.UpdateRepository(repo.RepositoryId, "changed", false)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if repo != nil {
        if repo.RepositoryId != dummyRepository.RepositoryId {
            t.Errorf("bad RepositoryId, got: %v want: %v", 
                     repo.RepositoryId, dummyRepository.RepositoryId)
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
    err = custodia.DeleteRepository(repo.RepositoryId, true)
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
    if repos[0].RepositoryId != dummyRepository.RepositoryId {
        t.Errorf("bad repository id, got: %v want: %v", 
            dummyRepository.RepositoryId, repos[0].RepositoryId)
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
            dummySchema.Structure[0].Default = 21    
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

    auth := common.NewClientAuth()  // auth is tested elsewhere
    client := common.NewClient(server.URL, auth)
    custodia := NewCustodiaAPIv1(client)

    // test CREATE: we submit an empty field list, since the response is mocked
    // and we will still get a working structure. The purpose here is to test
    // that the received data are correctly populating the objects
    structure := []SchemaField{}
    // we init a Repository with just the right id, don't need other data
    repo := Repository{RepositoryId: dummySchema.RepositoryId}
    schema, err := custodia.CreateSchema(&repo, "unittest", true, structure)

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if schema != nil {
        if schema.RepositoryId != repoId {
            t.Errorf("bad RepositoryId, got: %v want: %v", 
                schema.RepositoryId, repoId)
        }
        if schema.SchemaId != dummySchema.SchemaId {
            t.Errorf("bad SchemaId, got: %v want: %v", 
                schema.SchemaId, dummySchema.SchemaId)
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
    schema, err = custodia.UpdateSchema(schema.SchemaId, "changed", true, 
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
    err = custodia.DeleteSchema(schema.SchemaId, true, true)
    if err != nil {
        t.Errorf("error while listing schemas. Details: %v", err)
    }

    // test LIST
    schemas, err := custodia.ListSchemas(repoId)
    if err != nil {
        t.Errorf("error while listing schemas. Details: %v", err)
    }
    if len(schemas) != 1 {
        t.Errorf("bad schemas lenght, got: %v want: 1", len(schemas))
    }
    if schemas[0].SchemaId != dummySchema.SchemaId {
        t.Errorf("bad schema id, got: %v want: %v", 
        dummySchema.SchemaId, schemas[0].SchemaId)
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
    // this can be added to dummyDoc later
    dummyContent := map[string]interface{}{
        "integerField": 42,
        "flaotField": 3.14,
        "stringField": "antani",
        "textField": "this is not a very long string, but should be",
        "boolField": true,
        "dateField": "1970-01-01",
        "timeField": "00:01:30",
        "datetimeField": "1970-01-01T00:01:30",
        "base64Field": "VGhpcyBpcyBhIGJhc2UtNjQgZW5jb2RlZCBzdHJpbmcu",
        "jsonField": `{"success": true}`,
        "blobField": uuid.New().String(),
        "arrayIntegerField": `[0, 1, 1, 2, 3, 5]`,
        "arrayFloatField": `[1.1, 2.2, 3.3, 4.4]`,
        "arrayStringField": `["Hello", "world", "!"]`,
    }

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
            dummyDoc.Content = dummyContent
            writeDocResponse(w)
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

    // test CREATE: we submit no content, since the response is mocked
    // we init instead a Schema with just the right ids
    schema := Schema{
        RepositoryId: dummyDoc.RepositoryId, 
        SchemaId: dummyDoc.SchemaId,
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
        if (*document).DocumentId != dummyDoc.DocumentId {
            t.Errorf("bad DocumentId, got: %v want: %v", 
            document.DocumentId, dummyDoc.DocumentId)
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
        if doc.DocumentId != dummyDoc.DocumentId {
            t.Errorf("bad DocumentId, got: %v want: %v", doc.DocumentId,
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

        // check the content
        // check the content types  
        contentErrs := validateContent(doc.Content, schema.getStructureAsMap())
        if len(contentErrs) > 0 {
            e := fmt.Errorf("content errors: %w", errors.Join(contentErrs...))
            t.Errorf(fmt.Sprintf("%v", e))
        }
        
        // FIXME: add checks on the values (?)

    } else {
        t.Errorf("unexpected: both document and error are nil!")
    }

}