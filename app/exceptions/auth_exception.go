package exceptions

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	_ExceptionBaseCode_Auth ExceptionCode = AuthExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	AuthExceptionSubDomainCode ExceptionCode   = 31
	ExceptionBaseCode_Auth     ExceptionCode   = _ExceptionBaseCode_Auth + ReservedExceptionCode
	ExceptionPrefix_Auth       ExceptionPrefix = "Auth"
)

type AuthExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	APIExceptionDomain
	TypeExceptionDomain
	CommonExceptionDomain
}

var Auth = &AuthExceptionDomain{
	BaseCode: ExceptionBaseCode_Auth,
	Prefix:   ExceptionPrefix_Auth,
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Auth,
		_Prefix:   ExceptionPrefix_Auth,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Auth,
		_Prefix:   ExceptionPrefix_Auth,
	},
	CommonExceptionDomain: CommonExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Auth,
		_Prefix:   ExceptionPrefix_Auth,
	},
}

/* ============================== Handling Invalid From ============================== */

func (d *AuthExceptionDomain) WrongPassword() *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Message:        "The password is not match",
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) WrongAccessToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Message:        "The access token is not match or expired",
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) WrongRefreshToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Message:        "The refresh token is not match or expired",
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) WrongUserAgent() *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Message:        "The user agent is not match",
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) WrongAuthCode() *Exception {
	return &Exception{
		Code:           d.BaseCode + 5,
		Prefix:         d.Prefix,
		Message:        "The authentication code is not match",
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) FailedToExtractOrValidateAccessToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 6,
		Prefix:         d.Prefix,
		Message:        "Failed to get or validate the access token",
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) FailedToExtractOrValidateRefreshToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 7,
		Prefix:         d.Prefix,
		Message:        "Failed to get or validate the refresh token",
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) LoginBlockedDueToTryingTooManyTimes(blockedUntil time.Time) *Exception {
	return &Exception{
		Code:           d.BaseCode + 8,
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("Blocked the login procedure because user has tried too many times and require to wait until %v", blockedUntil),
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

/* ========================= Handling Permission Denied ========================= */

func (d *AuthExceptionDomain) PermissionDeniedDueToUserRole(userRole any) *Exception {
	return &Exception{
		Code:           d.BaseCode + 11,
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("The current user role of %v does not have access to this operation", userRole),
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) PermissionDeniedDueToUserPlan(userPlan any) *Exception {
	return &Exception{
		Code:           d.BaseCode + 12,
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("The current user plan of %v does not have access to this operation", userPlan),
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

/* ========================= Handling Internal Error ========================= */
// (this part is usually due to the developer, ex. Place the middleware in the wrong order)

func (d *AuthExceptionDomain) MissPlacingOrWrongMiddlewareOrder(optionalMessage ...string) *Exception {
	message := "Miss placing or placing the middleware in the wrong order"
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], "", " ")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d.BaseCode + 101,
		Prefix:         d.Prefix,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
