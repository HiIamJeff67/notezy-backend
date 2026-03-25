package exceptions

import (
	"net/http"
	traces "notezy-backend/app/traces"
)

const (
	_ExceptionBaseCode_UserAccount ExceptionCode = UserAccountExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	UserAccountExceptionSubDomainCode ExceptionCode   = 34
	ExceptionBaseCode_UserAccount     ExceptionCode   = _ExceptionBaseCode_UserAccount + ReservedExceptionCode
	ExceptionPrefix_UserAccount       ExceptionPrefix = "UserAccount"
)

type UserAccountExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
	APIExceptionDomain
	TypeExceptionDomain
}

var UserAccount = &UserAccountExceptionDomain{
	BaseCode: ExceptionBaseCode_UserAccount,
	Prefix:   ExceptionPrefix_UserAccount,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UserAccount,
		_Prefix:   ExceptionPrefix_UserAccount,
	},
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UserAccount,
		_Prefix:   ExceptionPrefix_UserAccount,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UserAccount,
		_Prefix:   ExceptionPrefix_UserAccount,
	},
}

/* ============================== Handling Conflict Account Settings Error ============================== */

func (d *UserAccountExceptionDomain) GoogleCredentialHasAlreadyBeenBinded() *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "GoogleCredentialIsAlreadyBinded",
		IsInternal:     false,
		Message:        "The current account has already been binded with google",
		HTTPStatusCode: http.StatusInternalServerError,
		LastTrace:      traces.GetTrace(1),
	}
}
