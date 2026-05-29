package exceptions

const (
	_ExceptionBaseCode_RoutineTag ExceptionCode = RoutineTagExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	RoutineTagExceptionSubDomainCode ExceptionCode   = 50
	ExceptionBaseCode_RoutineTag     ExceptionCode   = _ExceptionBaseCode_RoutineTag + ReservedExceptionCode
	ExceptionPrefix_RoutineTag       ExceptionPrefix = "RoutineTag"
)

type RoutineTagExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	APIExceptionDomain
	DatabaseExceptionDomain
	TypeExceptionDomain
	FileExceptionDomain
}

var RoutineTag = &RoutineTagExceptionDomain{
	BaseCode: ExceptionBaseCode_RoutineTag,
	Prefix:   ExceptionPrefix_RoutineTag,
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_RoutineTag,
		_Prefix:   ExceptionPrefix_RoutineTag,
	},
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_RoutineTag,
		_Prefix:   ExceptionPrefix_RoutineTag,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_RoutineTag,
		_Prefix:   ExceptionPrefix_RoutineTag,
	},
	FileExceptionDomain: FileExceptionDomain{
		_BaseCode: _ExceptionBaseCode_RoutineTag,
		_Prefix:   ExceptionPrefix_RoutineTag,
	},
}
