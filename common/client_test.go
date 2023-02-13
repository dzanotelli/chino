package common

import "testing"

func TestCallUrl(t *testing.T) {
	chinoAuth := NewClientAuth()
	chinoClient := NewClient("https://www.chino.io", chinoAuth)
	resp, err := chinoClient.Get("/")
	if err != nil {
		t.Errorf("Error while processing request: %s", err)
		return // stop execution here
	}

	if (resp.StatusCode != 200) {
		t.Errorf("Bad status code, got: %v want: 200", resp.StatusCode)
	}
}