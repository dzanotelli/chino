package custodia

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dzanotelli/chino/common"
	"github.com/simplereach/timeutils"
)


type Document struct {
	Id string `json:"document_id,omitempty"`
	SchemaId string `json:"schema_id,omitempty"`
	RepositoryId string `json:"repository_id,omitempty"`
	InsertDate timeutils.Time `json:"insert_date,omitempty"`
	LastUpdate timeutils.Time `json:"last_update,omitempty"`
	IsActive bool `json:"is_active"`
	Content map[string]interface{} `json:"content,omitempty"`
}

type DocumentEnvelope struct {
	Document *Document `json:"document"`
}

type DocumentsEnvelope struct {
	Documents []Document `json:"documents"`
}

// [C]reate a new document
func (ca *CustodiaAPIv1) CreateDocument(schema *Schema, isActive bool,
	content map[string]interface{}) (*Document, error) {
	if schema.Id == "" {
		return nil, fmt.Errorf("schema has no SchemaId, does it exist?")
	} else if !common.IsValidUUID(schema.Id) {
		return nil, fmt.Errorf("SchemaId is not a valid UUID: %s (it " +
			"should not be manually set)", schema.Id)
	}

	// validate document content
	contentErrors := validateContent(content, schema.getStructureAsMap())
	if len(contentErrors) > 0 {
		e := fmt.Errorf("content errors: %w", errors.Join(contentErrors...))
		return nil, e
	}

	doc := Document{IsActive: isActive, Content: content}
	url := fmt.Sprintf("/schemas/%s/documents", schema.Id)
	params := map[string]interface{}{
		"data": doc,
	}
	resp, err := ca.Call("POST", url, params)
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
func (ca *CustodiaAPIv1) ReadDocument(schema Schema, id string) (*Document,
	 error) {
	if !common.IsValidUUID(id) {
		return nil, errors.New("id is not a valid UUID: " + id)
	}

	url := fmt.Sprintf("/documents/%s", id)
	resp, err := ca.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	docEnvelope := DocumentEnvelope{}
	if err := json.Unmarshal([]byte(resp), &docEnvelope); err != nil {
		return nil, err
	}

	// convert values to concrete types
	converted, ee := convertData(docEnvelope.Document.Content, &schema)
	if len(ee) > 0 {
		err := fmt.Errorf("conversion errors: %w", errors.Join(ee...))
		return docEnvelope.Document, err
	}

	// all good, assign the new content to doc and return it
	docEnvelope.Document.Content = converted
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
	params := map[string]interface{}{
		"data": doc,
	}
	resp, err := ca.Call("PUT", url, params)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content and return a fresh document instance
	docEnvelope := DocumentEnvelope{}
	if err := json.Unmarshal([]byte(resp), &docEnvelope); err != nil {
		return nil, err
	}

	// PUT call returns the whole documents, along with its content
	return docEnvelope.Document, nil
}

// [D]elete an existent document
// if force=false document is just deactivated
// if consisten=true the operation is done sync (server waits to respond)
func (ca *CustodiaAPIv1) DeleteDocument(id string, force, consistent bool) (
	error) {
	url := fmt.Sprintf("/documents/%s", id)
	url += fmt.Sprintf("?force=%v&consistent=%v", force, consistent)

	_, err := ca.Call("DELETE", url, nil)
	if err != nil {
		return err
	}
	return nil
}

// [L]ist all the documents in a Schema
// url params:
//   full_document: bool
//   is_active: bool
//   insert_date__gt: time.Time
//   insert_date__lt: time.Time
//   last_update__gt: time.Time
//   last_update__lt: time.Time
func (ca *CustodiaAPIv1) ListDocuments(schemaId string,
	params map[string]interface{}) ([]*Document, error) {
	if !common.IsValidUUID(schemaId) {
		return nil, fmt.Errorf("schemaId is not a valid UUID: %s", schemaId)
	}

	url := fmt.Sprintf("/schemas/%s/documents", schemaId)
	if len(params) > 0 {
		url += "?"
	}

	availableParams := map[string]interface{}{
		"full_document": true,
		"is_active": true,
		"insert_date__gt": time.Time{},
		"insert_date__lt": time.Time{},
		"last_update__gt": time.Time{},
		"last_update__lt": time.Time{},
	}
	for param := range params {
		// check that param is legit
		_, ok := availableParams[param]
		if !ok {
			return nil, fmt.Errorf("got unexpected param '%s'", param)
		}
		value := params[param]

		switch param {
		case "full_document", "is_active":
			_, ok := value.(bool)
			if !ok {
				err := fmt.Errorf("param '%s': bad type: '%T', must be bool",
					param, value)
				return nil, err
			}
		case "insert_date__gt", "insert_date__lt", "last_update__gt",
			"last_update__lt":
			time_value, ok := value.(time.Time)
			if !ok {
				err := fmt.Errorf("param '%s': bad type: '%T', must be Time",
					param, value)
				return nil, err
			}
			value = time_value.Format(time.RFC3339)
		default:
			return nil, fmt.Errorf("got unexpected param '%s'", param)
		}

		url += fmt.Sprintf("%s=%v&", param, value)
	}

	url = strings.TrimRight(url, "&")
	resp, err := ca.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	docusEnvelope := DocumentsEnvelope{}
	if err := json.Unmarshal([]byte(resp), &docusEnvelope); err != nil {
		return nil, err
	}

	result := []*Document{}
	for _, doc := range docusEnvelope.Documents {
		result = append(result, &doc)
	}

	return result, nil
}
