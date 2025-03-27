package custodia

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/simplereach/timeutils"
)


type Document struct {
	Id uuid.UUID `json:"document_id,omitempty"`
	SchemaId uuid.UUID `json:"schema_id,omitempty"`
	RepositoryId uuid.UUID `json:"repository_id,omitempty"`
	InsertDate timeutils.Time `json:"insert_date,omitempty"`
	LastUpdate timeutils.Time `json:"last_update,omitempty"`
	IsActive bool `json:"is_active"`
	Content map[string]any `json:"content,omitempty"`
}

type DocumentEnvelope struct {
	Document *Document `json:"document"`
}

type DocumentsEnvelope struct {
	Documents []Document `json:"documents"`
}

// [C]reate a new document
func (ca *CustodiaAPIv1) CreateDocument(schema *Schema, isActive bool,
	content map[string]any) (*Document, error) {
	// validate document content
	contentErrors := validateContent(content, schema.getStructureAsMap())
	if len(contentErrors) > 0 {
		e := fmt.Errorf("content errors: %w", errors.Join(contentErrors...))
		return nil, e
	}

	doc := Document{IsActive: isActive, Content: content}
	url := fmt.Sprintf("/schemas/%s/documents", schema.Id)
	params := map[string]any{
		"_data": doc,
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

	// convert values to concrete types
	converted, ee := convertData(docEnvelope.Document.Content, schema)
	if len(ee) > 0 {
		err := fmt.Errorf("conversion errors: %w", errors.Join(ee...))
		return docEnvelope.Document, err
	}

	// all good, assign the new content to doc and return it
	docEnvelope.Document.Content = converted

	return docEnvelope.Document, nil
}

// [R]ead an existent document
func (ca *CustodiaAPIv1) ReadDocument(schema Schema, documentId uuid.UUID) (
	*Document, error) {
	url := fmt.Sprintf("/documents/%s", documentId)
	resp, err := ca.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	docEnvelope := DocumentEnvelope{}
	if err := json.Unmarshal([]byte(resp), &docEnvelope); err != nil {
		return nil, err
	}

	// FIXME: probably better to:
	//   1. remove this convertion here, and let Content be raw interfaces
	//      as returned by Unmarshal
	//   2. add a `GetContent` method to Document, which converts the
	//      underlying type of interfaces to the expected concrete types
	//      (e.g. strings to time.Time, or float to int)

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
func (ca *CustodiaAPIv1) UpdateDocument(schema Schema, documentId uuid.UUID,
	isActive bool, content map[string]any) (*Document, error) {
	url := fmt.Sprintf("/documents/%s", documentId)

	// create a doc with just the values we can send, and marshal it
	doc := Document{IsActive: isActive, Content: content}
	params := map[string]any{
		"_data": doc,
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

	// convert values to concrete types
	converted, ee := convertData(docEnvelope.Document.Content, &schema)
	if len(ee) > 0 {
		err := fmt.Errorf("conversion errors: %w", errors.Join(ee...))
		return docEnvelope.Document, err
	}

	// PUT call returns the whole documents, along with its content
	docEnvelope.Document.Content = converted
	return docEnvelope.Document, nil
}

// [D]elete an existent document
// if force=false document is just deactivated
// if consisten=true the operation is done sync (server waits to respond)
func (ca *CustodiaAPIv1) DeleteDocument(documentId uuid.UUID, force,
	consistent bool) (error) {
	url := fmt.Sprintf("/documents/%s", documentId)
	url += fmt.Sprintf("?force=%v&consistent=%v", force, consistent)

	_, err := ca.Call("DELETE", url, nil)
	if err != nil {
		return err
	}
	return nil
}

// [L]ist all the documents in a Schema
// url queryParams:
//   offset: int: number of items to skip from the beginning of the list
//   limit: int : maximum number of items to return in a single page
//   full_document: bool: return the full document
//   is_active: bool: filter by
//   insert_date__gt: time string (RFC3339, YYYY-MM-DDTHH:MM:SS): filter by
//   insert_date__lt: time string (RFC3339): filter by
//   last_update__gt: time string (RFC3339): filter by
//   last_update__lt: time string (RFC3339): filter by
func (ca *CustodiaAPIv1) ListDocuments(schema Schema,
	queryParams map[string]string) ([]*Document, error) {
	u, err := url.Parse(fmt.Sprintf("/schemas/%s/documents", schema.Id))
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %v", err)
	}

	// Adding query params
	q := u.Query()
	for k, v := range(queryParams) {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	resp, err := ca.Call("GET", u.String(), nil)
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
		converted, ee := convertData(doc.Content, &schema)
		if len(ee) > 0 {
			err := fmt.Errorf("conversion errors: %w", errors.Join(ee...))
			return nil, err
		}
		doc.Content = converted
		result = append(result, &doc)
	}

	return result, nil
}
