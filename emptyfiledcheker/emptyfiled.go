package emptyfieldcheker

import (
	"fmt"
	"reflect"
)
func EmptyField(obj any, excludedFields ...string) (bool, []string) {
	var empty bool
	var errors []string

	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	if val.Kind() != reflect.Struct {
		return false, []string{"Provided value is not a struct"}
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name

		
		isExcluded := false
		for _, excluded := range excludedFields {
			if fieldName == excluded {
				isExcluded = true
				break
			}
		}
		if isExcluded {
			continue
		}


		switch field.Kind() {
		case reflect.String:
			if field.String() == "" {
				errors = append(errors, fmt.Sprintf("Field: %v is empty", fieldName))
				empty = true
			}
		case reflect.Int:
			if field.Int() == 0 {
				errors = append(errors, fmt.Sprintf("Field: %v is empty", fieldName))
				empty = true
			}
		case reflect.Float64:
			if field.Float() == 0 {
				errors = append(errors, fmt.Sprintf("Field: %v is empty", fieldName))
				empty = true
			}
		case reflect.Struct:
	
			isEmpty, nestedErrors := EmptyField(field.Interface(), excludedFields...)
			if isEmpty {
				errors = append(errors, fmt.Sprintf("Nested structure field: %v has empty fields", fieldName))
				empty = true
			}
			errors = append(errors, nestedErrors...)

		default:

		}
	}

	return empty, errors
}
