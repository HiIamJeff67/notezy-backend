package kafka

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	franzkgo "github.com/twmb/franz-go/pkg/kgo"

	coreeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/events"
)

func TestDecodeEventEnvelopeRejectsMismatchedKafkaKey(t *testing.T) {
	aggregateId := uuid.New()
	payload, err := json.Marshal(coreeventscontract.EventEnvelope[json.RawMessage]{
		SchemaVersion: coreeventscontract.Version,
		EventId:       uuid.New(),
		EventType:     coreeventscontract.EventType_BlockPackAccessRevoked,
		AggregateType: coreeventscontract.AggregateType_BlockPack,
		AggregateId:   aggregateId,
		KafkaKey:      aggregateId.String(),
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = decodeEventEnvelope(&franzkgo.Record{
		Key:   []byte(uuid.NewString()),
		Value: payload,
	})
	if err == nil {
		t.Fatal("expected Kafka key mismatch to be rejected")
	}
}

func TestNewConsumerRejectsMissingGroupOrTopics(t *testing.T) {
	_, err := NewConsumer(ConsumerConfig{}, coreeventscontract.CoreLifecycleTopic.String())
	if err == nil {
		t.Fatal("expected Kafka consumer group to be required")
	}

	_, err = NewConsumer(ConsumerConfig{
		ConsumerGroup: "test-group",
	})
	if err == nil {
		t.Fatal("expected Kafka consumer topic to be required")
	}
}

func TestErrorClassificationDefaultsToTransient(t *testing.T) {
	if errorClassification(errors.New("temporary failure")) != ErrorClassification_Transient {
		t.Fatal("expected unclassified error to be retried")
	}
	if errorClassification(&ConsumerError{
		Classification: ErrorClassification_PoisonMessage,
		Origin:         errors.New("invalid event"),
	}) != ErrorClassification_PoisonMessage {
		t.Fatal("expected explicit Kafka consumer error classification")
	}
	if errorClassification(&ConsumerError{
		Classification: "Unknown",
		Origin:         errors.New("unknown classification"),
	}) != ErrorClassification_Transient {
		t.Fatal("expected unknown Kafka consumer classification to be retried")
	}
}

func TestRetryBackoffCapsAtConfiguredMaximum(t *testing.T) {
	consumerConfig := ConsumerConfig{
		InitialRetryBackoff: time.Second,
		MaximumRetryBackoff: 4 * time.Second,
	}
	if retryBackoff(consumerConfig, 4) != 4*time.Second {
		t.Fatal("expected retry backoff to cap at configured maximum")
	}
}

func TestDeadLetterTopic(t *testing.T) {
	if DeadLetterTopic(coreeventscontract.CoreLifecycleTopic.String()) != "notezy.core.lifecycle.v1.dlq" {
		t.Fatal("unexpected Kafka dead-letter topic name")
	}
}
