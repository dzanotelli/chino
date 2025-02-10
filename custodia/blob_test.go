package custodia

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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
    dummyUUID := uuid.New()

	blobCreateResp := map[string]interface{}{
		"blob": map[string]string{
			"upload_id": dummyUUID.String(),
			"expire_date": "2015-04-14T05:09:54.915Z",
		},
	}
	blobCommitResp := map[string]interface{}{
		"blob": map[string]string{
			"blob_id": dummyUUID.String(),
			"document_id": dummyUUID.String(),
			"sha1": "sha1",
			"md5": "md5",
		},
	}

	blobTokenResponse := map[string]interface{}{
		"token": "token",
		"expiration": "2015-04-14T05:09:54.915Z",
		"one_time": true,
	}

	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/blobs" && r.Method == "POST" {
			// mock CREATE blob
			w.WriteHeader(http.StatusOK)
			envelope.Data, _ = json.Marshal(blobCreateResp)
			out, _ := json.Marshal(envelope)
			w.Write(out)
		} else if r.URL.Path == "/api/v1/blobs/" + dummyUUID.String() &&
			r.Method == "PUT" {
			// mock upload a chunk
			w.WriteHeader(http.StatusOK)
			envelope.Data, _ = json.Marshal(blobCreateResp)
			out, _ := json.Marshal(envelope)
			w.Write(out)
		} else if r.URL.Path == "/api/v1/blobs/commit" && r.Method == "POST" {
			// mock commit a blob
			w.WriteHeader(http.StatusOK)
			envelope.Data, _ = json.Marshal(blobCommitResp)
			out, _ := json.Marshal(envelope)
			w.Write(out)
		} else if r.URL.Path == "/api/v1/blobs/" + dummyUUID.String() &&
			r.Method == "GET" {
			// mock GET blob
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("hello world!"))
		} else if r.URL.Path == "/api/v1/blobs/" + dummyUUID.String() &&
			r.Method == "DELETE" {
			// mock DELETE blob
			w.WriteHeader(http.StatusOK)
			envelope.Data = nil
			out, _ := json.Marshal(envelope)
			w.Write(out)
		} else if r.URL.Path == fmt.Sprintf("/api/v1/blobs/%s/generate",
			dummyUUID.String()) && r.Method == "POST" {
			// mock generate blob token
			w.WriteHeader(http.StatusOK)
			envelope.Data, _ = json.Marshal(blobTokenResponse)
			out, _ := json.Marshal(envelope)
			w.Write(out)
		} else if r.URL.Path == fmt.Sprintf("/api/v1/blobs/url/%s",
			dummyUUID) && r.Method == "GET" {
			// mock GET blob with token
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("hello token!"))
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

	// Test upload create
	upload, err := custodia.CreateBlob(dummyUUID, "field", "filename")
	if err != nil {
		t.Errorf("Error creating blob: %v", err)
	} else {
		var tests = []struct {
			want interface{}
			got  interface{}
		}{
			{dummyUUID, upload.Id},
			{2015, upload.ExpireDate.Year()},
			{4, int(upload.ExpireDate.Month())},
			{14, upload.ExpireDate.Day()},
			{5, upload.ExpireDate.Hour()},
			{9, upload.ExpireDate.Minute()},
			{54, upload.ExpireDate.Second()},
		}
		for i, test := range(tests) {
			if !reflect.DeepEqual(test.want, test.got) {
				t.Errorf("CreateBlob %d: expected %v, got %v", i, test.want,
					test.got)
			}
		}
	}

	// Test UploadChunk
	ub, err := custodia.CreateBlob(dummyUUID, "field", "filename")
	if err != nil {
		t.Errorf("Error creating blob: %v", err)
	}
	// we just check the UUID
	if (ub.Id != dummyUUID) {
		t.Errorf("CreateBlob: expected %v, got %v", dummyUUID, ub.Id)
	}

	data := []byte("hello world")
	ub, err = custodia.UploadChunk(dummyUUID, data, 11, 0)
	if err != nil {
		t.Errorf("Error uploading chunk: %v", err)
	} else {
		var tests = []struct {
			want interface{}
			got  interface{}
		}{
			{dummyUUID, ub.Id},
			{2015, ub.ExpireDate.Year()},
			{4, int(ub.ExpireDate.Month())},
			{14, ub.ExpireDate.Day()},
			{5, ub.ExpireDate.Hour()},
			{9, ub.ExpireDate.Minute()},
			{54, ub.ExpireDate.Second()},
		}
		for i, test := range(tests) {
			if !reflect.DeepEqual(test.want, test.got) {
				t.Errorf("UploadChunk %d: expected %v, got %v", i, test.want,
					test.got)
			}
		}
	}

	// Test CommitBlob
	blob, err := custodia.CommitBlob(dummyUUID)
	if err != nil {
		t.Errorf("Error committing blob: %v", err)
	} else {
		var tests = []struct {
			want interface{}
			got  interface{}
		}{
			{dummyUUID.String(), blob.Id.String()},
			{dummyUUID.String(), blob.DocumentId},
			{"sha1", blob.Sha1},
			{"md5", blob.Md5},
		}
		for i, test := range(tests) {
			if !reflect.DeepEqual(test.want, test.got) {
				t.Errorf("CommitBlob %d: expected %v (%t), got %v (%t)", i,
				test.want, test.want, test.got, test.got)
			}
		}
	}

	// Test GetBlobData
	stream, err := custodia.GetBlobData(dummyUUID)
	if err != nil {
		t.Errorf("Error getting blob data: %v", err)
	} else {
		data, err := io.ReadAll(stream)
		if err != nil {
			t.Errorf("Error reading blob data: %v", err)
		}
		if string(data) != "hello world!" {
			t.Errorf("GetBlobData: expected %v, got %v", []byte("hello world"),
				data)
		}
	}

	// Test DeleteBlob
	err = custodia.DeleteBlob(dummyUUID)
	if err != nil {
		t.Errorf("Error deleting blob: %v", err)
	}

	blobToken, err := custodia.GenerateBlobToken(dummyUUID, false, 0)
	if err != nil {
		t.Errorf("Error generating blob token: %v", err)
	} else {
		var tests = []struct {
			want interface{}
			got  interface{}
		}{
			{"token", blobToken.Token},
			{2015, blobToken.Expiration.Year()},
			{4, int(blobToken.Expiration.Month())},
			{14, blobToken.Expiration.Day()},
			{5, blobToken.Expiration.Hour()},
			{9, blobToken.Expiration.Minute()},
			{54, blobToken.Expiration.Second()},
			{true, blobToken.OneTime},
		}
		for i, test := range(tests) {
			if !reflect.DeepEqual(test.want, test.got) {
				t.Errorf("GenerateBlobToken %d: expected %v, got %v", i,
					test.want, test.got)
			}
		}
	}

	// Test GetBlobDataWithToken
	stream, err = custodia.GetBlobDataWithToken(dummyUUID, "token")
	if err != nil {
		t.Errorf("Error getting blob data with token: %v", err)
	} else {
		data, err := io.ReadAll(stream)
		if err != nil {
			t.Errorf("Error reading blob data  with token: %v", err)
		}
		if string(data) != "hello token!" {
			t.Errorf("GetBlobDataWithToken: expected %s, got %s",
				[]byte("hello world"), data)
		}
	}

	// Test upload a blob from a file
	file, err := os.CreateTemp("", "chino_unittest_*.txt")
	if err != nil {
		t.Errorf("Error creating temp file: %v", err)
	}
	defer os.Remove(file.Name())

	// try to upload the emtpy file: we should get an error
	_, err = custodia.CreateBlobFromFile(file.Name(), dummyUUID, "field", 0)
	if err == nil {
		t.Errorf("Expected error when uploading an empty file")
	} else if fmt.Sprint(err) != "file is empty" {
		t.Errorf("Expected error 'file is empty', got %v", err)
	}

	// write some content to the file
	_, err = file.WriteString("hello world!")
	if err != nil {
		t.Errorf("Error writing to temp file: %v", err)
	}
	file.Close()

	blob, err = custodia.CreateBlobFromFile(file.Name(), dummyUUID, "field", 0)
	if err != nil {
		t.Errorf("Error creating blob from file: %v", err)
	} else {
		// this is the result of commit blob
		var tests = []struct {
			want interface{}
			got  interface{}
		}{
			{dummyUUID.String(), blob.Id.String()},
			{dummyUUID.String(), blob.DocumentId},
			{"sha1", blob.Sha1},
			{"md5", blob.Md5},
		}
		for i, test := range(tests) {
			if !reflect.DeepEqual(test.want, test.got) {
				t.Errorf("CreateBlobFromFile %d: expected %v, got %v", i,
					test.want, test.got)
			}
		}
	}
}
