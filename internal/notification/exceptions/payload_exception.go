package exceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
)

type PayloadException struct {
	NotificationException
}

func NewPayloadException(domain string) PayloadException {
	return PayloadException{NotificationException: NewNotificationException(domain)}
}

func (e PayloadException) PayloadDecodeFailed(cause error) *exceptions.Exception {
	return exceptions.New("PayloadDecodeFailed", e.Domain, "DecodePayload", "Failed to decode the notification payload", http.StatusBadRequest).WithOrigin(cause)
}

func (e PayloadException) InvalidNewsPayload(cause error) *exceptions.Exception {
	return exceptions.New("InvalidNewsPayload", e.Domain, "ValidatePayload", "The news notification payload is invalid", http.StatusBadRequest).WithOrigin(cause)
}

func (e PayloadException) InvalidWarningPayload(cause error) *exceptions.Exception {
	return exceptions.New("InvalidWarningPayload", e.Domain, "ValidatePayload", "The warning notification payload is invalid", http.StatusBadRequest).WithOrigin(cause)
}

func (e PayloadException) InvalidImportantPayload(cause error) *exceptions.Exception {
	return exceptions.New("InvalidImportantPayload", e.Domain, "ValidatePayload", "The important notification payload is invalid", http.StatusBadRequest).WithOrigin(cause)
}

func (e PayloadException) ResponsePayloadDecodeFailed(cause error) *exceptions.Exception {
	exception := exceptions.New("ResponsePayloadDecodeFailed", e.Domain, "SearchPrivateNotifications", "Failed to decode a notification response payload", http.StatusInternalServerError, true)
	exception.Retryable = true
	return exception.WithOrigin(cause)
}
