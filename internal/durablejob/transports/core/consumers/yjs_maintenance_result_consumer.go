package coreconsumers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	coreeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/events"
	durablejobeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1/events"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/types"
	yjsworkereventscontract "github.com/HiIamJeff67/notezy-backend/contracts/yjsworker/v1/events"

	platformkafka "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka"
	logs "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/logs"
	metrics "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/metrics"

	corestrategies "github.com/HiIamJeff67/notezy-backend/internal/durablejob/transports/core/strategies"
)

type YjsMaintenanceResultConsumer struct {
	kafkaConfig platformkafka.ConsumerConfig
	strategy    *corestrategies.YjsMaintenanceStrategy
}

func NewYjsMaintenanceResultConsumer(
	strategy *corestrategies.YjsMaintenanceStrategy,
	kafkaConfig platformkafka.ConsumerConfig,
) *YjsMaintenanceResultConsumer {
	if strategy == nil {
		strategy = corestrategies.NewYjsMaintenanceStrategy()
	}
	return &YjsMaintenanceResultConsumer{
		kafkaConfig: kafkaConfig,
		strategy:    strategy,
	}
}

func (c *YjsMaintenanceResultConsumer) Start(ctx context.Context) func() {
	consumer, err := platformkafka.NewConsumer(
		c.kafkaConfig,
		durablejobeventscontract.DurableJobCoreYjsMaintenanceResultTopic.String(),
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
		event.AggregateId == uuid.Nil || event.KafkaKey != event.AggregateId.String() {
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
		(result.Operation != yjsworkereventscontract.YjsMaintenanceOperation_Compact && result.Operation != yjsworkereventscontract.YjsMaintenanceOperation_Project) {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_SchemaIncompatible,
			Origin:         errors.New("invalid Yjs maintenance result data"),
		}
	}

	if result.Success {
		request, exists := c.strategy.Complete(result.RequestId)
		if exists && result.Operation == yjsworkereventscontract.YjsMaintenanceOperation_Compact &&
			result.CompactedUntilSequence > request.Hint.CompactedUntilSequence {
			request.Hint.CompactedUntilSequence = result.CompactedUntilSequence
			request.Hint.UncompactedUpdateCount = request.Hint.LatestUpdateSequence - result.CompactedUntilSequence
			compactedAt := time.Now().UTC()
			request.Hint.LastCompactedAt = &compactedAt
			if result.ProjectedUntilSequence > request.Hint.ProjectedUntilSequence {
				request.Hint.ProjectedUntilSequence = result.ProjectedUntilSequence
			}
			if request.Hint.ProjectedUntilSequence < request.Hint.LatestUpdateSequence {
				c.retryHint(ctx, request.Hint)
			}
		}
		if metrics.NotezyMeter != nil {
			metrics.NotezyMeter.Count(ctx, "yjs.maintenance.result.success", 1)
		}
		return nil
	}

	if request, exists := c.strategy.Fail(result.RequestId); exists && request.Attempt < corestrategies.MaximumRequestAttempts {
		c.retryHint(ctx, request.Hint)
	}
	if metrics.NotezyMeter != nil {
		metrics.NotezyMeter.Count(ctx, "yjs.maintenance.result.failure", 1)
	}
	if logs.NotezyLogger != nil {
		logs.NotezyLogger.Error(ctx, errors.New(result.Error), "Yjs maintenance request failed")
	}

	return nil
}

func (c *YjsMaintenanceResultConsumer) retryHint(ctx context.Context, hint coreeventscontract.YjsMaintenanceHintData) {
	go func() {
		timer := time.NewTimer(time.Second)
		defer timer.Stop()
		select {
		case <-ctx.Done():
		case <-timer.C:
			if err := c.strategy.Enqueue(hint); err != nil && logs.NotezyLogger != nil {
				logs.NotezyLogger.Error(ctx, err, "failed to requeue Yjs maintenance hint")
			}
		}
	}()
}
