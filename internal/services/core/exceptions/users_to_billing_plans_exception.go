package apiexceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
)

type usersToBillingPlansExceptionDomain struct{}

var UsersToBillingPlans = usersToBillingPlansExceptionDomain{}

func (usersToBillingPlansExceptionDomain) NotFound() *exceptions.Exception {
	return exceptions.New(
		"NotFound",
		"UsersToBillingPlans",
		"Repository",
		"UsersToBillingPlans was not found",
		http.StatusNotFound,
	)
}

func (usersToBillingPlansExceptionDomain) FailedToCreate() *exceptions.Exception {
	return exceptions.New(
		"FailedToCreate",
		"UsersToBillingPlans",
		"Repository",
		"Failed to create UsersToBillingPlans",
		http.StatusInternalServerError,
		true,
	)
}

func (usersToBillingPlansExceptionDomain) FailedToUpdate() *exceptions.Exception {
	return exceptions.New(
		"FailedToUpdate",
		"UsersToBillingPlans",
		"Repository",
		"Failed to update UsersToBillingPlans",
		http.StatusInternalServerError,
		true,
	)
}

func (usersToBillingPlansExceptionDomain) FailedToDelete() *exceptions.Exception {
	return exceptions.New(
		"FailedToDelete",
		"UsersToBillingPlans",
		"Repository",
		"Failed to delete UsersToBillingPlans",
		http.StatusInternalServerError,
		true,
	)
}

func (usersToBillingPlansExceptionDomain) NoChanges() *exceptions.Exception {
	return exceptions.New(
		"NoChanges",
		"UsersToBillingPlans",
		"Repository",
		"No changes were applied to UsersToBillingPlans",
		http.StatusNotModified,
	)
}
