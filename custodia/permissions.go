package custodia

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dzanotelli/chino/common"
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
    err := json.Unmarshal(data, &value)
    if err!= nil {
        return err
    }
    intValue := indexOf(value, ps.Choices()) + 1  // enum starts from 1
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
    PermissionActionCreate PermissionType = iota + 1
	PermissionActionRead
    PermissionActionUpdate
    PermissionActionDelete
	PermissionActionList
    PermissionActionSearch
    PermissionActionAuthorize
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
	intValue := indexOf(value, pt.Choices()) + 1  // enum starts from 1
	if intValue < 1 {
        return fmt.Errorf("PermissionType: received unknown value '%v'", value)
    }

    *pt = PermissionType(intValue)
    return nil
}

type PermissionAction int

const (
	Grant PermissionAction = iota +1
	Revoke
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
    intValue := indexOf(value, pa.Choices()) + 1  // enum starts from 1
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
    return []string{"repositories", "schemas", "documents", "user_schemas",
        "users", "groups", "collections"}
}

func (rt ResourceType) String() string {
    return rt.Choices()[rt-1]
}


func (rt ResourceType) MarshalJSON() ([]byte, error) {
    return json.Marshal(rt.String())
}

func (rt* ResourceType) UnmarshalJSON(data []byte) error {
	var value string
    err := json.Unmarshal(data, &value)
    if err!= nil {
        return err
    }
    intValue := indexOf(value, rt.Choices()) + 1  // enum starts from 1
    if intValue < 1 {
        return fmt.Errorf("ResourceType: received unknown value '%v'", value)
    }

    *rt = ResourceType(intValue)
    return nil
}


type Resource struct {
	Access string `json:"access"`
	ParentId string `json:"parent_id,omitempty"`
	Id string `json:"resource_id"`
	Type ResourceType `json:"resource_type"`
	OwnerId string `json:"owner_id,omitempty"`
	OwnerType ResourceType `json:"owner_type,omitempty"`
	Permission map[PermissionScope][]PermissionType `json:"permission"`
}

// Grant or revoke permissions over resources of a specific type.
// It can be used only on Top Level resources.
func (ca *CustodiaAPIv1) PermissionOnResources(action PermissionAction,
	resourceType ResourceType, subjectType ResourceType, subjectId string,
	permissions map[PermissionScope][]PermissionType) (
	error) {
	if !common.IsValidUUID(subjectId) {
		return errors.New("subjectId is not a valid UUID: " + subjectId)
    }

	url := fmt.Sprintf("/perms/%s/%s/%s/%s", action.String(),
        resourceType.String(), subjectType.String(), subjectId)
	params := map[string]interface{}{"data": permissions}
	_, err := ca.Call("POST", url, params)
    if err!= nil {
        return err
    }
	return nil
}

// Grant or Revoke permissions over a specific resource.
// It can be called on all resources.
func (ca *CustodiaAPIv1) PermissionOnResource(action PermissionAction,
    resourceType ResourceType, resourceId string, subjectType ResourceType,
	subjectId string, permissions map[PermissionScope][]PermissionType) (
    error) {
	if !common.IsValidUUID(resourceId) {
        return errors.New("resourceId is not a valid UUID: " + resourceId)
    }
	if !common.IsValidUUID(subjectId) {
		return errors.New("subjectId is not a valid UUID: " + subjectId)
	}

    url := fmt.Sprintf("/perms/%s/%s/%s/%s/%s", action, resourceType,
		resourceId, subjectType, subjectId)
    params := map[string]interface{}{"data": permissions}
    _, err := ca.Call("POST", url, params)
    if err!= nil {
        return err
    }
    return nil
}

// Grant or Revoke permissions over all the children of a specific resource.
// It can be used only on resources that have a parent-child relationship.
func (ca *CustodiaAPIv1) PermissionOnResourceChildren(action PermissionAction,
    resourceType ResourceType, resourceId string,
	resourceChildType ResourceType, subjectType ResourceType,
	subjectId string, permissions map[PermissionScope][]PermissionType) (
	error) {
	if !common.IsValidUUID(resourceId) {
		return errors.New("resourceId is not a valid UUID: " + resourceId)
	}
	if !common.IsValidUUID(subjectId) {
		return errors.New("subjectId is not a valid UUID: " + subjectId)
	}

	url := fmt.Sprintf("/perms/%s/%s/%s/%s/%s/%s", action, resourceType,
	resourceId, resourceChildType, subjectType, subjectId)
	params := map[string]interface{}{"data": permissions}
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
func (ca *CustodiaAPIv1) ReadPermissionsOnDocument(documentId string) (
    []Resource, error) {
	if !common.IsValidUUID(documentId) {
		return nil, errors.New("documentId is not a valid UUID: " + documentId)
	}

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
func (ca *CustodiaAPIv1) ReadPermissionsOnUser(userId string) ([]Resource,
	error) {
    if !common.IsValidUUID(userId) {
        return nil, errors.New("userId is not a valid UUID: " + userId)
    }

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
func (ca *CustodiaAPIv1) ReadPermissionsOnGroup(groupId string) ([]Resource,
	error) {
	if!common.IsValidUUID(groupId) {
        return nil, errors.New("groupId is not a valid UUID: " + groupId)
    }

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