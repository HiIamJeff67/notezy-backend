package exceptions

const (
	_ExceptionBaseCode_Badge ExceptionCode = (DatabaseExceptionDomainCode*ExceptionDomainCodeShiftAmount +
		BadgeExceptionSubDomainCode*ExceptionSubDomainCodeShiftAmount)

	BadgeExceptionSubDomainCode ExceptionCode   = 6
	ExceptionBaseCode_Badge     ExceptionCode   = _ExceptionBaseCode_Badge + ReservedExceptionCode
	ExceptionPrefix_Badge       ExceptionPrefix = "Badge"
)

type BadgeExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
}

var Badge = &BadgeExceptionDomain{
	BaseCode: ExceptionBaseCode_Badge,
	Prefix:   ExceptionPrefix_Badge,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Badge,
		_Prefix:   ExceptionPrefix_Badge,
	},
}
