// app/util/partial_update_processor.go
package util

import (
	"reflect"
)

func PartialUpdatePreprocess[T any, S any](values T, setNull map[string]bool, existingValues S) (S, error) {
	existingReflect := reflect.ValueOf(&existingValues).Elem()
	existingType := reflect.TypeOf(existingValues)
	valuesReflect := reflect.ValueOf(values)
	valuesType := reflect.TypeOf(values)

	// handle pointer type for existing values
	if existingReflect.Kind() == reflect.Ptr {
		existingReflect = existingReflect.Elem()
		existingType = existingType.Elem()
	}

	// handle pointer type for new values
	if valuesReflect.Kind() == reflect.Ptr {
		valuesReflect = valuesReflect.Elem()
		valuesType = valuesType.Elem()
	}

	// iterate the existingField which is the entire struct of the database table
	for i := 0; i < existingReflect.NumField(); i++ {
		existingField := existingReflect.Field(i)
		existingFieldType := existingType.Field(i)
		fieldName := existingFieldType.Name

		// check if the field in existingField can be modified
		if !existingField.CanSet() {
			continue
		}

		// 1. check if there's a true value in setNull which indicating we need to set it to null
		if setNull != nil {
			if shouldSetNull, exists := setNull[fieldName]; exists && shouldSetNull {
				// setting the zero value（which is nil for pointer）
				existingField.Set(reflect.Zero(existingField.Type()))
				continue
			}
		}

		// 2. check if there's a corresponding field in valuesField
		valuesField, hasCorrespondingField := findFieldByName(valuesReflect, valuesType, fieldName)
		if hasCorrespondingField {
			if valuesField.Kind() == reflect.Ptr && !valuesField.IsNil() { // check if there's non-nil value in valuesField
				// new value in valuesField, setting it to the existingField
				if existingField.Kind() == reflect.Ptr {
					// both of them are pointers
					existingField.Set(valuesField)
				} else {
					// existingField is NOT a pointer, valuesField is a pointer
					existingField.Set(valuesField.Elem())
				}
				continue
			} else if valuesField.Kind() != reflect.Ptr { // if the valuesField is not pointer, set it directly
				if existingField.Kind() == reflect.Ptr {
					// existingField is a pointer, valuesField is NOT a pointer
					newValue := reflect.New(existingField.Type().Elem())
					newValue.Elem().Set(valuesField)
					existingField.Set(newValue)
				} else {
					// neither existingField nor valuesField are pointers
					existingField.Set(valuesField)
				}
				continue
			}
		}

		// 3. leave the existingField
	}

	return existingValues, nil
}

func findFieldByName(structReflect reflect.Value, structType reflect.Type, fieldName string) (reflect.Value, bool) {
	for i := 0; i < structReflect.NumField(); i++ {
		if structType.Field(i).Name == fieldName {
			return structReflect.Field(i), true
		}
	}
	return reflect.Value{}, false
}
