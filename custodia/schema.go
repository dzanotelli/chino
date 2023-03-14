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
	InsertDate timeutils.Time `json:"insert_date"`
	LastUpdate timeutils.Time `json:"last_update"`
	IsActive bool `json:"is_active"`
	Structure []SchemaField
}

type SchemaEnvelope struct {
	Schema *Schema `json:"schema"`
}

type SchemasEnvelope struct {
	Schemas []Schema `json:"schemas"`
}

// [C]reate a new schema
func (ca *CustodiaAPIv1) CreateSchema(r Repository, descritpion string, 
	fields []SchemaField, isActive bool) (*Schema, error) {
	schema := Schema{RepositoryId: r.RepositoryId, Description: descritpion,
		Structure: fields, IsActive: isActive}
	data, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}
	resp, err := ca.Call("POST", "/schemas", string(data))
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	schemaEnvelope := SchemaEnvelope{}
	if err := json.Unmarshal([]byte(resp), &schemaEnvelope); err != nil {
		return nil, err
	}

	return schemaEnvelope.Schema, nil
}

// [R]ead an existent schema
func (ca *CustodiaAPIv1) GetSchema(id string) (*Schema, error) {
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
	return schemaEnvelope.Schema, nil
}

// // [U]pdate an existent schema
// func (ca *CustodiaAPIv1) UpdateSchema(schema *Schema, description string,
// 	isActive bool, structure []SchemaField) {
// 		url := fmt.Sprintf("/schemas/%s", (*schema).SchemaId)

// 		// get a copy and update the values, so we can easily marshal it
// 		structure :=
// 		copy := *schema
// 		copy.Structure
// 	}