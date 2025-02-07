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

func TestGroupCRUDL(t *testing.T) {
    envelope := CustodiaEnvelope{
        Result: "success",
        ResultCode: 200,
        Message: nil,
    }

    dummyGroup := map[string]interface{}{
        "group_id": uuid.New().String(),
        "group_name": "unittest",
        "attributes": map[string]interface{}{"antani": 3.14},
        "is_active": true,
        "insert_date": "2015-02-07T12:14:46.754",
        "last_update": "2015-03-13T18:06:21.242",
    }
    gid, _ := dummyGroup["group_id"].(string)

    userSchemaId := uuid.New().String()
    dummyUser := map[string]interface{}{
        "username": "unittest",
        "schemas_id": userSchemaId,
        "user_id": uuid.New().String(),
        "insert_date": "2015-02-07T12:14:46.754",
        "last_update": "2015-03-13T18:06:21.242",
        "is_active": true,
        "attributes": map[string]interface{}{"antani": 3.14},
        "groups": []string{gid},
    }
    uid, _ := dummyUser["user_id"].(string)

    responseGroup := map[string]interface{}{
        "group": dummyGroup,
    }

    mockHandler := func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/api/v1/groups" && r.Method == "POST" {
            // mock CREATE response
            envelope.Data, _ = json.Marshal(responseGroup)
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/groups/%s", gid) &&
            // mock READ response
            r.Method == "GET" {
            data, _ := json.Marshal(responseGroup)
            envelope.Data = data
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/groups/%s", gid ) &&
            r.Method == "PUT" {
            // mock UPDATE response
            dummyGroup["group_name"] = "changed"
            dummyGroup["is_active"] = false
            dummyGroup["attributes"] = map[string]interface{}{
                "antani": 3.14, "something": "else"}
            envelope.Data, _ = json.Marshal(responseGroup)
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/groups/%s",
            gid) && r.Method == "DELETE" {
            // mock DELETE response
            data, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(data)
        } else if r.URL.Path == "/api/v1/groups" && r.Method == "GET" {
            // mock LIST response
            groupsResp := map[string]interface{}{
                "count": 1,
                "total_count": 1,
                "limit": 1,
                "offset": 0,
                "groups": []interface{}{dummyGroup},
            }
            data, _ := json.Marshal(groupsResp)
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            envelope.Data = data
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/groups/%s/users", gid) &&
            r.Method == "GET" {
            listResp := map[string]interface{}{
                "count": 1,
                "total_count": 1,
                "limit": 1,
                "offset": 0,
                "users": []map[string]interface{}{dummyUser},
            }
            data, _ := json.Marshal(listResp)
            envelope := CustodiaEnvelope{Result: "success", ResultCode: 200}
            envelope.Data = data
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/groups/%s/users/%s", gid,
            uid) {
            // we don't care the method. GET, POST or DELETE, they always
            // return just 200 success with empty message
            data, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(data)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/groups/%s/user_schemas/%s", gid, userSchemaId) {
            // we don't care the method. GET, POST or DELETE, they always
            // return just 200 success with empty message
            data, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(data)
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

    // init stuff
    client := common.NewClient(server.URL, common.GetFakeAuth())
    custodia := NewCustodiaAPIv1(client)

    // test CREATE
    group, err := custodia.CreateGroup("unittest", true,
        map[string]interface{}{})
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {gid, group.Id},
            {"unittest", group.Name},
            {map[string]interface{}{"antani": 3.14}, group.Attributes},
            {true, group.IsActive},
            {2015, group.InsertDate.Year()},
            {2, int(group.InsertDate.Month())},
            {7, int(group.InsertDate.Day())},
            {12, int(group.InsertDate.Hour())},
            {14, int(group.InsertDate.Minute())},
            {46, int(group.InsertDate.Second())},
            {2015, group.LastUpdate.Year()},
            {3, int(group.LastUpdate.Month())},
            {13, int(group.LastUpdate.Day())},
            {18, int(group.LastUpdate.Hour())},
            {6, int(group.LastUpdate.Minute())},
            {21, int(group.LastUpdate.Second())},
        }

        for _, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("Group Create: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    }

    // test UPDATE
    // response is mocked, so we don't need to pass the right data
    group, err = custodia.UpdateGroup(gid, "changed", false,
        map[string]interface{}{})
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {gid, group.Id},
            {"changed", group.Name},
            {false, group.IsActive},
            {map[string]interface{}{"antani": 3.14, "something": "else"},
                group.Attributes},
            {2015, group.InsertDate.Year()},
            {2, int(group.InsertDate.Month())},
            {7, int(group.InsertDate.Day())},
            {12, int(group.InsertDate.Hour())},
            {14, int(group.InsertDate.Minute())},
            {46, int(group.InsertDate.Second())},
            {2015, group.LastUpdate.Year()},
            {3, int(group.LastUpdate.Month())},
            {13, int(group.LastUpdate.Day())},
            {18, int(group.LastUpdate.Hour())},
            {6, int(group.LastUpdate.Minute())},
            {21, int(group.LastUpdate.Second())},
        }

        for _, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("Group Update: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    }

    // test DELETE
    err = custodia.DeleteGroup(gid, false)
    if err != nil {
        t.Errorf("error while deleting group: %v", err)
    }
    err = custodia.DeleteGroup(gid, true)
    if err != nil {
        t.Errorf("error while deleting group: %v", err)
    }

    // test LIST
    groups, err := custodia.ListGroups()
    if err != nil {
        t.Errorf("error while listing groups: %v", err)
    } else if reflect.TypeOf(groups) != reflect.TypeOf([]Group{}) {
        t.Errorf("groups is not list of Groups, got: %T want: %T",
            groups, []*Group{})
    }

    // Group members tests
    users, err := custodia.ListGroupUsers(gid)
    if err!= nil {
        t.Errorf("error while listing users in group: %v", err)
    } else {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {uid, users[0].Id},
            {"unittest", users[0].Username},
            {2015, int(users[0].InsertDate.Year())},
            {2, int(users[0].InsertDate.Month())},
            {7, int(users[0].InsertDate.Day())},
            {12, int(users[0].InsertDate.Hour())},
            {14, int(users[0].InsertDate.Minute())},
            {46, int(users[0].InsertDate.Second())},
            {2015, users[0].LastUpdate.Year()},
            {3, int(users[0].LastUpdate.Month())},
            {13, int(users[0].LastUpdate.Day())},
        }

        for _, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("Group Create: bad value, got: %v want: %v",
                    test.got, test.want)
            }
        }
    }

    // Add user to group
    err = custodia.AddUserToGroup(uid, gid)
    if err!= nil {
        t.Errorf("error while adding user to group: %v", err)
    }

    // Remove user from group
    err = custodia.RemoveUserFromGroup(uid, gid)
    if err!= nil {
        t.Errorf("error while removing user from group: %v", err)
    }

    // Add all users of schema to group
    err = custodia.AddUsersFromUserSchemaToGroup(userSchemaId, gid)
    if err!= nil {
        t.Errorf("error while adding users to group from schema: %v", err)
    }

    // Remove all users of schema from group
    err = custodia.RemoveUsersFromUserSchemaFromGroup(userSchemaId, gid)
    if err!= nil {
        t.Errorf("error while removing users from group from schema: %v", err)
    }
}
