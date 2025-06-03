package exceptions

import "net/http"

const (
	_ExceptionBaseCode_Auth ExceptionCode = (APIExceptionDomainCode*ExceptionDomainCodeShiftAmount +
		AuthExceptionSubDomainCode*ExceptionSubDomainCodeShiftAmount)

	AuthExceptionSubDomainCode ExceptionCode   = 3
	ExceptionBaseCode_Auth     ExceptionCode   = _ExceptionBaseCode_Auth + ReservedExceptionCode
	ExceptionPrefix_Auth       ExceptionPrefix = "Auth"
)

const (
	exceptionReason_WrongPassword     ExceptionReason = "Wrong_Password"
	exceptionReason_WrongAccessToken  ExceptionReason = "Wrong_AccessToken"
	exceptionReason_WrongRefreshToken ExceptionReason = "Wrong_RefreshToken"
)

type AuthExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	APIExceptionDomain
}

var Auth = &AuthExceptionDomain{
	BaseCode: ExceptionBaseCode_Auth,
	Prefix:   ExceptionPrefix_Auth,
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Auth,
		_Prefix:   ExceptionPrefix_Auth,
	},
}

/* ============================== Handling Invalid From ============================== */

func (d *AuthExceptionDomain) WrongPassword() *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         exceptionReason_WrongPassword,
		Message:        "The password is not match",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

func (d *AuthExceptionDomain) WrongAccessToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         exceptionReason_WrongAccessToken,
		Message:        "The access token is not match or expired",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

func (d *AuthExceptionDomain) WrongRefreshToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         exceptionReason_WrongRefreshToken,
		Message:        "The refresh token is not match or expired",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}
