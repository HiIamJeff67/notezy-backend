package exceptions

import (
	"fmt"
	"net/http"

	shared "notezy-backend/shared"
)

const (
	_ExceptionBaseCode_Cookie ExceptionCode = CookieExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	CookieExceptionSubDomainCode ExceptionCode   = 2
	ExceptionBaseCode_Cookie     ExceptionCode   = _ExceptionBaseCode_Cookie + ReservedExceptionCode
	ExceptionPrefix_Cookie       ExceptionPrefix = "Cookie"
)

type CookieExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	APIExceptionDomain
}

var Cookie = &CookieExceptionDomain{
	BaseCode: ExceptionBaseCode_Cookie,
	Prefix:   ExceptionPrefix_Cookie,
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Cookie,
		_Prefix:   ExceptionPrefix_Cookie,
	},
}

func (d *CookieExceptionDomain) NotFound(cookieName shared.ValidCookieName) *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_NotFound,
		Message:        fmt.Sprintf("Cannot find the %s in the cookie", convertCamelCaseToSentenceCase(cookieName.String())),
		HTTPStatusCode: http.StatusNotFound,
	}
}

func (d *CookieExceptionDomain) FailedToCreate(cookieName shared.ValidCookieName) *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToCreate,
		Message:        fmt.Sprintf("Failed to set the %s to the cache", convertCamelCaseToSentenceCase(cookieName.String())),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}
