package durablejobexceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
)

type searchExceptionDomain struct{}

var Search = searchExceptionDomain{}

func (searchExceptionDomain) FailedToDecode() *exceptions.Exception {
	return exceptions.New(
		"CursorDecodeFailed",
		"Search",
		"Cursor",
		"Failed to decode the search cursor",
		http.StatusInternalServerError,
		true,
	)
}

func (searchExceptionDomain) FailedToEncode() *exceptions.Exception {
	return exceptions.New(
		"CursorEncodeFailed",
		"Search",
		"Cursor",
		"Failed to encode the search cursor",
		http.StatusInternalServerError,
		true,
	)
}

func (searchExceptionDomain) FailedToUnmarshalSearchCursor() *exceptions.Exception {
	return exceptions.New(
		"CursorEncodingFailed",
		"Search",
		"Cursor",
		"Failed to encode the search cursor",
		http.StatusInternalServerError,
		true,
	)
}
