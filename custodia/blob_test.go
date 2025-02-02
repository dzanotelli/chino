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



func TestBlob(t *testing.T) {

    envelope := CustodiaEnvelope{
        Result: "success",
        ResultCode: 200,
        Message: nil,
    }
    dummyUUID := uuid.New().String()

	blobCreateResp := map[string]interface{}{
		"blob": map[string]string{
			"upload_id": dummyUUID,
			"expire_date": "2015-04-14T05:09:54.915Z",
		},
	}

	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/blobs" {
			w.WriteHeader(http.StatusOK)
			envelope.Data, _ = json.Marshal(blobCreateResp)
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

	// Test blob create

	blob, err := custodia.CreateBlob(dummyUUID, "field", "filename")
	if err != nil {
		t.Errorf("Error creating blob: %v", err)
	} else {
		var tests = []struct {
			want interface{}
			got  interface{}
		}{
			{dummyUUID, blob.Id},
			{2015, blob.ExpireDate.Year()},
			{4, int(blob.ExpireDate.Month())},
			{14, blob.ExpireDate.Day()},
			{5, blob.ExpireDate.Hour()},
			{9, blob.ExpireDate.Minute()},
			{54, blob.ExpireDate.Second()},
		}
		for i, test := range(tests) {
			if !reflect.DeepEqual(test.want, test.got) {
				t.Errorf("CreateBlob %d: expected %v, got %v", i, test.want,
					test.got)
			}
		}
	}
}
