package custodia

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dzanotelli/chino/common"
	"github.com/simplereach/timeutils"
)


type Document struct {
	DocumentId string `json:"document_id,omitempty"`
	SchemaId string `json:"schema_id,omitempty"`
	RepositoryId string `json:"repository_id,omitempty"`
	InsertDate timeutils.Time `json:"insert_date,omitempty"`
	LastUpdate timeutils.Time `json:"last_update,omitempty"`
	IsActive bool `json:"is_active"`
	Content  map[string]interface{} `json:"content,omitempty"`
}

type DocumentEnvelope struct {
	Document *Document `json:"document"`
}

type DocumentsEnvelope struct {
	Document []Document `json:"document"`
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
	data, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("/schemas/%s/documents", schema.SchemaId)
	resp, err := ca.Call("POST", url, string(data))
	if err != nil {
		return nil, err
	}
	// JSON: unmarshal resp content
	docEnvelope := DocumentEnvelope{}
	if err := json.Unmarshal([]byte(resp), &docEnvelope); err != nil {
		return nil, err
	}
	
	// if everything is ok, we can safely set the given content as the
	// returned document content, since the API doesn't return it
	docEnvelope.Document.Content = content

	return docEnvelope.Document, nil
}

// [R]ead an existent document
func (ca *CustodiaAPIv1) ReadDocument(id string) (*Document, error) {
	if !common.IsValidUUID(id) {
		return nil, errors.New("id is not a valid UUID: " + id)
	}

	url := fmt.Sprintf("/documents/%s", id)
	resp, err := ca.Call("GET", url)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	docEnvelope := DocumentEnvelope{}
	if err := json.Unmarshal([]byte(resp), &docEnvelope); err != nil {
		return nil, err
	}

	return docEnvelope.Document, nil
}

// [U]pdate an existent document
func (ca *CustodiaAPIv1) UpdateDocument(id string , isActive bool, 
	content map[string]interface{}) (*Document, error) {
		if !common.IsValidUUID(id) {
			return nil, errors.New("id is not a valid UUID: " + id)
		}	
	
	url := fmt.Sprintf("/documents/%s", id)

	// create a doc with just the values we can send, and marshal it
	doc := Document{IsActive: isActive, Content: content}
	data, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	resp, err := ca.Call("PUT", url, string(data))
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content overwriting the old document
	docEnvelope := DocumentEnvelope{}
	if err := json.Unmarshal([]byte(resp), &docEnvelope); err != nil {
		return nil, err
	}

	// PUT call returns the whole documents, along with its content
	return docEnvelope.Document, nil
}