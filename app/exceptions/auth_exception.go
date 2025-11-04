package exceptions

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	AuthExceptionSubDomainCode ExceptionCode   = 31
	_ExceptionBaseCode_Auth    ExceptionCode   = AuthExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount // the actual exception code, but we need to reserved spaces for some predefined exceptions
	ExceptionBaseCode_Auth     ExceptionCode   = _ExceptionBaseCode_Auth + ReservedExceptionCode
	ExceptionPrefix_Auth       ExceptionPrefix = "Auth"
)

type AuthExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	APIExceptionDomain
	TypeExceptionDomain
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
}

/* ============================== Handling Invalid From ============================== */

func (d *AuthExceptionDomain) WrongPassword() *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "WrongPassword",
		IsInternal:     false,
		Message:        "The password is not match",
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) WrongAccessToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         "WrongAccessToken",
		IsInternal:     true,
		Message:        "The access token is not match or expired",
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) WrongRefreshToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         "WrongRefreshToken",
		IsInternal:     true,
		Message:        "The refresh token is not match or expired",
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) WrongUserAgent() *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Reason:         "WrongUserAgent",
		IsInternal:     true,
		Message:        "The user agent is not match",
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) WrongAuthCode() *Exception {
	return &Exception{
		Code:           d.BaseCode + 5,
		Prefix:         d.Prefix,
		Reason:         "WrongAuthCode",
		IsInternal:     false,
		Message:        "The authentication code is not match",
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) LoginBlockedDueToTryingTooManyTimes(blockedUntil time.Time) *Exception {
	return &Exception{
		Code:           d.BaseCode + 6,
		Prefix:         d.Prefix,
		Reason:         "LoginBlockedDueToTryingTooManyTimes",
		IsInternal:     false,
		Message:        fmt.Sprintf("Blocked the login procedure because user has tried too many times and require to wait until %v", blockedUntil),
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) AuthCodeBlockedDueToTryingTooManyTimes(blockedUntil time.Time) *Exception {
	return &Exception{
		Code:           d.BaseCode + 7,
		Prefix:         d.Prefix,
		Reason:         "AuthCodeBlockedDueToTryingTooManyTimes",
		IsInternal:     false,
		Message:        fmt.Sprintf("Blocked the generating auth code procedure because user has tried too many times and require to wait until %v", blockedUntil),
		HTTPStatusCode: http.StatusTooManyRequests,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) NoClientIPOrReferenceToClient() *Exception {
	return &Exception{
		Code:           d.BaseCode + 8,
		Prefix:         d.Prefix,
		Reason:         "NoClientIPOrReferenceToClient",
		IsInternal:     true,
		Message:        "Cannot extract or find any reference to the client",
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) InvalidRateLimitTokenCount() *Exception {
	return &Exception{
		Code:           d.BaseCode + 9,
		Prefix:         d.Prefix,
		Reason:         "InvalidRateLimitTokenCount",
		IsInternal:     true,
		Message:        "The rate limit token count is invalid",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) RateLimitExceeded() *Exception {
	return &Exception{
		Code:           d.BaseCode + 10,
		Prefix:         d.Prefix,
		Reason:         "RateLimitExceeded",
		IsInternal:     false,
		Message:        "The rate limit token count is exceeded, please try again later",
		HTTPStatusCode: http.StatusBadRequest,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

/* ========================= Handling Permission Denied ========================= */

func (d *AuthExceptionDomain) PermissionDeniedDueToUserRole(userRole any) *Exception {
	return &Exception{
		Code:           d.BaseCode + 101,
		Prefix:         d.Prefix,
		Reason:         "PermissionDeniedDueToUserRole",
		IsInternal:     false,
		Message:        fmt.Sprintf("The current user role of %v does not have access to this operation", userRole),
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) PermissionDeniedDueToUserPlan(userPlan any) *Exception {
	return &Exception{
		Code:           d.BaseCode + 102,
		Prefix:         d.Prefix,
		Reason:         "PermissionDeniedDueToUserPlan",
		IsInternal:     false,
		Message:        fmt.Sprintf("The current user plan of %v does not have access to this operation", userPlan),
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) PermissionDeniedDueToInvalidRequestOriginDomain(origin string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 103,
		Prefix:         d.Prefix,
		Reason:         "PermissionDeniedDueToInvalidRequestOriginDomain",
		IsInternal:     false,
		Message:        fmt.Sprintf("The current request origin domain of %s is invalid", origin),
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *AuthExceptionDomain) PermissionDeniedDueToTooManyRequests() *Exception {
	return &Exception{
		Code:           d.BaseCode + 104,
		Prefix:         d.Prefix,
		Reason:         "PermissionDeniedDueToTooManyRequests",
		IsInternal:     false,
		Message:        "Too many requests, please wait for a while",
		HTTPStatusCode: http.StatusTooManyRequests,
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
		Code:           d.BaseCode + 201,
		Prefix:         d.Prefix,
		Reason:         "MissPlacingOrWrongMiddlewareOrder",
		IsInternal:     true,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
