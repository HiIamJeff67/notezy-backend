package yjsworkerconsumers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	durablejobeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1/events"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/events"
	platformkafka "github.com/HiIamJeff67/notezy-backend/internal/platform/kafka"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	durablejobproducers "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/durablejob/producers"
	"github.com/google/uuid"
)

type YjsMaintenanceResultConsumer struct {
	producer    *durablejobproducers.YjsMaintenanceResultProducer
	kafkaConfig platformkafka.ConsumerConfig
}

func NewYjsMaintenanceResultConsumer(
	producer *durablejobproducers.YjsMaintenanceResultProducer,
	kafkaConfig platformkafka.ConsumerConfig,
) *YjsMaintenanceResultConsumer {
	return &YjsMaintenanceResultConsumer{
		producer:    producer,
		kafkaConfig: kafkaConfig,
	}
}

func (c *YjsMaintenanceResultConsumer) Start(ctx context.Context) func() {
	consumer, err := platformkafka.NewConsumer(
		c.kafkaConfig,
		durablejobeventscontract.YjsWorkerCoreMaintenanceResultTopic.String(),
	)
	if err != nil {
		if logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(ctx, err, "failed to create Yjs maintenance result consumer")
		}

		return func() {}
	}

	workerCtx, cancel := context.WithCancel(ctx)
	go func() {
		if err := consumer.Run(workerCtx, c.consume); err != nil && workerCtx.Err() == nil && logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(workerCtx, err, "Yjs maintenance result consumer stopped")
		}
	}()

	return func() {
		cancel()
		consumer.Close()
	}
}

func (c *YjsMaintenanceResultConsumer) consume(
	ctx context.Context,
	_ platformkafka.ConsumerRecord,
	event eventcontract.EventEnvelope[json.RawMessage],
) error {
	if event.EventType != durablejobeventscontract.EventType_YjsMaintenanceCompleted ||
		event.AggregateType != durablejobeventscontract.AggregateType_BlockPack ||
		event.AggregateId == uuid.Nil ||
		event.KafkaKey != event.AggregateId.String() {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_SchemaIncompatible,
			Origin:         errors.New("invalid Yjs maintenance result envelope"),
		}
	}

	var result durablejobeventscontract.YjsMaintenanceResultData
	if err := json.Unmarshal(event.Data, &result); err != nil {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_SchemaIncompatible,
			Origin:         fmt.Errorf("decode Yjs maintenance result: %w", err),
		}
	}
	if result.RequestId == uuid.Nil || result.BlockPackId != event.AggregateId || result.DocumentId == uuid.Nil ||
		result.TargetSequence < 0 ||
		(result.Operation != durablejobeventscontract.YjsMaintenanceOperation_Compact && result.Operation != durablejobeventscontract.YjsMaintenanceOperation_Project) {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_SchemaIncompatible,
			Origin:         errors.New("invalid Yjs maintenance result data"),
		}
	}

	if err := c.producer.Produce(ctx, event, result); err != nil {
		return fmt.Errorf("produce forwarded Yjs maintenance result: %w", err)
	}

	return nil
}
