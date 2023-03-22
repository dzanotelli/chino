package custodia

import (
	"fmt"

	"github.com/dzanotelli/chino/common"
	"github.com/simplereach/timeutils"
)


type Document struct {
	DocumentId string `json:"document_id,omitempty"`
	SchemaId string `json:"schema_id,omitempty"`
	RepositoryId string `json:"repository_id,omitempty"`
	InsertDate timeutils.Time `json:"insert_date"`
	LastUpdate timeutils.Time `json:"last_update"`
	IsActive bool `json:"is_active"`
	Content  map[string]interface{} `json:"content"`
}

type DocumentEnvelope struct {
	Document *Document `json:"document"`
}

type DocumentsEnvelope struct {
	Document []Document `json:"documents"`
}

// [C]reate a new document
func (ca *CustodiaAPIv1) CreateDocument(schema *Schema, isActive bool, 
	content map[string]interface{}) (*Document, error) {
	if schema.SchemaId == "" {
		return nil, fmt.Errorf("schema has no SchemaId, does it exist?")
	} else if !common.IsValidUUID(schema.SchemaId) {
		return nil, fmt.Errorf("SchemaId is not a valid UUID: %s (it " +
			"should not be manually set)", schema.SchemaId)
	}

	// validate document content
	mappedStructure := schema.getStructureAsMap()
	validatedContent := make(map[string]interface{})
	for key, value := range content {
		field, ok := mappedStructure[key]
		if !ok {
			return nil, fmt.Errorf("given field '%s' not defined in " +
				"schema structure", key)
		}

		// field exist, check that is of the right type
		var val interface{}
		var err error
		switch field.Type {
		case "integer":
			val, ok = value.(int)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be int", key)
			}
		case "float":
			val, ok = value.(float64)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be float", key)
			}
		case "string", "text":
			val, ok = value.(string)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be string", key)
				break
			}
			converted := fmt.Sprintf("%v", val)
			if field.Type == "string" && len(converted) > 255 {
				ok = false
				err = fmt.Errorf("field '%s' exceeded max lenght of 255 chars", 
					key)
			}
		case "boolean":
			val, ok = value.(bool)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be bool", key)
			}
		// FIXME: missing other types
		}
		
		
		if !ok {
			return nil, err
		}
		validatedContent[key] = val

	}


	// doc := Document{IsActive: isActive, Content: content}


	return nil, nil
}