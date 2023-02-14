package custodia

import (
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestRepositoryCRUD(t *testing.T) {
	// mock calls
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		
	}
	server := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer server.Close()

	auth := NewClientAuth()  // auth is tested elsewhere (calls are mocked)
	chinoClient := chinoClient(server.URL, auth)
	

	// test READ


}