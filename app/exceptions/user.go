package exceptions

const (
	// bcs the general domain has some general exception that has be defined
	_ExceptionBaseCode_User ExceptionCode = 100000	// the actual base for the exceptions of user
	// if you need to code a custom exception of users, 
	// use the ExceptionBaseCode_User, instead of _ExceptionBaseCode_User
	// the exception codes that we can actually customize here is shifted with ReservedExceptionCode
	ExceptionBaseCode_User ExceptionCode = _ExceptionBaseCode_User + ReservedExceptionCode
	ExceptionPrefix_User ExceptionPrefix = "User"
)

type UserExceptionDomain struct { // as the down layer of APIExceptionDomain
	// so that we don't make methods for APIExceptionDomain
	// instead we make methods for UserExceptionDomain
	APIExceptionDomain
}

var User = &UserExceptionDomain{
	APIExceptionDomain{ BaseCode: _ExceptionBaseCode_User, Prefix: ExceptionPrefix_User }, 
}
