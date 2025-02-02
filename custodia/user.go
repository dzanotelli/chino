package custodia

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dzanotelli/chino/common"
	"github.com/simplereach/timeutils"
)


type User struct {
	Id string `json:"user_id,omitempty"`
	UserSchemaId string `json:"schema_id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	InsertDate timeutils.Time `json:"insert_date,omitempty"`
	LastUpdate timeutils.Time `json:"last_update,omitempty"`
	IsActive bool `json:"is_active"`
	Attributes map[string]interface{} `json:"content,omitempty"`
	Groups []string `jsong:"groups,omitempty"`
}

type UserEnvelope struct {
	User *User `json:"user"`
}

type UsersEnvelope struct {
	Users []User `json:"users"`
}

// [C]reate a new user
func (ca *CustodiaAPIv1) CreateUser(userSchema *UserSchema, isActive bool,
	attributes map[string]interface{}) (*User, error) {
	if userSchema.Id == "" {
		return nil, fmt.Errorf("user schema has no UserSchemaId, " +
			"does it exist?")
	} else if !common.IsValidUUID(userSchema.Id) {
		return nil, fmt.Errorf("SchemaId is not a valid UUID: %s (it " +
			"should not be manually set)", userSchema.Id)
	}

	// validate user content
	contentErrors := validateContent(attributes, userSchema.getStructureAsMap())
	if len(contentErrors) > 0 {
		e := fmt.Errorf("content errors: %w", errors.Join(contentErrors...))
		return nil, e
	}

	doc := User{IsActive: isActive, Attributes: attributes}
	url := fmt.Sprintf("/user_schemas/%s/users", userSchema.Id)
	params := map[string]interface{}{"_data": doc}
	resp, err := ca.Call("POST", url, params)
	if err != nil {
		return nil, err
	}
	// JSON: unmarshal resp content
	docEnvelope := UserEnvelope{}
	if err := json.Unmarshal([]byte(resp), &docEnvelope); err != nil {
		return nil, err
	}

	// if everything is ok, we can safely set the given content as the
	// returned user content, since the API doesn't return it
	docEnvelope.User.Attributes = attributes

	return docEnvelope.User, nil
}

// [R]ead an existent user
func (ca *CustodiaAPIv1) ReadUser(schema UserSchema, id string) (*User,
	 error) {
	if !common.IsValidUUID(id) {
		return nil, errors.New("id is not a valid UUID: " + id)
	}

	url := fmt.Sprintf("/users/%s", id)
	resp, err := ca.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	userEnvelope := UserEnvelope{}
	if err := json.Unmarshal([]byte(resp), &userEnvelope); err != nil {
		return nil, err
	}

	// convert values to concrete types
	converted, ee := convertData(userEnvelope.User.Attributes, &schema)
	if len(ee) > 0 {
		err := fmt.Errorf("conversion errors: %w", errors.Join(ee...))
		return userEnvelope.User, err
	}

	// all good, assign the new content to doc and return it
	userEnvelope.User.Attributes = converted
	return userEnvelope.User, nil
}

// [U]pdate an existent user
func (ca *CustodiaAPIv1) UpdateUser(id string , isActive bool,
	content map[string]interface{}) (*User, error) {
		if !common.IsValidUUID(id) {
			return nil, errors.New("id is not a valid UUID: " + id)
		}

	url := fmt.Sprintf("/users/%s", id)

	// create a user with just the values we can send, and marshal it
	user := User{IsActive: isActive, Attributes: content}
	params := map[string]interface{}{"_data": user}
	resp, err := ca.Call("PUT", url, params)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content and return a fresh user instance
	docEnvelope := UserEnvelope{}
	if err := json.Unmarshal([]byte(resp), &docEnvelope); err != nil {
		return nil, err
	}

	// PUT call returns the whole users, along with its content
	return docEnvelope.User, nil
}

// [D]elete an existent user
// if force=false user is just deactivated
// if consisten=true the operation is done sync (server waits to respond)
func (ca *CustodiaAPIv1) DeleteUser(id string, force, consistent bool) (
	error) {
	url := fmt.Sprintf("/users/%s", id)
	url += fmt.Sprintf("?force=%v&consistent=%v", force, consistent)

	_, err := ca.Call("DELETE", url, nil)
	if err != nil {
		return err
	}
	return nil
}

// [L]ist all the users in a Schema
// url params:
//   full_user: bool
//   is_active: bool
//   insert_date__gt: time.Time
//   insert_date__lt: time.Time
//   last_update__gt: time.Time
//   last_update__lt: time.Time
func (ca *CustodiaAPIv1) ListUsers(userSchemaId string,
	params map[string]interface{}) ([]*User, error) {
	if !common.IsValidUUID(userSchemaId) {
		return nil, fmt.Errorf("schemaId is not a valid UUID: %s",
			userSchemaId)
	}

	url := fmt.Sprintf("/user_schemas/%s/users", userSchemaId)
	if len(params) > 0 {
		url += "?"
	}

	availableParams := map[string]interface{}{
		"full_user": true,
		"is_active": true,
		"insert_date__gt": time.Time{},
		"insert_date__lt": time.Time{},
		"last_update__gt": time.Time{},
		"last_update__lt": time.Time{},
	}
	for param := range params {
		// check that param is legit
		_, ok := availableParams[param]
		if !ok {
			return nil, fmt.Errorf("got unexpected param '%s'", param)
		}
		value := params[param]

		switch param {
		case "full_user", "is_active":
			_, ok := value.(bool)
			if !ok {
				err := fmt.Errorf("param '%s': bad type: '%T', must be bool",
					param, value)
				return nil, err
			}
		case "insert_date__gt", "insert_date__lt", "last_update__gt",
			"last_update__lt":
			time_value, ok := value.(time.Time)
			if !ok {
				err := fmt.Errorf("param '%s': bad type: '%T', must be Time",
					param, value)
				return nil, err
			}
			value = time_value.Format(time.RFC3339)
		default:
			return nil, fmt.Errorf("got unexpected param '%s'", param)
		}

		url += fmt.Sprintf("%s=%v&", param, value)
	}

	url = strings.TrimRight(url, "&")
	resp, err := ca.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	docusEnvelope := UsersEnvelope{}
	if err := json.Unmarshal([]byte(resp), &docusEnvelope); err != nil {
		return nil, err
	}

	result := []*User{}
	for _, doc := range docusEnvelope.Users {
		result = append(result, &doc)
	}

	return result, nil
}
