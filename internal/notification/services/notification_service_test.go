package services

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	coreeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/events"
	notificationtypescontract "github.com/HiIamJeff67/notezy-backend/contracts/notification/v1/types"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/events"

	sharedvalidations "github.com/HiIamJeff67/notezy-backend/shared/validations"

	schemas "github.com/HiIamJeff67/notezy-backend/internal/notification/data/database/schemas"
	notificationvalidations "github.com/HiIamJeff67/notezy-backend/internal/notification/validations"
)

type notificationRepositoryStub struct {
	createCalls        int
	createErr          error
	deleteForUserCalls int
	deleteForUserErr   error
}

func (r *notificationRepositoryStub) CreateFromRequest(
	context.Context,
	eventcontract.EventEnvelope[coreeventscontract.NotificationRequestedData],
) error {
	r.createCalls++
	return r.createErr
}

func (r *notificationRepositoryStub) List(context.Context, uuid.UUID, *time.Time, int) ([]schemas.Notification, error) {
	return nil, nil
}

func (r *notificationRepositoryStub) CountUnread(context.Context, uuid.UUID) (int64, error) {
	return 0, nil
}

func (r *notificationRepositoryStub) MarkRead(context.Context, uuid.UUID, []uuid.UUID) (int64, error) {
	return 0, nil
}

func (r *notificationRepositoryStub) SoftDelete(context.Context, uuid.UUID, []uuid.UUID) (int64, error) {
	return 0, nil
}

func (r *notificationRepositoryStub) DeleteForUser(context.Context, uuid.UUID) (int64, error) {
	r.deleteForUserCalls++
	return 0, r.deleteForUserErr
}

func (r *notificationRepositoryStub) DeleteExpired(context.Context, time.Time, time.Duration) (int64, error) {
	return 0, nil
}

func (r *notificationRepositoryStub) ClaimOutbox(context.Context, string, int, time.Duration) ([]schemas.OutboxEvent, error) {
	return nil, nil
}

func (r *notificationRepositoryStub) MarkOutboxPublished(context.Context, string, []uuid.UUID) error {
	return nil
}

func (r *notificationRepositoryStub) MarkOutboxFailed(context.Context, string, []uuid.UUID, string, time.Time) error {
	return nil
}

func (r *notificationRepositoryStub) DeletePublishedOutbox(context.Context, time.Time) (int64, error) {
	return 0, nil
}

func newNotificationServiceForTest(repository *notificationRepositoryStub) *NotificationService {
	validate := validator.New()
	sharedvalidations.RegisterStringsValidation(validate)
	sharedvalidations.RegisterTimesValidation(validate)
	notificationvalidations.RegisterNotificationValidation(validate)
	notificationvalidations.RegisterNewsValidation(validate)
	notificationvalidations.RegisterWarningValidation(validate)
	notificationvalidations.RegisterImportantValidation(validate)

	return NewNotificationService(repository, validate)
}

func TestConsumeRequestedValidatesPayloadBeforePersisting(t *testing.T) {
	repository := &notificationRepositoryStub{}
	service := newNotificationServiceForTest(repository)
	recipientUserPublicId := uuid.New()
	payload, err := json.Marshal(notificationtypescontract.NewsPayload{
		Title:   "Release update",
		Summary: "A new release is available.",
		Body:    "Read the release notes for more details.",
	})
	if err != nil {
		t.Fatalf("failed to encode payload: %v", err)
	}

	event := eventcontract.EventEnvelope[coreeventscontract.NotificationRequestedData]{
		SchemaVersion: eventcontract.Version,
		EventId:       uuid.New(),
		EventType:     coreeventscontract.EventType_NotificationRequested,
		AggregateType: coreeventscontract.AggregateType_Notification,
		AggregateId:   recipientUserPublicId,
		KafkaKey:      recipientUserPublicId.String(),
		OccurredAt:    time.Now().UTC(),
		CorrelationId: "request-id",
		Data: coreeventscontract.NotificationRequestedData{
			RecipientUserPublicId: recipientUserPublicId,
			Type:                  coreeventscontract.NotificationType_News,
			Priority:              coreeventscontract.NotificationPriority_Normal,
			TemplateKey:           notificationtypescontract.TemplateKey_News,
			TemplateVersion:       1,
			Payload:               payload,
			DedupeKey:             "welcome:" + recipientUserPublicId.String(),
		},
	}

	if err := service.ConsumeRequested(context.Background(), event); err != nil {
		t.Fatalf("expected valid notification request, got %v", err)
	}
	if repository.createCalls != 1 {
		t.Fatalf("expected one repository call, got %d", repository.createCalls)
	}
}

func TestConsumeRequestedRejectsInvalidPayloadBeforePersisting(t *testing.T) {
	repository := &notificationRepositoryStub{}
	service := newNotificationServiceForTest(repository)
	recipientUserPublicId := uuid.New()
	payload, err := json.Marshal(notificationtypescontract.NewsPayload{Title: ""})
	if err != nil {
		t.Fatalf("failed to encode payload: %v", err)
	}

	event := eventcontract.EventEnvelope[coreeventscontract.NotificationRequestedData]{
		EventId:       uuid.New(),
		EventType:     coreeventscontract.EventType_NotificationRequested,
		AggregateType: coreeventscontract.AggregateType_Notification,
		AggregateId:   recipientUserPublicId,
		KafkaKey:      recipientUserPublicId.String(),
		OccurredAt:    time.Now().UTC(),
		Data: coreeventscontract.NotificationRequestedData{
			RecipientUserPublicId: recipientUserPublicId,
			Type:                  coreeventscontract.NotificationType_News,
			Priority:              coreeventscontract.NotificationPriority_Normal,
			TemplateKey:           notificationtypescontract.TemplateKey_News,
			TemplateVersion:       1,
			Payload:               payload,
			DedupeKey:             "invalid:" + recipientUserPublicId.String(),
		},
	}

	if err := service.ConsumeRequested(context.Background(), event); err == nil {
		t.Fatal("expected invalid notification payload to be rejected")
	}
	if repository.createCalls != 0 {
		t.Fatalf("expected no repository call for invalid payload, got %d", repository.createCalls)
	}
}

func TestDeleteForUserDelegatesToRepository(t *testing.T) {
	repository := &notificationRepositoryStub{}
	service := newNotificationServiceForTest(repository)

	if err := service.DeleteForUser(context.Background(), uuid.New()); err != nil {
		t.Fatalf("expected user notification deletion to succeed, got %v", err)
	}
	if repository.deleteForUserCalls != 1 {
		t.Fatalf("expected one user notification deletion call, got %d", repository.deleteForUserCalls)
	}
}
