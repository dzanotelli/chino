package custodia

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/simplereach/timeutils"
)

type Collection struct {
	Id   uuid.UUID `json:"collection_id"`
	Name string `json:"name"`
	InsertDate timeutils.Time `json:"insert_date"`
	LastUpdate timeutils.Time `json:"last_update"`
	IsActive bool `json:"is_active"`
}

type CollectionEnvelope struct {
	Collections []Collection `json:"collections"`
}

func (c Collection) String() string {
	return fmt.Sprintf("<Collection %s %s>", c.Name, c.Id)
}

// [C]reate a new collection
func (ca *CustodiaAPIv1) CreateCollection(name string) (*Collection, error) {
	url := "/collections"
	data := map[string]any{"name": name}
	params := map[string]any{"_data": data}
	resp, err := ca.Call("POST", url, params)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	collection := Collection{}
	if err := json.Unmarshal([]byte(resp), &collection); err != nil {
		return nil, err
	}

	return &collection, nil
}

// [R]ead an existent collection
func (ca *CustodiaAPIv1) ReadCollection(collectionId uuid.UUID) (*Collection,
	error) {
	url := fmt.Sprintf("/collections/%s", collectionId)
	resp, err := ca.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	collection := Collection{}
	if err := json.Unmarshal([]byte(resp), &collection); err != nil {
		return nil, err
	}

	return &collection, nil
}

// [U]pdate an existent collection
func (ca *CustodiaAPIv1) UpdateCollection(collectionId uuid.UUID,
	name string) (
	*Collection, error) {
	url := fmt.Sprintf("/collections/%s", collectionId)
	data := map[string]any{"name": name}
	params := map[string]any{"_data": data}
	resp, err := ca.Call("PUT", url, params)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	collection := Collection{}
	if err := json.Unmarshal([]byte(resp), &collection); err != nil {
		return nil, err
	}

	return &collection, nil
}

// [D]elete an existent collection
func (ca *CustodiaAPIv1) DeleteCollection(collectionId uuid.UUID, force bool) (
	error) {
	url := fmt.Sprintf("/collections/%s", collectionId)
	if force {
		url += "?force=true"
	}
	_, err := ca.Call("DELETE", url, nil)
	if err != nil {
		return err
	}
	return nil
}

// [L]ist collections
func (ca *CustodiaAPIv1) ListCollections() ([]*Collection, error) {
	url := "/collections"
	resp, err := ca.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	collectionEnvelope := CollectionEnvelope{}
	if err := json.Unmarshal([]byte(resp), &collectionEnvelope); err != nil {
		return nil, err
	}

	result := []*Collection{}
	for _, collection := range collectionEnvelope.Collections {
		result = append(result, &collection)
	}

	return result, nil
}

// List the collections of a document
func (ca *CustodiaAPIv1) ListDocumentCollections(documentId uuid.UUID) (
	[]*Collection, error) {
	url := fmt.Sprintf("/collections/documents/%s", documentId)
	resp, err := ca.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	collectionEnvelope := CollectionEnvelope{}
	if err := json.Unmarshal([]byte(resp), &collectionEnvelope); err != nil {
		return nil, err
	}

	result := []*Collection{}
	for _, collection := range collectionEnvelope.Collections {
		result = append(result, &collection)
	}
	return result, nil
}

// List the documents of a collection
func (ca *CustodiaAPIv1) ListCollectionDocuments(collectionId uuid.UUID,
	fullDocument bool) ([]*Document, error) {
	url := fmt.Sprintf("/collections/%s/documents", collectionId)
	if fullDocument {
		url += "?full_document=true"
	}

	resp, err := ca.Call("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	documentsEnvelope := DocumentsEnvelope{}
	if err := json.Unmarshal([]byte(resp), &documentsEnvelope); err != nil {
		return nil, err
	}

	// FIXME: the underlying type of interfaces is not the expected concrete
	//   type defined in the Schema structure, it's just the result of json
	//   unmarshalling. Check document/ReadDocument comments for more.
	result := []*Document{}
	for _, document := range documentsEnvelope.Documents {
		result = append(result, &document)
	}
	return result, nil
}

// Add a document to a collection
func (ca *CustodiaAPIv1) AddDocumentToCollection(documentId uuid.UUID,
	collectionId uuid.UUID) error {
	url := fmt.Sprintf("/collections/%s/documents/%s", collectionId,
		documentId)
	_, err := ca.Call("POST", url, nil)
	if err != nil {
		return err
	}
	return nil
}

// Remove a document from a collection
func (ca *CustodiaAPIv1) RemoveDocumentFromCollection(documentId uuid.UUID,
	collectionId uuid.UUID) error {
	url := fmt.Sprintf("/collections/%s/documents/%s", collectionId,
		documentId)
	_, err := ca.Call("DELETE", url, nil)
	if err != nil {
		return err
	}
	return nil
}

// Search a collection
func (ca *CustodiaAPIv1) SearchCollection(name string, contains bool) (
	[]*Collection, error) {
	url := "/collections/search"
	data := map[string]any{"name": name, "contains": contains}
	params := map[string]any{"_data": data}

	resp, err := ca.Call("POST", url, params)
	if err != nil {

	}

	// JSON: unmarshal resp content
	collectionEnvelope := CollectionEnvelope{}
	if err := json.Unmarshal([]byte(resp), &collectionEnvelope); err != nil {
		return nil, err
	}

	result := []*Collection{}
	for _, collection := range collectionEnvelope.Collections {
		result = append(result, &collection)
	}

	return result, nil
}
