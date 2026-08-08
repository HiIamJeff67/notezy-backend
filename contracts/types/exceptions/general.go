package exceptions

import (
	"net/http"
	"strings"
)

func InternalServerError(optionalMessage ...string) *Exception {
	message := "An internal server error occurred"
	if len(optionalMessage) > 0 && strings.TrimSpace(optionalMessage[0]) != "" {
		message = optionalMessage[0]
	}

	return New(
		"InternalServerError",
		"General",
		"Respond",
		message,
		http.StatusInternalServerError,
	)
}

func InvalidDto(domain string, optionalMessage ...string) *Exception {
	message := "Invalid " + domain + " request data"
	if len(optionalMessage) > 0 && strings.TrimSpace(optionalMessage[0]) != "" {
		message = optionalMessage[0]
	}

	return New(
		"InvalidDto",
		domain,
		"Bind",
		message,
		http.StatusBadRequest,
	)
}

func InvalidInput(domain string, optionalMessage ...string) *Exception {
	message := "Invalid " + domain + " request input"
	if len(optionalMessage) > 0 && strings.TrimSpace(optionalMessage[0]) != "" {
		message = optionalMessage[0]
	}

	return New(
		"InvalidInput",
		domain,
		"Bind",
		message,
		http.StatusBadRequest,
	)
}
