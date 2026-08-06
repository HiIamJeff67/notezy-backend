package coreproducers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	durablejobeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1/events"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/types"

	platformkafka "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka"

	durablejobroutinetask "github.com/HiIamJeff67/notezy-backend/internal/durablejob/routinetask"
)

type RoutineTaskResultProducer struct {
	producer *platformkafka.Producer
}

func NewRoutineTaskResultProducer(producer *platformkafka.Producer) *RoutineTaskResultProducer {
	return &RoutineTaskResultProducer{producer: producer}
}

func (p *RoutineTaskResultProducer) Produce(
	ctx context.Context,
	result durablejobroutinetask.RoutineTaskResult,
) error {
	eventType, ok := map[durablejobroutinetask.RoutineTaskResultKind]eventcontract.EventType{
		durablejobroutinetask.RoutineTaskResultKind_Completed: durablejobeventscontract.EventType_RoutineTasksCompleted,
		durablejobroutinetask.RoutineTaskResultKind_Failed:    durablejobeventscontract.EventType_RoutineTasksFailed,
	}[result.Kind]
	if !ok {
		return fmt.Errorf("unsupported routine task result kind: %s", result.Kind)
	}

	payload, err := json.Marshal(eventcontract.EventEnvelope[any]{
		SchemaVersion: eventcontract.Version,
		EventId:       uuid.New(),
		EventType:     eventType,
		AggregateType: durablejobeventscontract.AggregateType_DurableJobWorker,
		AggregateId:   result.WorkerId,
		KafkaKey:      result.WorkerId.String(),
		OccurredAt:    time.Now().UTC(),
		CorrelationId: result.CorrelationId,
		Data:          result.Data,
	})
	if err != nil {
		return err
	}

	return p.producer.Produce(
		ctx,
		durablejobeventscontract.CoreDurableJobRoutineTaskTopic.String(),
		result.WorkerId.String(),
		payload,
	)
}
