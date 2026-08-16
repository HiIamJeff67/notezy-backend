package apiexceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notegic-backend/contracts/types/exceptions"
)

type UsersToBillingPlansException struct {
	CoreException
}

func NewUsersToBillingPlansException() UsersToBillingPlansException {
	return UsersToBillingPlansException{CoreException: NewCoreException("UsersToBillingPlans")}
}

func (UsersToBillingPlansException) NotFound() *exceptions.Exception {
	return exceptions.New(
		"NotFound",
		"UsersToBillingPlans",
		"Repository",
		"UsersToBillingPlans was not found",
		http.StatusNotFound,
	)
}

func (UsersToBillingPlansException) FailedToCreate() *exceptions.Exception {
	return exceptions.New(
		"FailedToCreate",
		"UsersToBillingPlans",
		"Repository",
		"Failed to create UsersToBillingPlans",
		http.StatusInternalServerError,
		true,
	)
}

func (UsersToBillingPlansException) FailedToUpdate() *exceptions.Exception {
	return exceptions.New(
		"FailedToUpdate",
		"UsersToBillingPlans",
		"Repository",
		"Failed to update UsersToBillingPlans",
		http.StatusInternalServerError,
		true,
	)
}

func (UsersToBillingPlansException) FailedToDelete() *exceptions.Exception {
	return exceptions.New(
		"FailedToDelete",
		"UsersToBillingPlans",
		"Repository",
		"Failed to delete UsersToBillingPlans",
		http.StatusInternalServerError,
		true,
	)
}

func (UsersToBillingPlansException) NoChanges() *exceptions.Exception {
	return exceptions.New(
		"NoChanges",
		"UsersToBillingPlans",
		"Repository",
		"No changes were applied to UsersToBillingPlans",
		http.StatusNotModified,
	)
}
