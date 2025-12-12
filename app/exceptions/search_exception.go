package exceptions

import (
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

func (d *SearchExceptionDomain) FailedToDecode() *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "FailedToDecode",
		IsInternal:     true,
		Message:        "Failed to decode into search cursor",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *SearchExceptionDomain) FailedToEncode() *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         "FailedToEncode",
		IsInternal:     true,
		Message:        "Failed to encode into encoded search cursor",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *SearchExceptionDomain) FailedToMarshalSearchCursor() *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         "FailedToMarshalSearchCursor",
		IsInternal:     true,
		Message:        "Failed to marshal the search cursor",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *SearchExceptionDomain) FailedToUnmarshalSearchCursor() *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Reason:         "FailedToUnmarshalSearchCursor",
		IsInternal:     true,
		Message:        "Failed to unmarshal the search cursor",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
