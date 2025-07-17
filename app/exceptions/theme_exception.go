package exceptions

const (
	_ExceptionBaseCode_Theme ExceptionCode = ThemeExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	ThemeExceptionSubDomainCode ExceptionCode   = 38
	ExceptionBaseCode_Theme     ExceptionCode   = _ExceptionBaseCode_Theme + ReservedExceptionCode
	ExceptionPrefix_Theme       ExceptionPrefix = "Theme"
)

type ThemeExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
	APIExceptionDomain
	GraphQLExceptionDomain
	TypeExceptionDomain
	CommonExceptionDomain
}

var Theme = &ThemeExceptionDomain{
	BaseCode: ExceptionBaseCode_Theme,
	Prefix:   ExceptionPrefix_Theme,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Theme,
		_Prefix:   ExceptionPrefix_Theme,
	},
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Theme,
		_Prefix:   ExceptionPrefix_Theme,
	},
	GraphQLExceptionDomain: GraphQLExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Theme,
		_Prefix:   ExceptionPrefix_Theme,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Theme,
		_Prefix:   ExceptionPrefix_Theme,
	},
	CommonExceptionDomain: CommonExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Theme,
		_Prefix:   ExceptionPrefix_Theme,
	},
}
