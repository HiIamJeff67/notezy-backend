package inputs

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	coreeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/events"
)

type CreateOutboxEventInput struct {
	Id            uuid.UUID                        `json:"id"`
	AggregateType coreeventscontract.AggregateType `json:"aggregateType"`
	AggregateId   uuid.UUID                        `json:"aggregateId"`
	EventType     coreeventscontract.EventType     `json:"eventType"`
	Topic         coreeventscontract.Topic         `json:"topic"`
	KafkaKey      string                           `json:"kafkaKey"`
	Payload       json.RawMessage                  `json:"payload"`
	Metadata      json.RawMessage                  `json:"metadata"`
	AvailableAt   time.Time                        `json:"availableAt"`
}

type FailedOutboxEventInput struct {
	Id          uuid.UUID `json:"id"`
	LastError   string    `json:"lastError"`
	AvailableAt time.Time `json:"availableAt"`
}
