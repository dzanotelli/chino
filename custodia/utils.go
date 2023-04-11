package custodia

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dzanotelli/chino/common"
)


func validateContent(data map[string]interface{}, 
	structure map[string]SchemaField) []error {
	var errors []error
	var err error

	for key, value := range data {
		field, ok := structure[key]
		if !ok {
			err = fmt.Errorf("field '%s' not defined in given structure", key)
			errors = append(errors, err)
			continue
		}

		// field exist, check that is of the right type
		var val interface{}
		switch field.Type {
		case "integer":
			val, ok = value.(int)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be int", key)
			}
		case "float":
			val, ok = value.(float64)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be float", key)
			}
		case "string", "text":
			val, ok = value.(string)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be string", key)
				break
			}
			converted := fmt.Sprintf("%v", val)
			if field.Type == "string" && len(converted) > 255 {
				ok = false
				err = fmt.Errorf("field '%s' exceeded max lenght of 255 chars", 
					key)
			}
		case "boolean":
			val, ok = value.(bool)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be bool", key)
			}
		case "date", "time", "datetime":
			val, ok = value.(time.Time)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be time.Time", key)
			}
		case "base64":
			val, ok = value.(string)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be a string " +
					"(in base64 format)", key)
				break
			}
			converted := fmt.Sprintf("%v", val)
			_, err = base64.StdEncoding.DecodeString(converted)
			if err != nil {
				err = fmt.Errorf("field '%s' expected to be a valid base64 " +
					"string", key)
			}
		case "json":
			val, ok = value.(string)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be a string " +
					"(in json format)", key)
				break
			}
			converted := fmt.Sprintf("%v", val)
			if !json.Valid([]byte(converted)) {
				err = fmt.Errorf("field '%s' expected to be a valid json " +
					"string", key)
			}
		case "array[integer]":
			val, ok = value.([]int)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be a slice of int",
					key)
			}
		case "array[float]":
			val, ok = value.([]float64)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be a slice of " + 
				"float64", key)
			}
		case "array[string]":
			val, ok = value.([]string)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be a slice of " + 
				"string", key)
			}
		case "blob":
			val, ok = value.(string)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be a string " +
				"(UUID referencing a blob_id)", key)
				break
			}
			converted := fmt.Sprintf("%v", val)
			if len(converted) > 0 && !common.IsValidUUID(converted) {
				err = fmt.Errorf("field '%s' expected to be a valid " +
					"UUID (referencing a blob_id)", key)
			}
		default:			
			err = fmt.Errorf("unhandled type '%s' of field '%s'", 
				field.Type, key)
			panic(err)
		}
		
		// an error occurred, save it
		if !ok {
			errors = append(errors, err)
		}
	}
	return errors
}
