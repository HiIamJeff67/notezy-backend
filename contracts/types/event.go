package types

import (
	"time"

	"github.com/google/uuid"
)

const Version = "v1"

type Topic string

func (t Topic) String() string {
	return string(t)
}

type AggregateType string

type EventType string

type TraceMetadata struct {
	TraceParent string `json:"traceParent,omitempty"`
	TraceState  string `json:"traceState,omitempty"`
}

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
