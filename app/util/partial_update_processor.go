package util

import (
	"reflect"
	"strings"
)

func PartialUpdatePreprocess[T any, S any](values T, setNull *map[string]bool, existingValues S) (S, error) {
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

	// Step 1：handle setNull mark, set the corresponded fields to null
	if setNull != nil {
		for fieldName, shouldSetNull := range *setNull {
			if shouldSetNull {
				// find corresponded fields in existingValues
				existingField, found := findFieldByName(existingReflect, existingType, fieldName)
				if found && existingField.CanSet() {
					// if the pointer is nil, set it as to a zero value
					existingField.Set(reflect.Zero(existingField.Type()))
				}
			}
		}
	}

	// Step 2：handle non zero fields in values, replace them to existingValues
	for i := 0; i < valuesReflect.NumField(); i++ {
		valuesField := valuesReflect.Field(i)
		valuesFieldType := valuesType.Field(i)
		fieldName := valuesFieldType.Name

		// skip zero fields
		if isZeroValue(valuesField) {
			continue
		}

		// find corresponded fields in the existingValues
		existingField, found := findFieldByName(existingReflect, existingType, fieldName)
		if !found || !existingField.CanSet() {
			continue
		}

		if valuesField.Kind() == reflect.Ptr && !valuesField.IsNil() {
			// if values is non nil pointer
			if existingField.Kind() == reflect.Ptr {
				existingField.Set(valuesField)
			} else {
				existingField.Set(valuesField.Elem())
			}
		} else if valuesField.Kind() != reflect.Ptr {
			// if values is not a pointer
			if existingField.Kind() == reflect.Ptr {
				// construct new pointer
				newValue := reflect.New(existingField.Type().Elem())
				newValue.Elem().Set(valuesField)
				existingField.Set(newValue)
			} else {
				existingField.Set(valuesField)
			}
		}
	}

	return existingValues, nil
}

func findFieldByName(structReflect reflect.Value, structType reflect.Type, fieldName string) (reflect.Value, bool) {
	for i := 0; i < structReflect.NumField(); i++ {
		if structType.Field(i).Name == fieldName {
			return structReflect.Field(i), true
		}
	}

	capitalizedFieldName := strings.Title(fieldName)
	for i := 0; i < structReflect.NumField(); i++ {
		if structType.Field(i).Name == capitalizedFieldName {
			return structReflect.Field(i), true
		}
	}

	return reflect.Value{}, false
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	default:
		return v.IsZero()
	}
}

/* ============================== Helper functions for directly partial update operations ============================== */

func CheckSetNull(setNullMap *map[string]bool, fieldName string) bool {
	if setNullMap != nil {
		return (*setNullMap)[fieldName]
	}
	return false
}
