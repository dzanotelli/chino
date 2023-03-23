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

	errors := validateContent(content, schema.getStructureAsMap())
	if len(errors) > 0 {
		err := fmt.Errorf("content validation failed: ")
		for _, e := range errors {
			err = fmt.Errorf("%w %w", err, e)
		}

		return nil, err
	}

	doc := Document{IsActive: isActive, Content: content}
	return &doc, nil
}
