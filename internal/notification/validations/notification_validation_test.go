package validations

import (
	"testing"

	validator "github.com/go-playground/validator/v10"

	notificationtypescontract "github.com/HiIamJeff67/notezy-backend/contracts/notification/v1/types"
	sharedvalidations "github.com/HiIamJeff67/notezy-backend/shared/validations"
)

func TestRegisterNotificationValidations(t *testing.T) {
	validate := validator.New()
	sharedvalidations.RegisterStringsValidation(validate)
	sharedvalidations.RegisterTimesValidation(validate)
	RegisterNotificationValidation(validate)
	RegisterNewsValidation(validate)
	RegisterWarningValidation(validate)
	RegisterImportantValidation(validate)

	if err := validate.Struct(notificationtypescontract.NotificationMetadata{
		Type:            "news",
		Priority:        "normal",
		TemplateVersion: 1,
	}); err != nil {
		t.Fatalf("expected valid notification metadata, got %v", err)
	}
	if err := validate.Struct(notificationtypescontract.NotificationMetadata{
		Type:            "unknown",
		Priority:        "normal",
		TemplateVersion: 1,
	}); err == nil {
		t.Fatal("expected notification metadata validation error")
	}

	if err := validate.Struct(notificationtypescontract.NewsPayload{
		Title:   "Release update",
		Summary: "A new release is available.",
		Body:    "Read the release notes for more details.",
	}); err != nil {
		t.Fatalf("expected valid news payload, got %v", err)
	}
	if err := validate.Struct(notificationtypescontract.WarningPayload{Title: "Security warning"}); err == nil {
		t.Fatal("expected warning payload validation error")
	}
}
