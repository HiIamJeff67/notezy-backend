package exceptions

const (
	_ExceptionBaseCode_UserInfo ExceptionCode = (DatabaseExceptionDomainCode*ExceptionDomainCodeShiftAmount +
		UserInfoExceptionSubDomainCode*ExceptionSubDomainCodeShiftAmount)

	UserInfoExceptionSubDomainCode ExceptionCode   = 2
	ExceptionBaseCode_UserInfo     ExceptionCode   = _ExceptionBaseCode_UserInfo + ReservedExceptionCode
	ExceptionPrefix_UserInfo       ExceptionPrefix = "UserInfo"
)

type UserInfoExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
}

var UserInfo = &UserInfoExceptionDomain{
	BaseCode: ExceptionBaseCode_UserInfo,
	Prefix:   ExceptionPrefix_UserInfo,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UserInfo,
		_Prefix:   ExceptionPrefix_UserInfo,
	},
}
