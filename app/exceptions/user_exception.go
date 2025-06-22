package exceptions

const (
	// define this bcs the general domain has some general exception that has be defined

	_ExceptionBaseCode_User ExceptionCode = UserExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount // the actual base for the exceptions of user

	// if you need to code a custom exception of users,
	// use the ExceptionBaseCode_User, instead of _ExceptionBaseCode_User
	// the exception codes that we can actually customize here is shifted with ReservedExceptionCode

	UserExceptionSubDomainCode ExceptionCode   = 32
	ExceptionBaseCode_User     ExceptionCode   = _ExceptionBaseCode_User + ReservedExceptionCode
	ExceptionPrefix_User       ExceptionPrefix = "User"
)

type UserExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	// as the down layer of DatabaseExceptionDomain
	// so that we don't make methods for DatabaseExceptionDomain
	// instead we make methods for UserExceptionDomain
	DatabaseExceptionDomain
	APIExceptionDomain
	TypeExceptionDomain
	CommonExceptionDomain
}

var User = &UserExceptionDomain{
	BaseCode: ExceptionBaseCode_User,
	Prefix:   ExceptionPrefix_User,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_User,
		_Prefix:   ExceptionPrefix_User,
	},
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_User,
		_Prefix:   ExceptionPrefix_User,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_User,
		_Prefix:   ExceptionPrefix_User,
	},
	CommonExceptionDomain: CommonExceptionDomain{
		_BaseCode: _ExceptionBaseCode_User,
		_Prefix:   ExceptionPrefix_User,
	},
}
