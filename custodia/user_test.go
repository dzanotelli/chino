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
    // ResponseInnerUser will be included in responses
    type ResponseInnerUser struct {
        UserId string `json:"user_id"`
        UserSchemaId string `json:"schema_id"`
        Username string `json:"username"`
        InsertDate string `json:"insert_date"`
        LastUpdate string `json:"last_update"`
        IsActive bool `json:"is_active"`
        Attributes map[string]interface{} `json:"content,omitempty"`
        Groups []string `json:"groups"`
    }

    // UserResponse will be marshalled to create and API-like response
    type UserResponse struct {
        User ResponseInnerUser `json:"user"`
    }

    // UsersResponse will be marshalled to crete an API-like response
    type UsersResponse struct {
        Count int `json:"count"`
        TotalCount int `json:"total_count"`
        Limit int `json:"limit"`
        Offset int `json:"offset"`
        Users []ResponseInnerUser `json:"users"`
    }

    // init stuff
    dummyUser := ResponseInnerUser{
        UserId: uuid.New().String(),
        UserSchemaId: uuid.New().String(),
        Username: "unittest",
        InsertDate: "2015-02-24T21:48:16.332",
        LastUpdate: "2015-02-24T21:48:16.332",
        IsActive: false,
        Groups: []string{},
    }
    dummyAttrs := map[string]interface{}{
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
        "blobField": uuid.New().String(),
        "arrayIntegerField": `[0, 1, 1, 2, 3, 5]`,
        "arrayFloatField": `[1.1, 2.2, 3.3, 4.4]`,
        "arrayStringField": `["Hello", "world", "!"]`,
    }
    dummyUser.Attributes = dummyAttrs

    // shortcuts
    userSchemaId := dummyUser.UserSchemaId
    userId := dummyUser.UserId

    writeDocResponse := func(w http.ResponseWriter) {
        data, _ := json.Marshal(UserResponse{dummyUser})
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
        if r.URL.Path == fmt.Sprintf("/api/v1/user_schemas/%s/users",
            userSchemaId) && r.Method == "POST" {
            // mock CREATE response
            writeDocResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/users/%s",
            userId) && r.Method == "GET" {
            // mock READ response
            writeDocResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/users/%s",
            userId) && r.Method == "PUT" {
            // mock UPDATE response
            dummyUser.IsActive = true
            dummyAttrs["stringField"] = "brematurata"
            writeDocResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/users/%s",
            userId) && r.Method == "DELETE" {
            // mock DELETE response
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/user_schemas/%s/users",
            userSchemaId) && r.Method == "GET" {
            // mock LIST response
            usersResp := UsersResponse{
                Count: 1,
                TotalCount: 1,
                Limit: 100,
                Offset: 0,
                Users: []ResponseInnerUser{dummyUser},
            }
            data, _ := json.Marshal(usersResp)
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
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

    // test CREATE: we submit no content, since the response is mocked
    // we init instead a UserSchema with just the right ids
    userSchema := UserSchema{
        Id: dummyUser.UserSchemaId,
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
            {dummyUser.UserId, user.Id},
            {dummyUser.UserSchemaId, userSchema.Id},
            {dummyUser.Username, user.Username},
            {2015, user.InsertDate.Year()},
            {2, int(user.LastUpdate.Month())},
            {false, user.IsActive},
            {reflect.TypeOf(map[string]interface{}{}), reflect.TypeOf(user.Attributes)},
            {reflect.TypeOf([]string{}), reflect.TypeOf(user.Groups)},
        }
        for _, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("Users CREATE: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    } else {
        t.Errorf("unexpected: both user and error are nil!")
    }

    // test UPDATE
    user, err = custodia.UpdateUser(dummyUser.UserId, true, attributes)

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if user != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyUser.UserId, user.Id},
            {dummyUser.UserSchemaId, userSchema.Id},
            {dummyUser.Username, user.Username},
            {2015, user.InsertDate.Year()},
            {2, int(user.LastUpdate.Month())},
            {true, user.IsActive},
            {reflect.TypeOf(map[string]interface{}{}), reflect.TypeOf(user.Attributes)},
            {reflect.TypeOf([]string{}), reflect.TypeOf(user.Groups)},
        }
        for _, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("Users CREATE: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    } else {
        t.Errorf("unexpected: both user and error are nil!")
    }

    // test DELETE
    err = custodia.DeleteUser(dummyUser.UserId, false, false)
    if err != nil {
        t.Errorf("error while deleting user. Details: %v", err)
    }

    // test LIST
    // test we gave a wrong argument
    params := map[string]interface{}{"antani": 42}
    _, err = custodia.ListUsers(dummyUser.UserSchemaId, params)
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
    users, err := custodia.ListUsers(dummyUser.UserSchemaId, goodParams)

    if err != nil {
        t.Errorf("error while listing users: %v", err)
    } else if reflect.TypeOf(users) != reflect.TypeOf([]*User{}) {
        t.Errorf("users is not list of Users, got: %T want: %T",
            users, []*User{})
    }
}
