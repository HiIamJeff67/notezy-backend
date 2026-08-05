package yjsworkerproducers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	durablejobeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1/events"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/events"
	platformkafka "github.com/HiIamJeff67/notezy-backend/internal/platform/kafka"
)

type YjsMaintenanceCommandProducer struct {
	producer *platformkafka.Producer
}

func NewYjsMaintenanceCommandProducer(producer *platformkafka.Producer) *YjsMaintenanceCommandProducer {
	return &YjsMaintenanceCommandProducer{producer: producer}
}

func (p *YjsMaintenanceCommandProducer) Produce(
	ctx context.Context,
	source eventcontract.EventEnvelope[json.RawMessage],
	request durablejobeventscontract.YjsMaintenanceRequestData,
	documentId uuid.UUID,
	targetSequence int64,
) error {
	command := eventcontract.EventEnvelope[durablejobeventscontract.YjsMaintenanceCommandData]{
		SchemaVersion: eventcontract.Version,
		EventId:       uuid.New(),
		EventType:     durablejobeventscontract.EventType_YjsMaintenanceCommand,
		AggregateType: durablejobeventscontract.AggregateType_BlockPack,
		AggregateId:   request.BlockPackId,
		KafkaKey:      request.BlockPackId.String(),
		OccurredAt:    time.Now().UTC(),
		CorrelationId: request.CorrelationId,
		CausationId:   &source.EventId,
		Trace:         source.Trace,
		Data: durablejobeventscontract.YjsMaintenanceCommandData{
			RequestId:      request.RequestId,
			BlockPackId:    request.BlockPackId,
			DocumentId:     documentId,
			Operation:      request.Operation,
			TargetSequence: targetSequence,
			CorrelationId:  request.CorrelationId,
		},
	}
	payload, err := json.Marshal(command)
	if err != nil {
		return err
	}

	return p.producer.Produce(
		ctx,
		durablejobeventscontract.CoreYjsWorkerMaintenanceCommandTopic.String(),
		request.BlockPackId.String(),
		payload,
	)
}
