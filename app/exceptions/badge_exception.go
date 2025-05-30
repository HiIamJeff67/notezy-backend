package exceptions

const (
	_ExceptionBaseCode_Badge ExceptionCode = (DatabaseExceptionDomainCode*ExceptionDomainCodeShiftAmount +
		BadgeExceptionSubDomainCode*ExceptionSubDomainCodeShiftAmount)

	BadgeExceptionSubDomainCode ExceptionCode   = 6
	ExceptionBaseCode_Badge     ExceptionCode   = _ExceptionBaseCode_Badge + ReservedExceptionCode
	ExceptionPrefix_Badge       ExceptionPrefix = "Badge"
)

type BadgeExceptionDomain struct {
	DatabaseExceptionDomain
}

var Badge = &BadgeExceptionDomain{
	DatabaseExceptionDomain{BaseCode: _ExceptionBaseCode_Badge, Prefix: ExceptionPrefix_Badge},
}
