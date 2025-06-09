package exceptions

const (
	_ExceptionBaseCode_Theme ExceptionCode = (DatabaseExceptionDomainCode*ExceptionDomainCodeShiftAmount +
		ThemeExceptionSubDomainCode*ExceptionSubDomainCodeShiftAmount)

	ThemeExceptionSubDomainCode ExceptionCode   = 7
	ExceptionBaseCode_Theme     ExceptionCode   = _ExceptionBaseCode_Theme + ReservedExceptionCode
	ExceptionPrefix_Theme       ExceptionPrefix = "Theme"
)

type ThemeExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
}

var Theme = &ThemeExceptionDomain{
	BaseCode: ExceptionBaseCode_Theme,
	Prefix:   ExceptionPrefix_Theme,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Theme,
		_Prefix:   ExceptionPrefix_Theme},
}
