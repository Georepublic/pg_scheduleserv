/*GRP-GNU-AGPL******************************************************************

File: validator.go

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

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-multierror"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

var locationTags = map[string]bool{
	"location":       true,
	"p_location":     true,
	"d_location":     true,
	"start_location": true,
	"end_location":   true,
}

func NewValidator() *validator.Validate {
	validate := validator.New()
	// Get json tag name instead of actual struct field name
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return validate
}

// Verify that the type of input user struct is same as the required type
func ValidateInput(jsonStruct map[string]interface{}, originalStruct interface{}) error {
	var errors error
	resType := reflect.TypeOf(originalStruct)
	if resType.Kind() != reflect.Struct {
		logrus.Error("Bad type: requires a struct")
	}
	for i := 0; i < resType.NumField(); i++ {
		fieldType := resType.Field(i)
		tag := jsonTag(fieldType)

		// Ignore any nil fields in the input
		if jsonStruct[tag] == nil {
			continue
		}

		userType := reflect.TypeOf(jsonStruct[tag])
		requiredType := fieldType.Type.Elem()

		switch {
		case checkjsonTagField(fieldType, "string"):
			requiredType = reflect.TypeOf("string")
		}

		// Need to validate struct fields separately
		typ, ok := jsonStruct[tag].(map[string]interface{})
		if ok && requiredType.Kind() == reflect.Struct {
			// LocationParams Struct
			if _, locationTagsFound := locationTags[tag]; locationTagsFound {
				location := LocationParams{}
				if err := mapstructure.Decode(typ, &location); err != nil {
					return fmt.Errorf("Field 'latitude' and 'longitude' of type 'float64' is required")
				}
				validate := NewValidator()
				if err := validate.Struct(location); err != nil {
					return err
				}
				continue
			}
		}

		// Need to validate []int64 fields separately
		typ2, ok := jsonStruct[tag].([]interface{})
		if ok && requiredType.Kind() == reflect.Slice {
			convertible := true
			for i := 0; i < len(typ2); i++ {
				// validation for time_windows field: [][]string
				if tag == "time_windows" || tag == "p_time_windows" || tag == "d_time_windows" {
					if typ2[i] == nil {
						convertible = false
					} else {
						if typ3, ok := typ2[i].([]interface{}); !ok {
							convertible = false
							if len(typ3) != 2 {
								convertible = false
							}
							for j := 0; j < len(typ3); j++ {
								if (typ3[j] == nil) || (reflect.TypeOf(typ3[j]).Kind() != reflect.String) {
									convertible = false
								}
							}
						}
					}
				} else if typ2[i] == nil || !reflect.TypeOf(typ2[i]).ConvertibleTo(requiredType.Elem()) {
					convertible = false
				}
			}
			if !convertible {
				errors = multierror.Append(errors, fmt.Errorf(fmt.Sprintf("Field '%s' must be of '%s' type.", tag, requiredType)))
			}
			continue
		}

		if !userType.ConvertibleTo(requiredType) {
			errors = multierror.Append(errors, fmt.Errorf(fmt.Sprintf("Field '%s' must be of '%s' type.", tag, requiredType)))
		}
	}
	return errors
}
