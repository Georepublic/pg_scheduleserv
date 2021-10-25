package util

import (
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
)

type PartialSQL struct {
	Fields []string
	Args   []interface{}
}

var readOnlyFields = map[string]bool{
	"id":         true,
	"created_at": true,
	"updated_at": true,
}

var aliasFields = map[string]string{
	"location": "location_index",
}

type LocationParams struct {
	Latitude  *float64 `json:"latitude"  validate:"required"`
	Longitude *float64 `json:"longitude" validate:"required"`
}

// Get an SQL query with partial fields, excluding the read-only fields
// Takes a resource object as parameter, and returns its fields and arguments.
// (returns only those arguments which are not nil)
// See https://stackoverflow.com/questions/38206479/golang-rest-patch-and-building-an-update-query
func GetPartialSQL(resource interface{}) PartialSQL {
	var partialSQL PartialSQL
	resType := reflect.TypeOf(resource)
	resVal := reflect.ValueOf(resource)
	n := resType.NumField()

	partialSQL.Fields = make([]string, 0, n)
	partialSQL.Args = make([]interface{}, 0, n)

	for i := 0; i < n; i++ {
		fieldType := resType.Field(i)
		fieldVal := resVal.Field(i)
		tag := jsonTag(fieldType)

		// skip nil properties (not going to be patched), skip unexported fields or read only fields
		_, readOnlyFound := readOnlyFields[tag]
		if fieldVal.IsNil() || fieldType.PkgPath != "" || readOnlyFound == true {
			continue
		}

		// Change any alias fields
		alias, aliasFound := aliasFields[tag]
		if aliasFound == true {
			tag = alias
		}

		partialSQL.Fields = append(partialSQL.Fields, tag)

		var val reflect.Value
		if fieldVal.Kind() == reflect.Ptr {
			val = fieldVal.Elem()
		} else {
			val = fieldVal
		}

		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			partialSQL.Args = append(partialSQL.Args, val.Int())
		case reflect.String:
			partialSQL.Args = append(partialSQL.Args, val.String())
		case reflect.Bool:
			if val.Bool() {
				partialSQL.Args = append(partialSQL.Args, 1)
			} else {
				partialSQL.Args = append(partialSQL.Args, 0)
			}
		case reflect.Struct:
			location := val.Interface()
			if typ, ok := location.(LocationParams); ok {
				partialSQL.Args = append(partialSQL.Args, GetLocationIndex(*typ.Latitude, *typ.Longitude))
			}
		case reflect.Interface:
			partialSQL.Args = append(partialSQL.Args, val.Interface())
		case reflect.Slice:
			partialSQL.Args = append(partialSQL.Args, val.Interface())
		default:
			logrus.Error("Invalid value type: ", val.Kind(), tag)
		}
	}

	return partialSQL
}

// Get the json tag name from struct field
func jsonTag(field reflect.StructField) string {
	fieldName := field.Name
	name := fieldName

	if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
		parts := strings.Split(jsonTag, ",")
		name = parts[0]
		if name == "" {
			name = fieldName
		}
	}
	name = strings.ToLower(name)
	return name
}
