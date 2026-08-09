package apiexceptions

import (
	"fmt"
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
)

type UserException struct {
	CoreException
}

func NewUserException() UserException {
	return UserException{CoreException: NewCoreException("User")}
}

func (UserException) NotFound() *exceptions.Exception {
	return exceptions.New(
		"NotFound",
		"User",
		"Repository",
		"User was not found",
		http.StatusNotFound,
	)
}

func (UserException) FailedToCreate() *exceptions.Exception {
	return exceptions.New(
		"FailedToCreate",
		"User",
		"Repository",
		"Failed to create the user",
		http.StatusInternalServerError,
		true,
	)
}

func (UserException) FailedToUpdate() *exceptions.Exception {
	return exceptions.New(
		"FailedToUpdate",
		"User",
		"Repository",
		"Failed to update the user",
		http.StatusInternalServerError,
		true,
	)
}

func (UserException) FailedToDelete() *exceptions.Exception {
	return exceptions.New(
		"FailedToDelete",
		"User",
		"Repository",
		"Failed to delete the user",
		http.StatusInternalServerError,
		true,
	)
}

func (UserException) NoChanges() *exceptions.Exception {
	return exceptions.New(
		"NoChanges",
		"User",
		"Repository",
		"No changes were applied to the user",
		http.StatusNotModified,
	)
}

func (UserException) FailedToCommitTransaction() *exceptions.Exception {
	return exceptions.New(
		"FailedToCommitTransaction",
		"User",
		"Transaction",
		"Failed to commit the user transaction",
		http.StatusInternalServerError,
		true,
	)
}

func (UserException) InvalidInput() *exceptions.Exception {
	return exceptions.New(
		"InvalidInput",
		"User",
		"Validate",
		"Invalid user input",
		http.StatusBadRequest,
	)
}

func (UserException) DuplicateName(name string) *exceptions.Exception {
	return exceptions.New(
		"DuplicateName",
		"User",
		"Create",
		fmt.Sprintf("The name of %s is already in use", name),
		http.StatusConflict,
	)
}

func (UserException) DuplicateEmail(email string) *exceptions.Exception {
	return exceptions.New(
		"DuplicateEmail",
		"User",
		"Create",
		fmt.Sprintf("The email of %s is already in use", email),
		http.StatusConflict,
	)
}
