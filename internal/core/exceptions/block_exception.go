package apiexceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notegic-backend/contracts/types/exceptions"
)

type BlockException struct {
	CoreException
}

func NewBlockException() BlockException {
	return BlockException{CoreException: NewCoreException("Block")}
}

func (BlockException) NotFound() *exceptions.Exception {
	return exceptions.New(
		"NotFound",
		"Block",
		"Repository",
		"Block was not found",
		http.StatusNotFound,
	)
}

func (BlockException) FailedToCreate() *exceptions.Exception {
	return exceptions.New(
		"FailedToCreate",
		"Block",
		"Repository",
		"Failed to create the block",
		http.StatusInternalServerError,
		true,
	)
}

func (BlockException) FailedToUpdate() *exceptions.Exception {
	return exceptions.New(
		"FailedToUpdate",
		"Block",
		"Repository",
		"Failed to update the block",
		http.StatusInternalServerError,
		true,
	)
}

func (BlockException) NoChanges() *exceptions.Exception {
	return exceptions.New(
		"NoChanges",
		"Block",
		"Repository",
		"No changes were applied to the block",
		http.StatusNotModified,
	)
}

func (BlockException) FailedToCommitTransaction() *exceptions.Exception {
	return exceptions.New(
		"FailedToCommitTransaction",
		"Block",
		"Transaction",
		"Failed to commit the block transaction",
		http.StatusInternalServerError,
		true,
	)
}

func (BlockException) InvalidDto() *exceptions.Exception {
	return exceptions.New(
		"InvalidDto",
		"Block",
		"Validate",
		"Invalid block DTO",
		http.StatusBadRequest,
	)
}

func (BlockException) NoPermission(action string) *exceptions.Exception {
	return exceptions.New(
		"PermissionDenied",
		"Block",
		"Authorize",
		"Permission is denied to "+action,
		http.StatusForbidden,
	)
}
