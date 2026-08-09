package exceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
)

type EventException struct {
	NotificationException
}

func NewEventException(domain string) EventException {
	return EventException{NotificationException: NewNotificationException(domain)}
}

func (e EventException) UnsupportedEventType() *exceptions.Exception {
	return exceptions.New("UnsupportedEventType", e.Domain, "ConsumeEvent", "The notification event type is unsupported", http.StatusBadRequest)
}

func (e EventException) AggregateRecipientMismatch() *exceptions.Exception {
	return exceptions.New("AggregateRecipientMismatch", e.Domain, "ConsumeEvent", "The notification aggregate recipient does not match", http.StatusBadRequest)
}

func (e EventException) InvalidMetadata(cause error) *exceptions.Exception {
	return exceptions.New("InvalidNotificationMetadata", e.Domain, "ValidateEvent", "The notification metadata is invalid", http.StatusBadRequest).WithOrigin(cause)
}

func (e EventException) UnsupportedTemplateVersion(cause error) *exceptions.Exception {
	return exceptions.New("UnsupportedTemplateVersion", e.Domain, "ValidateEvent", "The notification template version is unsupported", http.StatusBadRequest).WithOrigin(cause)
}

func (e EventException) InvalidNewsTemplateKey() *exceptions.Exception {
	return exceptions.New("InvalidNewsTemplateKey", e.Domain, "ValidateEvent", "The news template key is invalid", http.StatusBadRequest)
}

func (e EventException) InvalidWarningTemplateKey() *exceptions.Exception {
	return exceptions.New("InvalidWarningTemplateKey", e.Domain, "ValidateEvent", "The warning template key is invalid", http.StatusBadRequest)
}

func (e EventException) InvalidImportantTemplateKey() *exceptions.Exception {
	return exceptions.New("InvalidImportantTemplateKey", e.Domain, "ValidateEvent", "The important template key is invalid", http.StatusBadRequest)
}

func (e EventException) UnsupportedType(cause error) *exceptions.Exception {
	return exceptions.New("UnsupportedNotificationType", e.Domain, "ValidateEvent", "The notification type is unsupported", http.StatusBadRequest).WithOrigin(cause)
}
