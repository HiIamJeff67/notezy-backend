package searchcursor

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	exceptions "notezy-backend/app/exceptions"
)

// we may name other search cursors based on their functionality.
// ex. ImprovableSearchCursor, AlphaNameSearchCursor, TimeBasedSearchCursor,
//     ClosestSeachCursor, ReachableSearchCursor, etc.

// This `SearchCursor` is a general and common search cursor used as default in most of the search scenarios
type SearchCursor[SearchCursorFieldType any] struct {
	Fields SearchCursorFieldType `json:"fields"`
}

func New[SearchCursorFieldType any](fields SearchCursorFieldType) *SearchCursor[SearchCursorFieldType] {
	return &SearchCursor[SearchCursorFieldType]{
		Fields: fields,
	}
}

func (sc *SearchCursor[SearchCursorFieldType]) Encode() (*string, *exceptions.Exception) {
	jsonData, err := json.Marshal(sc.Fields)
	if err != nil {
		return nil, exceptions.Search.FailedToMarshalSearchCursor().WithError(err)
	}

	encoded := base64.StdEncoding.EncodeToString(jsonData)
	return &encoded, nil
}

func Decode[SearchCursorFieldType any](encoded string) (*SearchCursor[SearchCursorFieldType], *exceptions.Exception) {
	if len(strings.ReplaceAll(encoded, " ", "")) == 0 {
		return nil, exceptions.Search.EmptyEncodedStringToDecodeSearchCursor()
	}

	jsonData, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, exceptions.Search.FailedToDecodeBase64String().WithError(err)
	}

	var fields SearchCursorFieldType
	if err := json.Unmarshal(jsonData, &fields); err != nil {
		return nil, exceptions.Search.FailedToUnmarshalSearchCursor().WithError(err)
	}

	return &SearchCursor[SearchCursorFieldType]{Fields: fields}, nil
}

func EncodeFromData[SearchCursorType any](data SearchCursorType) (*string, *exceptions.Exception) {
	cursor := New(data)
	return cursor.Encode()
}
