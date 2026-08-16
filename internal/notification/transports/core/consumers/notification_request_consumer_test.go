package consumers

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	coreeventscontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/events"
	eventcontract "github.com/HiIamJeff67/notegic-backend/contracts/types/events"

	platformkafka "github.com/HiIamJeff67/notegic-backend/shared/platform/kafka"

	services "github.com/HiIamJeff67/notegic-backend/internal/notification/services"
)

type notificationServiceStub struct {
	services.NotificationServiceInterface
	consumeCalls          int
	consumeErr            error
	deleteCalls           int
	deleteErr             error
	lastNotificationEvent eventcontract.EventEnvelope[coreeventscontract.NotificationRequestedData]
	lastDeletedUser       uuid.UUID
}

func (s *notificationServiceStub) ConsumeNotificationRequested(
	_ context.Context,
	event eventcontract.EventEnvelope[coreeventscontract.NotificationRequestedData],
) error {
	s.consumeCalls++
	s.lastNotificationEvent = event
	return s.consumeErr
}

func (s *notificationServiceStub) DeleteAllNotificationsForUser(
	_ context.Context,
	userPublicId uuid.UUID,
) error {
	s.deleteCalls++
	s.lastDeletedUser = userPublicId
	return s.deleteErr
}

func TestNotificationRequestConsumerConsumesNotificationRequestedEvent(t *testing.T) {
	recipientUserPublicId := uuid.New()
	service := &notificationServiceStub{}
	consumer := &NotificationRequestConsumer{service: service}
	event := eventcontract.EventEnvelope[json.RawMessage]{
		SchemaVersion: eventcontract.Version,
		EventId:       uuid.New(),
		EventType:     coreeventscontract.EventType_NotificationRequested,
		AggregateType: coreeventscontract.AggregateType_Notification,
		AggregateId:   recipientUserPublicId,
		KafkaKey:      recipientUserPublicId.String(),
		OccurredAt:    time.Now().UTC(),
		Data: mustMarshalNotificationRequest(t, coreeventscontract.NotificationRequestedData{
			RecipientUserPublicId: recipientUserPublicId,
			Type:                  coreeventscontract.NotificationType_News,
			Priority:              coreeventscontract.NotificationPriority_Normal,
			TemplateKey:           "news",
			TemplateVersion:       1,
			Payload:               json.RawMessage(`{"title":"Release"}`),
			DedupeKey:             "release:" + recipientUserPublicId.String(),
		}),
	}

	if err := consumer.consume(context.Background(), platformkafka.ConsumerRecord{}, event); err != nil {
		t.Fatalf("consume notification event: %v", err)
	}
	if service.consumeCalls != 1 {
		t.Fatalf("consume calls = %d, want 1", service.consumeCalls)
	}
	if service.lastNotificationEvent.Data.RecipientUserPublicId != recipientUserPublicId {
		t.Fatalf("recipient user public ID = %s, want %s", service.lastNotificationEvent.Data.RecipientUserPublicId, recipientUserPublicId)
	}
}

func TestNotificationRequestConsumerClassifiesServiceFailureAsTransient(t *testing.T) {
	service := &notificationServiceStub{consumeErr: errors.New("database unavailable")}
	consumer := &NotificationRequestConsumer{service: service}
	recipientUserPublicId := uuid.New()
	event := eventcontract.EventEnvelope[json.RawMessage]{
		SchemaVersion: eventcontract.Version,
		EventId:       uuid.New(),
		EventType:     coreeventscontract.EventType_NotificationRequested,
		AggregateType: coreeventscontract.AggregateType_Notification,
		AggregateId:   recipientUserPublicId,
		KafkaKey:      recipientUserPublicId.String(),
		Data: mustMarshalNotificationRequest(t, coreeventscontract.NotificationRequestedData{
			RecipientUserPublicId: recipientUserPublicId,
			Type:                  coreeventscontract.NotificationType_News,
			Priority:              coreeventscontract.NotificationPriority_Normal,
			TemplateKey:           "news",
			TemplateVersion:       1,
			Payload:               json.RawMessage(`{"title":"Release"}`),
			DedupeKey:             "release:" + recipientUserPublicId.String(),
		}),
	}

	resultErr := consumer.consume(context.Background(), platformkafka.ConsumerRecord{}, event)
	consumerErr, ok := resultErr.(*platformkafka.ConsumerError)
	if !ok {
		t.Fatalf("error type = %T, want *platformkafka.ConsumerError", resultErr)
	}
	if consumerErr.Classification != platformkafka.ErrorClassification_Transient {
		t.Fatalf("classification = %q, want %q", consumerErr.Classification, platformkafka.ErrorClassification_Transient)
	}
}

func TestNotificationRequestConsumerDeletesNotificationsForDeletedUser(t *testing.T) {
	userPublicId := uuid.New()
	service := &notificationServiceStub{}
	consumer := &NotificationRequestConsumer{service: service}
	event := eventcontract.EventEnvelope[json.RawMessage]{
		SchemaVersion: eventcontract.Version,
		EventId:       uuid.New(),
		EventType:     coreeventscontract.EventType_UserDeleted,
		AggregateType: coreeventscontract.AggregateType_User,
		AggregateId:   userPublicId,
		KafkaKey:      userPublicId.String(),
		Data:          mustMarshalUserDeleted(t, coreeventscontract.UserDeletedData{DeletedAt: time.Now().UTC()}),
	}

	if err := consumer.consume(context.Background(), platformkafka.ConsumerRecord{}, event); err != nil {
		t.Fatalf("consume user deleted event: %v", err)
	}
	if service.deleteCalls != 1 || service.lastDeletedUser != userPublicId {
		t.Fatalf("delete call = %d for %s, want one call for %s", service.deleteCalls, service.lastDeletedUser, userPublicId)
	}
}

func mustMarshalNotificationRequest(t *testing.T, data coreeventscontract.NotificationRequestedData) json.RawMessage {
	t.Helper()
	encoded, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshal notification request: %v", err)
	}
	return encoded
}

func mustMarshalUserDeleted(t *testing.T, data coreeventscontract.UserDeletedData) json.RawMessage {
	t.Helper()
	encoded, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshal user deleted event: %v", err)
	}
	return encoded
}
