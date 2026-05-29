package exceptions

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
