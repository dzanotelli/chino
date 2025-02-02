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
