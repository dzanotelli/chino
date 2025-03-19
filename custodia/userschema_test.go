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

func TestUserSchemaCRUDL(t *testing.T) {
    envelope := CustodiaEnvelope{
        Result: "success",
        ResultCode: 200,
        Message: nil,
    }
    dummyUUID := uuid.New()

    createResp := map[string]interface{}{
        "user_schema_id": dummyUUID.String(),
        "description": "unittest",
        "groups": []string{},
        "insert_date": "2015-04-24T21:48:16.332Z",
        "last_update": "2015-04-24T21:48:16.332Z",
        "is_active": false,
    }
    updateResp := map[string]interface{}{
        "user_schema_id": dummyUUID.String(),
        "description": "changed",
        "groups": []string{},
        "insert_date": "2015-04-24T21:48:16.332Z",
        "last_update": "2015-04-24T21:48:16.332Z",
        "is_active": true,
    }

    structure := []interface{}{
        map[string]interface{}{
            "name": "IntField",
            "type": "integer",
            "indexed": true,
            "default": 42,
        },
        map[string]interface{}{
            "name": "StrField",
            "type": "string",
            "indexed": true,
            "default": "asd",
        },
        map[string]interface{}{
            "name": "FloatField",
            "type": "number",
            "indexed": true,
            "default": 3.14,
        },
        map[string]interface{}{
            "name": "BoolField",
            "type": "boolean",
            "indexed": false,
        },
        map[string]interface{}{
            "name": "DateField",
            "type": "date",
            "default": "2023-03-15",
        },
        map[string]interface{}{
            "name": "TimeField",
            "type": "time",
            "default": "11:43:04.058",
        },
        map[string]interface{}{
            "name": "DatetimeField",
            "type": "datetime",
            "default": "2023-03-15T11:43:04.058",
        },
    }
    createResp["structure"] = structure
    updateResp["structure"] = structure


    // mock calls
    mockHandler := func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/api/v1/user_schemas" && r.Method == "POST" {
            // mock CREATE response
            w.WriteHeader(http.StatusCreated)
            data := map[string]any{"user_schema": createResp}
            envelope.Data, _ = json.Marshal(data)
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/user_schemas/%s", dummyUUID,
        ) && r.Method == "GET" {
            // mock READ response
            w.WriteHeader(http.StatusOK)
            data := map[string]any{"user_schema": createResp}
            envelope.Data, _ = json.Marshal(data)
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/user_schemas/%s", dummyUUID,
        ) && r.Method == "PUT" {
            // mock UPDATE response
            w.WriteHeader(http.StatusOK)
            data := map[string]any{"user_schema": updateResp}
            envelope.Data, _ = json.Marshal(data)
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/user_schemas/%s", dummyUUID,
        ) && r.Method == "DELETE" {
            // mock DELETE response
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == "/api/v1/user_schemas" &&  r.Method == "GET" {
            listResp := map[string]interface{}{
                "count": 1,
                "total_count": 1,
                "limit": 100,
                "offset": 0,
                "user_schemas": []interface{}{
                    createResp,
                },
            }
            envelope.Data, _ = json.Marshal(listResp)
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

    // test CREATE: we submit an empty field list, since the response is mocked
    // and we will still get a working structure. The purpose here is to test
    // that the received data are correctly populating the objects
    reqStruct := []SchemaField{}

    userSchema, err := custodia.CreateUserSchema("unittest", true, reqStruct);

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if userSchema != nil {
        var tests = []struct {
            want any
            got any
        }{
            {dummyUUID.String(), userSchema.Id.String()},
            {"unittest", userSchema.Description},
            {[]string(nil), userSchema.Groups},
            {2015, userSchema.InsertDate.Year()},
            {2015, userSchema.LastUpdate.Year()},
            {false, userSchema.IsActive},
        }
        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("CreateUserSchema #%d: bad value, got: %v (%T) " +
                    "want: %v (%T)", i, test.got, test.got, test.want,
                    test.want)
            }
        }
    } else {
        t.Errorf("unexpected: both userSchema and error are nil!")
    }

    // test READ
    userSchema, err = custodia.ReadUserSchema(dummyUUID)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if userSchema != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyUUID.String(), userSchema.Id.String()},
            {"unittest", userSchema.Description},
            {[]string(nil), userSchema.Groups},
            {2015, userSchema.InsertDate.Year()},
            {2015, userSchema.LastUpdate.Year()},
            {false, userSchema.IsActive},
        }
        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("ReadUserSchema #%d: bad value, got: %v want: %v", i,
                    test.got, test.want)
            }
        }
    } else {
        t.Errorf("unexpected: both userSchema and error are nil!")
    }

    // test UPDATE
    userSchema, err = custodia.UpdateUserSchema(dummyUUID, "antani2", true,
        reqStruct)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if userSchema != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyUUID.String(), userSchema.Id.String()},
            {"changed", userSchema.Description},
            {[]string(nil), userSchema.Groups},
            {2015, userSchema.InsertDate.Year()},
            {2015, userSchema.LastUpdate.Year()},
            {true, userSchema.IsActive},
        }
        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("UpdateUserSchema #%d: bad value, got: %v want: %v",
                    i, test.got, test.want)
            }
        }
    } else {
        t.Errorf("unexpected: both userSchema and error are nil!")
    }

    // test DELETE
    err = custodia.DeleteUserSchema(dummyUUID, false)
    if err != nil {
        t.Errorf("error while deleting user schema. Details: %v", err)
    }

    // test LIST
    uss, err := custodia.ListUserSchemas()
    if err != nil {
        t.Errorf("error while listing user schemas. Details: %v", err)
    } else if reflect.TypeOf(uss) != reflect.TypeOf([]*UserSchema{}) {
        t.Errorf("uss is not list of UserSchemas, got: %T want: %T",
            uss, []*UserSchema{})
    }
}
