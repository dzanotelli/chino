package custodia

import (
	"encoding/json"

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

type ScehamEnvelope struct {
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
	schemaEnvelope := ScehamEnvelope{}
	if err := json.Unmarshal([]byte(resp), &schemaEnvelope); err != nil {
		return nil, err
	}

	return schemaEnvelope.Schema, nil
}