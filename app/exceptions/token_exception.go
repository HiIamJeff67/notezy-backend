package exceptions

import (
	"net/http"
)

const (
	_ExceptionBaseCode_Token ExceptionCode = TokenExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	TokenExceptionSubDomainCode ExceptionCode   = 10
	ExceptionBaseCode_Token     ExceptionCode   = _ExceptionBaseCode_Token + ReservedExceptionCode
	ExceptionPrefix_Token       ExceptionPrefix = "Token"
)

type TokenExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	APIExceptionDomain
}

var Token = &TokenExceptionDomain{
	BaseCode: ExceptionBaseCode_Token,
	Prefix:   ExceptionPrefix_Token,
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Token,
		_Prefix:   ExceptionPrefix_Token,
	},
}

/* ============================== Handling Secret Key Environment Variable Not Found ============================== */

func (d *TokenExceptionDomain) AccessTokenSecretKeyNotFound() *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "AccessTokenSecretKeyNotFound",
		IsInternal:     true,
		Message:        "The environment variables of access token secret key is not found",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *TokenExceptionDomain) RefreshTokenSecretKeyNotFound() *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         "RefreshTokenSecretKeyNotFound",
		IsInternal:     true,
		Message:        "The environment variables of refresh token secret key is not found",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *TokenExceptionDomain) CSRFTokenSecretKeyNotFound() *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         "CSRFTokenSecretKeyNotFound",
		IsInternal:     true,
		Message:        "The environment variables of csrf token secret key is not found",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

/* ============================== Handling Generate Token Error ============================== */

func (d *TokenExceptionDomain) FailedToGenerateAccessToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 101,
		Prefix:         d.Prefix,
		Reason:         "FailedToGenerateAccessToken",
		IsInternal:     true,
		Message:        "Failed to generate the access token",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *TokenExceptionDomain) FailedToGenerateRefreshToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 102,
		Prefix:         d.Prefix,
		Reason:         "FailedToGenerateRefreshToken",
		IsInternal:     true,
		Message:        "Failed to generate the refresh token",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *TokenExceptionDomain) FailedToGenerateCSRFToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 103,
		Prefix:         d.Prefix,
		Reason:         "FailedToGenerateCSRFToken",
		IsInternal:     true,
		Message:        "Failed to generate the csrf token",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

/* ============================== Handling Parse Token Error ============================== */

func (d *TokenExceptionDomain) FailedToParseAccessToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 201,
		Prefix:         d.Prefix,
		Reason:         "FailedToParseAccessToken",
		IsInternal:     true,
		Message:        "Failed to parse the access token",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *TokenExceptionDomain) FailedToParseRefreshToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 202,
		Prefix:         d.Prefix,
		Reason:         "FailedToParseRefreshToken",
		IsInternal:     true,
		Message:        "Failed to parse the refresh token",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *TokenExceptionDomain) FailedToParseCSRFToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 203,
		Prefix:         d.Prefix,
		Reason:         "FailedToParseCSRFToken",
		IsInternal:     true,
		Message:        "Failed to parse the csrf token",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

/* ============================== Handling Token Expiration Error ============================== */

func (d *TokenExceptionDomain) AccessTokenExpired() *Exception {
	return &Exception{
		Code:           d.BaseCode + 301,
		Prefix:         d.Prefix,
		Reason:         "AccessTokenExpired",
		IsInternal:     true,
		Message:        "The given access token is expired",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *TokenExceptionDomain) RefreshTokenExpired() *Exception {
	return &Exception{
		Code:           d.BaseCode + 302,
		Prefix:         d.Prefix,
		Reason:         "RefreshTokenExpired",
		IsInternal:     true,
		Message:        "The given refresh token is expired",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *TokenExceptionDomain) CSRFTokenExpired() *Exception {
	return &Exception{
		Code:           d.BaseCode + 303,
		Prefix:         d.Prefix,
		Reason:         "CSRFTokenExpired",
		IsInternal:     true,
		Message:        "The given CSRF token is expired",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

/* ============================== Handling Invalid Type or Value in Token Claims ============================== */

func (d *TokenExceptionDomain) InvalidCSRFTokenSignature() *Exception {
	return &Exception{
		Code:           d.BaseCode + 401,
		Prefix:         d.Prefix,
		Reason:         "InvalidCSRFTokenSignature",
		IsInternal:     true,
		Message:        "The signature of CSRF token is invalid",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}

}

/* ============================== Extract or Validate Error ============================== */

func (d *TokenExceptionDomain) FailedToExtractOrValidateAccessToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 501,
		Prefix:         d.Prefix,
		Reason:         "FailedToExtractOrValidateAccessToken",
		IsInternal:     true,
		Message:        "Failed to get or validate the access token",
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *TokenExceptionDomain) FailedToExtractOrValidateRefreshToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 502,
		Prefix:         d.Prefix,
		Reason:         "FailedToExtractOrValidateRefreshToken",
		IsInternal:     true,
		Message:        "Failed to get or validate the refresh token",
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *TokenExceptionDomain) FailedToExtractOrValidateCSRFToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 503,
		Prefix:         d.Prefix,
		Reason:         "FailedToExtractOrValidateCSRFToken",
		IsInternal:     true,
		Message:        "Failed to get or validate the CSRF token",
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
