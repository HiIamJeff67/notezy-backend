package exceptions

const (
	_ExceptionBaseCode_UsersToBadges ExceptionCode = UsersToBadgesExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	UsersToBadgesExceptionSubDomainCode ExceptionCode   = 36
	ExceptionBaseCode_UsersToBadges     ExceptionCode   = _ExceptionBaseCode_UsersToBadges + ReservedExceptionCode
	ExceptionPrefix_UsersToBadges       ExceptionPrefix = "UsersToBadges"
)

type UsersToBadgesExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
	APIExceptionDomain
	TypeExceptionDomain
}

var UsersToBadges = &UsersToBadgesExceptionDomain{
	BaseCode: ExceptionBaseCode_UsersToBadges,
	Prefix:   ExceptionPrefix_UsersToBadges,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UsersToBadges,
		_Prefix:   ExceptionPrefix_UsersToBadges,
	},
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UsersToBadges,
		_Prefix:   ExceptionPrefix_UsersToBadges,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UsersToBadges,
		_Prefix:   ExceptionPrefix_UsersToBadges,
	},
}
