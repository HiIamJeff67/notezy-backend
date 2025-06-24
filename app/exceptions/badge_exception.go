package exceptions

const (
	_ExceptionBaseCode_Badge ExceptionCode = BadgeExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	BadgeExceptionSubDomainCode ExceptionCode   = 37
	ExceptionBaseCode_Badge     ExceptionCode   = _ExceptionBaseCode_Badge + ReservedExceptionCode
	ExceptionPrefix_Badge       ExceptionPrefix = "Badge"
)

type BadgeExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
	APIExceptionDomain
	TypeExceptionDomain
	CommonExceptionDomain
}

var Badge = &BadgeExceptionDomain{
	BaseCode: ExceptionBaseCode_Badge,
	Prefix:   ExceptionPrefix_Badge,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Badge,
		_Prefix:   ExceptionPrefix_Badge,
	},
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Badge,
		_Prefix:   ExceptionPrefix_Badge,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Badge,
		_Prefix:   ExceptionPrefix_Badge,
	},
	CommonExceptionDomain: CommonExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Badge,
		_Prefix:   ExceptionPrefix_Badge,
	},
}
