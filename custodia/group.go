package custodia

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/simplereach/timeutils"
)

type Group struct {
	Id uuid.UUID `json:"group_id,omitempty"`
	Name string `json:"group_name"`
	InsertDate timeutils.Time `json:"insert_date,omitempty"`
	LastUpdate timeutils.Time `json:"last_update,omitempty"`
	IsActive bool `json:"is_active"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

type GroupEnvelope struct {
	Group *Group `json:"group"`
}

type GroupsEnvelope struct {
	Groups []Group `json:"groups"`
}

// [C]reate a new group
func (ca *CustodiaAPIv1) CreateGroup(name string, isActive bool,
	attributes map[string]any) (*Group, error) {
	group := Group{Name: name, IsActive: isActive, Attributes: attributes}
	url := "/groups"
	params := map[string]any{"_data": group}
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
func (ca *CustodiaAPIv1) ReadGroup(groupId uuid.UUID) (*Group, error) {
	url := fmt.Sprintf("/groups/%s", groupId)
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
func (ca *CustodiaAPIv1) UpdateGroup(groupId uuid.UUID, name string,
	isActive bool, attributes map[string]any) (*Group, error) {
	group := Group{Name: name, IsActive: isActive, Attributes: attributes}
	url := fmt.Sprintf("/groups/%s", groupId)
	params := map[string]any{"_data": group}
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
func (ca *CustodiaAPIv1) DeleteGroup(groupId uuid.UUID, force bool) error {
	url := fmt.Sprintf("/groups/%s?force=%v", groupId, force)
	_, err := ca.Call("DELETE", url, nil)
	return err
}

// [L]ist all groups
func (ca *CustodiaAPIv1) ListGroups(queryParams map[string]string) (
	[]Group, error,
) {
	u, err := url.Parse("/groups")
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %v", err)
	}

	// Adding query params
	q := u.Query()
	for k, v := range queryParams {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	resp, err := ca.Call("GET", u.String(), nil)
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
func (ca *CustodiaAPIv1) ListGroupUsers(groupId uuid.UUID,
	queryParams map[string]string) ([]User, error,
) {
	u, err := url.Parse(fmt.Sprintf("/groups/%s/users", groupId))
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %v", err)
	}

	// Adding query params
	q := u.Query()
	for k, v := range queryParams {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	resp, err := ca.Call("GET", u.String(), nil)
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
func (ca *CustodiaAPIv1) AddUserToGroup(userId uuid.UUID, groupId uuid.UUID) (
	error) {
	url := fmt.Sprintf("/groups/%s/users/%s", groupId, userId)
	_, err := ca.Call("POST", url, nil)
	if err != nil {
		return err
	}
	return nil
}

// [C] Add all users of a UserSchema to the group
func (ca *CustodiaAPIv1) AddUsersFromUserSchemaToGroup(
    userSchemaId uuid.UUID, groupId uuid.UUID) error {
    url := fmt.Sprintf("/groups/%s/user_schemas/%s", groupId, userSchemaId)
    _, err := ca.Call("POST", url, nil)
    if err != nil {
        return err
    }
    return nil
}

// [D] Remove a user from the group
func (ca *CustodiaAPIv1) RemoveUserFromGroup(userId uuid.UUID,
	groupId uuid.UUID) (error) {
    url := fmt.Sprintf("/groups/%s/users/%s", groupId, userId)
    _, err := ca.Call("DELETE", url, nil)
    if err != nil {
        return err
    }
    return nil
}

// [D] Remove all users of a UserSchema from the group
func (ca *CustodiaAPIv1) RemoveUsersFromUserSchemaFromGroup(
    userSchemaId uuid.UUID, groupId uuid.UUID) error {
    url := fmt.Sprintf("/groups/%s/user_schemas/%s", groupId, userSchemaId)
    _, err := ca.Call("DELETE", url, nil)
    if err != nil {
        return err
    }
    return nil
}
