package exceptions

import (
	"errors"
	"net/http"
	"testing"
)

func TestExceptionToPublicMasksInternalException(t *testing.T) {
	exception := New(
		"DatabaseFailure",
		"API",
		"CreateRootShelf",
		"insert failed: duplicate key",
		http.StatusInternalServerError,
		true,
	).WithOrigin(errors.New("duplicate key value violates unique constraint"))

	publicException := exception.ToPublic()
	if publicException.isInternal {
		t.Fatal("expected a public exception")
	}
	if publicException.Message != "An internal server error occurred" {
		t.Fatalf("unexpected public message: %s", publicException.Message)
	}
	if publicException.origin != nil || publicException.details != nil {
		t.Fatal("public exception must not include internal diagnostic data")
	}
	if publicException.httpStatusCode != http.StatusInternalServerError {
		t.Fatalf("unexpected status: %d", publicException.httpStatusCode)
	}
}

func TestExceptionToPublicUsesSafeFallback(t *testing.T) {
	exception := New(
		"StorageUnavailable",
		"API",
		"UploadMaterial",
		"S3 returned an internal response",
		http.StatusInternalServerError,
		true,
	).WithPublicFallback(PublicFallback{
		Reason:         "ServiceUnavailable",
		Domain:         "Storage",
		Operation:      "UploadMaterial",
		Message:        "Storage is temporarily unavailable",
		HTTPStatusCode: http.StatusServiceUnavailable,
		Retryable:      true,
	})

	publicException := exception.ToPublic()
	if publicException.Reason != "ServiceUnavailable" || !publicException.Retryable {
		t.Fatalf("unexpected public fallback: %#v", publicException)
	}
	if publicException.Message != "Storage is temporarily unavailable" {
		t.Fatalf("unexpected fallback message: %s", publicException.Message)
	}
}

func TestNewUsesOnlyTheFirstInternalFlag(t *testing.T) {
	exception := New(
		"InvalidInput",
		"Gateway",
		"BindRequest",
		"The request is invalid",
		http.StatusBadRequest,
		true,
		false,
	)
	if !exception.IsInternal() {
		t.Fatal("expected the first internal flag to be used")
	}
}

func TestExceptionClonePreservesFieldsAndOverridesHTTPStatus(t *testing.T) {
	exception := New(
		"ServiceUnavailable",
		"API",
		"CreateRootShelf",
		"The Core service is unavailable",
		http.StatusServiceUnavailable,
	).WithPublicFallback(PublicFallback{
		Reason:         "ServiceUnavailable",
		Domain:         "API",
		Operation:      "CreateRootShelf",
		Message:        "The Core service is unavailable",
		HTTPStatusCode: http.StatusServiceUnavailable,
		Retryable:      true,
	})

	clone := exception.Clone(http.StatusBadGateway)
	if clone == exception {
		t.Fatal("expected a distinct exception clone")
	}
	if clone.HTTPStatusCode() != http.StatusBadGateway || clone.Reason != exception.Reason {
		t.Fatalf("unexpected clone: %#v", clone)
	}
	clone.publicFallback.Message = "changed"
	if exception.publicFallback.Message == "changed" {
		t.Fatal("expected the public fallback to be cloned")
	}
}
