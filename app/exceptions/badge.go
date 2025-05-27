package exceptions

const (
	_ExceptionBaseCode_Badge ExceptionCode = 600000
	ExceptionBaseCode_Badge ExceptionCode = _ExceptionBaseCode_Badge + ReservedExceptionCode
	ExceptionPrefix_Badge ExceptionPrefix = "Badge"
)

type BadgeExceptionDomain struct {
	APIExceptionDomain
}

var Badge = &BadgeExceptionDomain{
	APIExceptionDomain{ BaseCode: _ExceptionBaseCode_Badge, Prefix: ExceptionPrefix_Badge },
}