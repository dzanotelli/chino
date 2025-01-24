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
