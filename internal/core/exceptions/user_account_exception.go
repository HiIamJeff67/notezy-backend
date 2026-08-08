package apiexceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
)

type userAccountExceptionDomain struct{}

var UserAccount = userAccountExceptionDomain{}

func (userAccountExceptionDomain) NotFound() *exceptions.Exception {
	return exceptions.New(
		"NotFound",
		"UserAccount",
		"Repository",
		"User account was not found",
		http.StatusNotFound,
	)
}

func (userAccountExceptionDomain) FailedToCreate() *exceptions.Exception {
	return exceptions.New(
		"FailedToCreate",
		"UserAccount",
		"Repository",
		"Failed to create the user account",
		http.StatusInternalServerError,
		true,
	)
}

func (userAccountExceptionDomain) FailedToUpdate() *exceptions.Exception {
	return exceptions.New(
		"FailedToUpdate",
		"UserAccount",
		"Repository",
		"Failed to update the user account",
		http.StatusInternalServerError,
		true,
	)
}

func (userAccountExceptionDomain) NoChanges() *exceptions.Exception {
	return exceptions.New(
		"NoChanges",
		"UserAccount",
		"Repository",
		"No changes were applied to the user account",
		http.StatusNotModified,
	)
}
