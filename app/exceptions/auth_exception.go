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
	ExceptionReason_WrongPassword     ExceptionReason = "Wrong_Password"
	ExceptionReason_WrongAccessToken  ExceptionReason = "Wrong_AccessToken"
	ExceptionReason_WrongRefreshToken ExceptionReason = "Wrong_RefreshToken"
)

type AuthExceptionDomain struct {
	APIExceptionDomain
}

var Auth = &AuthExceptionDomain{
	APIExceptionDomain{BaseCode: _ExceptionBaseCode_Auth, Prefix: ExceptionPrefix_Auth},
}

/* ============================== Handling Invalid From ============================== */
func (d *AuthExceptionDomain) WrongPassword() *Exception {
	return &Exception{
		Code:           ExceptionBaseCode_Auth + 1,
		Prefix:         ExceptionPrefix_Auth,
		Reason:         ExceptionReason_WrongPassword,
		Message:        "The password is not match",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

func (d *AuthExceptionDomain) WrongAccessToken() *Exception {
	return &Exception{
		Code:           ExceptionBaseCode_Auth + 2,
		Prefix:         ExceptionPrefix_Auth,
		Reason:         ExceptionReason_WrongAccessToken,
		Message:        "The access token is not match or expired",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

func (d *AuthExceptionDomain) WrongRefreshToken() *Exception {
	return &Exception{
		Code:           ExceptionBaseCode_Auth + 3,
		Prefix:         ExceptionPrefix_Auth,
		Reason:         ExceptionReason_WrongRefreshToken,
		Message:        "The refresh token is not match or expired",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}
