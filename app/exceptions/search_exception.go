package exceptions

import (
	"fmt"
	"net/http"
)

const (
	_ExceptionBaseCode_Search ExceptionCode = SearchExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	SearchExceptionSubDomainCode ExceptionCode   = 7
	ExceptionBaseCode_Search     ExceptionCode   = _ExceptionBaseCode_Search + ReservedExceptionCode
	ExceptionPrefix_Search       ExceptionPrefix = "Search"
)

type SearchExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
	APIExceptionDomain
}

var Search = &SearchExceptionDomain{
	BaseCode: ExceptionBaseCode_Search,
	Prefix:   ExceptionPrefix_Search,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Search,
		_Prefix:   ExceptionPrefix_Search,
	},
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Email,
		_Prefix:   ExceptionPrefix_Email,
	},
}

func (d *SearchExceptionDomain) InvalidNilDataToEncodeSearchCursor() *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Reason:         "InvalidNilDataToEncodeSearchCursor",
		Prefix:         d.Prefix,
		Message:        "Invalid nil data to encode search cursor, data must be not nil",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *SearchExceptionDomain) InvalidNonMapToEncodeSearchCursor() *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Reason:         "InvalidNonMapToEncodeSearchCursor",
		Prefix:         d.Prefix,
		Message:        "Invalid non map data to encode search cursor, data must be map[string]interface{}",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *SearchExceptionDomain) FailedToMarshalSearchCursor() *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Reason:         "FailedToMarshalSearchCursor",
		Prefix:         d.Prefix,
		Message:        "Failed to marshal the search cursor",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *SearchExceptionDomain) FailedToUnmarshalSearchCursor() *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Reason:         "FailedToUnmarshalSearchCursor",
		Prefix:         d.Prefix,
		Message:        "Failed to unmarshal the search cursor",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *SearchExceptionDomain) EmptyEncodedStringToDecodeSearchCursor() *Exception {
	return &Exception{
		Code:           d.BaseCode + 5,
		Reason:         "EmptyEncodedStringToDecodeSearchCursor",
		Prefix:         d.Prefix,
		Message:        "Encoded string cannot be empty",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *SearchExceptionDomain) FailedToDecodeBase64String() *Exception {
	return &Exception{
		Code:           d.BaseCode + 6,
		Reason:         "FailedToDecodeBase64String",
		Prefix:         d.Prefix,
		Message:        "Failed to decode base64 string",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *SearchExceptionDomain) CannotFindFieldInEncodedSearchCursor(searchCursor string, fieldName string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 7,
		Reason:         "CannotFindFieldInEncodedSearchCursor",
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("Cannot find the field of %s in the search cursor: %s", fieldName, searchCursor),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
