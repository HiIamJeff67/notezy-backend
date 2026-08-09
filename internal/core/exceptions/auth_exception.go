package apiexceptions

import (
	"fmt"
	"net/http"
	"time"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
)

type AuthException struct {
	CoreException
}

func NewAuthException() AuthException {
	return AuthException{
		CoreException: NewCoreException("Auth"),
	}
}

func (AuthException) WrongPassword() *exceptions.Exception {
	return exceptions.New(
		"WrongPassword",
		"Auth",
		"Authenticate",
		"The password does not match",
		http.StatusUnauthorized,
	)
}

func (AuthException) WrongAuthCode() *exceptions.Exception {
	return exceptions.New(
		"WrongAuthCode",
		"Auth",
		"Authenticate",
		"The authentication code does not match",
		http.StatusUnauthorized,
	)
}

func (AuthException) LoginBlockedDueToTryingTooManyTimes(blockedUntil time.Time) *exceptions.Exception {
	return exceptions.New(
		"LoginBlockedDueToTryingTooManyTimes",
		"Auth",
		"Authenticate",
		fmt.Sprintf("Login is blocked until %v", blockedUntil),
		http.StatusTooManyRequests,
	)
}

func (AuthException) AuthCodeBlockedDueToTryingTooManyTimes(blockedUntil time.Time) *exceptions.Exception {
	return exceptions.New(
		"AuthCodeBlockedDueToTryingTooManyTimes",
		"Auth",
		"Authenticate",
		fmt.Sprintf("Auth code generation is blocked until %v", blockedUntil),
		http.StatusTooManyRequests,
	)
}
