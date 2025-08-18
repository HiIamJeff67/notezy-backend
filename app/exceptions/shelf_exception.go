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
	APIExceptionDomain
	DatabaseExceptionDomain
	TypeExceptionDomain
}

var Shelf = &ShelfExceptionDomain{
	BaseCode: ExceptionBaseCode_Shelf,
	Prefix:   ExceptionPrefix_Shelf,
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Shelf,
		_Prefix:   ExceptionPrefix_Shelf,
	},
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

func (d *ShelfExceptionDomain) MaximumWidthExceeded(currentWidth int, maxWidth int) *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "MaximumWidthExceeded",
		IsInternal:     false,
		Message:        fmt.Sprintf("The current width of %d is exceeded the limitation of %d", currentWidth, maxWidth),
		HTTPStatusCode: http.StatusBadRequest,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *ShelfExceptionDomain) MaximumDepthExceeded(currentDepth int, maxDepth int) *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         "MaximumDepthExceeded",
		IsInternal:     false,
		Message:        fmt.Sprintf("The current depth of %d is exceeded the limitation of %d", currentDepth, maxDepth),
		HTTPStatusCode: http.StatusBadRequest,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *ShelfExceptionDomain) MaximumTraverseCountExceeded(currentTraverseCount int, maxTraverseCount int) *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         "MaximumTraverseCountExceeded",
		IsInternal:     false,
		Message:        fmt.Sprintf("The current traverse count of %d is exceeded the limitation of %d", currentTraverseCount, maxTraverseCount),
		HTTPStatusCode: http.StatusBadRequest,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *ShelfExceptionDomain) CallingMethodsWithNilValue() *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Reason:         "CallingMethodWithNilValue",
		IsInternal:     true,
		Message:        "Nil value cannot call the methods of ShelfNode",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *ShelfExceptionDomain) FailedToEncode(node any) *Exception {
	return &Exception{
		Code:           d.BaseCode + 5,
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
		Code:           d.BaseCode + 6,
		Prefix:         d.Prefix,
		Reason:         "FailedToDecode",
		IsInternal:     true,
		Message:        fmt.Sprintf("Failed to decode encoded string of %v", data),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *ShelfExceptionDomain) InsertParentIntoItsChildren(destination any, target any) *Exception {
	return &Exception{
		Code:       d.BaseCode + 7,
		Prefix:     d.Prefix,
		Reason:     "InsertParentIntoItsChildren",
		IsInternal: false,
		Message: fmt.Sprintf(
			"Failed to insert %v into %v since %v is one of the child of %v, insert a parent node into its children is not allowed",
			target, destination, destination, target,
		),
		HTTPStatusCode: http.StatusBadRequest,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *ShelfExceptionDomain) FailedToConstructNewShelfNode(field string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 8,
		Prefix:         d.Prefix,
		Reason:         "FailedToConstructNewShelfNode",
		IsInternal:     false,
		Message:        fmt.Sprintf("The field of %s in ShelfNode is not pass by the validator", field),
		HTTPStatusCode: http.StatusBadRequest,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *ShelfExceptionDomain) CannotEncodeNonRootShelfNode(node any) *Exception {
	return &Exception{
		Code:           d.BaseCode + 9,
		Prefix:         d.Prefix,
		Reason:         "CannotEncodeNonRootShelfNode",
		IsInternal:     true,
		Message:        fmt.Sprintf("Cannot encoded the ShelfNode of %v which is not the root node", node),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *ShelfExceptionDomain) CircularChildrenDetectedInShelfNode() *Exception {
	return &Exception{
		Code:           d.BaseCode + 10,
		Prefix:         d.Prefix,
		Reason:         "CircularChildrenDetectedInShelfNode",
		IsInternal:     false,
		Message:        "Circular children detected in the given ShelfNode which is an invalid structure",
		HTTPStatusCode: http.StatusBadRequest,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *ShelfExceptionDomain) RepeatedShelfNodesDetected() *Exception {
	return &Exception{
		Code:           d.BaseCode + 11,
		Prefix:         d.Prefix,
		Reason:         "RepeatedShelfNodesDetected",
		IsInternal:     false,
		Message:        "Invalid ShelfNode structure with repeated shelf nodes detected in the same tree which is violating the uniqueness",
		HTTPStatusCode: http.StatusBadRequest,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *ShelfExceptionDomain) RepeatedMaterialIdsDetected() *Exception {
	return &Exception{
		Code:           d.BaseCode + 12,
		Prefix:         d.Prefix,
		Reason:         "RepeatedMaterialIdsDetected",
		IsInternal:     false,
		Message:        "Invalid ShelfNode structure with repeated material ids detected in the same tree which is violating the uniqueness",
		HTTPStatusCode: http.StatusBadRequest,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
