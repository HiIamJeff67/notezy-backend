package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"log/slog"
	"time"

	"github.com/google/uuid"
	franzkgo "github.com/twmb/franz-go/pkg/kgo"

	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/events"

	logs "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/logs"
	traces "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/traces"
)

type ErrorClassification string

const (
	ErrorClassification_Transient          ErrorClassification = "Transient"
	ErrorClassification_PoisonMessage      ErrorClassification = "PoisonMessage"
	ErrorClassification_SchemaIncompatible ErrorClassification = "SchemaIncompatible"
)

type ConsumerError struct {
	Classification ErrorClassification
	Origin         error
}

func (e *ConsumerError) Error() string {
	if e == nil || e.Origin == nil {
		return "Kafka consumer error"
	}

	return e.Origin.Error()
}

func (e *ConsumerError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Origin
}

type ConsumerRecord struct {
	Topic     string
	Partition int32
	Offset    int64
	Key       string
}

type ConsumerHandler func(
	ctx context.Context,
	record ConsumerRecord,
	envelope eventcontract.EventEnvelope[json.RawMessage],
) error

type DeadLetter struct {
	SchemaVersion   string              `json:"schemaVersion"`
	ConsumerGroup   string              `json:"consumerGroup"`
	SourceTopic     string              `json:"sourceTopic"`
	SourcePartition int32               `json:"sourcePartition"`
	SourceOffset    int64               `json:"sourceOffset"`
	Key             string              `json:"key"`
	EventId         *uuid.UUID          `json:"eventId,omitempty"`
	Classification  ErrorClassification `json:"classification"`
	Attempts        int                 `json:"attempts"`
	Error           string              `json:"error"`
	Value           []byte              `json:"value"`
	FailedAt        time.Time           `json:"failedAt"`
}

type Consumer struct {
	client        *franzkgo.Client
	consumerGroup string
	config        ConsumerConfig
}

func NewConsumer(
	kafkaConfig ConsumerConfig,
	topics ...string,
) (*Consumer, error) {
	if kafkaConfig.ConsumerGroup == "" {
		return nil, errors.New("Kafka consumer group is required")
	}
	if len(topics) == 0 {
		return nil, errors.New("at least one Kafka consumer topic is required")
	}

	options, err := newConnectionOptions(kafkaConfig.ClientConfig)
	if err != nil {
		return nil, err
	}
	options = append(options,
		franzkgo.ConsumerGroup(kafkaConfig.ConsumerGroup),
		franzkgo.ConsumeTopics(topics...),
		franzkgo.DisableAutoCommit(),
		franzkgo.BlockRebalanceOnPoll(),
	)
	client, err := franzkgo.NewClient(options...)
	if err != nil {
		return nil, fmt.Errorf("create Kafka consumer: %w", err)
	}

	return &Consumer{
		client:        client,
		consumerGroup: kafkaConfig.ConsumerGroup,
		config:        kafkaConfig,
	}, nil
}

/* ============================== Auxiliary Methods ============================== */

func (c *Consumer) consumeFetches(
	ctx context.Context,
	fetches franzkgo.Fetches,
	handler ConsumerHandler,
) {
	fetches.EachPartition(func(fetchTopicPartition franzkgo.FetchTopicPartition) {
		if fetchTopicPartition.Err != nil {
			c.recordFailure(ctx, fetchTopicPartition.Topic, fetchTopicPartition.Partition, 0, "Kafka consumer partition fetch failed", fetchTopicPartition.Err)
			return
		}

		recordsToCommit := make([]*franzkgo.Record, 0, len(fetchTopicPartition.Records))
		for _, record := range fetchTopicPartition.Records {
			RecordConsumerLag(
				ctx,
				record.Topic,
				c.consumerGroup,
				fetchTopicPartition.HighWatermark-record.Offset,
			)
			if !c.consumeRecord(ctx, record, handler) {
				break
			}
			recordsToCommit = append(recordsToCommit, record)
		}
		if len(recordsToCommit) == 0 {
			return
		}
		if err := c.client.CommitRecords(ctx, recordsToCommit[len(recordsToCommit)-1]); err != nil {
			c.recordFailure(
				ctx,
				fetchTopicPartition.Topic,
				fetchTopicPartition.Partition,
				recordsToCommit[len(recordsToCommit)-1].Offset,
				"Failed to commit Kafka consumer offset",
				err,
			)
		}
	})
}

func (c *Consumer) consumeRecord(
	ctx context.Context,
	record *franzkgo.Record,
	handler ConsumerHandler,
) bool {
	var envelope eventcontract.EventEnvelope[json.RawMessage]
	if err := json.Unmarshal(record.Value, &envelope); err != nil {
		return c.deadLetter(
			ctx,
			record,
			nil,
			ErrorClassification_SchemaIncompatible,
			1,
			fmt.Errorf("decode Kafka event envelope: %w", err),
		)
	}
	if envelope.SchemaVersion != eventcontract.Version {
		return c.deadLetter(
			ctx,
			record,
			nil,
			ErrorClassification_SchemaIncompatible,
			1,
			fmt.Errorf("unsupported Kafka event schema version %q", envelope.SchemaVersion),
		)
	}
	if envelope.EventId == uuid.Nil || envelope.EventType == "" || envelope.AggregateType == "" ||
		envelope.AggregateId == uuid.Nil || envelope.KafkaKey == "" {
		return c.deadLetter(
			ctx,
			record,
			nil,
			ErrorClassification_SchemaIncompatible,
			1,
			errors.New("Kafka event envelope is incomplete"),
		)
	}
	if envelope.KafkaKey != envelope.AggregateId.String() || envelope.KafkaKey != string(record.Key) {
		return c.deadLetter(
			ctx,
			record,
			nil,
			ErrorClassification_SchemaIncompatible,
			1,
			errors.New("Kafka event envelope key does not match the aggregate ID"),
		)
	}

	if envelope.Trace.TraceParent != "" {
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier{
			"traceparent": envelope.Trace.TraceParent,
			"tracestate":  envelope.Trace.TraceState,
		})
	}
	if traces.NotezyTracer != nil {
		consumerCtx, span := traces.NotezyTracer.Start(ctx, "kafka.consume")
		defer traces.NotezyTracer.End(span, nil)
		ctx = consumerCtx
	}

	var err error
	for attempt := 1; attempt <= c.config.MaximumAttempts; attempt++ {
		startedAt := time.Now()
		err = handler(ctx, ConsumerRecord{
			Topic:     record.Topic,
			Partition: record.Partition,
			Offset:    record.Offset,
			Key:       string(record.Key),
		}, envelope)
		if err == nil {
			RecordConsume(ctx, record.Topic, c.consumerGroup, time.Since(startedAt))
			return true
		}

		classification := ErrorClassification_Transient
		var consumerErr *ConsumerError
		if errors.As(err, &consumerErr) {
			switch consumerErr.Classification {
			case ErrorClassification_PoisonMessage, ErrorClassification_SchemaIncompatible:
				classification = consumerErr.Classification
			}
		}
		if classification != ErrorClassification_Transient {
			return c.deadLetter(ctx, record, &envelope.EventId, classification, attempt, err)
		}
		if attempt == c.config.MaximumAttempts {
			return c.deadLetter(ctx, record, &envelope.EventId, classification, attempt, err)
		}

		RecordRetry(ctx, record.Topic, c.consumerGroup)
		c.recordFailure(ctx, record.Topic, record.Partition, record.Offset, "Retrying Kafka consumer event", err)
		backoff := c.config.InitialRetryBackoff
		for index := 1; index < attempt && backoff < c.config.MaximumRetryBackoff; index++ {
			backoff *= 2
		}
		if backoff > c.config.MaximumRetryBackoff {
			backoff = c.config.MaximumRetryBackoff
		}
		select {
		case <-ctx.Done():
			return false
		case <-time.After(backoff):
		}
	}

	return false
}

func (c *Consumer) deadLetter(
	ctx context.Context,
	record *franzkgo.Record,
	eventId *uuid.UUID,
	classification ErrorClassification,
	attempts int,
	err error,
) bool {
	if traces.NotezyTracer != nil {
		traces.NotezyTracer.RecordError(ctx, err)
	}

	deadLetter := DeadLetter{
		SchemaVersion:   eventcontract.Version,
		ConsumerGroup:   c.consumerGroup,
		SourceTopic:     record.Topic,
		SourcePartition: record.Partition,
		SourceOffset:    record.Offset,
		Key:             string(record.Key),
		EventId:         eventId,
		Classification:  classification,
		Attempts:        attempts,
		Error:           err.Error(),
		Value:           record.Value,
		FailedAt:        time.Now().UTC(),
	}
	payload, marshalErr := json.Marshal(deadLetter)
	if marshalErr != nil {
		c.recordFailure(ctx, record.Topic, record.Partition, record.Offset, "Failed to serialize Kafka dead-letter record", marshalErr)
		return false
	}

	startedAt := time.Now()
	result := c.client.ProduceSync(ctx, &franzkgo.Record{
		Topic: DeadLetterTopic(record.Topic),
		Key:   record.Key,
		Value: payload,
	})
	if publishErr := result.FirstErr(); publishErr != nil {
		RecordPublish(ctx, DeadLetterTopic(record.Topic), time.Since(startedAt), publishErr)
		c.recordFailure(ctx, record.Topic, record.Partition, record.Offset, "Failed to publish Kafka dead-letter record", publishErr)
		return false
	}
	RecordPublish(ctx, DeadLetterTopic(record.Topic), time.Since(startedAt), nil)
	RecordDeadLetter(ctx, record.Topic, c.consumerGroup)
	if logs.NotezyLogger != nil {
		deadLetter.Value = nil
		_ = logs.NotezyLogger.JSON(ctx, slog.LevelError, "Kafka event sent to dead-letter topic", deadLetter)
	}

	return true
}

func (c *Consumer) recordFailure(
	ctx context.Context,
	topic string,
	partition int32,
	offset int64,
	message string,
	err error,
) {
	if traces.NotezyTracer != nil {
		traces.NotezyTracer.RecordError(ctx, err)
	}
	if logs.NotezyLogger != nil {
		logs.NotezyLogger.Error(ctx, err, message,
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination.name", topic),
			attribute.String("messaging.consumer.group.name", c.consumerGroup),
			attribute.Int("messaging.kafka.partition", int(partition)),
			attribute.Int64("messaging.kafka.offset", offset),
		)
	}
}

/* ============================== Consumer Methods ============================== */

func (c *Consumer) Run(ctx context.Context, handler ConsumerHandler) error {
	if handler == nil {
		return errors.New("Kafka consumer handler is required")
	}

	for ctx.Err() == nil {
		fetches := c.client.PollRecords(ctx, c.config.MaximumPollRecords)
		if err := fetches.Err(); err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("poll Kafka consumer records: %w", err)
		}

		c.consumeFetches(ctx, fetches, handler)
		c.client.AllowRebalance()
	}

	return nil
}

func (c *Consumer) Close() {
	c.client.Close()
}

/* ============================== Auxiliary Functions ============================== */

func DeadLetterTopic(topic string) string {
	return topic + ".dlq"
}
