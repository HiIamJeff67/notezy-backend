package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
)

type SearchCursorField struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

type SearchCursor struct {
	Fields []SearchCursorField `json:"fields"`
}

func EncodeSearchCursor(data interface{}, fieldNames ...string) (string, error) {
	if data == nil {
		return "", fmt.Errorf("data cannot be nil")
	}

	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return "", fmt.Errorf("data must be a struct")
	}

	searchCursor := SearchCursor{
		Fields: make([]SearchCursorField, 0, len(fieldNames)),
	}

	for _, fieldName := range fieldNames {
		field := val.FieldByName(fieldName)
		if !field.IsValid() {
			return "", fmt.Errorf("field %s not found", fieldName)
		}
		searhCursorField := SearchCursorField{
			Name:  fieldName,
			Value: field.Interface(),
			Type:  field.Type().String(),
		}
		searchCursor.Fields = append(searchCursor.Fields, searhCursorField)
	}

	jsonData, err := json.Marshal(searchCursor)
	if err != nil {
		return "", fmt.Errorf("failed to marshal search cursor: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(jsonData)
	return encoded, nil
}

func DecodeSearchCursor(encoded string) (*SearchCursor, error) {
	if encoded == "" {
		return nil, fmt.Errorf("encoded string cannot be empty")
	}

	jsonData, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	var cursor SearchCursor
	if err := json.Unmarshal(jsonData, &cursor); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cursor: %w", err)
	}

	return &cursor, nil
}

func (sc *SearchCursor) ConvertValue(value interface{}, typeName string) (interface{}, error) {
	switch typeName {
	case "time.Time":
		if str, ok := value.(string); ok {
			return time.Parse(time.RFC3339, str)
		}
	case "uuid.UUID":
		if str, ok := value.(string); ok {
			return uuid.Parse(str)
		}
	case "int8", "int", "int32", "int64":
		if num, ok := value.(float64); ok {
			return int64(num), nil
		}
	case "string":
		if str, ok := value.(string); ok {
			return str, nil
		}
	}

	return value, nil
}

func (sc *SearchCursor) GetValue(fieldName string) (interface{}, error) {
	for _, field := range sc.Fields {
		if field.Name == fieldName {
			return sc.ConvertValue(field.Value, field.Type)
		}
	}

	return nil, fmt.Errorf("field %s not found in search cursor", fieldName)
}
