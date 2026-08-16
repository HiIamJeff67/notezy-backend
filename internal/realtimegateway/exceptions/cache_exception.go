package exceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notegic-backend/contracts/types/exceptions"
)

type CacheException struct {
	RealtimeGatewayException
}

func NewCacheException(domain string) CacheException {
	return CacheException{RealtimeGatewayException: NewRealtimeGatewayException(domain)}
}

func (e CacheException) Unavailable(cause error) *exceptions.Exception {
	exception := exceptions.New("CacheClientUnavailable", e.Domain, "AccessCache", "The realtime cache is unavailable", http.StatusServiceUnavailable, true)
	exception.Retryable = true
	return exception.WithOrigin(cause)
}

func (e CacheException) NotFound(cause error) *exceptions.Exception {
	return exceptions.New("NotFound", e.Domain, "GetCache", "The realtime cache record was not found", http.StatusNotFound).WithOrigin(cause)
}

func (e CacheException) DeserializationFailed(cause error) *exceptions.Exception {
	return exceptions.New("DeserializationFailed", e.Domain, "ReadCache", "Failed to deserialize the realtime cache record", http.StatusInternalServerError, true).WithOrigin(cause)
}

func (e CacheException) SerializationFailed(cause error) *exceptions.Exception {
	return exceptions.New("SerializationFailed", e.Domain, "WriteCache", "Failed to serialize the realtime cache record", http.StatusInternalServerError, true).WithOrigin(cause)
}

func (e CacheException) CreateFailed(cause error) *exceptions.Exception {
	return exceptions.New("FailedToCreate", e.Domain, "CreateCache", "Failed to create the realtime cache record", http.StatusInternalServerError, true).WithOrigin(cause)
}

func (e CacheException) InvalidRateLimitTokenCount() *exceptions.Exception {
	return exceptions.New("InvalidRateLimitTokenCount", e.Domain, "SynchronizeCache", "The rate-limit token count is invalid", http.StatusBadRequest)
}

func (e CacheException) UpdateFailed(cause error) *exceptions.Exception {
	return exceptions.New("FailedToUpdate", e.Domain, "UpdateCache", "Failed to update the realtime cache record", http.StatusInternalServerError, true).WithOrigin(cause)
}

func (e CacheException) DeleteFailed(cause error) *exceptions.Exception {
	return exceptions.New("FailedToDelete", e.Domain, "DeleteCache", "Failed to delete the realtime cache record", http.StatusInternalServerError, true).WithOrigin(cause)
}
