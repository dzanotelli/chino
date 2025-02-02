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
    // ResponseInnerApp will be included in responses
    type ResponseInnerApp struct {
        AppSecret string `json:"app_secret"`
        ClientType string `json:"client_type"`
        GrantType string `json:"grant_type"`
        AppName string `json:"app_name"`
        RedirectUrl string `json:"redirect_url"`
        AppId string `json:"app_id"`
    }

    type ApplicationResponse struct {
        Application ResponseInnerApp `json:"application"`
    }

    type ApplicationsResponse struct {
        Count int `json:"count"`
        TotalCount int `json:"total_count"`
        Limit int `json:"limit"`
        Offset int `json:"offset"`
        Applications []ResponseInnerApp `json:"applications"`
    }

    // init stuff
    aid := "MyAppId42"
    dummyApp := ResponseInnerApp{
        AppId: aid,
        AppSecret: "123456",
        ClientType: "public",
        GrantType: "password",
        AppName: "antani",
        RedirectUrl: "",
    }

    writeAppResponse := func(w http.ResponseWriter) {
        data, _ := json.Marshal(ApplicationResponse{dummyApp})
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
        if r.URL.Path == "/api/v1/auth/applications" && r.Method == "POST" {
            // mock CREATE response
            writeAppResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/auth/applications/%s",
            aid) && r.Method == "GET" {
            // mock READ response
            writeAppResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/auth/applications/%s",
            aid) && r.Method == "PUT" {
            // mock UPDATE response
            dummyApp.GrantType = GrantAuthorizationCode.String()
            dummyApp.ClientType = ClientConfidential.String()
            dummyApp.RedirectUrl = "http://antani.org"
            writeAppResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/auth/applications/%s",
            aid) && r.Method == "DELETE" {
            // mock DELETE response
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == "/api/v1/auth/applications" &&
            r.Method == "GET" {
            // mock LIST response
            appsResp := ApplicationsResponse{
                Count: 1,
                TotalCount: 1,
                Limit: 100,
                Offset: 0,
                Applications: []ResponseInnerApp{dummyApp},
            }
            data, _ := json.Marshal(appsResp)
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

    // test CREATE
    app, err := custodia.CreateApplication("antani", GrantPassword,
        ClientConfidential, "")

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if app != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyApp.AppId, app.Id},
            {GrantPassword, app.GrantType},
            {dummyApp.AppName, app.Name},
            {dummyApp.AppSecret, app.Secret},

        }
        for _, test := range tests {
            if test.want != test.got {
                t.Errorf("Application CREATE: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    }

    // test READ
    app, err = custodia.ReadApplication(dummyApp.AppId)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if app != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyApp.AppId, app.Id},
            {GrantPassword, app.GrantType},
            {dummyApp.AppName, app.Name},
            {dummyApp.AppSecret, app.Secret},

        }
        for _, test := range tests {
            if test.want != test.got {
                t.Errorf("Application GET: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    }

    // test UPDATE
    app, err = custodia.UpdateApplication(dummyApp.AppId, "antani",
        GrantAuthorizationCode, ClientConfidential, "http://antani.org")
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if app != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyApp.AppId, app.Id},
            {ClientConfidential, app.ClientType},
            {GrantAuthorizationCode, app.GrantType},
            {"antani", app.Name},
            {dummyApp.AppSecret, app.Secret},

        }
        for _, test := range tests {
            if test.want != test.got {
                t.Errorf("Application GET: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    }

    // test DELETE
    err = custodia.DeleteApplication(dummyApp.AppId)
    if err != nil {
        t.Errorf("error while deleting application. Details: %v", err)
    }

    // test LIST
    apps, err := custodia.ListApplications()
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

