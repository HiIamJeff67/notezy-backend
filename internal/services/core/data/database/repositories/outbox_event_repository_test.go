package repositories

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	eventscontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/events"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas"
)

func TestNewCreateOutboxEventInputAndSerializePreserveEventContract(t *testing.T) {
	eventId := uuid.New()
	aggregateId := uuid.New()
	occurredAt := time.Now().UTC().Round(0)
	envelope := eventscontract.EventEnvelope[eventscontract.BlockPackAccessRevokedData]{
		SchemaVersion: eventscontract.Version,
		EventId:       eventId,
		EventType:     eventscontract.EventType_BlockPackAccessRevoked,
		AggregateType: eventscontract.AggregateType_BlockPack,
		AggregateId:   aggregateId,
		KafkaKey:      aggregateId.String(),
		OccurredAt:    occurredAt,
		CorrelationId: "request-123",
		Trace: eventscontract.TraceMetadata{
			TraceParent: "00-trace",
		},
		Data: eventscontract.BlockPackAccessRevokedData{},
	}

	createInput, err := newCreateOutboxEventInput(eventscontract.CoreLifecycleTopic, envelope)
	if err != nil {
		t.Fatalf("failed to create outbox input: %v", err)
	}
	payload, err := SerializeOutboxEvent(schemas.OutboxEvent{
		Id:            createInput.Id,
		AggregateType: createInput.AggregateType,
		AggregateId:   createInput.AggregateId,
		EventType:     createInput.EventType,
		Topic:         createInput.Topic,
		KafkaKey:      createInput.KafkaKey,
		Payload:       datatypes.JSON(createInput.Payload),
		Metadata:      datatypes.JSON(createInput.Metadata),
	})
	if err != nil {
		t.Fatalf("failed to serialize outbox event: %v", err)
	}

	var serialized eventscontract.EventEnvelope[eventscontract.BlockPackAccessRevokedData]
	if err := json.Unmarshal(payload, &serialized); err != nil {
		t.Fatalf("failed to decode serialized event: %v", err)
	}
	if serialized.EventId != eventId || serialized.AggregateId != aggregateId ||
		serialized.KafkaKey != aggregateId.String() || serialized.CorrelationId != "request-123" ||
		serialized.Trace.TraceParent != "00-trace" {
		t.Fatalf("serialized event lost contract fields: %#v", serialized)
	}
}

func TestNewCreateOutboxEventInputRejectsMismatchedKafkaKey(t *testing.T) {
	_, err := newCreateOutboxEventInput(
		eventscontract.CoreLifecycleTopic,
		eventscontract.EventEnvelope[eventscontract.UserSessionsRevokedData]{
			SchemaVersion: eventscontract.Version,
			EventId:       uuid.New(),
			EventType:     eventscontract.EventType_UserSessionsRevoked,
			AggregateType: eventscontract.AggregateType_User,
			AggregateId:   uuid.New(),
			KafkaKey:      "another-aggregate",
			OccurredAt:    time.Now(),
			CorrelationId: "request-123",
			Data:          eventscontract.UserSessionsRevokedData{},
		},
	)
	if err == nil {
		t.Fatal("expected mismatched Kafka key to be rejected")
	}
}
