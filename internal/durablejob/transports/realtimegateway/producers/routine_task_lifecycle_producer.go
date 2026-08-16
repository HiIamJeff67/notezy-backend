package realtimegatewayproducers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	durablejobeventscontract "github.com/HiIamJeff67/notegic-backend/contracts/durable-job/v1/events"
	routinetasktypes "github.com/HiIamJeff67/notegic-backend/contracts/durable-job/v1/types/routine-tasks"
	eventcontract "github.com/HiIamJeff67/notegic-backend/contracts/types/events"

	platformkafka "github.com/HiIamJeff67/notegic-backend/shared/platform/kafka"
)

type RoutineTaskLifecycleProducer struct {
	producer *platformkafka.Producer
}

func NewRoutineTaskLifecycleProducer(
	producer *platformkafka.Producer,
) *RoutineTaskLifecycleProducer {
	return &RoutineTaskLifecycleProducer{
		producer: producer,
	}
}

func (p *RoutineTaskLifecycleProducer) ProduceRoutineTaskRunning(
	ctx context.Context,
	assignment routinetasktypes.RoutineTaskAssignment,
) error {
	now := time.Now().UTC()
	payload, err := json.Marshal(eventcontract.EventEnvelope[durablejobeventscontract.RoutineTaskRunningData]{
		SchemaVersion: eventcontract.Version,
		EventId:       uuid.New(),
		EventType:     durablejobeventscontract.EventType_RoutineTaskRunning,
		AggregateType: durablejobeventscontract.AggregateType_RoutineTask,
		AggregateId:   assignment.RoutineTaskId,
		KafkaKey:      assignment.RoutineTaskId.String(),
		OccurredAt:    now,
		CorrelationId: assignment.RoutineTaskRecordId.String(),
		Data: durablejobeventscontract.RoutineTaskRunningData{
			RoutineTaskId:       assignment.RoutineTaskId,
			RoutineTaskRecordId: assignment.RoutineTaskRecordId,
			RoutineId:           assignment.RoutineId,
			ActorUserPublicId:   assignment.ActorUserPublicId,
			Purpose:             assignment.Purpose,
			Attempt:             assignment.Attempt,
			StartedAt:           now,
		},
	})
	if err != nil {
		return err
	}

	return p.producer.Produce(
		ctx,
		durablejobeventscontract.DurableJobRealtimeGatewayRoutineTaskLifecycleTopic.String(),
		assignment.RoutineTaskId.String(),
		payload,
	)
}
