package exceptions

const (
	_ExceptionBaseCode_UsersToBadges ExceptionCode = (DatabaseExceptionDomainCode*ExceptionDomainCodeShiftAmount +
		UsersToBadgesExceptionSubDomainCode*ExceptionSubDomainCodeShiftAmount)

	UsersToBadgesExceptionSubDomainCode ExceptionCode   = 5
	ExceptionBaseCode_UsersToBadges     ExceptionCode   = _ExceptionBaseCode_UsersToBadges + ReservedExceptionCode
	ExceptionPrefix_UsersToBadges       ExceptionPrefix = "UsersToBadges"
)

type UsersToBadgesExceptionDomain struct {
	DatabaseExceptionDomain
}

var UsersToBadges = &UsersToBadgesExceptionDomain{
	DatabaseExceptionDomain{BaseCode: _ExceptionBaseCode_UsersToBadges, Prefix: ExceptionPrefix_UsersToBadges},
}
