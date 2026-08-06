package apiexceptions

import (
	"fmt"
	"net/http"
	"strings"

	exceptions "github.com/HiIamJeff67/notezy-backend/shared/exceptions"
)

type domainException struct {
	domain string
}

func newDomainException(domain string) domainException {
	return domainException{
		domain: domain,
	}
}

func (d domainException) NotFound(optionalMessage ...string) *exceptions.Exception {
	message := d.domain + " was not found"
	if len(optionalMessage) > 0 && strings.TrimSpace(optionalMessage[0]) != "" {
		message = optionalMessage[0]
	}

	return exceptions.New(
		"NotFound",
		d.domain,
		"Repository",
		message,
		http.StatusNotFound,
	)
}

func (d domainException) FailedToCreate(optionalMessage ...string) *exceptions.Exception {
	message := "Failed to create " + d.domain
	if len(optionalMessage) > 0 && strings.TrimSpace(optionalMessage[0]) != "" {
		message = optionalMessage[0]
	}

	return exceptions.New(
		"FailedToCreate",
		d.domain,
		"Repository",
		message,
		http.StatusInternalServerError,
		true,
	)
}

func (d domainException) FailedToUpdate(optionalMessage ...string) *exceptions.Exception {
	message := "Failed to update " + d.domain
	if len(optionalMessage) > 0 && strings.TrimSpace(optionalMessage[0]) != "" {
		message = optionalMessage[0]
	}

	return exceptions.New(
		"FailedToUpdate",
		d.domain,
		"Repository",
		message,
		http.StatusInternalServerError,
		true,
	)
}

func (d domainException) FailedToDelete(optionalMessage ...string) *exceptions.Exception {
	message := "Failed to delete " + d.domain
	if len(optionalMessage) > 0 && strings.TrimSpace(optionalMessage[0]) != "" {
		message = optionalMessage[0]
	}

	return exceptions.New(
		"FailedToDelete",
		d.domain,
		"Repository",
		message,
		http.StatusInternalServerError,
		true,
	)
}

func (d domainException) NoChanges() *exceptions.Exception {
	return exceptions.New(
		"NoChanges",
		d.domain,
		"Repository",
		"No changes were applied to "+d.domain,
		http.StatusNotModified,
	)
}

func (d domainException) FailedToCommitTransaction() *exceptions.Exception {
	return exceptions.New(
		"FailedToCommitTransaction",
		d.domain,
		"Transaction",
		"Failed to commit the "+d.domain+" transaction",
		http.StatusInternalServerError,
		true,
	)
}

func (d domainException) NoPermission(action string) *exceptions.Exception {
	return exceptions.New(
		"PermissionDenied",
		d.domain,
		"Authorize",
		"Permission is denied to "+action,
		http.StatusForbidden,
	)
}

func (d domainException) InvalidInput(optionalMessage ...string) *exceptions.Exception {
	message := "Invalid " + d.domain + " input"
	if len(optionalMessage) > 0 && strings.TrimSpace(optionalMessage[0]) != "" {
		message = optionalMessage[0]
	}

	return exceptions.New(
		"InvalidInput",
		d.domain,
		"Validate",
		message,
		http.StatusBadRequest,
	)
}

func (d domainException) InvalidDto(optionalMessage ...string) *exceptions.Exception {
	message := "Invalid " + d.domain + " DTO"
	if len(optionalMessage) > 0 && strings.TrimSpace(optionalMessage[0]) != "" {
		message = optionalMessage[0]
	}

	return exceptions.New(
		"InvalidDto",
		d.domain,
		"Validate",
		message,
		http.StatusBadRequest,
	)
}

func (d domainException) InvalidType(value any) *exceptions.Exception {
	return exceptions.New(
		"InvalidType",
		d.domain,
		"Validate",
		"Invalid type in "+d.domain,
		http.StatusInternalServerError,
		true,
	).WithDetails(map[string]any{
		"actualType": fmt.Sprintf("%T", value),
		"value":      value,
	})
}

func (d domainException) FailedToCompileRegularExpression() *exceptions.Exception {
	return exceptions.New(
		"FailedToCompileRegularExpression",
		d.domain,
		"Validate",
		"Failed to compile regular expression",
		http.StatusInternalServerError,
		true,
	)
}

func (d domainException) CannotGetFileObjects() *exceptions.Exception {
	return exceptions.New(
		"CannotGetFileObjects",
		d.domain,
		"File",
		"Failed to get file objects",
		http.StatusInternalServerError,
		true,
	)
}

func (d domainException) FailedToMarshalData(data any) *exceptions.Exception {
	return exceptions.New(
		"FailedToMarshal",
		d.domain,
		"Marshal",
		fmt.Sprintf("Failed to marshal data of %v", data),
		http.StatusInternalServerError,
		true,
	)
}
