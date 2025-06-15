package exceptions

import (
	"fmt"
	"net/http"
)

const (
	_ExceptionBaseCode_Context ExceptionCode = (APIExceptionDomainCode*ExceptionDomainCodeShiftAmount +
		ContextExceptionSubDomainCode*ExceptionSubDomainCodeShiftAmount)

	ContextExceptionSubDomainCode ExceptionCode   = 5
	ExceptionBaseCode_Context     ExceptionCode   = _ExceptionBaseCode_Context + ReservedExceptionCode
	ExceptionPrefix_Context       ExceptionPrefix = "Context"
)

const (
	ExceptionReason_FailedToFetchContextFieldOfSpecificName   ExceptionReason = "Failed_To_Fetch_Context_Field_Of_Specific_Name"
	ExceptionReason_FailedToConvertContextFieldToSpecificType ExceptionReason = "Failed_To_Convert_Context_Field_To_Specific_Type"
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
		_Prefix:   ExceptionPrefix_Context},
}

func (d *ContextExceptionDomain) FailedToFetchContextFieldOfSpecificName(name string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToFetchContextFieldOfSpecificName,
		Message:        fmt.Sprintf("Failed to find and fetch the context field with name of %s since it is not exist in the current context", name),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *ContextExceptionDomain) FailedToConvertContextFieldToSpecificType(typeName string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToConvertContextFieldToSpecificType,
		Message:        fmt.Sprintf("Failed to convert context field from type of any to type of %s", typeName),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}
