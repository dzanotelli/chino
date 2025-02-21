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

func TestPermission(t *testing.T) {
    envelope := CustodiaEnvelope{
        Result: "success",
        ResultCode: 200,
        Message: nil,
    }

    dummyUUID := uuid.New()

    // dummy data to return from ReadAllPermissions
    allPermissions := []map[string]interface{}{
        {
            "access": "Structure",
            "parent_id": nil,
            "resource_type": "Repository",
            "owner_id": dummyUUID,
            "owner_type": "users",
            "permission": map[string][]string{
                "Manage": {
                  "R",
                },
            },
        },
        {
            "access": "Data",
            "resource_id": dummyUUID,
            "resource_type": "Schema",
            "owner_id": dummyUUID,
            "owner_type": "users",
            "permission": map[string][]string{
                "Authorize": {
                  "A",
                },
                "Manage": {
                  "R",
                  "U",
                  "L",
                },
            },
        },
    }
    userPermissions := []map[string]interface{}{
        {
            "access": "Structure",
            "parent_id": nil,
            "resource_type": "Repository",
            "owner_id": dummyUUID,
            "owner_type": "users",
            "permission": map[string][]string{
                "Manage": {
                  "R", "U",
                },
            },
        },
    }
    groupPermissions := []map[string]interface{}{
        {
            "access": "Structure",
            "parent_id": nil,
            "resource_type": "Repository",
            "owner_id": dummyUUID,
            "owner_type": "groups",
            "permission": map[string][]string{
                "Manage": {
                  "L", "R",
                },
            },
        },
    }

    mockHandler := func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == fmt.Sprintf(
            "/api/v1/perms/grant/repositories/users/%s", dummyUUID) &&
            r.Method == "POST" {
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/perms/grant/repositories/%s/groups/%s", dummyUUID,
            dummyUUID) && r.Method == "POST" {
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf(
            "/api/v1/perms/revoke/repositories/%s/schemas/groups/%s",
            dummyUUID, dummyUUID) && r.Method == "POST" {
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == "/api/v1/perms" && r.Method == "GET" {
            data := make(map[string]interface{})
            data["permissions"] = allPermissions
            envelope.Data, _ = json.Marshal(data)
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/perms/documents/%s",
            dummyUUID) && r.Method == "GET" {
            data := make(map[string]interface{})
            data["permissions"] = allPermissions  // reusing allPermissions
            envelope.Data, _ = json.Marshal(data)
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/perms/users/%s",
            dummyUUID) && r.Method == "GET" {
            data := make(map[string]interface{})
            data["permissions"] = userPermissions
            envelope.Data, _ = json.Marshal(data)
            out, _ := json.Marshal(envelope)
            w.WriteHeader(http.StatusOK)
            w.Write(out)
        } else if r.URL.Path == fmt.Sprintf("/api/v1/perms/groups/%s",
            dummyUUID) && r.Method == "GET" {
            data := make(map[string]interface{})
            data["permissions"] = groupPermissions
            envelope.Data, _ = json.Marshal(data)
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

    // Test Permission on Resources (multiple)
    perms := map[PermissionScope][]PermissionType{
        PermissionScopeManage: {PermissionTypeCreate, PermissionTypeList,
            PermissionTypeRead,},
        PermissionScopeAuthorize: {},
    }
    err := custodia.PermissionOnResources(PermissionActionGrant,
        ResourceRepository, ResourceUser, dummyUUID, perms)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }

    // Test Permission on Resource (single)
    err = custodia.PermissionOnResource(PermissionActionGrant,
        ResourceRepository, dummyUUID, ResourceGroup, dummyUUID, perms)
    if err!= nil {
        t.Errorf("unexpected error: %v", err)
    }

    // Test Permission on Resource children
    err = custodia.PermissionOnResourceChildren(PermissionActionRevoke,
        ResourceRepository, dummyUUID, ResourceSchema, ResourceGroup,
        dummyUUID, perms)
    if err!= nil {
        t.Errorf("unexpected error: %v", err)
    }

    // test ReadAllPermissions
    allPerms, err := custodia.ReadAllPermissions()
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {"Structure", allPerms[0].Access},
            {ResourceUser, allPerms[0].OwnerType},
            {dummyUUID, allPerms[0].OwnerId},
            {"", allPerms[0].ParentId},
            {map[PermissionScope][]PermissionType{
                PermissionScopeManage: {PermissionTypeRead,}},
                allPerms[0].Permission},
            {"Data", allPerms[1].Access},
            {ResourceUser, allPerms[1].OwnerType},
            {dummyUUID, allPerms[1].OwnerId},
            {dummyUUID, allPerms[1].Id},
            {"", allPerms[1].ParentId},
            {map[PermissionScope][]PermissionType{
                PermissionScopeAuthorize: {PermissionTypeAuthorize,},
                PermissionScopeManage: {PermissionTypeRead,
                    PermissionTypeUpdate, PermissionTypeList},

            }, allPerms[1].Permission},
        }
        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("ReadAllPermissions %d: bad value, got: %v want: %v",
                    i, test.got, test.want)
            }
        }
    }

    // Test ReadPermissionsOnDocument
    resources, err := custodia.ReadPermissionsOnDocument(dummyUUID)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {"Structure", resources[0].Access},
            {ResourceUser, resources[0].OwnerType},
            {dummyUUID, resources[0].OwnerId},
            {"", resources[0].ParentId},
            {map[PermissionScope][]PermissionType{
                PermissionScopeManage: {PermissionTypeRead,}},
                resources[0].Permission},
            {"Data", resources[1].Access},
            {ResourceUser, resources[1].OwnerType},
            {dummyUUID, resources[1].OwnerId},
            {dummyUUID, resources[1].Id},
            {"", resources[1].ParentId},
            {map[PermissionScope][]PermissionType{
                PermissionScopeAuthorize: {PermissionTypeAuthorize,},
                PermissionScopeManage: {PermissionTypeRead,
                    PermissionTypeUpdate, PermissionTypeList},

            }, resources[1].Permission},
        }
        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("ReadPermissionsOnDocument %d: bad value, got: " +
                    "%v want: %v", i, test.got, test.want)
            }
        }
    }

    // Test ReadPermissionsOnUser
    resources, err = custodia.ReadPermissionsOnUser(dummyUUID)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {"Structure", resources[0].Access},
            {ResourceUser, resources[0].OwnerType},
            {dummyUUID, resources[0].OwnerId},
            {"", resources[0].ParentId},
            {map[PermissionScope][]PermissionType{
                PermissionScopeManage: {PermissionTypeRead,
                    PermissionTypeUpdate}},
                resources[0].Permission},
        }
        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("ReadPermissionsOnUser %d: bad value, got: " +
                    "%v want: %v", i, test.got, test.want)
            }
        }
    }

    // Test ReadPermissionsOnGroup
    resources, err = custodia.ReadPermissionsOnGroup(dummyUUID)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    } else {
        var tests = []struct {
            want interface{}
            got interface{}
        }{
            {"Structure", resources[0].Access},
            {ResourceGroup, resources[0].OwnerType},
            {dummyUUID, resources[0].OwnerId},
            {"", resources[0].ParentId},
            {map[PermissionScope][]PermissionType{
                PermissionScopeManage: {PermissionTypeList,
                    PermissionTypeRead}},
                resources[0].Permission},
        }
        for i, test := range tests {
            if !reflect.DeepEqual(test.want, test.got) {
                t.Errorf("ReadPermissionsOnGroup %d: bad value, got: " +
                    "%v want: %v", i, test.got, test.want)
            }
        }
    }
}

