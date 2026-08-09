package apiexceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
)

type BadgeException struct {
	CoreException
}

func NewBadgeException() BadgeException {
	return BadgeException{CoreException: NewCoreException("Badge")}
}

func (BadgeException) NotFound() *exceptions.Exception {
	return exceptions.New(
		"NotFound",
		"Badge",
		"Repository",
		"Badge was not found",
		http.StatusNotFound,
	)
}
