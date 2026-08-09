package coreconsumers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	durablejobcontract "github.com/HiIamJeff67/notezy-backend/contracts/durable-job/v1"
	durablejobeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/durable-job/v1/events"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/events"

	platformkafka "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka"
	logs "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/logs"

	routinetask "github.com/HiIamJeff67/notezy-backend/internal/durablejob/routinetask"
)

type RoutineTaskAssignmentConsumer struct {
	routineTaskEngine *routinetask.Engine
	kafkaConfig       platformkafka.ConsumerConfig
}

func NewRoutineTaskAssignmentConsumer(
	routineTaskEngine *routinetask.Engine,
	kafkaConfig platformkafka.ConsumerConfig,
) *RoutineTaskAssignmentConsumer {
	return &RoutineTaskAssignmentConsumer{
		routineTaskEngine: routineTaskEngine,
		kafkaConfig:       kafkaConfig,
	}
}

func (c *RoutineTaskAssignmentConsumer) Start(ctx context.Context) func() {
	consumer, err := platformkafka.NewConsumer(
		c.kafkaConfig,
		durablejobeventscontract.CoreDurableJobRoutineTaskTopic.String(),
	)
	if err != nil {
		if logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(ctx, err, "Failed to create routine task assignment consumer")
		}

		return func() {}
	}

	workerCtx, cancel := context.WithCancel(ctx)
	go func() {
		if err := consumer.Run(workerCtx, c.consume); err != nil && workerCtx.Err() == nil && logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(workerCtx, err, "Routine task assignment consumer stopped")
		}
	}()

	return func() {
		cancel()
		consumer.Close()
	}
}

func (c *RoutineTaskAssignmentConsumer) consume(
	ctx context.Context,
	_ platformkafka.ConsumerRecord,
	event eventcontract.EventEnvelope[json.RawMessage],
) error {
	if event.EventType != durablejobeventscontract.EventType_RoutineTasksAssigned {
		return nil
	}
	if event.AggregateType != durablejobeventscontract.AggregateType_DurableJobWorker {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_SchemaIncompatible,
			Origin:         fmt.Errorf("DurableJob routine task assignment event is invalid"),
		}
	}

	var response durablejobcontract.ClaimRoutineTasksResponseDto
	if err := json.Unmarshal(event.Data, &response); err != nil {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_SchemaIncompatible,
			Origin:         fmt.Errorf("decode DurableJob routine task assignments: %w", err),
		}
	}
	if response.RequestId == uuid.Nil || response.WorkerId != event.AggregateId {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_SchemaIncompatible,
			Origin:         fmt.Errorf("DurableJob routine task assignment response is invalid"),
		}
	}
	for index, assignment := range response.Assignments {
		if assignment.RoutineTaskId == uuid.Nil || assignment.RoutineTaskRecordId == uuid.Nil ||
			assignment.ActorUserId == uuid.Nil || assignment.RoutineId == uuid.Nil ||
			assignment.Purpose == "" || len(assignment.Payload) == 0 || assignment.StartedAt.IsZero() {
			return &platformkafka.ConsumerError{
				Classification: platformkafka.ErrorClassification_SchemaIncompatible,
				Origin:         fmt.Errorf("DurableJob routine task assignment at index %d is invalid", index),
			}
		}
	}

	if err := c.routineTaskEngine.HandleRoutineTaskAssignments(ctx, response.Assignments); err != nil {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_Transient,
			Origin:         err,
		}
	}

	return nil
}
