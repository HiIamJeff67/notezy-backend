package util

import (
	"reflect"
)

func CopyNonNilFields(target interface{}, input interface{}) {
	targetVal := reflect.ValueOf(target).Elem()
	inputVal := reflect.ValueOf(input)
	inputType := inputVal.Type()

	for i := 0; i < inputVal.NumField(); i++ {
		inField := inputVal.Field(i)
		fieldName := inputType.Field(i).Name
		tarField := targetVal.FieldByName(fieldName)
		if !tarField.IsValid() || !tarField.CanSet() {
			continue
		}
		if inField.Kind() == reflect.Ptr && !inField.IsNil() {
			tarField.Set(inField.Elem())
		}
		if inField.Kind() == reflect.Ptr && tarField.Kind() == reflect.Ptr && !inField.IsNil() {
			tarField.Set(inField)
		}
	}
}
