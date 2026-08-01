package apiexceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
)

type blockExceptionDomain struct{}

var Block = blockExceptionDomain{}

func (blockExceptionDomain) NotFound() *exceptions.Exception {
	return exceptions.New(
		"NotFound",
		"Block",
		"Repository",
		"Block was not found",
		http.StatusNotFound,
	)
}

func (blockExceptionDomain) FailedToCreate() *exceptions.Exception {
	return exceptions.New(
		"FailedToCreate",
		"Block",
		"Repository",
		"Failed to create the block",
		http.StatusInternalServerError,
		true,
	)
}

func (blockExceptionDomain) FailedToUpdate() *exceptions.Exception {
	return exceptions.New(
		"FailedToUpdate",
		"Block",
		"Repository",
		"Failed to update the block",
		http.StatusInternalServerError,
		true,
	)
}

func (blockExceptionDomain) NoChanges() *exceptions.Exception {
	return exceptions.New(
		"NoChanges",
		"Block",
		"Repository",
		"No changes were applied to the block",
		http.StatusNotModified,
	)
}

func (blockExceptionDomain) FailedToCommitTransaction() *exceptions.Exception {
	return exceptions.New(
		"FailedToCommitTransaction",
		"Block",
		"Transaction",
		"Failed to commit the block transaction",
		http.StatusInternalServerError,
		true,
	)
}

func (blockExceptionDomain) InvalidDto() *exceptions.Exception {
	return exceptions.New(
		"InvalidDto",
		"Block",
		"Validate",
		"Invalid block DTO",
		http.StatusBadRequest,
	)
}

func (blockExceptionDomain) NoPermission(action string) *exceptions.Exception {
	return exceptions.New(
		"PermissionDenied",
		"Block",
		"Authorize",
		"Permission is denied to "+action,
		http.StatusForbidden,
	)
}
