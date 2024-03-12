package custodia

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	// "reflect"
	"strconv"
	"strings"
	"time"

	"github.com/dzanotelli/chino/common"
	"github.com/simplereach/timeutils"
	// "golang.org/x/text/cases"
)

const TypeInt, TypeArrayInt = "integer", "array[integer]"
const TypeFloat, TypeArrayFloat = "float", "array[float]"
const TypeStr, TypeText, TypeArrayStr = "string", "text", "array[string]"
const TypeBool = "boolean"
const TypeDate, TypeTime, TypeDateTime = "date", "time", "datetime"
const TypeBase64, TypeJson, TypeBlob = "base64", "json", "blob"


// func struct2json(fields *struct{}) string {
// 	out = "{"
// 	fieldTemplate :=


// 	for field := range fields {
// 		fName := reflect.TypeOf(field).Tags.get("name")
// 		fIndexed := reflect.TypeOf(field).Tags.get("indexed")
// 		fDefault := reflect.TypeOf(field).Tags.get("default")
// 		fType := ""

// 		switch reflect.TypeOf(field) {
// 		case int, int8, int16, int32, int64:
// 			fType = "integer"
// 		case string:
// 			if reflect.TypeOf(field).Tags.get("text") {
// 				fType = "text"
// 			} else {
// 				fType = "string"
// 			}
// 		}
// 		out += fmt.Sprintf("\"name\":\"%s\",\"type\":\"%s\"", &fName, fType)
// 		if fIndexed {
// 			out += fmt.Sprintf(",\"indexed\":true")
// 		}
// 		if fDefault {
// 			// FIXME: quotes or not depending on type
// 			out += fmt.Sprintf(",\"default\":true")
// 		}

// 	}

// 	return out
// }




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

func parseArrayString(arrayString string, itemType string) ([]interface{},
	error) {
	var result []interface{}
	var ee []error
	var err error

	// remove brackets and split values
	arrayString = strings.TrimLeft(arrayString, "[")
	arrayString = strings.TrimRight(arrayString, "]")

	splitted := strings.Split(arrayString, ",")
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
		// FIXME: remove
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
	// case "base64":
	// 	converted = fmt.Sprintf("%v", value)
	// 	_, err = base64.StdEncoding.DecodeString(converted.(string))
	// 	if err != nil {
	// 		e = fmt.Errorf("field '%s': not a valid base64 string, %w",
	// 			field.Name, err)
	// 	}
	case TypeArrayInt, TypeArrayFloat, TypeArrayStr:
		arrayStr := fmt.Sprintf("%v", value)
		converted, e = parseArrayString(arrayStr, field.Type)
	default:
		e := fmt.Errorf("field '%s': type '%s' not handled", field.Name,
			field.Type)
		panic(e)
	}

	return converted, e
}

func convertData(data map[string]interface{}, schema Schema) (
	map[string]interface{}, []error) {
	converted := map[string]interface{}{}
	errors := []error{}
	structure := schema.getStructureAsMap()

	for name, value := range data {
		var err error

		field, ok := structure[name]
		if !ok {
			e := fmt.Errorf("field '%s': not belonging to schema %s '%s'",
				name, schema.SchemaId, schema.Description)
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