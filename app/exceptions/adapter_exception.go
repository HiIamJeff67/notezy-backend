package exceptions

import (
	"fmt"
	"net/http"
)

const (
	_ExceptionBaseCode_Adapter ExceptionCode = AdapterExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	AdapterExceptionSubDomainCode ExceptionCode   = 9
	ExceptionBaseCode_Adapter     ExceptionCode   = _ExceptionBaseCode_Adapter + ReservedExceptionCode
	ExceptionPrefix_Adapter       ExceptionPrefix = "Adapter"
)

type AdapterExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
}

var Adapter = &AdapterExceptionDomain{
	BaseCode: ExceptionBaseCode_Adapter,
	Prefix:   ExceptionPrefix_Adapter,
}

/* ============================== Handling Multipart Adapter Errors ============================== */

func (d *AdapterExceptionDomain) InvalidMultipartForm() *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "InvalidMultipartForm",
		IsInternal:     false,
		Message:        "The multipart form in the context is missing or invalid",
		HTTPStatusCode: http.StatusForbidden,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AdapterExceptionDomain) FileTooLarge(size int64, maxSize int64) *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "FileTooLarge",
		IsInternal:     false,
		Message:        fmt.Sprintf("The size of the file in multipart form data is %d which is larger than the limit of %d", size, maxSize),
		HTTPStatusCode: http.StatusRequestEntityTooLarge,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
