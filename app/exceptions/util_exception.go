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
		Message:        "The environment variables of access token secret key is not found",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *UtilExceptionDomain) RefreshTokenSecretKeyNotFound() *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Message:        "The environment variables of refresh token secret key is not found",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *UtilExceptionDomain) FailedToGenerateAccessToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Message:        "Failed to generate the access token",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *UtilExceptionDomain) FailedToGenerateRefreshToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Message:        "Failed to generate the refresh token",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *UtilExceptionDomain) FailedToParseAccessToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 5,
		Prefix:         d.Prefix,
		Message:        "Failed to parse the access token",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *UtilExceptionDomain) FailedToParseRefreshToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 6,
		Prefix:         d.Prefix,
		Message:        "Failed to parse the refresh token",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

/* ============================== Handing Exception on Hash ============================== */

func (d *UtilExceptionDomain) FailedToGenerateHashValue() *Exception {
	return &Exception{
		Code:           d.BaseCode + 11,
		Prefix:         d.Prefix,
		Message:        "Failed to generate the hash value",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

/* ============================== Handling Wrapping Official Utility Error ============================== */

func (d *UtilExceptionDomain) FailedToReadFile() *Exception {
	return &Exception{
		Code:           d.BaseCode + 21,
		Prefix:         d.Prefix,
		Message:        "Failed to read the file",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *UtilExceptionDomain) FailedToPreprocessPartialUpdate(values interface{}, setNull *map[string]bool, existingValues interface{}) *Exception {
	return &Exception{
		Code:           d.BaseCode + 22,
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("Failed to preprocess partial update with value: %v, setNull: %v, and existingValues: %v", values, setNull, existingValues),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}
