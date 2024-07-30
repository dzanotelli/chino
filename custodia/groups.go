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
func (ca *CustodiaAPIv1) DeleteGroup(id string) error {
	if !common.IsValidUUID(id) {
		return errors.New("id is not a valid UUID: " + id)
	}
	url := fmt.Sprintf("/groups/%s", id)
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
