package kafka_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	franzkgo "github.com/twmb/franz-go/pkg/kgo"

	durablejobeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/durable-job/v1/events"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/events"
	platformkafka "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka"
	kafkatopics "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka/topics"
)

func TestCoreDurableJobKafkaBrokerFlow(t *testing.T) {
	if os.Getenv("NOTEZY_RUN_INTEGRATION") != "1" {
		t.Skip("set NOTEZY_RUN_INTEGRATION=1 to run Kafka broker integration tests")
	}

	brokers := configuredKafkaBrokers(t)
	producer, err := platformkafka.NewProducer(platformkafka.ClientConfig{
		ConnectionConfig: platformkafka.ConnectionConfig{
			Brokers:     brokers,
			DialTimeout: 10 * time.Second,
		},
		ClientId: "notezy-test-core-durablejob-producer",
	})
	if err != nil {
		t.Fatalf("create Kafka producer: %v", err)
	}
	t.Cleanup(producer.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	t.Cleanup(cancel)
	if err := producer.Ping(ctx); err != nil {
		t.Skipf("Kafka broker is unavailable: %v", err)
	}

	topic := durablejobeventscontract.CoreDurableJobRoutineTaskTopic.String()
	consumerConfig := platformkafka.ConsumerConfig{
		ClientConfig: platformkafka.ClientConfig{
			ConnectionConfig: platformkafka.ConnectionConfig{
				Brokers:     brokers,
				DialTimeout: 10 * time.Second,
			},
			ClientId: "notezy-test-core-durablejob-consumer",
		},
		ConsumerGroup:       "notezy-test-core-durablejob-" + uuid.NewString(),
		MaximumAttempts:     3,
		InitialRetryBackoff: 10 * time.Millisecond,
		MaximumRetryBackoff: 25 * time.Millisecond,
		MaximumPollRecords:  20,
	}

	var (
		mu       sync.Mutex
		attempts int
		sequence []int
		finished = make(chan struct{})
	)
	correlationId := uuid.NewString()
	consumer, err := platformkafka.NewConsumer(consumerConfig, topic)
	if err != nil {
		t.Fatalf("create Kafka consumer: %v", err)
	}
	t.Cleanup(consumer.Close)

	consumerContext, stopConsumer := context.WithCancel(ctx)
	t.Cleanup(stopConsumer)
	go func() {
		_ = consumer.Run(consumerContext, func(
			_ context.Context,
			_ platformkafka.ConsumerRecord,
			event eventcontract.EventEnvelope[json.RawMessage],
		) error {
			if event.CorrelationId != correlationId {
				return nil
			}

			var payload struct {
				Sequence int `json:"sequence"`
			}
			if err := json.Unmarshal(event.Data, &payload); err != nil {
				return &platformkafka.ConsumerError{
					Classification: platformkafka.ErrorClassification_SchemaIncompatible,
					Origin:         err,
				}
			}

			mu.Lock()
			attempts++
			if attempts < 3 {
				mu.Unlock()
				return &platformkafka.ConsumerError{
					Classification: platformkafka.ErrorClassification_Transient,
					Origin:         errors.New("simulated transient Core/DurableJob failure"),
				}
			}
			sequence = append(sequence, payload.Sequence)
			if len(sequence) == 3 {
				close(finished)
			}
			mu.Unlock()
			return nil
		})
	}()

	workerId := uuid.New()
	for index := 1; index <= 3; index++ {
		publishCoreDurableJobEvent(t, ctx, producer, topic, correlationId, workerId, index)
	}

	select {
	case <-finished:
	case <-ctx.Done():
		t.Fatalf("timed out waiting for Core/DurableJob events: %v", ctx.Err())
	}

	mu.Lock()
	defer mu.Unlock()
	if attempts != 5 {
		t.Fatalf("consumer attempts = %d, want 5 (3 retries for the first event, then two ordered events)", attempts)
	}
	if want := []int{1, 2, 3}; fmt.Sprint(sequence) != fmt.Sprint(want) {
		t.Fatalf("event sequence = %v, want %v", sequence, want)
	}
}

func TestCoreDurableJobKafkaSchemaIncompatibleEventGoesToDLQ(t *testing.T) {
	if os.Getenv("NOTEZY_RUN_INTEGRATION") != "1" {
		t.Skip("set NOTEZY_RUN_INTEGRATION=1 to run Kafka broker integration tests")
	}

	brokers := configuredKafkaBrokers(t)
	producer, err := platformkafka.NewProducer(platformkafka.ClientConfig{
		ConnectionConfig: platformkafka.ConnectionConfig{
			Brokers:     brokers,
			DialTimeout: 10 * time.Second,
		},
		ClientId: "notezy-test-schema-producer",
	})
	if err != nil {
		t.Fatalf("create Kafka producer: %v", err)
	}
	t.Cleanup(producer.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	t.Cleanup(cancel)
	if err := producer.Ping(ctx); err != nil {
		t.Skipf("Kafka broker is unavailable: %v", err)
	}

	topic := durablejobeventscontract.CoreDurableJobRoutineTaskTopic.String()
	deadLetterTopic := platformkafka.DeadLetterTopic(topic)
	sourceConsumer, err := platformkafka.NewConsumer(platformkafka.ConsumerConfig{
		ClientConfig: platformkafka.ClientConfig{
			ConnectionConfig: platformkafka.ConnectionConfig{
				Brokers:     brokers,
				DialTimeout: 10 * time.Second,
			},
			ClientId: "notezy-test-schema-source-consumer",
		},
		ConsumerGroup:       "notezy-test-schema-source-" + uuid.NewString(),
		MaximumAttempts:     1,
		InitialRetryBackoff: time.Millisecond,
		MaximumRetryBackoff: time.Millisecond,
		MaximumPollRecords:  10,
	}, topic)
	if err != nil {
		t.Fatalf("create Kafka source consumer: %v", err)
	}
	t.Cleanup(sourceConsumer.Close)

	dlqConsumer, err := franzkgo.NewClient(
		franzkgo.SeedBrokers(brokers...),
		franzkgo.ClientID("notezy-test-schema-dlq-consumer"),
		franzkgo.ConsumerGroup("notezy-test-schema-dlq-"+uuid.NewString()),
		franzkgo.ConsumeTopics(deadLetterTopic),
		franzkgo.DialTimeout(10*time.Second),
	)
	if err != nil {
		t.Fatalf("create Kafka DLQ consumer: %v", err)
	}
	t.Cleanup(dlqConsumer.Close)

	dlqReceived := make(chan platformkafka.DeadLetter, 1)
	workerId := uuid.New()
	consumerContext, stopConsumer := context.WithCancel(ctx)
	t.Cleanup(stopConsumer)
	go func() {
		_ = sourceConsumer.Run(consumerContext, func(
			_ context.Context,
			_ platformkafka.ConsumerRecord,
			_ eventcontract.EventEnvelope[json.RawMessage],
		) error {
			return errors.New("the incompatible event must be dead-lettered before reaching the handler")
		})
	}()
	go func() {
		for consumerContext.Err() == nil {
			fetches := dlqConsumer.PollRecords(consumerContext, 1)
			if fetches.Err() != nil {
				return
			}
			fetches.EachRecord(func(record *franzkgo.Record) {
				var deadLetter platformkafka.DeadLetter
				if err := json.Unmarshal(record.Value, &deadLetter); err == nil && deadLetter.Key == workerId.String() {
					dlqReceived <- deadLetter
				}
			})
		}
	}()

	payload, err := json.Marshal(eventcontract.EventEnvelope[map[string]any]{
		SchemaVersion: "v999",
		EventId:       uuid.New(),
		EventType:     durablejobeventscontract.EventType_RoutineTasksCompleted,
		AggregateType: durablejobeventscontract.AggregateType_DurableJobWorker,
		AggregateId:   workerId,
		KafkaKey:      workerId.String(),
		OccurredAt:    time.Now().UTC(),
		Data:          map[string]any{},
	})
	if err != nil {
		t.Fatalf("marshal incompatible event: %v", err)
	}
	if err := producer.Produce(ctx, topic, workerId.String(), payload); err != nil {
		t.Fatalf("publish incompatible event: %v", err)
	}

	select {
	case deadLetter := <-dlqReceived:
		if deadLetter.SourceTopic != topic {
			t.Fatalf("DLQ source topic = %q, want %q", deadLetter.SourceTopic, topic)
		}
		if deadLetter.Classification != platformkafka.ErrorClassification_SchemaIncompatible {
			t.Fatalf("DLQ classification = %q, want %q", deadLetter.Classification, platformkafka.ErrorClassification_SchemaIncompatible)
		}
		if deadLetter.Attempts != 1 {
			t.Fatalf("DLQ attempts = %d, want 1", deadLetter.Attempts)
		}
	case <-ctx.Done():
		t.Fatalf("timed out waiting for Kafka DLQ record: %v", ctx.Err())
	}
}

func configuredKafkaBrokers(t *testing.T) []string {
	t.Helper()

	values := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	brokers := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			brokers = append(brokers, value)
		}
	}
	if len(brokers) > 0 {
		ensureKafkaTopics(t, brokers)
		return brokers
	}

	t.Skip("KAFKA_BROKERS is not set; start the integration Compose stack first")
	return nil
}

func ensureKafkaTopics(t *testing.T, brokers []string) {
	t.Helper()

	provisioner, err := platformkafka.NewTopicProvisioner(platformkafka.ClientConfig{
		ConnectionConfig: platformkafka.ConnectionConfig{
			Brokers:     brokers,
			DialTimeout: 10 * time.Second,
		},
		ClientId: "notezy-test-kafka-topic-bootstrap",
	})
	if err != nil {
		t.Fatalf("create Kafka topic provisioner: %v", err)
	}
	t.Cleanup(provisioner.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	if err := provisioner.EnsureTopics(ctx, kafkatopics.All()); err != nil {
		t.Fatalf("ensure Kafka topics: %v", err)
	}
}

func publishCoreDurableJobEvent(
	t *testing.T,
	ctx context.Context,
	producer *platformkafka.Producer,
	topic string,
	correlationId string,
	workerId uuid.UUID,
	sequence int,
) {
	t.Helper()

	payload, err := json.Marshal(eventcontract.EventEnvelope[map[string]int]{
		SchemaVersion: eventcontract.Version,
		EventId:       uuid.New(),
		EventType:     durablejobeventscontract.EventType_RoutineTasksCompleted,
		AggregateType: durablejobeventscontract.AggregateType_DurableJobWorker,
		AggregateId:   workerId,
		KafkaKey:      workerId.String(),
		OccurredAt:    time.Now().UTC(),
		CorrelationId: correlationId,
		Data:          map[string]int{"sequence": sequence},
	})
	if err != nil {
		t.Fatalf("marshal Core/DurableJob event: %v", err)
	}
	if err := producer.Produce(ctx, topic, workerId.String(), payload); err != nil {
		t.Fatalf("publish Core/DurableJob event: %v", err)
	}
}
