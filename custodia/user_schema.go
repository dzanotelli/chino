package custodia

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/simplereach/timeutils"
)

type UserSchema struct {
	Id uuid.UUID `json:"user_schema_id,omitempty"`
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
	params := map[string]any{"_data": user_schema}
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
func (ca *CustodiaAPIv1) ReadUserSchema(userSchemaId uuid.UUID) (*UserSchema,
	error) {
	url := fmt.Sprintf("/user_schemas/%s", userSchemaId)
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
func (ca *CustodiaAPIv1) UpdateUserSchema(userSchemaId uuid.UUID,
	description string, isActive bool, structure []SchemaField) (*UserSchema,
	error) {
	// isActive bool, structure json.RawMessage) (*UserSchema, error) {
	url := fmt.Sprintf("/user_schemas/%s", userSchemaId)

	// UserSchema with just the data to send, so we can easily marshal it
	schema := UserSchema{
		Description: description,
		IsActive: isActive,
		Structure: structure,
	}
	params := map[string]any{"_data": schema}
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
func (ca *CustodiaAPIv1) DeleteUserSchema(userSchemaId uuid.UUID, force bool) (
	error) {
	url := fmt.Sprintf("/user_schemas/%s", userSchemaId)
	if force {
		url += "?force=true"
	}

	_, err := ca.Call("DELETE", url, nil)
	if err != nil {
		return err
	}
	return nil
}

// [L]ist all the schemas in a UserSchema
// queryParams (optional):
//   offset: int: number of items to skip from the beginning of the list
//   limit: int : maximum number of items to return in a single page
func (ca *CustodiaAPIv1) ListUserSchemas(queryParams map[string]string) (
	[]*UserSchema, error,
) {
	u, err := url.Parse("/user_schemas")
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
