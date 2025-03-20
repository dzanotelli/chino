package custodia

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type SearchResponse struct {
	Documents []*Document `json:"documents,omitempty"`
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
func (ca *CustodiaAPIv1) SearchDocuments(schemaId uuid.UUID,
	resultType ResultType, query map[string]any,
	sort map[string]any) (*SearchResponse, error) {
	url := fmt.Sprintf("/search/documents/%s", schemaId)
	data := map[string]any{"result_type": resultType.String(),
		"query": query}
	if sort != nil {
		data["sort"] = sort
	}
	params := map[string]any{"_data": true}
	resp, err := ca.Call("POST", url, params)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	searchResponse := &SearchResponse{}
	if err := json.Unmarshal([]byte(resp), searchResponse); err != nil {
		return nil, err
	}

	// FIXME: golang unmarshals returns always float64 for numbers
	//   dunno if let the user do this, or force convertion of the
	//   underlying type in SearchDocuments (but we need the schema!)
	//   Check `ReadDocument` for more (there we apply conversion)

	return searchResponse, nil
}

// Search users
func (ca *CustodiaAPIv1) SearchUsers(userSchemaId uuid.UUID,
	resultType ResultType, query map[string]any,
	sort map[string]any) (*SearchResponse, error) {
	url := fmt.Sprintf("/search/users/%s", userSchemaId)
	data := map[string]any{"result_type": resultType.String(),
		"query": query}
	if sort != nil {
		data["sort"] = sort
	}
	params := map[string]any{"_data": data}
	resp, err := ca.Call("POST", url, params)
	if err != nil {
		return nil, err
	}

	// JSON: unmarshal resp content
	searchResponse := &SearchResponse{}
	if err := json.Unmarshal([]byte(resp), searchResponse); err != nil {
		return nil, err
	}

	return searchResponse, nil
}
