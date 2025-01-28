package custodia

import (
	"encoding/json"
	"fmt"
	"time"
)


type UploadBlob struct {
	Id string `json:"upload_id"`
	ExpireDate time.Time `json:"expire_date"`

}

type Blob struct {
	Id string `json:"blob_id"`
	DocumentId string `json:"document_id"`
	Sha1 string `json:"sha1"`
	Md5 string `json:"md5"`
}

type BlobEnvelope struct {
	Blob map[string]interface{} `json:"blob"`
}

// Upload a new blob
// This is the starting point to begin to upload a new blob
// It returns an UploadBlob object which is used later to upload
// data to the server
func (ca *CustodiaAPIv1) CreateBlob(documentId string, fieldName string,
	fileName string) (*UploadBlob, error) {


	data := map[string]interface{}{"document_id": documentId,
		"field": fieldName, "file_name": fileName}
	resp, err := ca.Call("POST", "/blobs/" , data)
	if err != nil {
		return nil, err
	}
	blobEnvelope := BlobEnvelope{}
	if err := json.Unmarshal([]byte(resp), &blobEnvelope); err != nil {
		return nil, err
	}

	ub := &UploadBlob{
		Id: blobEnvelope.Blob["blob_id"].(string),
		ExpireDate: blobEnvelope.Blob["expire_date"].(time.Time),
	}

	return ub, nil
}

// Upload a chunk
// Args:
// - ub: the UploadBlob returned by CreateBlob
// - data: the chunk of data
// - length: the total length of the file. Must be passed as header
// - offset: the offset of this chunk in the file. Must be passed as header
func (ca *CustodiaAPIv1) UploadChunk(ub *UploadBlob, data []byte,
	length int, offset int) (*UploadBlob, error) {
	url := fmt.Sprintf("/blobs/%s", ub.Id)
	params := map[string]interface{}{
		"Content-Type": "application/octet-stream",
		"Length": fmt.Sprint(length),
		"Offset": fmt.Sprint(offset),
		"body": data,
	}

	resp, err := ca.Call("PUT", url, params)
	if err != nil {
		return nil, err
	}

	blobEnvelope := BlobEnvelope{}
	if err := json.Unmarshal([]byte(resp), &blobEnvelope); err != nil {
		return nil, err
	}

	ub = &UploadBlob{
		Id: blobEnvelope.Blob["upload_id"].(string),
		ExpireDate: blobEnvelope.Blob["expire_date"].(time.Time),
	}

	return ub, nil
}

// Commit a blob
func (ca *CustodiaAPIv1) CommitBlob(ub *UploadBlob) (*Blob, error) {
	url := "/blobs/commit"
	data := map[string]interface{}{"upload_id": ub.Id}
	resp, err := ca.Call("POST", url, data)
	if err != nil {
		return nil, err
	}
	blobEnvelope := BlobEnvelope{}
	if err := json.Unmarshal([]byte(resp), &blobEnvelope); err != nil {
		return nil, err
	}
	blob := &Blob{
		Id: blobEnvelope.Blob["blob_id"].(string),
		DocumentId: blobEnvelope.Blob["document_id"].(string),
		Sha1: blobEnvelope.Blob["sha1"].(string),
		Md5: blobEnvelope.Blob["md5"].(string),
	}
	return blob, nil
}

// Download a blob
func (ca *CustodiaAPIv1) GetBlobData(blobId string) ([]byte, error) {
	url := fmt.Sprintf("/blobs/%s", blobId)
	params := map[string]interface{}{"raw_response": true}
	resp, err := ca.Call("GET", url, params)
	if err != nil {
		return nil, err
	}

	return []byte(resp), nil
}

// Delete a blob
func (ca *CustodiaAPIv1) DeleteBlob(blobId string) error {
	url := fmt.Sprintf("/blobs/%s", blobId)
	_, err := ca.Call("DELETE", url, nil)
	return err
}


// Generate a blob url
func (ca *CustodiaAPIv1) GetBlobUrl(blobId string, oneTime bool,
	duration int) (map[string]string, error) {
	url := fmt.Sprintf("/blobs/%s/generate", blobId)
	data := map[string]interface{}{"one_time": oneTime, "duration": duration}
	resp, err := ca.Call("POST", url, data)
	if err != nil {
		return nil, err
	}

	// response returns a map with toekn, expiration, and one_time
	// since one_time is a bool, but we already know it (it's a func param)
	// we return just the map with two strings
	respData := map[string]interface{}{}
	if err := json.Unmarshal([]byte(resp), &respData); err != nil {
		return nil, err
	}

	out := map[string]string{
		"token": respData["token"].(string),
		"expiration": respData["expiration"].(string),
	}

	return out, nil
}