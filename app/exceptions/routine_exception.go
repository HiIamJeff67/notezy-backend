package exceptions

import (
	"net/http"

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
}

/* ============================== Handling Linking Error of Routines ============================== */

func (d *RoutineExceptionDomain) FailedToLinkRoutineTags() *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "FailedToLinkRoutineTags",
		IsInternal:     true,
		Message:        "Cannot link the given routine tags to the target routine",
		HTTPStatusCode: http.StatusInternalServerError,
		LastTrace:      traces.GetTrace(1),
	}
}

func (d *RoutineExceptionDomain) FailedToLinkRoutineTasks() *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         "FailedToLinkRoutineTasks",
		IsInternal:     true,
		Message:        "Cannot link the given routine tasks to the target routine",
		HTTPStatusCode: http.StatusInternalServerError,
		LastTrace:      traces.GetTrace(1),
	}
}

func (d *RoutineExceptionDomain) FailedToLinkItems() *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         "FailedToLinkItems",
		IsInternal:     true,
		Message:        "Cannot link the given items to the target routine",
		HTTPStatusCode: http.StatusInternalServerError,
		LastTrace:      traces.GetTrace(1),
	}
}
