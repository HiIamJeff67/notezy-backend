package exceptions

const (
	// define this bcs the general domain has some general exception that has be defined

	_ExceptionBaseCode_User ExceptionCode = (DatabaseExceptionDomainCode*ExceptionDomainCodeShiftAmount +
		UserExceptionSubDomainCode*ExceptionSubDomainCodeShiftAmount) // the actual base for the exceptions of user

	// if you need to code a custom exception of users,
	// use the ExceptionBaseCode_User, instead of _ExceptionBaseCode_User
	// the exception codes that we can actually customize here is shifted with ReservedExceptionCode

	UserExceptionSubDomainCode ExceptionCode   = 1
	ExceptionBaseCode_User     ExceptionCode   = _ExceptionBaseCode_User + ReservedExceptionCode
	ExceptionPrefix_User       ExceptionPrefix = "User"
)

type UserExceptionDomain struct {
	// as the down layer of DatabaseExceptionDomain
	// so that we don't make methods for DatabaseExceptionDomain
	// instead we make methods for UserExceptionDomain
	DatabaseExceptionDomain
}

var User = &UserExceptionDomain{
	DatabaseExceptionDomain{BaseCode: _ExceptionBaseCode_User, Prefix: ExceptionPrefix_User},
}
