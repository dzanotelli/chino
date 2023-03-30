package custodia

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dzanotelli/chino/common"
	"github.com/simplereach/timeutils"
)

type SchemaField struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Indexed bool `json:"bool,omitempty"`
	Default interface{} `json:"default,omitempty"`
}

type Schema struct {
	SchemaId string `json:"schema_id,omitempty"`
	RepositoryId string `json:"repository_id,omitempty"`
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
// automatically converted to float and we want int instead)
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
func (ca *CustodiaAPIv1) CreateSchema(repository *Repository, descritpion string, 
	isActive bool, fields []SchemaField) (*Schema, error) {
	if repository.RepositoryId == "" {
		return nil, fmt.Errorf("repository has no RepositoryId, " + 
			"does it exist?")
	} else if !common.IsValidUUID(repository.RepositoryId) {
		return nil, fmt.Errorf("RepositoryId is not a valid UUID: %s (it " +
			"should not be manually set)", repository.RepositoryId)
	}

	// FIXME: missing field type validation, and indexed property validation
	//   and insensitive property
	
	schema := Schema{RepositoryId: repository.RepositoryId, 
		Description: descritpion, Structure: fields, IsActive: isActive}
	data, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("/repositories/%s/schemas", repository.RepositoryId)
	resp, err := ca.Call("POST", url, string(data))
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
func (ca *CustodiaAPIv1) ReadSchema(id string) (*Schema, error) {
	if !common.IsValidUUID(id) {
		return nil, errors.New("id is not a valid UUID: " + id)
	}

	url := fmt.Sprintf("/schemas/%s", id)
	resp, err := ca.Call("GET", url)
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
func (ca *CustodiaAPIv1) UpdateSchema(id string, description string,
	isActive bool, structure []SchemaField) (*Schema, error) {
		url := fmt.Sprintf("/schemas/%s", id)

		// Schema with just the data to send, so we can easily marshal it
		schema := Schema{
			Description: description, 
			IsActive: isActive, 
			Structure: structure,
		}
		data, err := json.Marshal(schema)
		if err != nil {
			return nil, err
		}
		resp, err := ca.Call("PUT", url, string(data))
		if err != nil {
			return nil, err
		}

		// JSON: unmarshal resp content overwriting the old schema
		schemaEnvelope := SchemaEnvelope{}
		if err := json.Unmarshal([]byte(resp), &schemaEnvelope); err != nil {
			return nil, err
		}
		schemaEnvelope.Schema.adjustDefaultTypes()

		return schemaEnvelope.Schema, nil
}

// [D]elete and existent schema
// if force=true the schema is deleted, else it's just deactivated
// if all_content=true the schema content is deleted too (it also sets 
// automatically force=true)
func (ca *CustodiaAPIv1) DeleteSchema(id string, force, allContent bool) (
	error) {
	url := fmt.Sprintf("/schemas/%s", id)

	// allContent requires force=true
	if allContent {
		url += "?force=true&all_content=true"
	} else if force {
		url += "?force=true"
	}

	_, err := ca.Call("DELETE", url)
	if err != nil {
		return err
	}
	return nil
}

// [L]ist all the schemas in a repository
func (ca *CustodiaAPIv1) ListSchemas(repositoryId string) ([]*Schema, error) {
	if !common.IsValidUUID(repositoryId) {
		err := fmt.Errorf("repositoryId is not a valid UUID: %v", repositoryId)
		return nil, err
	}
	url := fmt.Sprintf("/repositories/%s/schemas", repositoryId)
	resp, err := ca.Call("GET", url)
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