package apiexceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/shared/exceptions"
)

type badgeExceptionDomain struct{}

var Badge = badgeExceptionDomain{}

func (badgeExceptionDomain) NotFound() *exceptions.Exception {
	return exceptions.New(
		"NotFound",
		"Badge",
		"Repository",
		"Badge was not found",
		http.StatusNotFound,
	)
}
