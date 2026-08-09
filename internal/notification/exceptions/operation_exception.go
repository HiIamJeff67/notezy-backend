package exceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
)

type OperationException struct {
	NotificationException
}

func NewOperationException(domain string) OperationException {
	return OperationException{NotificationException: NewNotificationException(domain)}
}

func (e OperationException) CreateFailed(cause error) *exceptions.Exception {
	exception := exceptions.New("CreateFailed", e.Domain, "CreateNotification", "Failed to create the notification", http.StatusInternalServerError, true)
	exception.Retryable = true
	return exception.WithOrigin(cause)
}

func (e OperationException) ListFailed(cause error) *exceptions.Exception {
	exception := exceptions.New("ListFailed", e.Domain, "ListNotifications", "Failed to list notifications", http.StatusInternalServerError, true)
	exception.Retryable = true
	return exception.WithOrigin(cause)
}

func (e OperationException) CountUnreadFailed(cause error) *exceptions.Exception {
	exception := exceptions.New("CountUnreadFailed", e.Domain, "CountUnreadNotifications", "Failed to count notifications", http.StatusInternalServerError, true)
	exception.Retryable = true
	return exception.WithOrigin(cause)
}

func (e OperationException) MarkReadFailed(cause error) *exceptions.Exception {
	exception := exceptions.New("MarkReadFailed", e.Domain, "MarkNotificationsRead", "Failed to mark notifications as read", http.StatusInternalServerError, true)
	exception.Retryable = true
	return exception.WithOrigin(cause)
}

func (e OperationException) DeleteFailed(cause error) *exceptions.Exception {
	exception := exceptions.New("DeleteFailed", e.Domain, "DeleteNotifications", "Failed to delete notifications", http.StatusInternalServerError, true)
	exception.Retryable = true
	return exception.WithOrigin(cause)
}

func (e OperationException) HardDeleteFailed(cause error) *exceptions.Exception {
	exception := exceptions.New("HardDeleteFailed", e.Domain, "HardDeleteNotifications", "Failed to hard-delete notifications", http.StatusInternalServerError, true)
	exception.Retryable = true
	return exception.WithOrigin(cause)
}

func (e OperationException) DeleteForUserFailed(cause error) *exceptions.Exception {
	exception := exceptions.New("DeleteForUserFailed", e.Domain, "DeleteForUser", "Failed to delete notifications for the user", http.StatusInternalServerError, true)
	exception.Retryable = true
	return exception.WithOrigin(cause)
}
