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


func ConvertSliceItemsOLD[T any](inputList interface{}) ([]T, error) {
	var output []T

	// input is a slice
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
func ConvertSliceItems[T any](input interface{}) ([]T, error) {
	list, ok := input.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected []interface{}, got: %T", input)
	}

	var output []T
	for _, item := range list {
		var converted T
		switch v := any(item).(type) {
		case float64:
			// Se T è un int, dobbiamo convertire il float64 in int
			if _, isInt := any(converted).(int); isInt {
				converted = any(int(v)).(T)
			} else if _, isInt64 := any(converted).(int64); isInt64 {
				converted = any(int64(v)).(T)
			} else if _, isFloat := any(converted).(float64); isFloat {
				converted = any(v).(T) // È già un float64, va bene
			} else {
				return nil, fmt.Errorf("unsupported number conversion for type %T", converted)
			}
		case string:
			if _, isString := any(converted).(string); isString {
				converted = any(v).(T)
			} else {
				return nil, fmt.Errorf("expected string but got different type")
			}
		default:
			// Tentiamo il cast diretto
			var ok bool
			converted, ok = any(v).(T)
			if !ok {
				return nil, fmt.Errorf("failed to convert item: %v", item)
			}
		}
		output = append(output, converted)
	}
	return output, nil
}