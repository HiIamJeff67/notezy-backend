package apiexceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notegic-backend/contracts/types/exceptions"
)

type SearchException struct {
	CoreException
}

func NewSearchException() SearchException {
	return SearchException{CoreException: NewCoreException("Search")}
}

func (SearchException) FailedToDecode() *exceptions.Exception {
	return exceptions.New(
		"CursorDecodeFailed",
		"Search",
		"Cursor",
		"Failed to decode the search cursor",
		http.StatusInternalServerError,
		true,
	)
}

func (SearchException) FailedToEncode() *exceptions.Exception {
	return exceptions.New(
		"CursorEncodeFailed",
		"Search",
		"Cursor",
		"Failed to encode the search cursor",
		http.StatusInternalServerError,
		true,
	)
}

func (SearchException) FailedToUnmarshalSearchCursor() *exceptions.Exception {
	return exceptions.New(
		"CursorEncodingFailed",
		"Search",
		"Cursor",
		"Failed to encode the search cursor",
		http.StatusInternalServerError,
		true,
	)
}
