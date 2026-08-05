package coreproducers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	durablejobeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1/events"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/events"
	platformkafka "github.com/HiIamJeff67/notezy-backend/internal/platform/kafka"
)

type YjsMaintenanceRequestProducer struct {
	producer *platformkafka.Producer
}

func NewYjsMaintenanceRequestProducer(producer *platformkafka.Producer) *YjsMaintenanceRequestProducer {
	return &YjsMaintenanceRequestProducer{producer: producer}
}

func (p *YjsMaintenanceRequestProducer) Produce(
	ctx context.Context,
	hint durablejobeventscontract.YjsMaintenanceHintData,
	operation durablejobeventscontract.YjsMaintenanceOperation,
	targetSequence int64,
	requestId uuid.UUID,
) error {
	correlationId := uuid.NewString()
	request := eventcontract.EventEnvelope[durablejobeventscontract.YjsMaintenanceRequestData]{
		SchemaVersion: eventcontract.Version,
		EventId:       uuid.New(),
		EventType:     durablejobeventscontract.EventType_YjsMaintenanceRequested,
		AggregateType: durablejobeventscontract.AggregateType_BlockPack,
		AggregateId:   hint.BlockPackId,
		KafkaKey:      hint.BlockPackId.String(),
		OccurredAt:    time.Now().UTC(),
		CorrelationId: correlationId,
		Data: durablejobeventscontract.YjsMaintenanceRequestData{
			RequestId:      requestId,
			BlockPackId:    hint.BlockPackId,
			DocumentId:     hint.DocumentId,
			Operation:      operation,
			TargetSequence: targetSequence,
			CorrelationId:  correlationId,
		},
	}
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}
	if err := p.producer.Produce(
		ctx,
		durablejobeventscontract.DurableJobCoreYjsMaintenanceRequestTopic.String(),
		hint.BlockPackId.String(),
		payload,
	); err != nil {
		return err
	}

	return nil
}
