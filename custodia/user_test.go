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
	"github.com/google/uuid"
)

func TestUserCRUDL(t *testing.T) {
    envelope := CustodiaEnvelope{
        Result: "success",
        ResultCode: 200,
        Message: nil,
    }
    dummyUUID := uuid.New()

    createResp := map[string]interface{}{
        "user_id": dummyUUID.String(),
        "schema_id": dummyUUID.String(),
        "username": "unittest",
        "insert_date": "2015-04-24T21:48:16.332Z",
        "last_update": "2015-04-24T21:48:16.332Z",
        "is_active": false,
        "groups": []string{},
    }

    updateResp := map[string]interface{}{
        "user_id": dummyUUID.String(),
        "schema_id": dummyUUID.String(),
        "username": "unittest",
        "insert_date": "2015-04-24T21:48:16.332Z",
        "last_update": "2015-04-24T21:48:16.332Z",
        "is_active": true,
        "groups": []string{},
        "attributes": createResp["attributes"],
    }

    dummyAttributes := map[string]interface{}{
        "integerField": 42,
        "flaotField": 3.14,
        "stringField": "antani",
        "textField": "this is not a very long string, but should be",
        "boolField": true,
        "dateField": "1970-01-01",
        "timeField": "00:01:30",
        "datetimeField": "2001-03-08T23:31:42",
        "base64Field": "VGhpcyBpcyBhIGJhc2UtNjQgZW5jb2RlZCBzdHJpbmcu",
        "jsonField": `{"success": true}`,
        "blobField": dummyUUID.String(),
        "arrayIntegerField": `[0, 1, 1, 2, 3, 5]`,
        "arrayFloatField": `[1.1, 2.2, 3.3, 4.4]`,
        "arrayStringField": `["Hello", "world", "!"]`,
    }

    createResp["attributes"] = dummyAttributes
    updateResp["attributes"] = dummyAttributes

    // mock calls
    mockHandler := func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == fmt.Sprintf(
            "/api/v1/user_schemas/%s/users", dummyUUID,
        ) && r.Method == "POST" {
            // mock CREATE response
            w.WriteHeader(http.StatusCreated)
            data, _ := json.Marshal(createResp)
            envelope.Data = data
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/users/%s", dummyUUID) &&
            r.Method == "GET" {
            // mock READ response
            w.WriteHeader(http.StatusOK)
            data, _ := json.Marshal(createResp)
            envelope.Data = data
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/users/%s", dummyUUID) &&
            r.Method == "PUT" {
            // mock UPDATE response
            dummyAttributes["stringField"] = "brematurata"
            w.WriteHeader(http.StatusOK)
            data, _ := json.Marshal(updateResp)
            envelope.Data = data
            out, _ := json.Marshal(envelope)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/users/%s", dummyUUID) &&
            r.Method == "DELETE" {
            // mock DELETE response
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/user_schemas/%s/users", dummyUUID,
        ) && r.Method == "GET" {
            // mock LIST response
            usersResp := map[string]interface{}{
                "count": 1,
                "total_count": 1,
                "limit": 100,
                "offset": 0,
                "users": []map[string]interface{}{
                    createResp,
                },
            }
            envelope.Data, _ = json.Marshal(usersResp)
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

    // test CREATE: we submit no content, since the response is mocked
    // we init instead a UserSchema with just the right ids
    userSchema := UserSchema{
        Id: dummyUUID,
        Description: "unittest",
        IsActive: true,
        Structure: []SchemaField{},
    }
    attributes := map[string]interface{}{}
    user, err := custodia.CreateUser(&userSchema, false, attributes)

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if user != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyUUID.String(), user.Id.String()},
            {dummyUUID.String(), userSchema.Id.String()},
            {"unittest", user.Username},
            {2015, user.InsertDate.Year()},
            {2, int(user.LastUpdate.Month())},
            {false, user.IsActive},
            {reflect.TypeOf(map[string]interface{}{}), reflect.TypeOf(user.Attributes)},
            {reflect.TypeOf([]string{}), reflect.TypeOf(user.Groups)},
        }
        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("CreateUser #%d: bad value, got: %v want: %v", i,
                    test.got, test.want)
            }
        }
    } else {
        t.Errorf("unexpected: both user and error are nil!")
    }

    // test UPDATE
    user, err = custodia.UpdateUser(dummyUUID, true, dummyAttributes)

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if user != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyUUID.String(), user.Id.String()},
            {dummyUUID.String(), userSchema.Id.String()},
            {"unittest", user.Username},
            {2015, user.InsertDate.Year()},
            {2, int(user.LastUpdate.Month())},
            {true, user.IsActive},
            {reflect.TypeOf(map[string]interface{}{}), reflect.TypeOf(user.Attributes)},
            {reflect.TypeOf([]string{}), reflect.TypeOf(user.Groups)},
        }
        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("UpdateUser #%d: bad value, got: %v want: %v", i,
                    test.got, test.want)
            }
        }
    } else {
        t.Errorf("unexpected: both user and error are nil!")
    }

    // test DELETE
    err = custodia.DeleteUser(dummyUUID, false, false)
    if err != nil {
        t.Errorf("error while deleting user. Details: %v", err)
    }

    // test LIST
    // test we gave a wrong argument
    params := map[string]interface{}{"antani": 42}
    _, err = custodia.ListUsers(dummyUUID, params)
    if err == nil {
        t.Errorf("ListUsers is not giving error with wrong param %v",
            params)
    }

    // test that all the other params are accepted instead
    goodParams := map[string]interface{}{
        "full_user": true,
        "is_active": true,
        "insert_date__gt": time.Time{},
        "insert_date__lt": time.Time{},
        "last_update__gt": time.Time{},
        "last_update__lt": time.Time{},
    }
    users, err := custodia.ListUsers(dummyUUID, goodParams)

    if err != nil {
        t.Errorf("error while listing users: %v", err)
    } else if reflect.TypeOf(users) != reflect.TypeOf([]*User{}) {
        t.Errorf("users is not list of Users, got: %T want: %T",
            users, []*User{})
    }
}
