package util

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	exceptions "notezy-backend/app/exceptions"
)

type SearchCursor[SearchableCursorFieldType any] struct {
	Fields SearchableCursorFieldType `json:"fields"`
}

func NewSearchCursor[SearchableCursorFieldType any](fields SearchableCursorFieldType) *SearchCursor[SearchableCursorFieldType] {
	return &SearchCursor[SearchableCursorFieldType]{
		Fields: fields,
	}
}

func (sc *SearchCursor[SearchableCursorFieldType]) EncodeSearchCursor() (*string, *exceptions.Exception) {
	jsonData, err := json.Marshal(sc.Fields)
	if err != nil {
		return nil, exceptions.Searchable.FailedToMarshalSearchCursor().WithError(err)
	}

	encoded := base64.StdEncoding.EncodeToString(jsonData)
	return &encoded, nil
}

func DecodeSearchCursor[SearchableCursorFieldType any](encoded string) (*SearchCursor[SearchableCursorFieldType], *exceptions.Exception) {
	if len(strings.ReplaceAll(encoded, " ", "")) == 0 {
		return nil, exceptions.Searchable.EmptyEncodedStringToDecodeSearchCursor()
	}

	jsonData, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, exceptions.Searchable.FailedToDecodeBase64String().WithError(err)
	}

	var fields SearchableCursorFieldType
	if err := json.Unmarshal(jsonData, &fields); err != nil {
		return nil, exceptions.Searchable.FailedToUnmarshalSearchCursor().WithError(err)
	}

	return &SearchCursor[SearchableCursorFieldType]{Fields: fields}, nil
}

func EncodeSearchCursorFromData[SearchCursorType any](data SearchCursorType) (*string, *exceptions.Exception) {
	cursor := NewSearchCursor(data)
	return cursor.EncodeSearchCursor()
}
