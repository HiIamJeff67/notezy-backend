package apiexceptions

import (
	"fmt"
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
)

type StorageException struct {
	CoreException
}

func NewStorageException() StorageException {
	return StorageException{CoreException: NewCoreException("Storage")}
}

func (StorageException) FailedToReadObjectBytes() *exceptions.Exception {
	return exceptions.New(
		"FailedToReadObjectBytes",
		"Storage",
		"Read",
		"Failed to read object bytes",
		http.StatusInternalServerError,
		true,
	)
}

func (StorageException) FailedToPutObject(object any) *exceptions.Exception {
	return exceptions.New(
		"FailedToPutObject",
		"Storage",
		"Put",
		fmt.Sprintf("Failed to put object %v", object),
		http.StatusInternalServerError,
		true,
	)
}

func (StorageException) FailedToPresignedGetObject(object any) *exceptions.Exception {
	return exceptions.New(
		"FailedToPresignedGetObject",
		"Storage",
		"PresignGet",
		fmt.Sprintf("Failed to presign object %v", object),
		http.StatusInternalServerError,
		true,
	)
}
