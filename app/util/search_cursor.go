package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"

	exceptions "notezy-backend/app/exceptions"
)

type SearchCursorField struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

type SearchCursor struct {
	Fields []SearchCursorField `json:"fields"`
}

func EncodeSearchCursor(data interface{}) (*string, *exceptions.Exception) {
	if data == nil {
		return nil, exceptions.Searchable.InvalidNilDataToEncodeSearchCursor()
	}

	var searchCursor SearchCursor

	if mapData, ok := data.(map[string]interface{}); ok {
		searchCursor.Fields = make([]SearchCursorField, 0, len(mapData))

		for fieldName, fieldValue := range mapData {
			fieldType := "nil"
			if fieldValue != nil {
				fieldType = reflect.TypeOf(fieldValue).String()
			}
			field := SearchCursorField{
				Name:  fieldName,
				Value: fieldValue,
				Type:  fieldType,
			}
			searchCursor.Fields = append(searchCursor.Fields, field)
		}
	} else {
		return nil, exceptions.Searchable.InvalidNonMapToEncodeSearchCursor()
	}

	jsonData, err := json.Marshal(searchCursor)
	if err != nil {
		return nil, exceptions.Searchable.FailedToMarshalSearchCursor().WithError(err)
	}

	encoded := base64.StdEncoding.EncodeToString(jsonData)
	return &encoded, nil
}

func DecodeSearchCursor(encoded string) (*SearchCursor, *exceptions.Exception) {
	if len(strings.ReplaceAll(encoded, " ", "")) == 0 {
		return nil, exceptions.Searchable.EmptyEncodedStringToDecodeSearchCursor()
	}

	jsonData, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, exceptions.Searchable.FailedToDecodeBase64String().WithError(err)
	}

	var cursor SearchCursor
	if err := json.Unmarshal(jsonData, &cursor); err != nil {
		return nil, exceptions.Searchable.FailedToUnMarshalSearchCursor().WithError(err)
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
