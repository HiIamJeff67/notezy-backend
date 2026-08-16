package apiexceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notegic-backend/contracts/types/exceptions"
)

type ThemeException struct {
	CoreException
}

func NewThemeException() ThemeException {
	return ThemeException{CoreException: NewCoreException("Theme")}
}

func (ThemeException) NotFound() *exceptions.Exception {
	return exceptions.New(
		"NotFound",
		"Theme",
		"Repository",
		"Theme was not found",
		http.StatusNotFound,
	)
}

func (ThemeException) FailedToCreate() *exceptions.Exception {
	return exceptions.New(
		"FailedToCreate",
		"Theme",
		"Repository",
		"Failed to create Theme",
		http.StatusInternalServerError,
		true,
	)
}

func (ThemeException) FailedToUpdate() *exceptions.Exception {
	return exceptions.New(
		"FailedToUpdate",
		"Theme",
		"Repository",
		"Failed to update Theme",
		http.StatusInternalServerError,
		true,
	)
}

func (ThemeException) FailedToDelete() *exceptions.Exception {
	return exceptions.New(
		"FailedToDelete",
		"Theme",
		"Repository",
		"Failed to delete Theme",
		http.StatusInternalServerError,
		true,
	)
}

func (ThemeException) NoChanges() *exceptions.Exception {
	return exceptions.New(
		"NoChanges",
		"Theme",
		"Repository",
		"No changes were applied to Theme",
		http.StatusNotModified,
	)
}
