package custodia

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dzanotelli/chino/common"
	"github.com/simplereach/timeutils"
)

type UserSchema struct {
	Id string `json:"user_schema_id,omitempty"`
	Description string `json:"description"`
	Groups []string `json:"gropus,omitempty"`
	InsertDate timeutils.Time `json:"insert_date,omitempty"`
	LastUpdate timeutils.Time `json:"last_update,omitempty"`
	IsActive bool `json:"is_active"`
	Structure []SchemaField `json:"structure"`
}

type UserSchemaEnvelope struct {
	UserSchema *UserSchema `json:"user_schema"`
}

type UserSchemasEnvelope struct {
	UserSchemas []UserSchema `json:"user_schemas"`
}

// adjustDefaultTypes for each field calls adjustDefaultType
func (s *UserSchema) adjustDefaultTypes() {
	for i := range s.Structure {
		s.Structure[i].adjustDefaultType()
	}
}

// [C]reate a new user schema
func (ca *CustodiaAPIv1) CreateUserSchema(descritpion string, isActive bool,
	 fields []SchemaField) (*UserSchema, error) {
	// FIXME: missing field type validation, and indexed property validation
	//   and insensitive property

	user_schema := UserSchema{Description: descritpion, Structure: fields,
		 IsActive: isActive}
	url := "/user_schemas"
	params := map[string]interface{}{"data": user_schema}
	resp, err := ca.Call("POST", url, params)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	schemaEnvelope := UserSchemaEnvelope{}
	if err := json.Unmarshal([]byte(resp), &schemaEnvelope); err != nil {
		return nil, err
	}
	schemaEnvelope.UserSchema.adjustDefaultTypes()

	return schemaEnvelope.UserSchema, nil
}

// [R]ead an existent user schema
func (ca *CustodiaAPIv1) ReadUserSchema(id string) (*UserSchema, error) {
	if !common.IsValidUUID(id) {
		return nil, errors.New("id is not a valid UUID: " + id)
	}

	url := fmt.Sprintf("/user_schemas/%s", id)
	resp, err := ca.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	schemaEnvelope := UserSchemaEnvelope{}
	if err := json.Unmarshal([]byte(resp), &schemaEnvelope); err != nil {
		return nil, err
	}
	schemaEnvelope.UserSchema.adjustDefaultTypes()

	return schemaEnvelope.UserSchema, nil
}

// [U]pdate an existent user schema
func (ca *CustodiaAPIv1) UpdateUserSchema(id string, description string,
	isActive bool, structure []SchemaField) (*UserSchema, error) {
	// isActive bool, structure json.RawMessage) (*UserSchema, error) {
		url := fmt.Sprintf("/user_schemas/%s", id)

		// UserSchema with just the data to send, so we can easily marshal it
		schema := UserSchema{
			Description: description,
			IsActive: isActive,
			Structure: structure,
		}
		params := map[string]interface{}{"data": schema}
		resp, err := ca.Call("PUT", url, params)
		if err != nil {
			return nil, err
		}

		// JSON: unmarshal resp content and return new user schema
		schemaEnvelope := UserSchemaEnvelope{}
		if err := json.Unmarshal([]byte(resp), &schemaEnvelope); err != nil {
			return nil, err
		}
		schemaEnvelope.UserSchema.adjustDefaultTypes()

		return schemaEnvelope.UserSchema, nil
}

// [D]elete an existent user schema
// if force=true the user schema is deleted, else it's just deactivated
func (ca *CustodiaAPIv1) DeleteUserSchema(id string, force bool) (
	error) {
	url := fmt.Sprintf("/user_schemas/%s", id)
	if force {
		url += "?force=true"
	}

	_, err := ca.Call("DELETE", url, nil)
	if err != nil {
		return err
	}
	return nil
}

// [L]ist all the user schemas
func (ca *CustodiaAPIv1) ListUserSchemas() ([]*UserSchema, error) {
	url := "/user_schemas"
	resp, err := ca.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	schemasEnvelope := UserSchemasEnvelope{}
	if err := json.Unmarshal([]byte(resp), &schemasEnvelope); err != nil {
		return nil, err
	}

	result := []*UserSchema{}
	for _, schema := range schemasEnvelope.UserSchemas {
		result = append(result, &schema)
		schema.adjustDefaultTypes()
	}

	return result, nil
}

// getStructureAsMap returns the list of fields in a map using the Name
// as key for quick access
func (us *UserSchema) getStructureAsMap() map[string]SchemaField {
	result := make(map[string]SchemaField)
	for _, field := range us.Structure {
		result[field.Name] = field
	}
	return result
}