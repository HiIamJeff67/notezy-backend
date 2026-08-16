package apiexceptions

import (
	"fmt"
	"net/http"
	"strings"

	exceptions "github.com/HiIamJeff67/notegic-backend/contracts/types/exceptions"
)

type CoreException struct {
	Domain string
}

func NewCoreException(domain string) CoreException {
	return CoreException{Domain: domain}
}

func (e CoreException) NotFound(optionalMessage ...string) *exceptions.Exception {
	message := e.Domain + " was not found"
	if len(optionalMessage) > 0 && strings.TrimSpace(optionalMessage[0]) != "" {
		message = optionalMessage[0]
	}

	return exceptions.New("NotFound", e.Domain, "Repository", message, http.StatusNotFound)
}

func (e CoreException) FailedToCreate(optionalMessage ...string) *exceptions.Exception {
	message := "Failed to create " + e.Domain
	if len(optionalMessage) > 0 && strings.TrimSpace(optionalMessage[0]) != "" {
		message = optionalMessage[0]
	}

	return exceptions.New("FailedToCreate", e.Domain, "Repository", message, http.StatusInternalServerError, true)
}

func (e CoreException) FailedToUpdate(optionalMessage ...string) *exceptions.Exception {
	message := "Failed to update " + e.Domain
	if len(optionalMessage) > 0 && strings.TrimSpace(optionalMessage[0]) != "" {
		message = optionalMessage[0]
	}

	return exceptions.New("FailedToUpdate", e.Domain, "Repository", message, http.StatusInternalServerError, true)
}

func (e CoreException) FailedToDelete(optionalMessage ...string) *exceptions.Exception {
	message := "Failed to delete " + e.Domain
	if len(optionalMessage) > 0 && strings.TrimSpace(optionalMessage[0]) != "" {
		message = optionalMessage[0]
	}

	return exceptions.New("FailedToDelete", e.Domain, "Repository", message, http.StatusInternalServerError, true)
}

func (e CoreException) NoChanges() *exceptions.Exception {
	return exceptions.New("NoChanges", e.Domain, "Repository", "No changes were applied to "+e.Domain, http.StatusNotModified)
}

func (e CoreException) FailedToCommitTransaction() *exceptions.Exception {
	return exceptions.New("FailedToCommitTransaction", e.Domain, "Transaction", "Failed to commit the "+e.Domain+" transaction", http.StatusInternalServerError, true)
}

func (e CoreException) NoPermission(action string) *exceptions.Exception {
	return exceptions.New("PermissionDenied", e.Domain, "Authorize", "Permission is denied to "+action, http.StatusForbidden)
}

func (e CoreException) InvalidInput(optionalMessage ...string) *exceptions.Exception {
	message := "Invalid " + e.Domain + " input"
	if len(optionalMessage) > 0 && strings.TrimSpace(optionalMessage[0]) != "" {
		message = optionalMessage[0]
	}

	return exceptions.New("InvalidInput", e.Domain, "Validate", message, http.StatusBadRequest)
}

func (e CoreException) InvalidDto(optionalMessage ...string) *exceptions.Exception {
	message := "Invalid " + e.Domain + " DTO"
	if len(optionalMessage) > 0 && strings.TrimSpace(optionalMessage[0]) != "" {
		message = optionalMessage[0]
	}

	return exceptions.New("InvalidDto", e.Domain, "Validate", message, http.StatusBadRequest)
}

func (e CoreException) InvalidType(value any) *exceptions.Exception {
	return exceptions.New("InvalidType", e.Domain, "Validate", "Invalid type in "+e.Domain, http.StatusInternalServerError, true).WithDetails(map[string]any{
		"actualType": fmt.Sprintf("%T", value),
		"value":      value,
	})
}

func (e CoreException) FailedToCompileRegularExpression() *exceptions.Exception {
	return exceptions.New("FailedToCompileRegularExpression", e.Domain, "Validate", "Failed to compile regular expression", http.StatusInternalServerError, true)
}

func (e CoreException) CannotGetFileObjects() *exceptions.Exception {
	return exceptions.New("CannotGetFileObjects", e.Domain, "File", "Failed to get file objects", http.StatusInternalServerError, true)
}

func (e CoreException) FailedToMarshalData(data any) *exceptions.Exception {
	return exceptions.New("FailedToMarshal", e.Domain, "Marshal", fmt.Sprintf("Failed to marshal data of %v", data), http.StatusInternalServerError, true)
}
