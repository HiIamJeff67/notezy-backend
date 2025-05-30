package exceptions

const (
	_ExceptionBaseCode_UserAccount ExceptionCode = (DatabaseExceptionDomainCode*ExceptionDomainCodeShiftAmount +
		UserAccountExceptionSubDomainCode*ExceptionSubDomainCodeShiftAmount)

	UserAccountExceptionSubDomainCode ExceptionCode   = 3
	ExceptionBaseCode_UserAccount     ExceptionCode   = _ExceptionBaseCode_UserAccount + ReservedExceptionCode
	ExceptionPrefix_UserAccount       ExceptionPrefix = "UserAccount"
)

type UserAccountExceptionDomain struct {
	DatabaseExceptionDomain
}

var UserAccount = &UserAccountExceptionDomain{
	DatabaseExceptionDomain{BaseCode: _ExceptionBaseCode_UserAccount, Prefix: ExceptionPrefix_UserAccount},
}
