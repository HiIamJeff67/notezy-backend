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
	DurableJobExceptionDomain
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
	DurableJobExceptionDomain: DurableJobExceptionDomain{
		_BaseCode: _ExceptionBaseCode_RoutineTask,
		_Prefix:   ExceptionPrefix_RoutineTask,
	},
}
