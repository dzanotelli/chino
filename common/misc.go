package common

import (
	"fmt"

	"github.com/google/uuid"
)

// IsValidUUID checks that u is a valid UUID4 string
func IsValidUUID(u string) bool {
    _, err := uuid.Parse(u)
    return err == nil
 }



func GetFakeAuth() *ClientAuth {
	fakeAuth := NewClientAuth(map[string]interface{}{
		"customerId": "00000000-0000-0000-0000-000000000000",
		"customerKey": "00000000-0000-0000-0000-000000000000",
	})
	fakeAuth.SwitchTo(CustomerAuth)

	return fakeAuth
}


func ConvertSliceItems[T any](inputList interface{}) ([]T, error) {
	var output []T

	// input is a list
	list, ok := inputList.([]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to convert slice: %v", inputList)
	}

	for _, item := range list {
		converted, ok := item.(T)
		if !ok {
			return nil, fmt.Errorf("failed to convert slice item: %v", item)
		}
		output = append(output, converted)
	}
	return output, nil
}