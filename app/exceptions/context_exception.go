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

func (d *ContextExceptionDomain) FailedToGetContextFieldOfSpecificName(name string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "FailedToGetContextFieldOfSpecificName",
		IsInternal:     true,
		Message:        fmt.Sprintf("Failed to find and fetch the context field with name of %s since it is not exist in the current context", name),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *ContextExceptionDomain) FailedToConvertContextFieldToSpecificType(typeName string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         "FailedToConvertContextFieldToSpecificType",
		IsInternal:     true,
		Message:        fmt.Sprintf("Failed to convert context field from type of any to type of %s", typeName),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *ContextExceptionDomain) FailedToGetCorrectContextValue(v interface{}) *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         "FailedToGetCorrectContextValue",
		IsInternal:     true,
		Message:        fmt.Sprintf("Failed to get correct context value, got %v instead", v),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *ContextExceptionDomain) FailedToConvertContextToGinContext() *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Reason:         "FailedToConvertContextToGinContext",
		IsInternal:     true,
		Message:        "Failed to convert from context.Context to gin.Context",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *ContextExceptionDomain) FailedToConvertGinContextToContext() *Exception {
	return &Exception{
		Code:           d.BaseCode + 5,
		Prefix:         d.Prefix,
		Reason:         "FailedToConvertGinContextToContext",
		IsInternal:     true,
		Message:        "Failed to convert from gin.Context to context.Context",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *ContextExceptionDomain) MaxContextBodySizeExceeded(sizeKiloBytes int64, maxSizeKiloByte int64) *Exception {
	return &Exception{
		Code:           d.BaseCode + 6,
		Prefix:         d.Prefix,
		Reason:         "MaxContextBodySizeExceeded",
		IsInternal:     true,
		Message:        fmt.Sprintf("The context body size of %d KB is larger than the maximum of %d KB", sizeKiloBytes, maxSizeKiloByte),
		HTTPStatusCode: http.StatusRequestEntityTooLarge,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
