package common

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCallUrl(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		// accept: json header needed
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Accept header, got: %v want: application/json", 
					 r.Header.Get("Accept"))
		}

		// mock endpoint responses
		if r.URL.Path == "/my/client/test" {
			w.WriteHeader(http.StatusOK)
       		w.Write([]byte(`{"value":"ok"}`))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"unsupported url"}`))
		}
	}
	server := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer server.Close()


	chinoAuth := NewClientAuth()
	chinoClient := NewClient(server.URL, chinoAuth)

	// test good URL
	resp, err := chinoClient.Get("/my/client/test")
	if err != nil {
		t.Errorf("error while processing request: %s", err)
		return // stop execution here
	}

	if (resp.StatusCode != 200) {
		t.Errorf("bad status code, got: %v want: 200", resp.StatusCode)
	}

	// test bad URL
	resp, _ = chinoClient.Get("/my/bad/url")
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("got wrong status, got: %v , want: %d", resp.StatusCode,
				 http.StatusBadRequest)
	}
}