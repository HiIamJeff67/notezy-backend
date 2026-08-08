package apiexceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
)

type itemExceptionDomain struct{}

var Item = itemExceptionDomain{}

func (itemExceptionDomain) NotFound() *exceptions.Exception {
	return exceptions.New(
		"NotFound",
		"Item",
		"Repository",
		"Item was not found",
		http.StatusNotFound,
	)
}

func (itemExceptionDomain) NoPermission(_ ...string) *exceptions.Exception {
	return exceptions.New(
		"PermissionDenied",
		"Item",
		"Authorize",
		"Permission is denied",
		http.StatusBadRequest,
	)
}
