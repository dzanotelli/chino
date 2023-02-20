package storage

import (
	"net/http"

	"github.com/dzanotelli/chino/common"
)

type StorageAPIv1 struct {
	client *common.Client
}

// NewCustodiaAPI returns a new CustodiaAPI object to interact
// with the Custodia APIs
func NewStorageAPIv1(client *common.Client) *StorageAPIv1 {
	capi := &StorageAPIv1{}
	capi.client = client
	return capi
}

func (ca *StorageAPIv1) Get(path string) (*http.Response, error) {
	return ca.client.Get("/api/v1" + path)
}

func (ca *StorageAPIv1) Post(path, payload string) (*http.Response, error) {
	return ca.client.Post("/api/v1" + path, payload)
}

func (ca *StorageAPIv1) Put(path, payload string) (*http.Response, error) {
	return ca.client.Put("/api/v1" + path, payload)
}

func (ca *StorageAPIv1) Patch(path, payload string) (*http.Response, error) {
	return ca.client.Patch("/api/v1" + path, payload)
}

func (ca *StorageAPIv1) Delete(path string) (*http.Response, error) {
	return ca.client.Delete("/api/v1" + path)
}