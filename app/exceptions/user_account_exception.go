package exceptions

const (
	_ExceptionBaseCode_UserAccount ExceptionCode = UserAccountExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	UserAccountExceptionSubDomainCode ExceptionCode   = 34
	ExceptionBaseCode_UserAccount     ExceptionCode   = _ExceptionBaseCode_UserAccount + ReservedExceptionCode
	ExceptionPrefix_UserAccount       ExceptionPrefix = "UserAccount"
)

type UserAccountExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
	APIExceptionDomain
	TypeExceptionDomain
}

var UserAccount = &UserAccountExceptionDomain{
	BaseCode: ExceptionBaseCode_UserAccount,
	Prefix:   ExceptionPrefix_UserAccount,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UserAccount,
		_Prefix:   ExceptionPrefix_UserAccount,
	},
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UserAccount,
		_Prefix:   ExceptionPrefix_UserAccount,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UserAccount,
		_Prefix:   ExceptionPrefix_UserAccount,
	},
}
