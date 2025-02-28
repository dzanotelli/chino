package custodia

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"strconv"
	"strings"
	"time"

	"github.com/dzanotelli/chino/common"
	"github.com/google/uuid"
	"github.com/simplereach/timeutils"
)

const TypeInt, TypeArrayInt = "integer", "array[integer]"
const TypeFloat, TypeArrayFloat = "float", "array[float]"
const TypeStr, TypeText, TypeArrayStr = "string", "text", "array[string]"
const TypeBool = "boolean"
const TypeDate, TypeTime, TypeDateTime = "date", "time", "datetime"
const TypeBase64, TypeJson, TypeBlob = "base64", "json", "blob"

// Return the index of the first found occurence of word in data
// or -1 if not found
func indexOf(word string, data []string) (int) {
    for k, v := range data {
        if word == v {
            return k
        }
    }
    return -1
}

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
		case TypeInt:
			val, ok = value.(int64)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be int64", key)
			}
		case TypeFloat:
			val, ok = value.(float64)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be float64", key)
			}
		case TypeStr, TypeText:
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
		case TypeBool:
			val, ok = value.(bool)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be bool", key)
			}
		case TypeDate, TypeTime, TypeDateTime:
			val, ok = value.(time.Time)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be time.Time", key)
			}
		case TypeBase64:
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
		case TypeJson:
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
		case TypeArrayInt:
			val, ok = value.([]int64)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be a slice of int64",
					key)
			}
		case TypeArrayFloat:
			val, ok = value.([]float64)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be a slice of " +
				"float64", key)
			}
		case TypeArrayStr:
			val, ok = value.([]string)
			if !ok {
				err = fmt.Errorf("field '%s' expected to be a slice of " +
				"string", key)
			}
		case TypeBlob:
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

// Manually parse a JSON array of int, floats or strings
func parseJSONArray(strArray string, itemType string) ([]interface{},
	error) {
	var result []interface{}
	var ee []error
	var err error

	// remove brackets and split values
	strArray = strings.TrimLeft(strArray, "[")
	strArray = strings.TrimRight(strArray, "]")

	splitted := strings.Split(strArray, ",")
	for i, v := range splitted {
		splitted[i] = strings.Trim(v, " ")
	}

	switch itemType {
	case TypeArrayInt:
		for i, v := range splitted {
			converted, e := strconv.ParseInt(v, 10, 64)
			if e != nil {
				ee = append(ee, fmt.Errorf("%d: ParseInt error", i))
				result = append(result, nil)
			} else {
				result = append(result, converted)
			}
		}
	case TypeArrayFloat:
		for i, v := range splitted {
			converted, e := strconv.ParseFloat(v, 64)
			if e != nil {
				ee = append(ee, fmt.Errorf("%d: ParseFloat error", i))
				result = append(result, nil)
			} else {
				result = append(result, converted)
			}
		}
	case TypeArrayStr:
		// since we are manually parsing the JSON array of values, we need to
		// remove the double quotes around each value
		for _, v := range splitted {
			result = append(result, strings.Trim(v, "\""))
		}
	default:
		panic(fmt.Sprintf("unhandled type '%s'", itemType))
	}

	if len(ee) > 0 {
		err = errors.Join(ee...)
	}
	return result, err
}

func convertField(value interface{}, field SchemaField) (interface{}, error) {
	var converted interface{}
	var e, err error
	var ok bool

	switch field.Type {
	case TypeInt:
		// json.Unmarshall always returns float64 for numbers
		c, ok := value.(float64)
		if !ok {
			e = fmt.Errorf("field '%s': cannot convert to int64", field.Name)
		}
		converted = int64(c)
	case TypeFloat:
		converted, ok = value.(float64)
		if !ok {
			e = fmt.Errorf("field '%s': cannot convert to float64", field.Name)
		}
	case TypeStr, TypeText, TypeBase64, TypeJson, TypeBlob:
		converted = fmt.Sprintf("%v", value)
	case TypeBool:
		converted, ok = value.(bool)
		if !ok {
			e = fmt.Errorf("field '%s': cannot convert to bool", field.Name)
		}
	case TypeDate, TypeTime, TypeDateTime:
		dateStr := fmt.Sprintf("%v", value)
		converted, err = timeutils.ParseDateString(dateStr)
		if err != nil {
			e = fmt.Errorf("field '%s': error while converting to " +
				"time.Time, %w", field.Name, err)
		}
	case TypeArrayInt, TypeArrayFloat, TypeArrayStr:
		arrayStr := fmt.Sprintf("%v", value)
		converted, e = parseJSONArray(arrayStr, field.Type)
	default:
		e := fmt.Errorf("field '%s': type '%s' not handled", field.Name,
			field.Type)
		panic(e)
	}

	return converted, e
}

type StructureMapper interface {
	getStructureAsMap() map[string]SchemaField
}

func convertData(data map[string]interface{}, schema StructureMapper) (
	map[string]interface{}, []error) {
	converted := map[string]interface{}{}
	errors := []error{}
	structure := schema.getStructureAsMap()

	// get id and description of schema/userschema
	var id uuid.UUID
	var descr string
	switch concreteSchema := schema.(type) {
	case *UserSchema:
		id = concreteSchema.Id
		descr = concreteSchema.Description
	case *Schema:
		id = concreteSchema.Id
		descr = concreteSchema.Description
	default:
		panic(fmt.Sprintf("unhandled type '%T'", schema))
	}

	for name, value := range data {
		var err error

		field, ok := structure[name]
		if !ok {
			e := fmt.Errorf("field '%s': not belonging to schema %s '%s'",
				name, id, descr)
			errors = append(errors, e)
			continue
		}

		converted[name], err = convertField(value, field)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return converted, errors
}
