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

func TestApplicationCRUDL(t *testing.T) {
    envelope := CustodiaEnvelope{
        Result: "success",
        ResultCode: 200,
        Message: nil,
    }
    aid := "MyAppId42"
    dummyApp := map[string]any{
        "app_id": aid,
        "app_secret": "123456",
        "client_type": "public",
        "grant_type": "password",
        "app_name": "antani",
        "redirect_url": "",
    }

    // mock calls
    mockHandler := func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/api/v1/auth/applications" && r.Method == "POST" {
            // mock CREATE response
            data, _ := json.Marshal(map[string]any{
                "application": dummyApp,
            })
            envelope.Data = data
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusCreated)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/auth/applications/%s", aid,
        ) && r.Method == "GET" {
            // mock READ response
            data, _ := json.Marshal(map[string]any{
                "application": dummyApp,
            })
            envelope.Data = data
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/auth/applications/%s", aid,
        ) && r.Method == "PUT" {
            // mock UPDATE response
            dummyApp["grant_type"] = GrantAuthorizationCode.String()
            dummyApp["client_type"] = ClientConfidential.String()
            dummyApp["redirect_url"] = "http://antani.org"
            data, _ := json.Marshal(map[string]any{
                "application": dummyApp,
            })
            envelope.Data = data
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/auth/applications/%s", aid,
        ) && r.Method == "DELETE" {
            // mock DELETE response
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == "/api/v1/auth/applications" &&
            r.Method == "GET" {
            // mock LIST response
            data := map[string]any{
                "count": 1,
                "total_count": 1,
                "limit": 100,
                "offset": 0,
                "applications": []any{dummyApp},
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

    // test CREATE
    app, err := custodia.CreateApplication("antani", GrantPassword,
        ClientConfidential, "")

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if app != nil {
        var tests = []struct {
            want any
            got any
        }{
            {aid, app.Id},
            {GrantPassword, app.GrantType},
            {"antani", app.Name},
            {"123456", app.Secret},

        }
        for _, test := range tests {
            if test.want != test.got {
                t.Errorf("Application CREATE: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    }

    // test READ
    app, err = custodia.ReadApplication(aid)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if app != nil {
        var tests = []struct {
            want any
            got any
        }{
            {aid, app.Id},
            {GrantPassword, app.GrantType},
            {"antani", app.Name},
            {"123456", app.Secret},

        }
        for _, test := range tests {
            if test.want != test.got {
                t.Errorf("ReadApplication: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    }

    // test UPDATE
    app, err = custodia.UpdateApplication(aid, "antani",
        GrantAuthorizationCode, ClientConfidential, "http://antani.org")
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if app != nil {
        var tests = []struct {
            want any
            got any
        }{
            {aid, app.Id},
            {ClientConfidential, app.ClientType},
            {GrantAuthorizationCode, app.GrantType},
            {"antani", app.Name},
            {"123456", app.Secret},

        }
        for _, test := range tests {
            if test.want != test.got {
                t.Errorf("UpdateApplication: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    }

    // test DELETE
    err = custodia.DeleteApplication(aid)
    if err != nil {
        t.Errorf("error while deleting application. Details: %v", err)
    }

    // test LIST
    queryParams := map[string]string{"offset": "0", "limit": "100"}
    apps, err := custodia.ListApplications(queryParams)
    if err != nil {
        t.Errorf("error while listing applications. Details: %v", err)
    } else if reflect.TypeOf(apps) != reflect.TypeOf([]*Application{}) {
        t.Errorf("apps is not list of Applications, got: %T want: %T",
            apps, []*Application{})
    }
}

func TestOAuth(t *testing.T) {
    envelope := CustodiaEnvelope{
        Result: "success",
        ResultCode: 200,
        Message: nil,
    }
    responseLogin := map[string]any{
        "access_token": "ans2fN08sliGpIOLMGg3fv4BpPhWRq",
        "token_type": "Bearer",
        "expires_in": 36000,
        "refresh_token": "vL0durAhdhNNYFI27F3zGGHXeNLwcO",
        "scope": "read write",
    }
    responseRefresh := map[string]any{
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
            want any
            got any
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
            want any
            got any
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
