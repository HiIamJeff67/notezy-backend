package events

import (
	"time"

	"github.com/google/uuid"
)

type EventEnvelope[D any] struct {
	SchemaVersion string        `json:"schemaVersion"`
	EventId       uuid.UUID     `json:"eventId"`
	EventType     EventType     `json:"eventType"`
	AggregateType AggregateType `json:"aggregateType"`
	AggregateId   uuid.UUID     `json:"aggregateId"`
	KafkaKey      string        `json:"kafkaKey"`
	OccurredAt    time.Time     `json:"occurredAt"`
	CorrelationId string        `json:"correlationId"`
	CausationId   *uuid.UUID    `json:"causationId,omitempty"`
	Trace         TraceMetadata `json:"trace"`
	Data          D             `json:"data"`
}
