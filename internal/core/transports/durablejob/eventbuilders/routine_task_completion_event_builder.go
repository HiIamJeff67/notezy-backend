package eventbuilders

import (
	"time"

	"github.com/google/uuid"

	coreeventscontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/events"
	routinetasktypes "github.com/HiIamJeff67/notegic-backend/contracts/durable-job/v1/types/routine-tasks"
	eventcontract "github.com/HiIamJeff67/notegic-backend/contracts/types/events"
)

type RoutineTaskCompletionEventBuilder struct{}

func NewRoutineTaskCompletionEventBuilder() *RoutineTaskCompletionEventBuilder {
	return &RoutineTaskCompletionEventBuilder{}
}

func (b *RoutineTaskCompletionEventBuilder) Build(
	completedTask routinetasktypes.CompletedRoutineTask,
	workerId uuid.UUID,
	occurredAt time.Time,
) eventcontract.EventEnvelope[coreeventscontract.RoutineTaskCompletedData] {
	return eventcontract.EventEnvelope[coreeventscontract.RoutineTaskCompletedData]{
		SchemaVersion: eventcontract.Version,
		EventId:       uuid.New(),
		EventType:     coreeventscontract.EventType_RoutineTaskCompleted,
		AggregateType: coreeventscontract.AggregateType_RoutineTask,
		AggregateId:   completedTask.RoutineTaskId,
		KafkaKey:      completedTask.RoutineTaskId.String(),
		OccurredAt:    occurredAt,
		CorrelationId: workerId.String(),
		Data: coreeventscontract.RoutineTaskCompletedData{
			RoutineTaskId:       completedTask.RoutineTaskId,
			RoutineTaskRecordId: completedTask.RoutineTaskRecordId,
			RoutineId:           completedTask.PreparedTask.RoutineId,
			ActorUserPublicId:   completedTask.PreparedTask.ActorUserPublicId,
			Purpose:             completedTask.PreparedTask.Purpose,
			WorkerId:            workerId,
			Attempt:             completedTask.PreparedTask.Attempt,
			CompletedAt:         completedTask.CompletedAt,
		},
	}
}
