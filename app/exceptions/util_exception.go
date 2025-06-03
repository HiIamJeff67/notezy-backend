package exceptions

import "net/http"

const (
	_ExceptionBaseCode_Util ExceptionCode = (APIExceptionDomainCode*ExceptionDomainCodeShiftAmount +
		UtilExceptionSubDomainCode*ExceptionSubDomainCodeShiftAmount)

	UtilExceptionSubDomainCode ExceptionCode   = 1
	ExceptionBaseCode_Util     ExceptionCode   = _ExceptionBaseCode_Util + ReservedExceptionCode
	ExceptionPrefix_Util       ExceptionPrefix = "Util"
)

const (
	exceptionReason_FailedToGenerateHashValue     ExceptionReason = "Failed_To_Generate_Hash_Value"
	exceptionReason_AccessTokenSecretKeyNotFound  ExceptionReason = "Access_Token_Secret_Key_Not_Found"
	exceptionReason_RefreshTokenSecretKeyNotFound ExceptionReason = "Refresh_Token_Secret_Key_Not_Found"
	exceptionReason_FailedToGenerateAccessToken   ExceptionReason = "Failed_To_Generate_Access_Token"
	exceptionReason_FailedToGenerateRefreshToken  ExceptionReason = "Failed_To_Generate_Refresh_Token"
	exceptionReason_FailedToParseAccessToken      ExceptionReason = "Failed_To_Parse_Access_Token"
	exceptionReason_FailedToParseRefreshToken     ExceptionReason = "Failed_To_Parse_Refresh_Token"
)

type UtilExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	APIExceptionDomain
}

var Util = &UtilExceptionDomain{
	BaseCode: ExceptionBaseCode_Util,
	Prefix:   ExceptionPrefix_Util,
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Util,
		_Prefix:   ExceptionPrefix_Util,
	},
}

/* ============================== Handing Exception on Hash ============================== */

func (d *UtilExceptionDomain) FailedToGenerateHashValue() *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         exceptionReason_FailedToGenerateHashValue,
		Message:        "Failed to generate the hash value",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

/* ============================== Handling Exception on Json Web Tokens ============================== */

func (d *UtilExceptionDomain) AccessTokenSecretKeyNotFound() *Exception {
	return &Exception{
		Code:           d.BaseCode + 11,
		Prefix:         d.Prefix,
		Reason:         exceptionReason_AccessTokenSecretKeyNotFound,
		Message:        "The environment variables of access token secret key is not found",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *UtilExceptionDomain) RefreshTokenSecretKeyNotFound() *Exception {
	return &Exception{
		Code:           d.BaseCode + 12,
		Prefix:         d.Prefix,
		Reason:         exceptionReason_RefreshTokenSecretKeyNotFound,
		Message:        "The environment variables of refresh token secret key is not found",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *UtilExceptionDomain) FailedToGenerateAccessToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 13,
		Prefix:         d.Prefix,
		Reason:         exceptionReason_FailedToGenerateAccessToken,
		Message:        "Failed to generate the access token",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *UtilExceptionDomain) FailedToGenerateRefreshToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 14,
		Prefix:         d.Prefix,
		Reason:         exceptionReason_FailedToGenerateRefreshToken,
		Message:        "Failed to generate the refresh token",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *UtilExceptionDomain) FailedToParseAccessToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 15,
		Prefix:         d.Prefix,
		Reason:         exceptionReason_FailedToParseAccessToken,
		Message:        "Failed to parse the access token",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *UtilExceptionDomain) FailedToParseRefreshToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 16,
		Prefix:         d.Prefix,
		Reason:         exceptionReason_FailedToParseRefreshToken,
		Message:        "Failed to parse the refresh token",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}
