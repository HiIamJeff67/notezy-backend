package adaptersproducers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	durablejobeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/durable-job/v1/events"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/events"
	yjsworkereventscontract "github.com/HiIamJeff67/notezy-backend/contracts/yjs-worker/v1/events"

	platformkafka "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka"
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
	command := eventcontract.EventEnvelope[yjsworkereventscontract.YjsMaintenanceCommandData]{
		SchemaVersion: eventcontract.Version,
		EventId:       uuid.New(),
		EventType:     yjsworkereventscontract.EventType_YjsMaintenanceCommand,
		AggregateType: yjsworkereventscontract.AggregateType_BlockPack,
		AggregateId:   request.BlockPackId,
		KafkaKey:      request.BlockPackId.String(),
		OccurredAt:    time.Now().UTC(),
		CorrelationId: request.CorrelationId,
		CausationId:   &source.EventId,
		Trace:         source.Trace,
		Data: yjsworkereventscontract.YjsMaintenanceCommandData{
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
		yjsworkereventscontract.YjsWorkerCoreMaintenanceCommandTopic.String(),
		request.BlockPackId.String(),
		payload,
	)
}
