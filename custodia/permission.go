package custodia

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

// Define PermissionContext
type PermissionScope int

const (
    PermissionScopeManage PermissionScope = iota + 1
    PermissionScopeAuthorize
)


func (ps PermissionScope) Choices() []string {
    return []string{"manage", "authorize"}
}

func (ps PermissionScope) String() string {
    return ps.Choices()[ps-1]
}

func (ps PermissionScope) MarshalJSON() ([]byte, error) {
    return json.Marshal(ps.String())
}

func (ps* PermissionScope) UnmarshalJSON(data []byte) error {
    var value string

    // NOTE: the API is inconsistent, sometimes returns a camel-cased string
    //   sometimes just all lower. We need to lowerize our value always.
    rawValue := strings.Trim(string(data), "\"")
    byteValue := []byte(fmt.Sprintf(`"%s"`, strings.ToLower(rawValue)))

    err := json.Unmarshal(byteValue, &value)
    if err!= nil {
        return err
    }
    intValue := slices.Index(ps.Choices(), value) + 1   // enum starts from 1
    if intValue < 1 {
        return fmt.Errorf("PermissionScope: received unknown value '%v'",
            value)
    }

    *ps = PermissionScope(intValue)
    return nil
}

// Define PermissionType
type PermissionType int

const (
    PermissionTypeCreate PermissionType = iota + 1
	PermissionTypeRead
    PermissionTypeUpdate
    PermissionTypeDelete
	PermissionTypeList
    PermissionTypeSearch
    PermissionTypeAuthorize
)

func (pt PermissionType) Choices() []string {
	return []string{"C", "R", "U", "D", "L", "S", "A"}
}

func (pt PermissionType) String() string {
    return pt.Choices()[pt-1]
}

func (pt PermissionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(pt.String())
}

func (pt* PermissionType) UnmarshalJSON(data []byte) error {
	var value string
    err := json.Unmarshal(data, &value)
    if err!= nil {
        return err
    }
    intValue := slices.Index(pt.Choices(), value) + 1   // enum starts from 1
	if intValue < 1 {
        return fmt.Errorf("PermissionType: received unknown value '%v'", value)
    }

    *pt = PermissionType(intValue)
    return nil
}

type PermissionAction int

const (
	PermissionActionGrant PermissionAction = iota +1
	PermissionActionRevoke
)

func (pa PermissionAction) Choices() []string {
	return []string{"grant", "revoke"}
}

func (pa PermissionAction) String() string {
	return pa.Choices()[pa-1]
}

func (pa PermissionAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(pa.String())
}

func (pa* PermissionAction) UnmarshalJSON(data []byte) error {
	var value string
    err := json.Unmarshal(data, &value)
    if err!= nil {
        return err
    }
    intValue := slices.Index(pa.Choices(), value) + 1   // enum starts from 1
    if intValue < 1 {
        return fmt.Errorf("PermissionAction: received unknown value '%v'", value)
    }

    *pa = PermissionAction(intValue)
    return nil
}

type ResourceType int

const (
	ResourceRepository ResourceType = iota + 1
	ResourceSchema
	ResourceDocument
	ResourceUserSchema
	ResourceUser
	ResourceGroup
	ResourceCollection
)


func (rt ResourceType) Choices() []string {
    return []string{"Repository", "Schema", "Document", "UserSchema",
        "User", "Group", "Collection"}

}

func (rt ResourceType) String() string {
    return rt.Choices()[rt-1]
}

func (rt ResourceType) UrlChoices() []string {
    return []string{"repositories", "schemas", "documents", "user_schemas",
        "users", "groups", "collections"}
}

func (rt ResourceType) UrlString() string {
    return rt.UrlChoices()[rt-1]
}

func (rt ResourceType) MarshalJSON() ([]byte, error) {
    return json.Marshal(rt.String())
}

func (rt* ResourceType) UnmarshalJSON(data []byte) error {
    // NOTE: the API is not consistent. Sometimes returns the singular
    //   camel-case version of the resource type, sometimes the lower plural
    //   we need to handle both cases.
    rawValue := strings.Trim(string(data), "\"")

    // single camle-case format
    choices := rt.Choices()
    if slices.Contains(choices, rawValue) {
        val := slices.Index(choices, rawValue)
        if val != -1 {
            *rt = ResourceType(val+1)
        }
        return nil
    }
    // plural format
    choices = rt.UrlChoices()
    if slices.Contains(choices, rawValue) {
        val := slices.Index(choices, rawValue)
        if val != -1 {
            *rt = ResourceType(val+1)
        }
        return nil
    }

    return fmt.Errorf("ResourceType: received unknown value '%v'", rawValue)

	// var value string
    // err := json.Unmarshal(data, &value)
    // if err!= nil {
    //     return err
    // }
    // intValue := indexOf(value, rt.Choices()) + 1  // enum starts from 1
    // if intValue < 1 {
    //     return fmt.Errorf("ResourceType: received unknown value '%v'", value)
    // }

    // *rt = ResourceType(intValue)
    // return nil
}


type Resource struct {
	Access string `json:"access"`
	ParentId uuid.UUID `json:"parent_id,omitempty"`
	Id uuid.UUID `json:"resource_id"`
	Type ResourceType `json:"resource_type"`
	OwnerId uuid.UUID `json:"owner_id,omitempty"`
	OwnerType ResourceType `json:"owner_type,omitempty"`
	Permission map[PermissionScope][]PermissionType `json:"permission"`
}

func (r* Resource) UnmarshalJSON(data []byte) error {
    // create an alias of Resource in order to transform the map Permission
    // with string keys
	type Alias Resource   // inherits from Resource
	aux := struct {
		Permission map[string][]PermissionType `json:"permission"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

    // unmarshal using the alias, all resource rules are preserved
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

    r.Permission = make(map[PermissionScope][]PermissionType)
    for scope, types := range aux.Permission {
        found := slices.Index(PermissionScope(1).Choices(),
            strings.ToLower(scope))
        ps := PermissionScope(found + 1)
        r.Permission[ps] = types
    }
    return nil
}

// Grant or revoke permissions over resources of a specific type.
// It can be used only on Top Level resources.
func (ca *CustodiaAPIv1) PermissionOnResources(action PermissionAction,
	resourceType ResourceType, subjectType ResourceType, subjectId uuid.UUID,
	permissions map[PermissionScope][]PermissionType) (
	error) {
	url := fmt.Sprintf("/perms/%s/%s/%s/%s", action, resourceType.UrlString(),
        subjectType.UrlString(), subjectId.String())
	params := map[string]any{"data": permissions}
	_, err := ca.Call("POST", url, params)
    if err!= nil {
        return err
    }
	return nil
}

// Grant or Revoke permissions over a specific resource.
// It can be called on all resources.
func (ca *CustodiaAPIv1) PermissionOnResource(action PermissionAction,
    resourceType ResourceType, resourceId uuid.UUID, subjectType ResourceType,
	subjectId uuid.UUID, permissions map[PermissionScope][]PermissionType) (
    error) {
    url := fmt.Sprintf("/perms/%s/%s/%s/%s/%s", action,
        resourceType.UrlString(), resourceId.String(), subjectType.UrlString(),
        subjectId.String())
    params := map[string]any{"data": permissions}
    _, err := ca.Call("POST", url, params)
    if err!= nil {
        return err
    }
    return nil
}

// Grant or Revoke permissions over all the children of a specific resource.
// It can be used only on resources that have a parent-child relationship.
func (ca *CustodiaAPIv1) PermissionOnResourceChildren(action PermissionAction,
    resourceType ResourceType, resourceId uuid.UUID,
	resourceChildType ResourceType, subjectType ResourceType,
	subjectId uuid.UUID, permissions map[PermissionScope][]PermissionType) (
	error) {
	url := fmt.Sprintf("/perms/%s/%s/%s/%s/%s/%s", action,
        resourceType.UrlString(), resourceId.String(),
        resourceChildType.UrlString(), subjectType.UrlString(),
        subjectId.String())
	params := map[string]any{"data": permissions}
	_, err := ca.Call("POST", url, params)
	if err!= nil {
		return err
	}
	return nil
}

// Read permissions on all resources
func (ca *CustodiaAPIv1) ReadAllPermissions() ([]Resource, error) {
	url := "/perms"
    resp, err := ca.Call("GET", url, nil)
    if err!= nil {
        return nil, err
    }

    // JSON: unmarshal resp content
    resourcesEnvelope := map[string][]Resource{}
    if err := json.Unmarshal([]byte(resp), &resourcesEnvelope); err!= nil {
        return nil, err
    }

	permissions, ok := resourcesEnvelope["permissions"]
	if !ok {
		return nil, fmt.Errorf("missing 'permissions' key in response")
	}
    return permissions, nil
}

// Read permissions over a document
func (ca *CustodiaAPIv1) ReadPermissionsOnDocument(documentId uuid.UUID) (
    []Resource, error) {
	url := fmt.Sprintf("/perms/documents/%s", documentId)
	resp, err := ca.Call("GET", url, nil)
	if err!= nil {
        return nil, err
    }
	// JSON: unmarshal resp content
	resourcesEnvelope := map[string][]Resource{}
    if err := json.Unmarshal([]byte(resp), &resourcesEnvelope); err!= nil {
        return nil, err
    }

    permissions, ok := resourcesEnvelope["permissions"]
    if!ok {
        return nil, fmt.Errorf("missing 'permissions' key in response")
    }
    return permissions, nil
}

// Read permissions over a user.
// List all the permissions that the user has on Resources.
func (ca *CustodiaAPIv1) ReadPermissionsOnUser(userId uuid.UUID) ([]Resource,
	error) {
    url := fmt.Sprintf("/perms/users/%s", userId)
    resp, err := ca.Call("GET", url, nil)
    if err!= nil {
        return nil, err
    }
    // JSON: unmarshal resp content
    resourcesEnvelope := map[string][]Resource{}
    if err := json.Unmarshal([]byte(resp), &resourcesEnvelope); err!= nil {
        return nil, err
    }

    permissions, ok := resourcesEnvelope["permissions"]
    if !ok {
        return nil, fmt.Errorf("missing 'permissions' key in response")
    }
    return permissions, nil
}

// Read permissions over a group.
func (ca *CustodiaAPIv1) ReadPermissionsOnGroup(groupId uuid.UUID) ([]Resource,
	error) {
    url := fmt.Sprintf("/perms/groups/%s", groupId)
    resp, err := ca.Call("GET", url, nil)
    if err!= nil {
        return nil, err
    }
    // JSON: unmarshal resp content
    resourcesEnvelope := map[string][]Resource{}
    if err := json.Unmarshal([]byte(resp), &resourcesEnvelope); err!= nil {
        return nil, err
    }

    permissions, ok := resourcesEnvelope["permissions"]
    if!ok {
        return nil, fmt.Errorf("missing 'permissions' key in response")
    }
    return permissions, nil
}