package custodia

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dzanotelli/chino/common"
)

// Define PermissionType
type PermissionType int

const (
    PermissionCreate PermissionType = iota + 1
	PermissionRead
    PermissionUpdate
    PermissionDelete
	PermissionList
    PermissionSearch
    PermissionAuthorize
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
	ResourceRepository ResourceType = iota
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
	Permission map[PermissionAction][]PermissionType `json:"permission"`
}

// Grant or revoke permissions over resources of a specific type.
// It can be used only on Top Level resources.
func (ca *CustodiaAPIv1) PermissionOnResources(action PermissionAction,
	resourceType ResourceType, subjectType ResourceType, subjectId string,
	permissions map[PermissionAction][]PermissionType) (
	error) {
	if !common.IsValidUUID(subjectId) {
		return errors.New("subjectId is not a valid UUID: " + subjectId)
    }

	url := fmt.Sprintf("/perms/%s/%s/%s/%s", action, resourceType, subjectType,
		subjectId)
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
	subjectId string, permissions map[PermissionAction][]PermissionType) (
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

