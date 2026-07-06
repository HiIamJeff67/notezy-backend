package exceptions

import (
	"fmt"
	"net/http"
	"time"

	"github.com/HiIamJeff67/notezy-backend/app/monitor/traces"
)

const (
	_ExceptionBaseCode_Routine ExceptionCode = RoutineExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	RoutineExceptionSubDomainCode ExceptionCode   = 49
	ExceptionBaseCode_Routine     ExceptionCode   = _ExceptionBaseCode_Routine + ReservedExceptionCode
	ExceptionPrefix_Routine       ExceptionPrefix = "Routine"
)

type RoutineExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	APIExceptionDomain
	DatabaseExceptionDomain
	TypeExceptionDomain
	FileExceptionDomain
	MarshalerExceptionDomain
}

var Routine = &RoutineExceptionDomain{
	BaseCode: ExceptionBaseCode_Routine,
	Prefix:   ExceptionPrefix_Routine,
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Routine,
		_Prefix:   ExceptionPrefix_Routine,
	},
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Routine,
		_Prefix:   ExceptionPrefix_Routine,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Routine,
		_Prefix:   ExceptionPrefix_Routine,
	},
	FileExceptionDomain: FileExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Routine,
		_Prefix:   ExceptionPrefix_Routine,
	},
	MarshalerExceptionDomain: MarshalerExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Routine,
		_Prefix:   ExceptionPrefix_Routine,
	},
}

/* ============================== Handling Time Range Validation ============================== */

func (d *RoutineExceptionDomain) QueriedTimeRangeTooLarge(from time.Time, to time.Time) *Exception {
	return &Exception{
		Code:           d.BaseCode + 11,
		Prefix:         d.Prefix,
		Reason:         "QueriedTimeRangeTooLarge",
		IsInternal:     false,
		Message:        fmt.Sprintf("Cannot querying with time range from %s to %s, which is too large", from.String(), to.String()),
		HTTPStatusCode: http.StatusBadRequest,
		LastTrace:      traces.GetTrace(1),
	}
}

/* ============================== Handling Linking Error of Routines ============================== */

func (d *RoutineExceptionDomain) FailedToLinkRoutineTags() *Exception {
	return &Exception{
		Code:           d.BaseCode + 21,
		Prefix:         d.Prefix,
		Reason:         "FailedToLinkRoutineTags",
		IsInternal:     false,
		Message:        "Cannot link the given routine tags to the target routine",
		HTTPStatusCode: http.StatusBadRequest,
		LastTrace:      traces.GetTrace(1),
	}
}

func (d *RoutineExceptionDomain) FailedToLinkItems() *Exception {
	return &Exception{
		Code:           d.BaseCode + 23,
		Prefix:         d.Prefix,
		Reason:         "FailedToLinkItems",
		IsInternal:     false,
		Message:        "Cannot link the given items to the target routine",
		HTTPStatusCode: http.StatusBadRequest,
		LastTrace:      traces.GetTrace(1),
	}
}
