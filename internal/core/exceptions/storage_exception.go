package apiexceptions

import (
	"fmt"
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
)

type storageExceptionDomain struct{}

var Storage = storageExceptionDomain{}

func (storageExceptionDomain) FailedToReadObjectBytes() *exceptions.Exception {
	return exceptions.New(
		"FailedToReadObjectBytes",
		"Storage",
		"Read",
		"Failed to read object bytes",
		http.StatusInternalServerError,
		true,
	)
}

func (storageExceptionDomain) FailedToPutObject(object any) *exceptions.Exception {
	return exceptions.New(
		"FailedToPutObject",
		"Storage",
		"Put",
		fmt.Sprintf("Failed to put object %v", object),
		http.StatusInternalServerError,
		true,
	)
}

func (storageExceptionDomain) FailedToPresignedGetObject(object any) *exceptions.Exception {
	return exceptions.New(
		"FailedToPresignedGetObject",
		"Storage",
		"PresignGet",
		fmt.Sprintf("Failed to presign object %v", object),
		http.StatusInternalServerError,
		true,
	)
}
