package apiexceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notegic-backend/contracts/types/exceptions"
)

type UserAccountException struct {
	CoreException
}

func NewUserAccountException() UserAccountException {
	return UserAccountException{CoreException: NewCoreException("UserAccount")}
}

func (UserAccountException) NotFound() *exceptions.Exception {
	return exceptions.New(
		"NotFound",
		"UserAccount",
		"Repository",
		"User account was not found",
		http.StatusNotFound,
	)
}

func (UserAccountException) FailedToCreate() *exceptions.Exception {
	return exceptions.New(
		"FailedToCreate",
		"UserAccount",
		"Repository",
		"Failed to create the user account",
		http.StatusInternalServerError,
		true,
	)
}

func (UserAccountException) FailedToUpdate() *exceptions.Exception {
	return exceptions.New(
		"FailedToUpdate",
		"UserAccount",
		"Repository",
		"Failed to update the user account",
		http.StatusInternalServerError,
		true,
	)
}

func (UserAccountException) NoChanges() *exceptions.Exception {
	return exceptions.New(
		"NoChanges",
		"UserAccount",
		"Repository",
		"No changes were applied to the user account",
		http.StatusNotModified,
	)
}
