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
    // ResponseInnerUserSchema will be included in responses
    type ResponseInnerUserSchema struct {
        UserSchemaId string `json:"user_schema_id"`
        Description string `json:"description"`
        Groups []string `json:"groups"`
        InsertDate string `json:"insert_date"`
        LastUpdate string `json:"last_update"`
        IsActive bool `json:"is_active"`
        Structure []SchemaField `json:"structure"`
    }

    // SchemaResponse will be marshalled to create an API-like response
    type UserSchemaResponse struct {
        Schema ResponseInnerUserSchema `json:"user_schema"`
    }

    // SchemasResponse will be marshalled to create an API-like response
    type UserSchemasResponse struct {
        Count int `json:"count"`
        TotalCount int `json:"total_count"`
        Limit int `json:"limit"`
        Offset int `json:"offset"`
        UserSchemas []ResponseInnerUserSchema `json:"user_schemas"`
    }

    dummyUserSchema := ResponseInnerUserSchema{
        UserSchemaId: uuid.New().String(),
        Description: "unittest",
        InsertDate: "2015-02-24T21:48:16.332",
        LastUpdate: "2015-02-24T21:48:16.332",
        IsActive: false,
        // Structure: json.RawMessage{},
        Structure: []SchemaField{
            {Name: "IntField", Type: "integer", Indexed: true, Default: 42},
            {Name: "StrField", Type: "string", Indexed: true, Default: "asd"},
            {Name: "FloatField", Type: "float", Indexed: false, Default: 3.14},
            {Name: "BoolField", Type: "bool", Indexed: false},
            {Name: "DateField", Type: "date", Default: "2023-03-15"},
            {Name: "TimeField", Type: "time", Default: "11:43:04.058"},
            {Name: "DateTimeField", Type: "datetime",
                Default: "2023-03-15T11:43:04.058"},
        },
    }

    writeUserSchemaResponse := func(w http.ResponseWriter) {
        data, _ := json.Marshal(UserSchemaResponse{dummyUserSchema})
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
        if r.URL.Path == "/api/v1/user_schemas" && r.Method == "POST" {
            // mock CREATE response
            writeUserSchemaResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/user_schemas/%s",
            dummyUserSchema.UserSchemaId) && r.Method == "GET" {
            // mock READ response
            writeUserSchemaResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/user_schemas/%s",
            dummyUserSchema.UserSchemaId) && r.Method == "PUT" {
            // mock UPDATE response
            dummyUserSchema.Description = "changed"
            dummyUserSchema.IsActive = true
            // dummyUserSchema.Structure[0].Default = 21
            writeUserSchemaResponse(w)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/user_schemas/%s",
            dummyUserSchema.UserSchemaId) && r.Method == "DELETE" {
            // mock DELETE response
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == "/api/v1/user_schemas" &&  r.Method == "GET" {
            // mock LIST response
            schemasResp := UserSchemasResponse{
                Count: 1,
                TotalCount: 1,
                Limit: 100,
                Offset: 0,
                UserSchemas: []ResponseInnerUserSchema{dummyUserSchema},
            }
            data, _ := json.Marshal(schemasResp)
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

    // test CREATE: we submit an empty field list, since the response is mocked
    // and we will still get a working structure. The purpose here is to test
    // that the received data are correctly populating the objects
    structure := []SchemaField{}

    userSchema, err := custodia.CreateUserSchema("unittest", true, structure);

    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if userSchema != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyUserSchema.UserSchemaId, userSchema.Id},
            {dummyUserSchema.Description, userSchema.Description},
            {dummyUserSchema.Groups, userSchema.Groups},
            {2015, userSchema.InsertDate.Year()},
            {2015, userSchema.LastUpdate.Year()},
            {false, userSchema.IsActive},
        }
        for _, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("UserSchema CREATE: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    } else {
        t.Errorf("unexpected: both userSchema and error are nil!")
    }

    // test READ
    userSchema, err = custodia.ReadUserSchema(dummyUserSchema.UserSchemaId)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if userSchema != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyUserSchema.UserSchemaId, userSchema.Id},
            {dummyUserSchema.Description, userSchema.Description},
            {dummyUserSchema.Groups, userSchema.Groups},
            {2015, userSchema.InsertDate.Year()},
            {2015, userSchema.LastUpdate.Year()},
            {false, userSchema.IsActive},
        }
        for _, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("UserSchema CREATE: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    } else {
        t.Errorf("unexpected: both userSchema and error are nil!")
    }

    // test UPDATE
    userSchema, err = custodia.UpdateUserSchema(dummyUserSchema.UserSchemaId,
        "antani2", true, dummyUserSchema.Structure)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else if userSchema != nil {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {dummyUserSchema.UserSchemaId, userSchema.Id},
            {dummyUserSchema.Description, userSchema.Description},
            {dummyUserSchema.Groups, userSchema.Groups},
            {2015, userSchema.InsertDate.Year()},
            {2015, userSchema.LastUpdate.Year()},
            {true, userSchema.IsActive},
        }
        for _, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("UserSchema CREATE: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    } else {
        t.Errorf("unexpected: both userSchema and error are nil!")
    }

    // test DELETE
    err = custodia.DeleteUserSchema(dummyUserSchema.UserSchemaId, false)
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
