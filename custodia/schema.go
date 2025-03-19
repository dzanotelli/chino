package custodia

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/simplereach/timeutils"
)

// SchemaField is used by both Schema and UserSchema
type SchemaField struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Indexed bool `json:"indexed,omitempty"`
	Default interface{} `json:"default,omitempty"`
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
	case "integer":
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
	params := map[string]interface{}{"_data": schema}
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
		params := map[string]interface{}{"_data": schema}
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

// [L]ist all the schemas in a repository
func (ca *CustodiaAPIv1) ListSchemas(repoId uuid.UUID) ([]*Schema, error) {
	url := fmt.Sprintf("/repositories/%s/schemas", repoId)
	resp, err := ca.Call("GET", url, nil)
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