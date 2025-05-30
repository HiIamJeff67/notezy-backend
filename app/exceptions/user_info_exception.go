package exceptions

const (
	_ExceptionBaseCode_UserInfo ExceptionCode = (DatabaseExceptionDomainCode*ExceptionDomainCodeShiftAmount +
		UserInfoExceptionSubDomainCode*ExceptionSubDomainCodeShiftAmount)

	UserInfoExceptionSubDomainCode ExceptionCode   = 2
	ExceptionBaseCode_UserInfo     ExceptionCode   = _ExceptionBaseCode_UserInfo + ReservedExceptionCode
	ExceptionPrefix_UserInfo       ExceptionPrefix = "UserInfo"
)

type UserInfoExceptionDomain struct {
	DatabaseExceptionDomain
}

var UserInfo = &UserInfoExceptionDomain{
	DatabaseExceptionDomain{BaseCode: _ExceptionBaseCode_UserInfo, Prefix: ExceptionPrefix_UserInfo},
}
