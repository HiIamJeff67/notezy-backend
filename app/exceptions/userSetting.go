package exceptions

const (
	_ExceptionBaseCode_UserSetting ExceptionCode = 400000
	ExceptionBaseCode_UserSetting ExceptionCode = _ExceptionBaseCode_UserSetting + ReservedExceptionCode
	ExceptionPrefix_UserSetting ExceptionPrefix = "UserSetting"
)

type UserSettingExceptionDomain struct {
	APIExceptionDomain
}

var UserSetting = &UserSettingExceptionDomain{
	APIExceptionDomain{ BaseCode: _ExceptionBaseCode_UserSetting, Prefix: ExceptionPrefix_UserSetting },
}