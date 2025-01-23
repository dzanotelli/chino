package custodia

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dzanotelli/chino/common"
	"github.com/simplereach/timeutils"
)

type Group struct {
	Id string `json:"group_id,omitempty"`
	Name string `json:"group_name"`
	InsertDate timeutils.Time `json:"insert_date,omitempty"`
	LastUpdate timeutils.Time `json:"last_update,omitempty"`
	IsActive bool `json:"is_active"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

type GroupEnvelope struct {
	Group *Group `json:"group"`
}

type GroupsEnvelope struct {
	Groups []Group `json:"groups"`
}

// [C]reate a new group
func (ca *CustodiaAPIv1) CreateGroup(name string, isActive bool,
	attributes map[string]interface{}) (*Group, error) {
	group := Group{Name: name, IsActive: isActive, Attributes: attributes}
	url := "/groups"
	params := map[string]interface{}{"data": group}
	resp, err := ca.Call("POST", url, params)
	if err != nil {
		return nil, err
	}
	// JSON: unmarshal resp content
	groupEnvelope := GroupEnvelope{}
	if err := json.Unmarshal([]byte(resp), &groupEnvelope); err != nil {
		return nil, err
	}
	return groupEnvelope.Group, nil
}

// [R]ead an existent group
func (ca *CustodiaAPIv1) ReadGroup(id string) (*Group, error) {
	if !common.IsValidUUID(id) {
		return nil, errors.New("id is not a valid UUID: " + id)
	}

	url := fmt.Sprintf("/groups/%s", id)
	resp, err := ca.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// JSON: unmarshal resp content
	groupEnvelope := GroupEnvelope{}
	if err := json.Unmarshal([]byte(resp), &groupEnvelope); err != nil {
		return nil, err
	}
	return groupEnvelope.Group, nil
}

// [U]pdate an existent group
func (ca *CustodiaAPIv1) UpdateGroup(id string, name string, isActive bool,
	attributes map[string]interface{}) (*Group, error) {
	if !common.IsValidUUID(id) {
		return nil, errors.New("id is not a valid UUID: " + id)
	}
	group := Group{Name: name, IsActive: isActive, Attributes: attributes}
	url := fmt.Sprintf("/groups/%s", id)
	params := map[string]interface{}{"data": group}
	resp, err := ca.Call("PUT", url, params)
	if err != nil {
		return nil, err
	}
	// JSON: unmarshal resp content
	groupEnvelope := GroupEnvelope{}
	if err := json.Unmarshal([]byte(resp), &groupEnvelope); err != nil {
		return nil, err
	}
	return groupEnvelope.Group, nil
}

// [D]elete an existent group
func (ca *CustodiaAPIv1) DeleteGroup(id string, force bool) error {
	if !common.IsValidUUID(id) {
		return errors.New("id is not a valid UUID: " + id)
	}
	url := fmt.Sprintf("/groups/%s?force=%v", id, force)
	_, err := ca.Call("DELETE", url, nil)
	return err
}

// [L]ist all groups
func (ca *CustodiaAPIv1) ListGroups() ([]Group, error) {
	url := "/groups"
	resp, err := ca.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// JSON: unmarshal resp content
	groupsEnvelope := GroupsEnvelope{}
	if err := json.Unmarshal([]byte(resp), &groupsEnvelope); err != nil {
		return nil, err
	}
	return groupsEnvelope.Groups, nil
}

// Group Members

// [L]ist group's users
func (ca *CustodiaAPIv1) ListGroupUsers(groupId string) ([]User, error) {
	if !common.IsValidUUID(groupId) {
		return nil, errors.New("groupId is not a valid UUID: " + groupId)
	}

	url := fmt.Sprintf("/groups/%s/users", groupId)
	resp, err := ca.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// JSON: unmarshal resp content
	usersEnvelope := UsersEnvelope{}
	if err := json.Unmarshal([]byte(resp), &usersEnvelope); err != nil {
		return nil, err
	}
	return usersEnvelope.Users, nil
}

// [C] Add a user to the group
func (ca *CustodiaAPIv1) AddUserToGroup(userId string, groupId string) (
	error) {
	if !common.IsValidUUID(userId) {
		return errors.New("userId is not a valid UUID: " + userId)
	}
	if !common.IsValidUUID(groupId) {
		return errors.New("groupId is not a valid UUID: " + groupId)
	}

	url := fmt.Sprintf("/groups/%s/users/%s", groupId, userId)
	_, err := ca.Call("POST", url, nil)
	if err != nil {
		return err
	}
	return nil
}

// [C] Add all users of a UserSchema to the group
func (ca *CustodiaAPIv1) AddUsersFromUserSchemaToGroup(
    userSchemaId string, groupId string) error {
    if !common.IsValidUUID(userSchemaId) {
        return errors.New("userSchemaId is not a valid UUID: " +
            userSchemaId)
    }
    if !common.IsValidUUID(groupId) {
        return errors.New("groupId is not a valid UUID: " + groupId)
    }

    url := fmt.Sprintf("/groups/%s/user_schemas/%s", groupId, userSchemaId)
    _, err := ca.Call("POST", url, nil)
    if err != nil {
        return err
    }
    return nil
}

// [D] Remove a user from the group
func (ca *CustodiaAPIv1) RemoveUserFromGroup(userId string, groupId string) (
    error) {
    if !common.IsValidUUID(userId) {
        return errors.New("userId is not a valid UUID: " + userId)
    }
    if !common.IsValidUUID(groupId) {
        return errors.New("groupId is not a valid UUID: " + groupId)
    }

    url := fmt.Sprintf("/groups/%s/users/%s", groupId, userId)
    _, err := ca.Call("DELETE", url, nil)
    if err != nil {
        return err
    }
    return nil
}

// [D] Remove all users of a UserSchema from the group
func (ca *CustodiaAPIv1) RemoveUsersFromUserSchemaFromGroup(
    userSchemaId string, groupId string) error {
    if !common.IsValidUUID(userSchemaId) {
        return errors.New("userSchemaId is not a valid UUID: " +
            userSchemaId)
    }
    if !common.IsValidUUID(groupId) {
        return errors.New("groupId is not a valid UUID: " + groupId)
    }

    url := fmt.Sprintf("/groups/%s/user_schemas/%s", groupId, userSchemaId)
    _, err := ca.Call("DELETE", url, nil)
    if err != nil {
        return err
    }
    return nil
}