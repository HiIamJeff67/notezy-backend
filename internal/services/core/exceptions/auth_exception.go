package apiexceptions

import (
	"fmt"
	"net/http"
	"time"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
)

type authExceptionDomain struct {
	domainException
}

var Auth = authExceptionDomain{
	domainException: newDomainException("Auth"),
}

func (authExceptionDomain) WrongPassword() *exceptions.Exception {
	return exceptions.New(
		"WrongPassword",
		"Auth",
		"Authenticate",
		"The password does not match",
		http.StatusUnauthorized,
	)
}

func (authExceptionDomain) WrongAuthCode() *exceptions.Exception {
	return exceptions.New(
		"WrongAuthCode",
		"Auth",
		"Authenticate",
		"The authentication code does not match",
		http.StatusUnauthorized,
	)
}

func (authExceptionDomain) LoginBlockedDueToTryingTooManyTimes(blockedUntil time.Time) *exceptions.Exception {
	return exceptions.New(
		"LoginBlockedDueToTryingTooManyTimes",
		"Auth",
		"Authenticate",
		fmt.Sprintf("Login is blocked until %v", blockedUntil),
		http.StatusTooManyRequests,
	)
}

func (authExceptionDomain) AuthCodeBlockedDueToTryingTooManyTimes(blockedUntil time.Time) *exceptions.Exception {
	return exceptions.New(
		"AuthCodeBlockedDueToTryingTooManyTimes",
		"Auth",
		"Authenticate",
		fmt.Sprintf("Auth code generation is blocked until %v", blockedUntil),
		http.StatusTooManyRequests,
	)
}
