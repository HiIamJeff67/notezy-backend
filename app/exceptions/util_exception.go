package exceptions

import (
	"fmt"
	"net/http"
)

const (
	_ExceptionBaseCode_Util ExceptionCode = UtilExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	UtilExceptionSubDomainCode ExceptionCode   = 1
	ExceptionBaseCode_Util     ExceptionCode   = _ExceptionBaseCode_Util + ReservedExceptionCode
	ExceptionPrefix_Util       ExceptionPrefix = "Util"
)

type UtilExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	CommonExceptionDomain
}

var Util = &UtilExceptionDomain{
	BaseCode: ExceptionBaseCode_Util,
	Prefix:   ExceptionPrefix_Util,
	CommonExceptionDomain: CommonExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Util,
		_Prefix:   ExceptionPrefix_Util,
	},
}

/* ============================== Handling Exception on Json Web Tokens ============================== */

func (d *UtilExceptionDomain) AccessTokenSecretKeyNotFound() *Exception {
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

func (d *UtilExceptionDomain) RefreshTokenSecretKeyNotFound() *Exception {
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

func (d *UtilExceptionDomain) FailedToGenerateAccessToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         "FailedToGenerateAccessToken",
		IsInternal:     true,
		Message:        "Failed to generate the access token",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *UtilExceptionDomain) FailedToGenerateRefreshToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Reason:         "FailedToGenerateRefreshToken",
		IsInternal:     true,
		Message:        "Failed to generate the refresh token",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *UtilExceptionDomain) FailedToParseAccessToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 5,
		Prefix:         d.Prefix,
		Reason:         "FailedToParseAccessToken",
		IsInternal:     true,
		Message:        "Failed to parse the access token",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *UtilExceptionDomain) FailedToParseRefreshToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 6,
		Prefix:         d.Prefix,
		Reason:         "FailedToParseRefreshToken",
		IsInternal:     true,
		Message:        "Failed to parse the refresh token",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

/* ============================== Handing Exception on Hash ============================== */

func (d *UtilExceptionDomain) FailedToGenerateHashValue() *Exception {
	return &Exception{
		Code:           d.BaseCode + 11,
		Prefix:         d.Prefix,
		Reason:         "FailedToGenerateHashValue",
		IsInternal:     true,
		Message:        "Failed to generate the hash value",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

/* ============================== Handling Wrapping Official Utility Error ============================== */

func (d *UtilExceptionDomain) FailedToReadFile() *Exception {
	return &Exception{
		Code:           d.BaseCode + 21,
		Prefix:         d.Prefix,
		Reason:         "FailedToReadFile",
		IsInternal:     true,
		Message:        "Failed to read the file",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *UtilExceptionDomain) FailedToPreprocessPartialUpdate(values interface{}, setNull *map[string]bool, existingValues interface{}) *Exception {
	return &Exception{
		Code:           d.BaseCode + 22,
		Prefix:         d.Prefix,
		Reason:         "FailedToPreprocessPartialUpdate",
		IsInternal:     true,
		Message:        fmt.Sprintf("Failed to preprocess partial update with value: %v, setNull: %v, and existingValues: %v", values, setNull, existingValues),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
