/*GRP-GNU-AGPL******************************************************************

File: query.go

Copyright (C) 2021  Team Georepublic <info@georepublic.de>

Developer(s):
Copyright (C) 2021  Ashish Kumar <ashishkr23438@gmail.com>

-----

This file is part of pg_scheduleserv.

pg_scheduleserv is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

pg_scheduleserv is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with pg_scheduleserv.  If not, see <https://www.gnu.org/licenses/>.

******************************************************************GRP-GNU-AGPL*/

package util

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
)

type PartialSQL struct {
	Fields []string
	Args   []interface{}
}

type LocationParams struct {
	Latitude  *float64 `json:"latitude"  validate:"required,max=90,min=-90" example:"2.0365"`
	Longitude *float64 `json:"longitude" validate:"required,max=180,min=-180" example:"48.6113"`
}

// Get an SQL query with partial fields
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

		// skip nil properties (not going to be patched), skip unexported fields
		if fieldVal.IsNil() || fieldType.PkgPath != "" {
			continue
		}

		// Change any alias fields
		alias, aliasFound := AliasFields[tag]
		if aliasFound {
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
		case reflect.Float32, reflect.Float64:
			partialSQL.Args = append(partialSQL.Args, val.Float())
		case reflect.String:
			partialSQL.Args = append(partialSQL.Args, val.String())
		case reflect.Bool:
			if val.Bool() {
				partialSQL.Args = append(partialSQL.Args, 1)
			} else {
				partialSQL.Args = append(partialSQL.Args, 0)
			}
		case reflect.Struct:
			value := val.Interface()
			if typ, ok := value.(LocationParams); ok {
				partialSQL.Args = append(partialSQL.Args, GetLocationId(*typ.Latitude, *typ.Longitude))
			}
			if typ, ok := value.(string); ok {
				partialSQL.Args = append(partialSQL.Args, typ)
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

func GetOutputFields(resourceStruct interface{}) (sql string) {
	val := reflect.ValueOf(resourceStruct)
	for i := 0; i < val.Type().NumField(); i++ {
		if i != 0 {
			sql += ","
		}
		field := val.Type().Field(i)
		fieldName := jsonTag(field)

		if aliasField, aliasFieldFound := AliasFields[fieldName]; aliasFieldFound {
			fieldName = aliasField
		}

		if _, intervalFieldFound := IntervalFields[fieldName]; intervalFieldFound {
			fieldName = fmt.Sprintf("to_char(%s, 'HH24:MI:SS')", fieldName)
		}
		if _, timestampFieldFound := TimestampFields[fieldName]; timestampFieldFound {
			fieldName = fmt.Sprintf("to_char(%s, 'YYYY-MM-DD') || 'T' || to_char(%s, 'HH24:MI:SS')", fieldName, fieldName)
		}
		sql += " " + fieldName
	}
	return sql
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

// Checks whether the field tag of json matches a particular field
func checkjsonTagField(field reflect.StructField, matchField string) bool {
	if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
		parts := strings.Split(jsonTag, ",")
		for i := 1; i < len(parts); i++ {
			if parts[i] == matchField {
				return true
			}
		}
	}
	return false
}
