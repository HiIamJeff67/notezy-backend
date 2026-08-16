package exceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notegic-backend/contracts/types/exceptions"
)

type DeliveryException struct {
	EmailException
}

func NewDeliveryException(domain string) DeliveryException {
	return DeliveryException{EmailException: NewEmailException(domain)}
}

func (e DeliveryException) DeliveryFailed(cause error) *exceptions.Exception {
	exception := exceptions.New("DeliveryFailed", e.Domain, "SendEmail", "Failed to deliver the email", http.StatusBadGateway, true)
	exception.Retryable = true
	return exception.WithOrigin(cause)
}

func (e DeliveryException) EnqueueFailed(cause error) *exceptions.Exception {
	exception := exceptions.New("EnqueueFailed", e.Domain, "EnqueueEmail", "Failed to enqueue the email", http.StatusInternalServerError, true)
	exception.Retryable = true
	return exception.WithOrigin(cause)
}
