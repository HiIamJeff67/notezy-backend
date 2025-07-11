package exceptions

import "net/http"

const (
	_ExceptionBaseCode_Searchable ExceptionCode = SearchableExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	SearchableExceptionSubDomainCode ExceptionCode   = 7
	ExceptionBaseCode_Searchable     ExceptionCode   = _ExceptionBaseCode_Searchable + ReservedExceptionCode
	ExceptionPrefix_Searchable       ExceptionPrefix = "Searchable"
)

type SearchableExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
	APIExceptionDomain
}

var Searchable = &SearchableExceptionDomain{
	BaseCode: ExceptionBaseCode_Searchable,
	Prefix:   ExceptionPrefix_Searchable,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Searchable,
		_Prefix:   ExceptionPrefix_Searchable,
	},
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Email,
		_Prefix:   ExceptionPrefix_Email,
	},
}

func (d *SearchableExceptionDomain) InvalidNilDataToEncodeSearchCursor() *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Message:        "Invalid nil data to encode search cursor, data must be not nil",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *SearchableExceptionDomain) InvalidNonMapToEncodeSearchCursor() *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Message:        "Invalid non map data to encode search cursor, data must be map[string]interface{}",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *SearchableExceptionDomain) FailedToMarshalSearchCursor() *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Message:        "Failed to marshal the search cursor",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *SearchableExceptionDomain) FailedToUnMarshalSearchCursor() *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Message:        "Failed to unmarshal the search cursor",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *SearchableExceptionDomain) EmptyEncodedStringToDecodeSearchCursor() *Exception {
	return &Exception{
		Code:           d.BaseCode + 5,
		Prefix:         d.Prefix,
		Message:        "Encoded string cannot be empty",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *SearchableExceptionDomain) FailedToDecodeBase64String() *Exception {
	return &Exception{
		Code:           d.BaseCode + 6,
		Prefix:         d.Prefix,
		Message:        "Failed to decode base64 string",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}
