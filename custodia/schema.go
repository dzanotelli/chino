package custodia

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/simplereach/timeutils"
	"golang.org/x/exp/slices"
)

// Define `type` enum
type FieldType int

const (
	TypeInteger FieldType = iota + 1
	TypeArryInteger
	TypeFloat
	TypeArrayFloat
	TypeString
	TypeArrayString
	TypeText FieldType
	TypeBoolean FieldType
	TypeDate FieldType
	TypeTime FieldType
	TypeDateTime FieldType
	TypeBase64 FieldType
	TypeJson FieldType
	TypeBlob FieldType
)

func (ft FieldType) Choices() []string {
	return []string{
		"integer",
		"array[integer]",
		"float",
		"array[float]",
		"string",
		"array[string]",
		"text",
		"boolean",
		"date",
		"time",
		"datetime",
		"base64",
		"json",
		"blob",
	}
}

func (ft FieldType) String() string {
	return ft.Choices()[ft-1]
}

func (ft FieldType) MarshalJSON() ([]byte, error) {
	return json.Marshal(ft.String())
}

func (ft *FieldType) UnmarshalJSON(data []byte) (err error) {
	var str string
	if err = json.Unmarshal(data, &str); err != nil {
		return err
	}

	intValue := slices.Index(ft.Choices(), str) + 1   // enum starts from 1
	if intValue < 1 {
		return fmt.Errorf("FieldType: received unknown value '%v'", str)
	}

	*ft = FieldType(intValue)
	return nil
}


// SchemaField is used by both Schema and UserSchema
type SchemaField struct {
	Name string `json:"name"`
	Type FieldType `json:"type"`
	Indexed bool `json:"indexed,omitempty"`
	Default any `json:"default,omitempty"`
	Insensitive bool `json:"insensitive,omitempty"`
}

type Schema struct {
	Id uuid.UUID `json:"schema_id,omitempty"`
	RepositoryId uuid.UUID `json:"repository_id,omitempty"`
	Description string `json:"description"`
	InsertDate timeutils.Time `json:"insert_date,omitempty"`
	LastUpdate timeutils.Time `json:"last_update,omitempty"`
	IsActive bool `json:"is_active"`
	Structure []SchemaField `json:"structure"`
}

type SchemaEnvelope struct {
	Schema *Schema `json:"schema"`
}

type SchemasEnvelope struct {
	Schemas []Schema `json:"schemas"`
}

// adjustDefaultType fixes the automatic interface-to-type conversion done
// by json.Unmarshal to the desired type (e.g. json int values are
// automatically decoded to float64 and we want int instead)
func (f *SchemaField) adjustDefaultType() {
	if f.Default == nil {
		return
	}

	switch f.Type {
	case TypeInteger:
		floatVal, _ := f.Default.(float64)
		f.Default = int(floatVal)
	}
}

// adjustDefaultTypes for each field calls adjustDefaultType
func (s *Schema) adjustDefaultTypes() {
	for i := range s.Structure {
		s.Structure[i].adjustDefaultType()
	}
}

// [C]reate a new schema
func (ca *CustodiaAPIv1) CreateSchema(repoId uuid.UUID, descritpion string,
	isActive bool, fields []SchemaField) (*Schema, error) {
	// FIXME: missing field type validation, and indexed property validation
	//   and insensitive property

	schema := Schema{RepositoryId: repoId, Description: descritpion,
		Structure: fields, IsActive: isActive}
	url := fmt.Sprintf("/repositories/%s/schemas", repoId)
	params := map[string]any{"_data": schema}
	resp, err := ca.Call("POST", url, params)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	schemaEnvelope := SchemaEnvelope{}
	if err := json.Unmarshal([]byte(resp), &schemaEnvelope); err != nil {
		return nil, err
	}
	schemaEnvelope.Schema.adjustDefaultTypes()

	return schemaEnvelope.Schema, nil
}

// [R]ead an existent schema
func (ca *CustodiaAPIv1) ReadSchema(schemaId uuid.UUID) (*Schema, error) {
	url := fmt.Sprintf("/schemas/%s", schemaId)
	resp, err := ca.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	schemaEnvelope := SchemaEnvelope{}
	if err := json.Unmarshal([]byte(resp), &schemaEnvelope); err != nil {
		return nil, err
	}
	schemaEnvelope.Schema.adjustDefaultTypes()

	return schemaEnvelope.Schema, nil
}

// [U]pdate an existent schema
func (ca *CustodiaAPIv1) UpdateSchema(schemaId uuid.UUID, description string,
	isActive bool, structure []SchemaField) (*Schema, error) {
	// isActive bool, structure json.RawMessage) (*Schema, error) {
		url := fmt.Sprintf("/schemas/%s", schemaId)

		// Schema with just the data to send, so we can easily marshal it
		schema := Schema{
			Description: description,
			IsActive: isActive,
			Structure: structure,
		}
		params := map[string]any{"_data": schema}
		resp, err := ca.Call("PUT", url, params)
		if err != nil {
			return nil, err
		}

		// JSON: unmarshal resp content and return new schema
		schemaEnvelope := SchemaEnvelope{}
		if err := json.Unmarshal([]byte(resp), &schemaEnvelope); err != nil {
			return nil, err
		}
		schemaEnvelope.Schema.adjustDefaultTypes()

		return schemaEnvelope.Schema, nil
}

// [D]elete an existent schema
// if force=true the schema is deleted, else it's just deactivated
// if all_content=true the schema content is deleted too (it also sets
// automatically force=true)
func (ca *CustodiaAPIv1) DeleteSchema(schemaId uuid.UUID, force bool,
	allContent bool) (error) {
	url := fmt.Sprintf("/schemas/%s", schemaId)

	// allContent requires force=true
	if allContent {
		url += "?force=true&all_content=true"
	} else if force {
		url += "?force=true"
	}

	_, err := ca.Call("DELETE", url, nil)
	if err != nil {
		return err
	}
	return nil
}

// [L]ist all the schemas in a Schema
// queryParams (optional):
//   offset: int: number of items to skip from the beginning of the list
//   limit: int : maximum number of items to return in a single page
func (ca *CustodiaAPIv1) ListSchemas(repoId uuid.UUID,
	queryParams map[string]string) ([]*Schema, error,
) {
	u, err := url.Parse(fmt.Sprintf("/repositories/%s/schemas", repoId))
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
	schemasEnvelope := SchemasEnvelope{}
	if err := json.Unmarshal([]byte(resp), &schemasEnvelope); err != nil {
		return nil, err
	}

	result := []*Schema{}
	for _, schema := range schemasEnvelope.Schemas {
		result = append(result, &schema)
		schema.adjustDefaultTypes()
	}

	return result, nil
}

// getStructureAsMap returns the list of fields in a map using the Name
// as key for quick access
func (s *Schema) getStructureAsMap() map[string]SchemaField {
	result := make(map[string]SchemaField)
	for _, field := range s.Structure {
		result[field.Name] = field
	}
	return result
}
