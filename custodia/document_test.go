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
    envelope := CustodiaEnvelope{
        Result: "success",
        ResultCode: 200,
        Message: nil,
    }
    dummyUUID := uuid.New()

    docCreateResponse := map[string]interface{}{
        "document_id": dummyUUID.String(),
        "schema_id": dummyUUID.String(),
        "repository_id": dummyUUID.String(),
        "insert_date": "2015-04-14T05:09:54.915Z",
        "last_update": "2015-04-14T05:09:54.915Z",
        "is_active": true,
    }
    docUpdateResponse := map[string]interface{}{
        "document_id": dummyUUID.String(),
        "schema_id": dummyUUID.String(),
        "repository_id": dummyUUID.String(),
        "insert_date": "2025-04-14T05:09:54.915Z",
        "last_update": "2025-04-14T05:09:54.915Z",
        "is_active": false,
    }    // // ResponseInnerDocument will be included in responses
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
        "blobField": dummyUUID.String(),
        "arrayIntegerField": `[0, 1, 2, 3, 5]`,
        "arrayFloatField": `[1.1, 2.2, 3.3, 4.4]`,
        "arrayStringField": `["Hello", "world", "!"]`,
    }
    docCreateResponse["content"] = dummyContent
    docUpdateResponse["content"] = dummyContent

    // mock calls
    mockHandler := func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == fmt.Sprintf("/api/v1/schemas/%s/documents",
            dummyUUID) && r.Method == "POST" {
            // mock CREATE response
            w.WriteHeader(http.StatusOK)
            data := map[string]interface{}{
                "document": docCreateResponse,
            }
            envelope.Data, _ = json.Marshal(data)
			out, _ := json.Marshal(envelope)
			w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/documents/%s",
            dummyUUID) && r.Method == "GET" {
            // mock READ response
            w.WriteHeader(http.StatusOK)
            data := map[string]interface{}{
                "document": docCreateResponse,
            }
            envelope.Data, _ = json.Marshal(data)
			out, _ := json.Marshal(envelope)
			w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/documents/%s",
            dummyUUID) && r.Method == "PUT" {
            // mock UPDATE response
            w.WriteHeader(http.StatusOK)
            dummyContent["stringField"] = "brematurata"
            data := map[string]interface{}{
                "document": docUpdateResponse,
            }
            envelope.Data, _ = json.Marshal(data)
			out, _ := json.Marshal(envelope)
			w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/documents/%s",
            dummyUUID) && r.Method == "DELETE" {
            // mock DELETE response
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            envelope.Data = nil
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/schemas/%s/documents",
            dummyUUID) && r.Method == "GET" {
            // mock LIST response
            data := map[string]interface{}{
                "count": 1,
                "total_count": 1,
                "limit": 100,
                "offset": 0,
                "documents": []map[string]interface{}{
                    docCreateResponse,
                    docUpdateResponse,
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

    // test CREATE: we submit no content, since the response is mocked
    // we init instead a Schema with just the right ids
    schema := Schema{
        RepositoryId: dummyUUID,
        Id: dummyUUID,
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
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyUUID.String(), document.RepositoryId.String()},
            {dummyUUID.String(), document.SchemaId.String()},
            {dummyUUID.String(), document.Id.String()},
            {2015, document.InsertDate.Year()},
            {2015, document.LastUpdate.Year()},
            {true, document.IsActive},
        }
        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("CreateDocument %d: bad value, got: %v want: %v",
                    i, test.got, test.want)
            }
        }
    } else {
        t.Errorf("unexpected: both document and error are nil!")
    }

    // test READ
    doc, err := custodia.ReadDocument(schema, dummyUUID)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if doc != nil {
        // some values need type conversion
        convArrayString, _ := common.ConvertSliceItems[string](
            doc.Content["arrayStringField"],
        )
        convArrayInteger, _ := common.ConvertSliceItems[int](
            doc.Content["arrayIntegerField"],
        )
        convArrayFloat, _ := common.ConvertSliceItems[float64](
            doc.Content["arrayFloatField"],
        )

        // setup tests now
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyUUID.String(), document.RepositoryId.String()},
            {dummyUUID.String(), document.SchemaId.String()},
            {dummyUUID.String(), document.Id.String()},
            {2015, document.InsertDate.Year()},
            {2015, document.LastUpdate.Year()},
            {true, document.IsActive},
            // check the content
            {int64(42), doc.Content["integerField"]},
            {3.14, doc.Content["flaotField"]},
            {"antani", doc.Content["stringField"]},
            {"this is not a very long string, but should be",
                doc.Content["textField"]},
            {true, doc.Content["boolField"]},
            {time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
                doc.Content["dateField"]},
            {time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
                doc.Content["timeField"]},
            {time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
                doc.Content["datetimeField"]},
            {"VGhpcyBpcyBhIGJhc2UtNjQgZW5jb2RlZCBzdHJpbmcu",
                doc.Content["base64Field"]},
            {`{"success": true}`, doc.Content["jsonField"]},
            {dummyUUID.String(), doc.Content["blobField"].(string)},
            {[]int{0, 1, 2, 3, 4, 5}, convArrayInteger},
            {[]float64{1.1, 2.2, 3.3, 4.4}, convArrayFloat},
            {[]string{"Hello", "world", "!"}, convArrayString},
            // for date/time/datetime we check all
            {1970, doc.Content["dateField"].(time.Time).Year()},
            {1, int(doc.Content["dateField"].(time.Time).Month())},
            {1, doc.Content["dateField"].(time.Time).Day()},
            {0, doc.Content["timeField"].(time.Time).Hour()},
            {1, doc.Content["timeField"].(time.Time).Minute()},
            {30, doc.Content["timeField"].(time.Time).Second()},
            {2001, doc.Content["datetimeField"].(time.Time).Year()},
            {3, int(doc.Content["datetimeField"].(time.Time).Month())},
            {8, doc.Content["datetimeField"].(time.Time).Day()},
            {23, doc.Content["datetimeField"].(time.Time).Hour()},
            {31, doc.Content["datetimeField"].(time.Time).Minute()},
            {42, doc.Content["datetimeField"].(time.Time).Second()},
        }
        // add othere tests that need type conversion before

        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("CreateDocument %d: bad value, got: %v want: %v",
                    i, test.got, test.want)
            }
        }





        // // test the content
        // if doc.Content["integerField"] != int64(42) {
        //     t.Errorf("content: bad integerField, got: %v want: %v",
        //         doc.Content["integerField"], int64(42))
        // }
        // if doc.Content["flaotField"] != 3.14 {
        //     t.Errorf("content: bad flaotField, got: %v want: %v",
        //         doc.Content["flaotField"], 3.14 )
        // }
        // if doc.Content["stringField"] != "antani" {
        //     t.Errorf("content: bad stringField, got: %v want: %v",
        //         doc.Content["stringField"], "antani")
        // }
        // if doc.Content["textField"] !=
        //     "this is not a very long string, but should be" {
        //     t.Errorf("content: bad textField, got: %v want: %v",
        //         doc.Content["textField"],
        //         "this is not a very long string, but should be")
        // }
        // if doc.Content["boolField"] != true {
        //     t.Errorf("content: bad boolField, got: %v want: %v",
        //         doc.Content["boolField"], true)
        // }

        // // for date we check just the yyyy-mm-dd part (ignoring time)
        // dateField, _ := doc.Content["dateField"].(time.Time)
        // if dateField.Year() != 1970 {
        //     t.Errorf("content: bad dateField year, got: %v want: 1970",
        //         dateField.Year())
        // }
        // if dateField.Month() != 1 {
        //     t.Errorf("content: bad dateField month, got: %v want: 1",
        //         dateField.Month())
        // }
        // if dateField.Day() != 1 {
        //     t.Errorf("content: bad dateField day, got: %v want: 1",
        //         dateField.Day())
        // }

        // // for time we check just the HH:MM:SS part (ignoring date)
        // timeField, _ := doc.Content["timeField"].(time.Time)
        // if timeField.Hour() != 0 {
        //     t.Errorf("content: bad timeField hour, got: %v want: 0",
        //         timeField.Hour())
        // }
        // if timeField.Minute() != 1 {
        //     t.Errorf("content: bad timeField minute, got: %v want: 1",
        //         timeField.Minute())
        // }
        // if timeField.Second() != 30 {
        //     t.Errorf("content: bad timeField second, got: %v want: 30",
        //         timeField.Second())
        // }

        // // for datetime we check all
        // dateTimeField, _ := doc.Content["datetimeField"].(time.Time)
        // if dateTimeField.Year() != 2001 {
        //     t.Errorf("content: bad dateTimeField year, got: %v want: 2001",
        //         dateTimeField.Year())
        // }
        // if dateTimeField.Month() != 3 {
        //     t.Errorf("content: bad dateTimeField month, got: %v want: 3",
        //         dateTimeField.Month())
        // }
        // if dateTimeField.Day() != 8 {
        //     t.Errorf("content: bad dateTimeField day, got: %v want: 8",
        //         dateTimeField.Day())
        // }
        // if dateTimeField.Hour() != 23 {
        //     t.Errorf("content: bad dateTimeField hour, got: %v want: 23",
        //         dateTimeField.Hour())
        // }
        // if dateTimeField.Minute() != 31 {
        //     t.Errorf("content: bad dateTimeField minute, got: %v want: 31",
        //         dateTimeField.Minute())
        // }
        // if dateTimeField.Second() != 42 {
        //     t.Errorf("content: bad dateTimeField second, got: %v want: 42",
        //         dateTimeField.Second())
        // }

        // base64, json, and blob are just string fields. We don't care if the
        // string in it actually converts to real data. This is a problem of
        // the user, his duty to encode/decode data correctly
        // if doc.Content["base64Field"] !=
        //     "VGhpcyBpcyBhIGJhc2UtNjQgZW5jb2RlZCBzdHJpbmcu" {
        //     t.Errorf("content: bad base64Field, got: %v want: %v",
        //     doc.Content["base64Field"],
        //     "VGhpcyBpcyBhIGJhc2UtNjQgZW5jb2RlZCBzdHJpbmcu")
        // }
        // if doc.Content["jsonField"] != `{"success": true}` {
        //     t.Errorf("content: bad jsonField, got: %v want: {\"success\": true}",
        //         doc.Content["jsonField"])
        // }
        // wantBlobField, _ := dummyContent["blobField"].(string)
        // if doc.Content["blobField"] != wantBlobField {
        //     t.Errorf("content: bad blobField, got: %v want: %v",
        //         doc.Content["blobField"], wantBlobField)
        // }

        // // test array fields
        // if reflect.DeepEqual(doc.Content["arrayIntegerField"],
        //     []int{0, 1, 2, 3, 4, 5}) {
        //     t.Errorf("content: bad arrayIntegerField, got: %v want: %v",
        //         doc.Content["arrayIntegerField"], []int{0, 1, 2, 3, 4, 5})
        // }
        // if reflect.DeepEqual(doc.Content["arrayFloatField"],
        //     []float64{1.1, 2.2, 3.3, 4.4}) {
        //     t.Errorf("content: bad arrayIntegerField, got: %v want: %v",
        //         doc.Content["arrayFloatField"], []float64{1.1, 2.2, 3.3, 4.4})
        // }
        // if reflect.DeepEqual(doc.Content["arrayStringField"],
        //     []string{"Hello", "world", "!"}) {
        //     t.Errorf("content: bad arrayIntegerField, got: %v want: %v",
        //         doc.Content["arrayStringField"],
        //         []string{"Hello", "world", "!"})
        // }
    } else {
        t.Errorf("unexpected: both document and error are nil!")
    }

    // test UPDATE
    // we still pass empty content, non influential, response is mocked
    doc, err = custodia.UpdateDocument(doc.Id, true, content)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if doc != nil {
        var tests = []struct {
            want interface{}
            got  interface{}
        }{
            {dummyUUID.String(), document.RepositoryId.String()},
            {dummyUUID.String(), document.SchemaId.String()},
            {dummyUUID.String(), document.Id.String()},
            {2025, document.InsertDate.Year()},  // changed
            {2025, document.LastUpdate.Year()},  // changed
            {want: false, got: doc.IsActive},    // changed
            // check the content
            {want: "brematurata", got: doc.Content["stringField"]}, // changed
            {int64(42), doc.Content["integerField"]},
            {3.14, doc.Content["flaotField"]},
            {"antani", doc.Content["stringField"]},
            {"this is not a very long string, but should be",
                doc.Content["textField"]},
            {true, doc.Content["boolField"]},
            {time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
                doc.Content["dateField"]},
            {time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
                doc.Content["timeField"]},
            {time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
                doc.Content["datetimeField"]},
            {"VGhpcyBpcyBhIGJhc2UtNjQgZW5jb2RlZCBzdHJpbmcu",
                doc.Content["base64Field"]},
            {`{"success": true}`, doc.Content["jsonField"]},
            {dummyUUID.String(), doc.Content["blobField"].(string)},
            {[]int{0, 1, 2, 3, 4, 5}, doc.Content["arrayIntegerField"]},
            {[]float64{1.1, 2.2, 3.3, 4.4}, doc.Content["arrayFloatField"]},
            {[]string{"Hello", "world", "!"}, doc.Content["arrayStringField"]},
            // for date/time/datetime we check all
            {1970, doc.Content["dateField"].(time.Time).Year()},
            {1, doc.Content["dateField"].(time.Time).Month()},
            {1, doc.Content["dateField"].(time.Time).Day()},
            {0, doc.Content["timeField"].(time.Time).Hour()},
            {1, doc.Content["timeField"].(time.Time).Minute()},
            {30, doc.Content["timeField"].(time.Time).Second()},
            {2001, doc.Content["datetimeField"].(time.Time).Year()},
            {3, doc.Content["datetimeField"].(time.Time).Month()},
            {8, doc.Content["datetimeField"].(time.Time).Day()},
            {23, doc.Content["datetimeField"].(time.Time).Hour()},
            {31, doc.Content["datetimeField"].(time.Time).Minute()},
            {42, doc.Content["datetimeField"].(time.Time).Second()},
        }
        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("UpdateDocument %d: bad value, got: %v want: %v",
                    i, test.got, test.want)
            }
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
    } else {
        if reflect.TypeOf(documents) != reflect.TypeOf([]*Document{}) {
            t.Errorf("documents is not list of Documents, got: %T want: %T",
                documents, []*Document{})
        }

        // we don't check the content of every single doc, just some values
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {true, documents[0].IsActive},
            {dummyUUID.String(), documents[0].SchemaId},
            {dummyUUID.String(), documents[0].Id},
            {2015, documents[0].InsertDate.Year()},
            {2015, documents[0].LastUpdate.Year()},
            {int(42), documents[0].Content["integerField"]},
            {true, documents[1].IsActive},
            {dummyUUID.String(), documents[0].SchemaId},
            {dummyUUID.String(), documents[0].Id},
            {2025, documents[1].InsertDate.Year()},
            {2025, documents[1].LastUpdate.Year()},
            {int(42), documents[1].Content["integerField"]},
        }
        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("ListDocuments %d: bad value, got: %v want: %v",
                    i, test.got, test.want)
            }
        }
    }
}
