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
	ExceptionReason__WrongPassword                         ExceptionReason = "Wrong_Password"
	ExceptionReason__WrongAccessToken                      ExceptionReason = "Wrong_AccessToken"
	ExceptionReason__WrongRefreshToken                     ExceptionReason = "Wrong_RefreshToken"
	ExceptionReason__FailedToExtractOrValidateAccessToken  ExceptionReason = "Failed_To_Extract_Or_Validate_AccessToken"
	ExceptionReason__FailedToExtractOrValidateRefreshToken ExceptionReason = "Failed_To_Extract_Or_Validate_RefreshToken"
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
		Reason:         ExceptionReason__WrongPassword,
		Message:        "The password is not match",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

func (d *AuthExceptionDomain) WrongAccessToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason__WrongAccessToken,
		Message:        "The access token is not match or expired",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

func (d *AuthExceptionDomain) WrongRefreshToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason__WrongRefreshToken,
		Message:        "The refresh token is not match or expired",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

func (d *AuthExceptionDomain) FailedToExtractOrValidateAccessToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason__FailedToExtractOrValidateAccessToken,
		Message:        "Failed to get or validate the access token",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

func (d *AuthExceptionDomain) FailedToExtractOrValidateRefreshToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 5,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason__FailedToExtractOrValidateRefreshToken,
		Message:        "Failed to get or validate the refresh token",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}
