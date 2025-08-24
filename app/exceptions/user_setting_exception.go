package exceptions

const (
	_ExceptionBaseCode_UserSetting ExceptionCode = UserSettingExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	UserSettingExceptionSubDomainCode ExceptionCode   = 35
	ExceptionBaseCode_UserSetting     ExceptionCode   = _ExceptionBaseCode_UserSetting + ReservedExceptionCode
	ExceptionPrefix_UserSetting       ExceptionPrefix = "UserSetting"
)

type UserSettingExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
	APIExceptionDomain
	TypeExceptionDomain
}

var UserSetting = &UserSettingExceptionDomain{
	BaseCode: ExceptionBaseCode_UserSetting,
	Prefix:   ExceptionPrefix_UserSetting,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UserSetting,
		_Prefix:   ExceptionPrefix_UserSetting,
	},
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UserSetting,
		_Prefix:   ExceptionPrefix_UserSetting,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UserSetting,
		_Prefix:   ExceptionPrefix_UserSetting,
	},
}
