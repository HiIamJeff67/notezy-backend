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
}

var Util = &UtilExceptionDomain{
	BaseCode: ExceptionBaseCode_Util,
	Prefix:   ExceptionPrefix_Util,
}

/* ============================== Handing Exception on Hash ============================== */

func (d *UtilExceptionDomain) FailedToGenerateHashValue() *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
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

/* ============================== Handling Get Block Time Error ============================== */

func (d *UtilExceptionDomain) InvalidLoginCount(loginCount int32) *Exception {
	return &Exception{
		Code:           d.BaseCode + 31,
		Prefix:         d.Prefix,
		Reason:         "InvalidLoginCount",
		IsInternal:     true,
		Message:        fmt.Sprintf("The given loginCount is invalid: %d", loginCount),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *UtilExceptionDomain) InvalidAuthCodeRequestTimes(authCodeRequestTimes int32) *Exception {
	return &Exception{
		Code:           d.BaseCode + 32,
		Prefix:         d.Prefix,
		Reason:         "InvalidAuthCodeRequestTimes",
		IsInternal:     true,
		Message:        fmt.Sprintf("The given loginCount is invalid: %d", authCodeRequestTimes),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *UtilExceptionDomain) NotRequiredToBlockLogin(loginCount int32) *Exception {
	return &Exception{
		Code:           d.BaseCode + 33,
		Prefix:         d.Prefix,
		Reason:         "NotRequiredToBlock",
		IsInternal:     true,
		Message:        fmt.Sprintf("The loginCount of %d is no need to block", loginCount),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *UtilExceptionDomain) NotRequiredToBlockAuthCode(authCodeRequestTime int32) *Exception {
	return &Exception{
		Code:           d.BaseCode + 34,
		Prefix:         d.Prefix,
		Reason:         "NotRequiredToBlockAuthCode",
		IsInternal:     true,
		Message:        fmt.Sprintf("The authCodeRequestTime of %d is no need to block", authCodeRequestTime),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
