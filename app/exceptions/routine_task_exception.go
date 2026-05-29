package exceptions

const (
	_ExceptionBaseCode_RoutineTask ExceptionCode = RoutineTaskExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	RoutineTaskExceptionSubDomainCode ExceptionCode   = 51
	ExceptionBaseCode_RoutineTask     ExceptionCode   = _ExceptionBaseCode_RoutineTask + ReservedExceptionCode
	ExceptionPrefix_RoutineTask       ExceptionPrefix = "RoutineTask"
)

type RoutineTaskExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	APIExceptionDomain
	DatabaseExceptionDomain
	TypeExceptionDomain
	FileExceptionDomain
}

var RoutineTask = &RoutineTaskExceptionDomain{
	BaseCode: ExceptionBaseCode_RoutineTask,
	Prefix:   ExceptionPrefix_RoutineTask,
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_RoutineTask,
		_Prefix:   ExceptionPrefix_RoutineTask,
	},
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_RoutineTask,
		_Prefix:   ExceptionPrefix_RoutineTask,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_RoutineTask,
		_Prefix:   ExceptionPrefix_RoutineTask,
	},
	FileExceptionDomain: FileExceptionDomain{
		_BaseCode: _ExceptionBaseCode_RoutineTask,
		_Prefix:   ExceptionPrefix_RoutineTask,
	},
}
