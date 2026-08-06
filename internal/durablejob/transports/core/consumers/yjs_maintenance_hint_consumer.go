package coreconsumers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"sync"
	"time"

	"github.com/google/uuid"

	coreeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/events"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/types"
	yjsworkercontract "github.com/HiIamJeff67/notezy-backend/contracts/yjs-worker/v1"
	yjsworkereventscontract "github.com/HiIamJeff67/notezy-backend/contracts/yjs-worker/v1/events"

	platformkafka "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka"
	logs "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/logs"
	metrics "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/metrics"

	coreproducers "github.com/HiIamJeff67/notezy-backend/internal/durablejob/transports/core/producers"
	corestrategies "github.com/HiIamJeff67/notezy-backend/internal/durablejob/transports/core/strategies"
)

type YjsMaintenanceHintConsumer struct {
	producer    *coreproducers.YjsMaintenanceRequestProducer
	kafkaConfig platformkafka.ConsumerConfig
	strategy    *corestrategies.YjsMaintenanceStrategy
	slots       chan struct{}
}

func NewYjsMaintenanceHintConsumer(
	producer *coreproducers.YjsMaintenanceRequestProducer,
	strategy *corestrategies.YjsMaintenanceStrategy,
	kafkaConfig platformkafka.ConsumerConfig,
) *YjsMaintenanceHintConsumer {
	if strategy == nil {
		strategy = corestrategies.NewYjsMaintenanceStrategy()
	}
	return &YjsMaintenanceHintConsumer{
		producer:    producer,
		kafkaConfig: kafkaConfig,
		strategy:    strategy,
		slots:       make(chan struct{}, corestrategies.MaximumDispatchWorkers),
	}
}

func (c *YjsMaintenanceHintConsumer) Start(ctx context.Context) func() {
	consumer, err := platformkafka.NewConsumer(
		c.kafkaConfig,
		coreeventscontract.CoreDurableJobYjsMaintenanceHintTopic.String(),
	)
	if err != nil {
		if logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(ctx, err, "failed to create Yjs maintenance hint consumer")
		}
		return func() {}
	}

	workerCtx, cancel := context.WithCancel(ctx)
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	go func() {
		defer waitGroup.Done()
		if err := consumer.Run(workerCtx, c.consume); err != nil && workerCtx.Err() == nil && logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(workerCtx, err, "Yjs maintenance hint consumer stopped")
		}
	}()
	go func() {
		defer waitGroup.Done()
		c.dispatch(workerCtx)
	}()

	return func() {
		cancel()
		consumer.Close()
		waitGroup.Wait()
	}
}

func (c *YjsMaintenanceHintConsumer) consume(
	ctx context.Context,
	_ platformkafka.ConsumerRecord,
	event eventcontract.EventEnvelope[json.RawMessage],
) error {
	if event.EventType != coreeventscontract.EventType_YjsMaintenanceHint ||
		event.AggregateType != coreeventscontract.AggregateType_BlockPack ||
		event.AggregateId == uuid.Nil || event.KafkaKey != event.AggregateId.String() {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_SchemaIncompatible,
			Origin:         errors.New("invalid Yjs maintenance hint envelope"),
		}
	}

	var hint coreeventscontract.YjsMaintenanceHintData
	if err := json.Unmarshal(event.Data, &hint); err != nil {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_SchemaIncompatible,
			Origin:         fmt.Errorf("decode Yjs maintenance hint: %w", err),
		}
	}
	if hint.BlockPackId != event.AggregateId || hint.DocumentId == uuid.Nil || hint.LatestUpdateSequence < 0 ||
		hint.CompactedUntilSequence < 0 || hint.ProjectedUntilSequence < -1 {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_SchemaIncompatible,
			Origin:         errors.New("invalid Yjs maintenance hint data"),
		}
	}

	if err := c.strategy.Enqueue(hint); err != nil {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_Transient,
			Origin:         err,
		}
	}
	if metrics.NotezyMeter != nil {
		metrics.NotezyMeter.Value(ctx, "yjs.maintenance.queue.size", int64(c.strategy.PendingCount()))
	}

	return nil
}

func (c *YjsMaintenanceHintConsumer) dispatch(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.strategy.Notify():
			c.dispatchPending(ctx)
		}
	}
}

func (c *YjsMaintenanceHintConsumer) dispatchPending(ctx context.Context) {
	for {
		hints := c.strategy.DequeueBatch(corestrategies.MaximumDispatchBatch)
		if len(hints) == 0 {
			return
		}
		var waitGroup sync.WaitGroup
		for _, hint := range hints {
			c.slots <- struct{}{}
			waitGroup.Add(1)
			go func(hint coreeventscontract.YjsMaintenanceHintData) {
				defer waitGroup.Done()
				defer func() { <-c.slots }()
				if err := c.dispatchHint(ctx, hint); err != nil {
					if logs.NotezyLogger != nil {
						logs.NotezyLogger.Error(ctx, err, "failed to dispatch Yjs maintenance request")
					}
					c.retryHint(ctx, hint)
					if metrics.NotezyMeter != nil {
						metrics.NotezyMeter.Count(ctx, "yjs.maintenance.request.failure", 1)
					}
				}
			}(hint)
		}
		waitGroup.Wait()
	}
}

func (c *YjsMaintenanceHintConsumer) retryHint(ctx context.Context, hint coreeventscontract.YjsMaintenanceHintData) {
	go func() {
		timer := time.NewTimer(time.Second)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			if err := c.strategy.Enqueue(hint); err != nil && logs.NotezyLogger != nil {
				logs.NotezyLogger.Error(ctx, err, "failed to requeue Yjs maintenance hint")
			}
		}
	}()
}

func (c *YjsMaintenanceHintConsumer) dispatchHint(ctx context.Context, hint coreeventscontract.YjsMaintenanceHintData) error {
	startedAt := time.Now()
	operation := yjsworkereventscontract.YjsMaintenanceOperation_Project
	if hint.UncompactedUpdateCount >= yjsworkercontract.YjsCompactionUpdateThreshold ||
		(hint.CompactedUntilSequence < hint.LatestUpdateSequence && hint.LastCompactedAt == nil) {
		operation = yjsworkereventscontract.YjsMaintenanceOperation_Compact
	}
	if operation == yjsworkereventscontract.YjsMaintenanceOperation_Project && hint.ProjectedUntilSequence >= hint.LatestUpdateSequence {
		if metrics.NotezyMeter != nil {
			metrics.NotezyMeter.Duration(ctx, "yjs.maintenance.dispatch.duration", time.Since(startedAt),
				attribute.String("operation", string(operation)), attribute.String("outcome", "noop"))
		}
		return nil
	}

	requestId := uuid.New()
	c.strategy.Track(requestId, hint)
	if err := c.producer.Produce(ctx, hint, operation, hint.LatestUpdateSequence, requestId); err != nil {
		c.strategy.Complete(requestId)
		return err
	}
	if metrics.NotezyMeter != nil {
		metrics.NotezyMeter.Duration(ctx, "yjs.maintenance.dispatch.duration", time.Since(startedAt),
			attribute.String("operation", string(operation)), attribute.String("outcome", "success"))
		metrics.NotezyMeter.Value(ctx, "yjs.maintenance.queue.size", int64(c.strategy.PendingCount()))
	}

	return nil
}
