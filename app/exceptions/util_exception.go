package exceptions

import (
	"fmt"
	"net/http"
)

const (
	_ExceptionBaseCode_Util ExceptionCode = (APIExceptionDomainCode*ExceptionDomainCodeShiftAmount +
		UtilExceptionSubDomainCode*ExceptionSubDomainCodeShiftAmount)

	UtilExceptionSubDomainCode ExceptionCode   = 1
	ExceptionBaseCode_Util     ExceptionCode   = _ExceptionBaseCode_Util + ReservedExceptionCode
	ExceptionPrefix_Util       ExceptionPrefix = "Util"
)

const (
	ExceptionReason_AccessTokenSecretKeyNotFound    ExceptionReason = "Access_Token_Secret_Key_Not_Found"
	ExceptionReason_RefreshTokenSecretKeyNotFound   ExceptionReason = "Refresh_Token_Secret_Key_Not_Found"
	ExceptionReason_FailedToGenerateAccessToken     ExceptionReason = "Failed_To_Generate_Access_Token"
	ExceptionReason_FailedToGenerateRefreshToken    ExceptionReason = "Failed_To_Generate_Refresh_Token"
	ExceptionReason_FailedToParseAccessToken        ExceptionReason = "Failed_To_Parse_Access_Token"
	ExceptionReason_FailedToParseRefreshToken       ExceptionReason = "Failed_To_Parse_Refresh_Token"
	ExceptionReason_FailedToGenerateHashValue       ExceptionReason = "Failed_To_Generate_Hash_Value"
	ExceptionReason_FailedToReadFile                ExceptionReason = "Failed_To_Read_File"
	ExceptionReason_FailedToPreprocessPartialUpdate ExceptionReason = "Failed_To_Preprocess_Partial_Update"
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

/* ============================== Handling Exception on Json Web Tokens ============================== */

func (d *UtilExceptionDomain) AccessTokenSecretKeyNotFound() *Exception {
	return &Exception{
		Code:           d.BaseCode + 11,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_AccessTokenSecretKeyNotFound,
		Message:        "The environment variables of access token secret key is not found",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *UtilExceptionDomain) RefreshTokenSecretKeyNotFound() *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_RefreshTokenSecretKeyNotFound,
		Message:        "The environment variables of refresh token secret key is not found",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *UtilExceptionDomain) FailedToGenerateAccessToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToGenerateAccessToken,
		Message:        "Failed to generate the access token",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *UtilExceptionDomain) FailedToGenerateRefreshToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToGenerateRefreshToken,
		Message:        "Failed to generate the refresh token",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *UtilExceptionDomain) FailedToParseAccessToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 5,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToParseAccessToken,
		Message:        "Failed to parse the access token",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *UtilExceptionDomain) FailedToParseRefreshToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 6,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToParseRefreshToken,
		Message:        "Failed to parse the refresh token",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

/* ============================== Handing Exception on Hash ============================== */

func (d *UtilExceptionDomain) FailedToGenerateHashValue() *Exception {
	return &Exception{
		Code:           d.BaseCode + 11,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToGenerateHashValue,
		Message:        "Failed to generate the hash value",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

/* ============================== Handling Wrapping Official Utility Error ============================== */

func (d *UtilExceptionDomain) FailedToReadFile() *Exception {
	return &Exception{
		Code:           d.BaseCode + 21,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToReadFile,
		Message:        "Failed to read the file",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *UtilExceptionDomain) FailedToPreprocessPartialUpdate(values interface{}, setNull map[string]bool, existingValues interface{}) *Exception {
	return &Exception{
		Code:           d.BaseCode + 22,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToPreprocessPartialUpdate,
		Message:        fmt.Sprintf("Failed to preprocess partial update with value: %v, setNull: %v, and existingValues: %v", values, setNull, existingValues),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}
