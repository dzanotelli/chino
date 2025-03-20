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

func TestSearch(t *testing.T) {
	envelope := CustodiaEnvelope{
		Result: "success",
		ResultCode: 200,
		Message: nil,
	}
	dummyUUID := uuid.New()

	docsResponse := map[string]any{
		"documents": []any{
			map[string]any{
				"document_id": dummyUUID.String(),
				"schema_id": dummyUUID.String(),
				"repository_id": dummyUUID.String(),
				"insert_date": "2015-02-07T12:14:46.754",
				"last_update": "2015-03-13T18:06:21.242",
				"is_active": true,
				"content": map[string]any{
					"antani": 42,
				},
			},
		},
		"count": 1,
		"total_count": 1,
		"limit": 1,
		"offset": 0,
	}

	usersResponse := map[string]any{
		"users": []any{
			map[string]any{
				"user_id": dummyUUID.String(),
				"schema_id": dummyUUID.String(),
				"username": "unittest",
				"insert_date": "2015-02-07T12:14:46.754",
				"last_update": "2015-03-13T18:06:21.242",
				"is_active": true,
				"attributes": map[string]any{"antani": 3.14},
			},
		},
		"count": 1,
		"total_count": 1,
		"limit": 1,
		"offset": 0,
	}

	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == fmt.Sprintf(
			"/api/v1/search/documents/%s", dummyUUID,
		) && r.Method == "POST" {
			data, _ := json.Marshal(docsResponse)
			envelope.Data = data
			out, _ := json.Marshal(envelope)
			w.WriteHeader(http.StatusOK)
			w.Write(out)
		} else if r.URL.Path == fmt.Sprintf(
			"/api/v1/search/users/%s", dummyUUID,
		) && r.Method == "POST" {
			data, _ := json.Marshal(usersResponse)
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

	// Test search documents
	query := map[string]any{
		"or": []map[string]any{
			{
				"field": "antani",
				"type": "eq",
				"value": 42,
			},
			{
				"field": "other",
				"type": "eq",
				"value": 21,
			},
		},
	}
	resp, err := custodia.SearchDocuments(dummyUUID, FullContent, query, nil)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		if len(resp.Documents) != 1 {
			t.Errorf("expected 1 document, got %d", len(resp.Documents))
		}
		var tests = []struct {
			want any
			got any
		}{
			{1, resp.Count},
			{1, resp.TotalCount},
			{1, resp.Limit},
			{0, resp.Offset},
			{dummyUUID.String(), resp.Documents[0].Id.String()},
			{dummyUUID.String(), resp.Documents[0].SchemaId.String()},
			{dummyUUID.String(), resp.Documents[0].RepositoryId.String()},
			{2015, resp.Documents[0].InsertDate.Year()},
			{2, int(resp.Documents[0].InsertDate.Month())},
			{7, resp.Documents[0].InsertDate.Day()},
			{2015, resp.Documents[0].LastUpdate.Year()},
			{3, int(resp.Documents[0].LastUpdate.Month())},
			{13, resp.Documents[0].LastUpdate.Day()},
			// FIXME: golang unmarshals returns always float64 (not ints)
			//   dunno if let the user do this, or force convertion of the
			//   underlying type in SearchDocuments (but we need the schema!)
			{42.0, resp.Documents[0].Content["antani"].(float64)},
		}

		for i, test := range tests {
			if !reflect.DeepEqual(test.want, test.got) {
				t.Errorf("SearchDocuments %d: bad value, got %d, want %d",
					i, test.got, test.want)
			}
		}
	}

	// test SearchUsers
	query = map[string]any{
		"field": "antani", "type": "gte", "value": 3.14,
	}
	sort := map[string]any{
		"field": "antani", "order": "asc",
	}
	resp, err = custodia.SearchUsers(dummyUUID, FullContent, query, sort)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	} else {
		if len(resp.Users) != 1 {
			t.Errorf("expected 1 user, got %d", len(resp.Users))
		}
		var tests = []struct {
			want any
			got any
		}{
			{1, resp.Count},
			{1, resp.TotalCount},
			{1, resp.Limit},
			{0, resp.Offset},
			{dummyUUID.String(), resp.Users[0].Id.String()},
			{dummyUUID.String(), resp.Users[0].UserSchemaId.String()},
			{"unittest", resp.Users[0].Username},
			{2015, resp.Users[0].InsertDate.Year()},
			{2, int(resp.Users[0].InsertDate.Month())},
			{7, resp.Users[0].InsertDate.Day()},
			{2015, resp.Users[0].LastUpdate.Year()},
			{3, int(resp.Users[0].LastUpdate.Month())},
			{13, resp.Users[0].LastUpdate.Day()},
			{3.14, resp.Users[0].Attributes["antani"].(float64)},
		}

		for i, test := range tests {
			if !reflect.DeepEqual(test.want, test.got) {
				t.Errorf("SearchUsers %d: bad value, got %d, want %d",
					i, test.got, test.want)
			}
		}
	}
}
