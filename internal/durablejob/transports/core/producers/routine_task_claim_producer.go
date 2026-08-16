package coreproducers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	durablejobcontract "github.com/HiIamJeff67/notegic-backend/contracts/durable-job/v1"
	durablejobeventscontract "github.com/HiIamJeff67/notegic-backend/contracts/durable-job/v1/events"
	eventcontract "github.com/HiIamJeff67/notegic-backend/contracts/types/events"

	platformkafka "github.com/HiIamJeff67/notegic-backend/shared/platform/kafka"
)

type RoutineTaskClaimProducer struct {
	producer *platformkafka.Producer
}

func NewRoutineTaskClaimProducer(producer *platformkafka.Producer) *RoutineTaskClaimProducer {
	return &RoutineTaskClaimProducer{producer: producer}
}

func (p *RoutineTaskClaimProducer) Produce(
	ctx context.Context,
	request durablejobcontract.ClaimRoutineTasksRequestDto,
) error {
	payload, err := json.Marshal(eventcontract.EventEnvelope[durablejobcontract.ClaimRoutineTasksRequestDto]{
		SchemaVersion: eventcontract.Version,
		EventId:       uuid.New(),
		EventType:     durablejobeventscontract.EventType_RoutineTaskClaimRequested,
		AggregateType: durablejobeventscontract.AggregateType_DurableJobWorker,
		AggregateId:   request.WorkerId,
		KafkaKey:      request.WorkerId.String(),
		OccurredAt:    time.Now().UTC(),
		CorrelationId: request.RequestId.String(),
		Data:          request,
	})
	if err != nil {
		return err
	}

	return p.producer.Produce(
		ctx,
		durablejobeventscontract.CoreDurableJobRoutineTaskTopic.String(),
		request.WorkerId.String(),
		payload,
	)
}
