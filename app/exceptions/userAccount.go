package exceptions

const (
	_ExceptionBaseCode_UserAccount ExceptionCode = 300000
	ExceptionBaseCode_UserAccount ExceptionCode = _ExceptionBaseCode_UserAccount + ReservedExceptionCode
	ExceptionPrefix_UserAccount ExceptionPrefix = "UserAccount"
)

type UserAccountExceptionDomain struct {
	APIExceptionDomain
}

var UserAccount = &UserAccountExceptionDomain{
	APIExceptionDomain{ BaseCode: _ExceptionBaseCode_UserAccount, Prefix: ExceptionPrefix_UserAccount },
}