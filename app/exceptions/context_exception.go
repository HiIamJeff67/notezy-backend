package exceptions

import (
	"fmt"
	"net/http"
)

const (
	_ExceptionBaseCode_Context ExceptionCode = ContextExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	ContextExceptionSubDomainCode ExceptionCode   = 4
	ExceptionBaseCode_Context     ExceptionCode   = _ExceptionBaseCode_Context + ReservedExceptionCode
	ExceptionPrefix_Context       ExceptionPrefix = "Context"
)

type ContextExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	APIExceptionDomain
}

var Context = &ContextExceptionDomain{
	BaseCode: ExceptionBaseCode_Context,
	Prefix:   ExceptionPrefix_Context,
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Context,
		_Prefix:   ExceptionPrefix_Context,
	},
}

func (d *ContextExceptionDomain) FailedToFetchContextFieldOfSpecificName(name string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("Failed to find and fetch the context field with name of %s since it is not exist in the current context", name),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *ContextExceptionDomain) FailedToConvertContextFieldToSpecificType(typeName string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("Failed to convert context field from type of any to type of %s", typeName),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
