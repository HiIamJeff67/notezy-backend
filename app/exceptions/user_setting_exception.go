package exceptions

const (
	_ExceptionBaseCode_UserSetting ExceptionCode = (DatabaseExceptionDomainCode*ExceptionDomainCodeShiftAmount +
		UserSettingExceptionSubDomainCode*ExceptionSubDomainCodeShiftAmount)

	UserSettingExceptionSubDomainCode ExceptionCode   = 4
	ExceptionBaseCode_UserSetting     ExceptionCode   = _ExceptionBaseCode_UserSetting + ReservedExceptionCode
	ExceptionPrefix_UserSetting       ExceptionPrefix = "UserSetting"
)

type UserSettingExceptionDomain struct {
	DatabaseExceptionDomain
}

var UserSetting = &UserSettingExceptionDomain{
	DatabaseExceptionDomain{BaseCode: _ExceptionBaseCode_UserSetting, Prefix: ExceptionPrefix_UserSetting},
}
