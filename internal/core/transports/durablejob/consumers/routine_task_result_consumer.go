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

type DurableJobRoutineTaskResultConsumer struct {
	routineTaskService          routineservices.RoutineTaskServiceInterface
	routineTaskExecutionService routineservices.RoutineTaskExecutionServiceInterface
	kafkaConfig                 platformkafka.ConsumerConfig
}

func NewDurableJobRoutineTaskResultConsumer(
	routineTaskService routineservices.RoutineTaskServiceInterface,
	kafkaConfig platformkafka.ConsumerConfig,
	routineTaskExecutionService routineservices.RoutineTaskExecutionServiceInterface,
) *DurableJobRoutineTaskResultConsumer {
	return &DurableJobRoutineTaskResultConsumer{
		routineTaskService:          routineTaskService,
		routineTaskExecutionService: routineTaskExecutionService,
		kafkaConfig:                 kafkaConfig,
	}
}

func (c *DurableJobRoutineTaskResultConsumer) Start(ctx context.Context) func() {
	consumer, err := platformkafka.NewConsumer(
		c.kafkaConfig,
		durablejobeventscontract.CoreDurableJobRoutineTaskTopic.String(),
	)
	if err != nil {
		if logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(ctx, err, "Failed to create DurableJob routine task result consumer")
		}
		return func() {}
	}

	workerCtx, cancel := context.WithCancel(ctx)
	go func() {
		if err := consumer.Run(workerCtx, c.consume); err != nil && workerCtx.Err() == nil && logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(workerCtx, err, "DurableJob routine task result consumer stopped")
		}
	}()

	return func() {
		cancel()
		consumer.Close()
	}
}

func (c *DurableJobRoutineTaskResultConsumer) consume(
	ctx context.Context,
	_ platformkafka.ConsumerRecord,
	event eventcontract.EventEnvelope[json.RawMessage],
) error {
	if event.EventType != durablejobeventscontract.EventType_RoutineTasksCompleted && event.EventType != durablejobeventscontract.EventType_RoutineTasksFailed {
		return nil
	}

	if event.EventId == uuid.Nil || event.AggregateType != durablejobeventscontract.AggregateType_DurableJobWorker || event.AggregateId == uuid.Nil || event.KafkaKey != event.AggregateId.String() {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_SchemaIncompatible,
			Origin:         fmt.Errorf("DurableJob routine task result aggregate is invalid"),
		}
	}

	switch event.EventType {
	case durablejobeventscontract.EventType_RoutineTasksCompleted:
		var request durablejobcontract.MarkCompletedRoutineTasksRequestDto
		if err := json.Unmarshal(event.Data, &request); err != nil {
			return &platformkafka.ConsumerError{
				Classification: platformkafka.ErrorClassification_SchemaIncompatible,
				Origin:         fmt.Errorf("decode DurableJob routine task completed result: %w", err),
			}
		}
		if request.WorkerId != event.AggregateId {
			return &platformkafka.ConsumerError{Classification: platformkafka.ErrorClassification_SchemaIncompatible, Origin: fmt.Errorf("DurableJob routine task completed worker does not match aggregate")}
		}
		if exception := c.routineTaskExecutionService.ApplyPreparedRoutineTasks(ctx, event.EventId, &request); exception != nil {
			return &platformkafka.ConsumerError{Classification: platformkafka.ErrorClassification_Transient, Origin: exception}
		}
	case durablejobeventscontract.EventType_RoutineTasksFailed:
		var request durablejobcontract.MarkFailedRoutineTasksRequestDto
		if err := json.Unmarshal(event.Data, &request); err != nil {
			return &platformkafka.ConsumerError{
				Classification: platformkafka.ErrorClassification_SchemaIncompatible,
				Origin:         fmt.Errorf("decode DurableJob routine task failed result: %w", err),
			}
		}
		if request.WorkerId != event.AggregateId {
			return &platformkafka.ConsumerError{Classification: platformkafka.ErrorClassification_SchemaIncompatible, Origin: fmt.Errorf("DurableJob routine task failed worker does not match aggregate")}
		}
		if exception := c.routineTaskService.MarkFailedRoutineTasks(ctx, event.EventId, &request); exception != nil {
			return &platformkafka.ConsumerError{Classification: platformkafka.ErrorClassification_Transient, Origin: exception}
		}
	}

	return nil
}
