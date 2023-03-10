package custodia

import (
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