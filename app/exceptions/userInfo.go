package exceptions

const (
	_ExceptionBaseCode_UserInfo ExceptionCode = 200000
	ExceptionBaseCode_UserInfo ExceptionCode = _ExceptionBaseCode_UserInfo + ReservedExceptionCode
	ExceptionPrefix_UserInfo ExceptionPrefix = "UserInfo"
)

type UserInfoExceptionDomain struct {
	APIExceptionDomain
}

var UserInfo = &UserInfoExceptionDomain{
	APIExceptionDomain{ BaseCode: _ExceptionBaseCode_UserInfo, Prefix: ExceptionPrefix_UserInfo },
}