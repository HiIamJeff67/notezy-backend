package exceptions

const (
	_ExceptionBaseCode_UserInfo ExceptionCode = UserInfoExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	UserInfoExceptionSubDomainCode ExceptionCode   = 33
	ExceptionBaseCode_UserInfo     ExceptionCode   = _ExceptionBaseCode_UserInfo + ReservedExceptionCode
	ExceptionPrefix_UserInfo       ExceptionPrefix = "UserInfo"
)

type UserInfoExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
	APIExceptionDomain
	GraphQLExceptionDomain
	TypeExceptionDomain
	CommonExceptionDomain
}

var UserInfo = &UserInfoExceptionDomain{
	BaseCode: ExceptionBaseCode_UserInfo,
	Prefix:   ExceptionPrefix_UserInfo,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UserInfo,
		_Prefix:   ExceptionPrefix_UserInfo,
	},
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UserInfo,
		_Prefix:   ExceptionPrefix_UserInfo,
	},
	GraphQLExceptionDomain: GraphQLExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UserInfo,
		_Prefix:   ExceptionPrefix_UserInfo,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UserInfo,
		_Prefix:   ExceptionPrefix_UserInfo,
	},
	CommonExceptionDomain: CommonExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UserInfo,
		_Prefix:   ExceptionPrefix_UserInfo,
	},
}
