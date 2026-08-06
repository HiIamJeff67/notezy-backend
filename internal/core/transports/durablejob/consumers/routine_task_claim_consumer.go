package durablejobconsumers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	durablejobcontract "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1"
	durablejobeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1/events"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/types"

	platformkafka "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka"
	logs "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/logs"

	routineservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/routines"
)

type DurableJobRoutineTaskClaimConsumer struct {
	routineTaskService routineservices.RoutineTaskServiceInterface
	kafkaConfig        platformkafka.ConsumerConfig
}

func NewDurableJobRoutineTaskClaimConsumer(
	routineTaskService routineservices.RoutineTaskServiceInterface,
	kafkaConfig platformkafka.ConsumerConfig,
) *DurableJobRoutineTaskClaimConsumer {
	return &DurableJobRoutineTaskClaimConsumer{
		routineTaskService: routineTaskService,
		kafkaConfig:        kafkaConfig,
	}
}

func (c *DurableJobRoutineTaskClaimConsumer) Start(ctx context.Context) func() {
	consumer, err := platformkafka.NewConsumer(
		c.kafkaConfig,
		durablejobeventscontract.CoreDurableJobRoutineTaskTopic.String(),
	)
	if err != nil {
		if logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(ctx, err, "Failed to create DurableJob routine task claim consumer")
		}

		return func() {}
	}

	workerCtx, cancel := context.WithCancel(ctx)
	go func() {
		if err := consumer.Run(workerCtx, c.consume); err != nil && workerCtx.Err() == nil && logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(workerCtx, err, "DurableJob routine task claim consumer stopped")
		}
	}()

	return func() {
		cancel()
		consumer.Close()
	}
}

func (c *DurableJobRoutineTaskClaimConsumer) consume(
	ctx context.Context,
	_ platformkafka.ConsumerRecord,
	event eventcontract.EventEnvelope[json.RawMessage],
) error {
	if event.EventType != durablejobeventscontract.EventType_RoutineTaskClaimRequested {
		return nil
	}

	var request durablejobcontract.ClaimRoutineTasksRequestDto
	if err := json.Unmarshal(event.Data, &request); err != nil {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_SchemaIncompatible,
			Origin:         fmt.Errorf("decode DurableJob routine task claim request: %w", err),
		}
	}
	if event.AggregateType != durablejobeventscontract.AggregateType_DurableJobWorker ||
		request.RequestId == uuid.Nil ||
		request.WorkerId == uuid.Nil ||
		request.WorkerId != event.AggregateId {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_SchemaIncompatible,
			Origin:         fmt.Errorf("DurableJob routine task claim request is invalid"),
		}
	}

	if _, exception := c.routineTaskService.ClaimRoutineTasks(ctx, event.EventId, &request); exception != nil {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_Transient,
			Origin:         exception,
		}
	}

	return nil
}
