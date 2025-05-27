package exceptions

const (
	_ExceptionBaseCode_UsersToBadges ExceptionCode = 500000
	ExceptionBaseCode_UsersToBadges ExceptionCode = _ExceptionBaseCode_UsersToBadges + ReservedExceptionCode
	ExceptionPrefix_UsersToBadges ExceptionPrefix = "UsersToBadges"
)

type UsersToBadgesExceptionDomain struct {
	APIExceptionDomain
}

var UsersToBadges = &UsersToBadgesExceptionDomain{
	APIExceptionDomain{ BaseCode: _ExceptionBaseCode_UsersToBadges, Prefix: ExceptionPrefix_UsersToBadges },
}