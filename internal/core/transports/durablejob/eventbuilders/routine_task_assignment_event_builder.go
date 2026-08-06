package eventbuilders

import (
	"time"

	"github.com/google/uuid"

	durablejobcontract "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1"
	durablejobeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1/events"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/types"
)

type RoutineTaskAssignmentEventBuilder struct{}

func NewRoutineTaskAssignmentEventBuilder() *RoutineTaskAssignmentEventBuilder {
	return &RoutineTaskAssignmentEventBuilder{}
}

func (b *RoutineTaskAssignmentEventBuilder) Build(
	response durablejobcontract.ClaimRoutineTasksResponseDto,
	occurredAt time.Time,
) eventcontract.EventEnvelope[durablejobcontract.ClaimRoutineTasksResponseDto] {
	return eventcontract.EventEnvelope[durablejobcontract.ClaimRoutineTasksResponseDto]{
		SchemaVersion: eventcontract.Version,
		EventId:       uuid.New(),
		EventType:     durablejobeventscontract.EventType_RoutineTasksAssigned,
		AggregateType: durablejobeventscontract.AggregateType_DurableJobWorker,
		AggregateId:   response.WorkerId,
		KafkaKey:      response.WorkerId.String(),
		OccurredAt:    occurredAt,
		CorrelationId: response.RequestId.String(),
		Data:          response,
	}
}
