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

	// 第一步：處理 setNull 標記，將對應欄位設為 nil
	if setNull != nil {
		for fieldName, shouldSetNull := range *setNull {
			if shouldSetNull {
				// 在 existingValues 中找到對應欄位
				existingField, found := findFieldByName(existingReflect, existingType, fieldName)
				if found && existingField.CanSet() {
					// 設置為零值（對於 pointer 就是 nil）
					existingField.Set(reflect.Zero(existingField.Type()))
				}
			}
		}
	}

	// 第二步：處理 values 中的非零值，替換到 existingValues
	for i := 0; i < valuesReflect.NumField(); i++ {
		valuesField := valuesReflect.Field(i)
		valuesFieldType := valuesType.Field(i)
		fieldName := valuesFieldType.Name

		// 跳過零值欄位
		if isZeroValue(valuesField) {
			continue
		}

		// 在 existingValues 中找到對應欄位
		existingField, found := findFieldByName(existingReflect, existingType, fieldName)
		if !found || !existingField.CanSet() {
			continue
		}

		// 根據類型進行適當的設置
		if valuesField.Kind() == reflect.Ptr && !valuesField.IsNil() {
			// values 欄位是非 nil 的 pointer
			if existingField.Kind() == reflect.Ptr {
				existingField.Set(valuesField)
			} else {
				existingField.Set(valuesField.Elem())
			}
		} else if valuesField.Kind() != reflect.Ptr {
			// values 欄位不是 pointer，直接設置
			if existingField.Kind() == reflect.Ptr {
				// 創建新的 pointer
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
