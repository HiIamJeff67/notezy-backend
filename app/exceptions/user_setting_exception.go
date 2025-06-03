package exceptions

const (
	_ExceptionBaseCode_UserSetting ExceptionCode = (DatabaseExceptionDomainCode*ExceptionDomainCodeShiftAmount +
		UserSettingExceptionSubDomainCode*ExceptionSubDomainCodeShiftAmount)

	UserSettingExceptionSubDomainCode ExceptionCode   = 4
	ExceptionBaseCode_UserSetting     ExceptionCode   = _ExceptionBaseCode_UserSetting + ReservedExceptionCode
	ExceptionPrefix_UserSetting       ExceptionPrefix = "UserSetting"
)

type UserSettingExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
}

var UserSetting = &UserSettingExceptionDomain{
	BaseCode: ExceptionBaseCode_UserSetting,
	Prefix:   ExceptionPrefix_UserSetting,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UserSetting,
		_Prefix:   ExceptionPrefix_UserSetting,
	},
}
