package apiexceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/shared/exceptions"
)

type themeExceptionDomain struct{}

var Theme = themeExceptionDomain{}

func (themeExceptionDomain) NotFound() *exceptions.Exception {
	return exceptions.New(
		"NotFound",
		"Theme",
		"Repository",
		"Theme was not found",
		http.StatusNotFound,
	)
}

func (themeExceptionDomain) FailedToCreate() *exceptions.Exception {
	return exceptions.New(
		"FailedToCreate",
		"Theme",
		"Repository",
		"Failed to create Theme",
		http.StatusInternalServerError,
		true,
	)
}

func (themeExceptionDomain) FailedToUpdate() *exceptions.Exception {
	return exceptions.New(
		"FailedToUpdate",
		"Theme",
		"Repository",
		"Failed to update Theme",
		http.StatusInternalServerError,
		true,
	)
}

func (themeExceptionDomain) FailedToDelete() *exceptions.Exception {
	return exceptions.New(
		"FailedToDelete",
		"Theme",
		"Repository",
		"Failed to delete Theme",
		http.StatusInternalServerError,
		true,
	)
}

func (themeExceptionDomain) NoChanges() *exceptions.Exception {
	return exceptions.New(
		"NoChanges",
		"Theme",
		"Repository",
		"No changes were applied to Theme",
		http.StatusNotModified,
	)
}
