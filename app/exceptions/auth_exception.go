package exceptions

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	_ExceptionBaseCode_Auth ExceptionCode = (APIExceptionDomainCode*ExceptionDomainCodeShiftAmount +
		AuthExceptionSubDomainCode*ExceptionSubDomainCodeShiftAmount)

	AuthExceptionSubDomainCode ExceptionCode   = 3
	ExceptionBaseCode_Auth     ExceptionCode   = _ExceptionBaseCode_Auth + ReservedExceptionCode
	ExceptionPrefix_Auth       ExceptionPrefix = "Auth"
)

const (
	ExceptionReason_WrongPassword                         ExceptionReason = "Wrong_Password"
	ExceptionReason_WrongAccessToken                      ExceptionReason = "Wrong_Access_Token"
	ExceptionReason_WrongRefreshToken                     ExceptionReason = "Wrong_Refresh_Token"
	ExceptionReason_WrongUserAgent                        ExceptionReason = "Wrong_User_Agent"
	ExceptionReason_WrongAuthCode                         ExceptionReason = "Wrong_Authentication_Code"
	ExceptionReason_FailedToExtractOrValidateAccessToken  ExceptionReason = "Failed_To_Extract_Or_Validate_Access_Token"
	ExceptionReason_FailedToExtractOrValidateRefreshToken ExceptionReason = "Failed_To_Extract_Or_Validate_Refresh_Token"
	ExceptionReason_PermissionDeniedDueToUserRole         ExceptionReason = "Permission_Denied_Due_To_User_Role"
	ExceptionReason_PermissionDeniedDueToUserPlan         ExceptionReason = "Permission_Denied_Due_To_User_Plan"
	ExceptionReason_MissPlacingOrWrongMiddlewareOrder     ExceptionReason = "Miss_Placing_Or_Wrong_Middleware_Order"
	ExceptionReason_LoginBlockedDueToTryingTooManyTimes   ExceptionReason = "Login_Blocked_Due_To_Trying_Too_Many_Times"
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
		Reason:         ExceptionReason_WrongPassword,
		Message:        "The password is not match",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

func (d *AuthExceptionDomain) WrongAccessToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_WrongAccessToken,
		Message:        "The access token is not match or expired",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

func (d *AuthExceptionDomain) WrongRefreshToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_WrongRefreshToken,
		Message:        "The refresh token is not match or expired",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

func (d *AuthExceptionDomain) WrongUserAgent() *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_WrongUserAgent,
		Message:        "The user agent is not match",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

func (d *AuthExceptionDomain) WrongAuthCode() *Exception {
	return &Exception{
		Code:           d.BaseCode + 5,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_WrongAuthCode,
		Message:        "The authentication code is not match",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

func (d *AuthExceptionDomain) FailedToExtractOrValidateAccessToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 6,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToExtractOrValidateAccessToken,
		Message:        "Failed to get or validate the access token",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

func (d *AuthExceptionDomain) FailedToExtractOrValidateRefreshToken() *Exception {
	return &Exception{
		Code:           d.BaseCode + 7,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToExtractOrValidateRefreshToken,
		Message:        "Failed to get or validate the refresh token",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

func (d *AuthExceptionDomain) LoginBlockedDueToTryingTooManyTimes(blockedUntil time.Time) *Exception {
	return &Exception{
		Code:           d.BaseCode + 8,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_LoginBlockedDueToTryingTooManyTimes,
		Message:        fmt.Sprintf("Blocked the login procedure because user has tried too many time and require to wait until %v", blockedUntil),
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

/* ========================= Handling Permission Denied ========================= */

func (d *AuthExceptionDomain) PermissionDeniedDueToUserRole(userRole any) *Exception {
	return &Exception{
		Code:           d.BaseCode + 11,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_PermissionDeniedDueToUserRole,
		Message:        fmt.Sprintf("The current user role of %v does not have access to this operation", userRole),
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

func (d *AuthExceptionDomain) PermissionDeniedDueToUserPlan(userPlan any) *Exception {
	return &Exception{
		Code:           d.BaseCode + 12,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_PermissionDeniedDueToUserPlan,
		Message:        fmt.Sprintf("The current user plan of %v does not have access to this operation", userPlan),
		HTTPStatusCode: http.StatusUnauthorized,
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
		Reason:         ExceptionReason_MissPlacingOrWrongMiddlewareOrder,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
	}
}
