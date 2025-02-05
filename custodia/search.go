package custodia

import (
	"encoding/json"
	"fmt"

	"github.com/dzanotelli/chino/common"
)

type SearchDocumentsResponse struct {
	Documents []*Document `json:"documents,omitempty"`
	Count int `json:"count"`
	TotalCount int `json:"total_count"`
	Limit int `json:"limit"`
	Offset int `json:"offset"`
}

type SearchUsersResponse struct {
	Users []*User `json:"users,omitempty"`
	Count int `json:"count"`
	TotalCount int `json:"total_count"`
	Limit int `json:"limit"`
	Offset int `json:"offset"`
}

type ResultType int

const (
	FullContent ResultType = iota + 1
	NoContent
	OnlyId
	Count
	Exists
	UsernameExists
)

func (rt ResultType) Choices() []string {
	return []string{"FULL_CONTENT", "NO_CONTENT", "ONLY_ID", "COUNT", "EXISTS",
		"USERNAME_EXISTS"}
}

func (rt ResultType) String() string {
    return rt.Choices()[rt-1]
}

func (rt ResultType) MarshalJSON() ([]byte, error) {
	return json.Marshal(rt.String())
}

func (rt* ResultType) UnmarshalJSON(data []byte) error {
	var value string
	err := json.Unmarshal(data, &value)
	if err!= nil {
		return err
	}
	intValue := indexOf(value, rt.Choices()) + 1  // enum starts from 1
	if intValue < 1 {
		return fmt.Errorf("ResultType: received unknown value '%v'", value)
	}

	*rt = ResultType(intValue)
	return nil
}

// Search documents
func (ca *CustodiaAPIv1) SearchDocuments(schemaId string,
	resultType ResultType, query map[string]interface{},
	sort map[string]interface{}) (*SearchDocumentsResponse, error) {
	if !common.IsValidUUID(schemaId) {
		return nil, fmt.Errorf("schemaId is not a valid UUID: %s", schemaId)
	}

	url := fmt.Sprintf("/search/documents/%s", schemaId)
	data := map[string]interface{}{"result_type": resultType, "query": query}
	if sort != nil {
		data["sort"] = sort
	}
	resp, err := ca.Call("POST", url, data)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	searchResponse := &SearchDocumentsResponse{}
	if err := json.Unmarshal([]byte(resp), searchResponse); err != nil {
		return nil, err
	}

	return searchResponse, nil
}

// Search users
func (ca *CustodiaAPIv1) SearchUsers(userSchemaId string,
	resultType ResultType, query map[string]interface{},
	sort map[string]interface{}) (*SearchUsersResponse, error) {
	if !common.IsValidUUID(userSchemaId) {
		return nil, fmt.Errorf("userSchemaId is not a valid UUID: %s",
			userSchemaId)
	}

	url := fmt.Sprintf("/search/users/%s", userSchemaId)
	data := map[string]interface{}{"result_type": resultType, "query": query}
	if sort != nil {
		data["sort"] = sort
	}
	resp, err := ca.Call("POST", url, data)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	searchResponse := &SearchUsersResponse{}
	if err := json.Unmarshal([]byte(resp), searchResponse); err != nil {
		return nil, err
	}

	return searchResponse, nil
}
