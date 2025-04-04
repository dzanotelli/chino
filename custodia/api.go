package custodia

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dzanotelli/chino/common"
)

// CustodiaEnvelope is the enveloped response, with data in subobject "data"
type CustodiaEnvelope struct {
	Result string `json:"result"`
	ResultCode uint64 `json:"result_code"`
	Message json.RawMessage `json:"message"`
	Data json.RawMessage `json:"data"`
}

type CustodiaAPIv1 struct {
	client *common.Client
	RawResponse *http.Response
}

// NewCustodiaAPI returns a new CustodiaAPI object to interact
// with the Custodia APIs
func NewCustodiaAPIv1(client *common.Client) *CustodiaAPIv1 {
	capi := &CustodiaAPIv1{}
	capi.client = client
	return capi
}

func (ca *CustodiaAPIv1) Call(method, path string,
	params map[string]interface{}) (string, error) {
	rawResponse, ok := params["_rawResponse"].(bool)
	if !ok {
		rawResponse = false
	}

	httpResp, err := ca.client.Call(method, "/api/v1" + path, params)

	// save the response for further inspection on need
	ca.RawResponse = httpResp

	if err != nil || rawResponse {
		return "", err
	}
	defer httpResp.Body.Close()

	resp := CustodiaEnvelope{}
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return "", err
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode > 299 {
		err := fmt.Errorf("error %v: %s", resp.ResultCode, resp.Message)
		return "", err
	}

	return string(resp.Data), nil
}