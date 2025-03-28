package custodia

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/simplereach/timeutils"
)


type User struct {
	Id uuid.UUID `json:"user_id,omitempty"`
	UserSchemaId uuid.UUID `json:"schema_id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	InsertDate timeutils.Time `json:"insert_date,omitempty"`
	LastUpdate timeutils.Time `json:"last_update,omitempty"`
	IsActive bool `json:"is_active"`
	Attributes map[string]any `json:"attributes,omitempty"`
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
	attributes map[string]any) (*User, error) {
	// validate user content
	contentErrors := validateContent(attributes,
		userSchema.getStructureAsMap())
	if len(contentErrors) > 0 {
		e := fmt.Errorf("content errors: %w", errors.Join(contentErrors...))
		return nil, e
	}

	doc := User{IsActive: isActive, Attributes: attributes}
	url := fmt.Sprintf("/user_schemas/%s/users", userSchema.Id)
	params := map[string]any{"_data": doc}
	resp, err := ca.Call("POST", url, params)
	if err != nil {
		return nil, err
	}
	// JSON: unmarshal resp content
	userEnvelope := UserEnvelope{}
	if err := json.Unmarshal([]byte(resp), &userEnvelope); err != nil {
		return nil, err
	}

	// if everything is ok, we can safely set the given content as the
	// returned user content, since the API doesn't return it
	userEnvelope.User.Attributes = attributes

	return userEnvelope.User, nil
}

// [R]ead an existent user
func (ca *CustodiaAPIv1) ReadUser(userSchema UserSchema, userId uuid.UUID) (
	*User, error) {
	url := fmt.Sprintf("/users/%s", userId)
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
	converted, ee := convertData(userEnvelope.User.Attributes, &userSchema)
	if len(ee) > 0 {
		err := fmt.Errorf("conversion errors: %w", errors.Join(ee...))
		return userEnvelope.User, err
	}

	// all good, assign the new content to doc and return it
	userEnvelope.User.Attributes = converted
	return userEnvelope.User, nil
}

// [U]pdate an existent user
func (ca *CustodiaAPIv1) UpdateUser(userId uuid.UUID , isActive bool,
	content map[string]any) (*User, error) {
	url := fmt.Sprintf("/users/%s", userId)

	// create a user with just the values we can send, and marshal it
	user := User{IsActive: isActive, Attributes: content}
	params := map[string]any{"_data": user}
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
func (ca *CustodiaAPIv1) DeleteUser(userId uuid.UUID, force, consistent bool) (
	error) {
	url := fmt.Sprintf("/users/%s", userId)
	url += fmt.Sprintf("?force=%v&consistent=%v", force, consistent)

	_, err := ca.Call("DELETE", url, nil)
	if err != nil {
		return err
	}
	return nil
}

// [L]ist all the users in a Schema
// url queryParams:
//   offset: int: number of items to skip from the beginning of the list
//   limit: int : maximum number of items to return in a single page
//   full_user: bool: return the full document
//   is_active: bool: filter by
//   insert_date__gt: time string (RFC3339, YYYY-MM-DDTHH:MM:SS): filter by
//   insert_date__lt: time string (RFC3339): filter by
//   last_update__gt: time string (RFC3339): filter by
//   last_update__lt: time string (RFC3339): filter by
func (ca *CustodiaAPIv1) ListUsers(userSchemaId uuid.UUID,
	queryParams map[string]string) ([]*User, error,
) {
	u, err := url.Parse(fmt.Sprintf("/user_schemas/%s/users", userSchemaId))
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %v", err)
	}

	// Adding query params
	q := u.Query()
	for k, v := range(queryParams) {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	resp, err := ca.Call("GET", u.String(), nil)
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
