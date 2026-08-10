package exceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
)

type RequestException struct {
	NotificationException
}

func NewRequestException(domain string) RequestException {
	return RequestException{NotificationException: NewNotificationException(domain)}
}

func (e RequestException) RecipientRequired() *exceptions.Exception {
	return exceptions.New("RecipientUserPublicIdRequired", e.Domain, "ValidateRequest", "The recipient user public ID is required", http.StatusBadRequest)
}

func (e RequestException) InvalidSearchRequest(cause error) *exceptions.Exception {
	return exceptions.New("InvalidSearchRequest", e.Domain, "SearchPrivateNotifications", "The private notification search request is invalid", http.StatusBadRequest).WithOrigin(cause)
}

func (e RequestException) InvalidCountRequest(cause error) *exceptions.Exception {
	return exceptions.New("InvalidCountRequest", e.Domain, "CountUnreadNotifications", "The unread notification count request is invalid", http.StatusBadRequest).WithOrigin(cause)
}

func (e RequestException) InvalidMarkReadRequest(cause error) *exceptions.Exception {
	return exceptions.New("InvalidMarkReadRequest", e.Domain, "MarkNotificationsRead", "The mark-read request is invalid", http.StatusBadRequest).WithOrigin(cause)
}

func (e RequestException) InvalidDeleteRequest(cause error) *exceptions.Exception {
	return exceptions.New("InvalidDeleteRequest", e.Domain, "DeleteNotifications", "The delete notifications request is invalid", http.StatusBadRequest).WithOrigin(cause)
}

func (e RequestException) UserRequired() *exceptions.Exception {
	return exceptions.New("UserPublicIdRequired", e.Domain, "DeleteAllNotificationsForUser", "The user public ID is required", http.StatusBadRequest)
}
