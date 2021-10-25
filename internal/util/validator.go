package util

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-multierror"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

// Verify that the type of input user struct is same as the required type
func ValidateInput(jsonStruct map[string]interface{}, originalStruct interface{}) error {
	var errors error
	resType := reflect.TypeOf(originalStruct)
	if resType.Kind() != reflect.Struct {
		logrus.Debug(resType.Kind())
		panic("Bad type: requires a struct")
	}
	for i := 0; i < resType.NumField(); i++ {
		fieldType := resType.Field(i)
		tag := jsonTag(fieldType)

		// Ignore any nil fields in the input
		if jsonStruct[tag] == nil {
			logrus.Debug("Skipping ", tag, " ", jsonStruct)
			continue
		}

		userType := reflect.TypeOf(jsonStruct[tag])
		requiredType := fieldType.Type.Elem()

		// Need to validate struct fields separately
		typ, ok := jsonStruct[tag].(map[string]interface{})
		if ok && requiredType.Kind() == reflect.Struct {
			// LocationParams Struct
			if tag == "location" {
				location := LocationParams{}
				mapstructure.Decode(typ, &location)
				validate := validator.New()
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
				if !reflect.TypeOf(typ2[i]).ConvertibleTo(requiredType.Elem()) {
					convertible = false
				}
			}
			if !convertible {
				errors = multierror.Append(errors, fmt.Errorf(fmt.Sprintf("Field '%s' must be of %s type. Found type: %s", tag, requiredType, userType)))
			}
			continue
		}

		if !userType.ConvertibleTo(requiredType) {
			errors = multierror.Append(errors, fmt.Errorf(fmt.Sprintf("Field '%s' must be of %s type. Found type: %s", tag, requiredType, userType)))
		}
	}
	return errors
}
