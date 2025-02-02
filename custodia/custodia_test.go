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
)







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


