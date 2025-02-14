package custodia

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dzanotelli/chino/common"
	"github.com/google/uuid"
	"github.com/simplereach/timeutils"
)


type UploadBlob struct {
	Id uuid.UUID `json:"upload_id"`
	ExpireDate timeutils.Time `json:"expire_date"`
}

type UploadBlobEnvelope struct {
	UploadBlob UploadBlob `json:"blob"`
}

type Blob struct {
	Id uuid.UUID `json:"blob_id"`
	DocumentId string `json:"document_id"`
	Sha1 string `json:"sha1"`
	Md5 string `json:"md5"`
}

type BlobToken struct {
	Token string `json:"token"`
	Expiration timeutils.Time `json:"expiration"`
	OneTime bool `json:"one_time"`
}

type BlobEnvelope struct {
	Blob Blob `json:"blob"`
}

// Upload a new blob
// This is the starting point to begin to upload a new blob
// It returns an UploadBlob object which is used later to upload
// data to the server
func (ca *CustodiaAPIv1) CreateBlob(documentId uuid.UUID, fieldName string,
	fileName string) (*UploadBlob, error) {
	data := map[string]interface{}{"document_id": documentId.String(),
		"field": fieldName, "file_name": fileName}
	params := map[string]interface{}{"_data": data}
	resp, err := ca.Call("POST", "/blobs" , params)
	if err != nil {
		return nil, err
	}
	blobEnvelope := BlobEnvelope{}
	if err := json.Unmarshal([]byte(resp), &blobEnvelope); err != nil {
		return nil, err
	}

	uploadBlobEnvelope := &UploadBlobEnvelope{}
	if err := json.Unmarshal([]byte(resp), &uploadBlobEnvelope); err != nil {
		return nil, err
	}

	return &uploadBlobEnvelope.UploadBlob, nil
}

// Upload a chunk
// Args:
// - uploadId: the ID of the UploadBlob returned by CreateBlob
// - data: the chunk of data to upload
// - length: the total length of the file. Must be passed as header
// - offset: the offset of this chunk in the file. Must be passed as header
func (ca *CustodiaAPIv1) UploadChunk(uploadId uuid.UUID, data []byte,
	length int, offset int) (*UploadBlob, error) {
	url := fmt.Sprintf("/blobs/%s", uploadId)
	params := map[string]interface{}{
		"Content-Type": "application/octet-stream",
		"Length": fmt.Sprint(length),
		"Offset": fmt.Sprint(offset),
		"_data": data,
	}

	resp, err := ca.Call("PUT", url, params)
	if err != nil {
		return nil, err
	}

	blobEnvelope := BlobEnvelope{}
	if err := json.Unmarshal([]byte(resp), &blobEnvelope); err != nil {
		return nil, err
	}

	uploadBlobEnvelope := &UploadBlobEnvelope{}
	if err := json.Unmarshal([]byte(resp), &uploadBlobEnvelope); err != nil {
		return nil, err
	}

	return &uploadBlobEnvelope.UploadBlob, nil
}

// Commit a blob
func (ca *CustodiaAPIv1) CommitBlob(uploadId uuid.UUID) (*Blob, error) {
	url := "/blobs/commit"
	data := map[string]interface{}{"upload_id": uploadId.String()}
	params := map[string]interface{}{"_data": data}
	resp, err := ca.Call("POST", url, params)
	if err != nil {
		return nil, err
	}
	blobEnvelope := BlobEnvelope{}
	if err := json.Unmarshal([]byte(resp), &blobEnvelope); err != nil {
		return nil, err
	}

	return &blobEnvelope.Blob, nil
}

// Download a blob
func (ca *CustodiaAPIv1) GetBlobData(blobId uuid.UUID) (io.Reader, error) {
	url := fmt.Sprintf("/blobs/%s", blobId)
	params := map[string]interface{}{"_rawResponse": true}
	_, err := ca.Call("GET", url, params)
	if err != nil {
		return nil, err
	}

	return ca.RawResponse.Body, nil
}

// Delete a blob
func (ca *CustodiaAPIv1) DeleteBlob(blobId uuid.UUID) error {
	url := fmt.Sprintf("/blobs/%s", blobId)
	_, err := ca.Call("DELETE", url, nil)
	return err
}

// Generate a blob token used later to authenticate blob download
func (ca *CustodiaAPIv1) GenerateBlobToken(blobId uuid.UUID, oneTime bool,
	duration int) (*BlobToken, error) {
	url := fmt.Sprintf("/blobs/%s/generate", blobId)
	data := map[string]interface{}{"one_time": oneTime, "duration": duration}
	resp, err := ca.Call("POST", url, data)
	if err != nil {
		return nil, err
	}

	// response returns a map with toekn, expiration, and one_time
	// since one_time is a bool, but we already know it (it's a func param)
	// we return just the map with two strings
	blobToken := &BlobToken{}
	if err := json.Unmarshal([]byte(resp), blobToken); err != nil {
		return nil, err
	}

	return blobToken, nil
}

// download a blob with a token
func (ca *CustodiaAPIv1) GetBlobDataWithToken(blobId uuid.UUID, token string) (
	io.Reader, error) {
	url := fmt.Sprintf("/blobs/url/%s?token=%s", blobId, token)
	params := map[string]interface{}{"_rawResponse": true}

	ca.client.GetAuth().SwitchTo(common.NoAuth)

	_, err := ca.Call("GET", url, params)
	if err != nil {
		return nil, err
	}

	ca.client.GetAuth().SwitchBack()

	return ca.RawResponse.Body, nil
}

// Upload a blob from a file
//
// Args:
// - filePath: the path to the file to upload
// - documentId: the id of the document that will be associated with the blob
// - fieldName: the name of the field in the document where to store the blob
//
// Returns:
// - the created blob
func (ca *CustodiaAPIv1) CreateBlobFromFile(filePath string,
	documentId uuid.UUID, fieldName string, chunkSize int64) (*Blob, error) {
	if chunkSize == 0 {
		chunkSize = 1024 * 1024 // 1MB
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get the file name from the path
	fileName := filepath.Base(filePath)

	// Create a new blob
	uploadBlob, err := ca.CreateBlob(documentId, fieldName, fileName)
	if err != nil {
		return nil, err
	}

	// Get file size
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if fileInfo.Size() == 0 {
		return nil, fmt.Errorf("file is empty")
	}

	// Upload the file in chunks
	buffer := make([]byte, chunkSize)
	var offset int64 = 0

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}

		_, err = ca.UploadChunk(uploadBlob.Id, buffer[:n], int(chunkSize),
			int(offset))
		if err != nil {
			return nil, err
		}

		offset += int64(n)
	}

	// Commit the blob
	blob, err := ca.CommitBlob(uploadBlob.Id)
	if err != nil {
		return nil, err
	}

	return blob, nil
}

// Save a blob to a file
//
// Args:
// - blobId: the id of the blob to download
// - filePath: the path where to save the file
//
// Returns:
// - nil if successful, error otherwise
func (ca *CustodiaAPIv1) GetBlobToFile(blobId uuid.UUID, filePath string) (
	error) {
	// Create a new file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get the blob data
	data, err := ca.GetBlobData(blobId)
	if err != nil {
		return err
	}

	// Copy the blob data to the file
	_, err = io.Copy(file, data)
	if err != nil {
		return err
	}
	return nil
}
