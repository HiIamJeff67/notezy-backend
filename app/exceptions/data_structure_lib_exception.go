package exceptions

import (
	"net/http"
)

const (
	_ExceptionBaseCode_DataStructureLib ExceptionCode = DataStructureLibExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	DataStructureLibExceptionSubDomainCode ExceptionCode   = 11
	ExceptionBaseCode_DataStructureLib     ExceptionCode   = _ExceptionBaseCode_DataStructureLib + ReservedExceptionCode
	ExceptionPrefix_DataStructureLib       ExceptionPrefix = "DataStructureLib"
)

type DataStructureLibExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	TypeExceptionDomain
}

var DataStructureLib = &DataStructureLibExceptionDomain{
	BaseCode: ExceptionBaseCode_DataStructureLib,
	Prefix:   ExceptionPrefix_DataStructureLib,
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_DataStructureLib,
		_Prefix:   ExceptionPrefix_DataStructureLib,
	},
}

/* ============================== Queue Error In Service ============================== */

func (d *DataStructureLibExceptionDomain) FailedToManipulateQueue() *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "FailedToManipulateQueue",
		IsInternal:     true,
		Message:        "Failed to manipulate with the queue due to some reason or data structure error",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
