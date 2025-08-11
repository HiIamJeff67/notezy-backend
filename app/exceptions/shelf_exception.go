package exceptions

import (
	"fmt"
	"net/http"
)

const (
	_ExceptionBaseCode_Shelf ExceptionCode = ShelfExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	ShelfExceptionSubDomainCode ExceptionCode   = 40
	ExceptionBaseCode_Shelf     ExceptionCode   = _ExceptionBaseCode_Shelf + ReservedExceptionCode
	ExceptionPrefix_Shelf       ExceptionPrefix = "Shelf"
)

type ShelfExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
	TypeExceptionDomain
}

var Shelf = &ShelfExceptionDomain{
	BaseCode: ExceptionBaseCode_Shelf,
	Prefix:   ExceptionPrefix_Shelf,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Shelf,
		_Prefix:   ExceptionPrefix_Shelf,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Shelf,
		_Prefix:   ExceptionPrefix_Shelf,
	},
}

/* ============================== Handling Structure Error of ShelfNode ============================== */

func (d *ShelfExceptionDomain) CallingMethodsWithNilValue() *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "CallingMethodWithNilValue",
		IsInternal:     true,
		Message:        "Nil value cannot call the methods of ShelfNode",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *ShelfExceptionDomain) FailedToEncode(node interface{}) *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         "FailedToEncode",
		IsInternal:     true,
		Message:        fmt.Sprintf("Failed to encode %v", node),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *ShelfExceptionDomain) FailedToDecode(data []byte) *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         "FailedToDecode",
		IsInternal:     true,
		Message:        fmt.Sprintf("Failed to decode encoded string of %v", data),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *ShelfExceptionDomain) InsertParentIntoItsChildren(destination interface{}, target interface{}) *Exception {
	return &Exception{
		Code:       d.BaseCode + 4,
		Prefix:     d.Prefix,
		Reason:     "InsertParentIntoItsChildren",
		IsInternal: false,
		Message: fmt.Sprintf(
			"Failed to insert %v into %v since %v is one of the child of %v, insert a parent node into its children is not allowed",
			target, destination, destination, target,
		),
		HTTPStatusCode: http.StatusConflict,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
