package apiexceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notegic-backend/contracts/types/exceptions"
)

type ItemException struct {
	CoreException
}

func NewItemException() ItemException {
	return ItemException{CoreException: NewCoreException("Item")}
}

func (ItemException) NotFound() *exceptions.Exception {
	return exceptions.New(
		"NotFound",
		"Item",
		"Repository",
		"Item was not found",
		http.StatusNotFound,
	)
}

func (ItemException) NoPermission(_ ...string) *exceptions.Exception {
	return exceptions.New(
		"PermissionDenied",
		"Item",
		"Authorize",
		"Permission is denied",
		http.StatusBadRequest,
	)
}
