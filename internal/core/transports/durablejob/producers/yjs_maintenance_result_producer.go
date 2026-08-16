package durablejobproducers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	durablejobeventscontract "github.com/HiIamJeff67/notegic-backend/contracts/durable-job/v1/events"
	eventcontract "github.com/HiIamJeff67/notegic-backend/contracts/types/events"
	yjsworkereventscontract "github.com/HiIamJeff67/notegic-backend/contracts/yjs-worker/v1/events"

	platformkafka "github.com/HiIamJeff67/notegic-backend/shared/platform/kafka"
)

type YjsMaintenanceResultProducer struct {
	producer *platformkafka.Producer
}

func NewYjsMaintenanceResultProducer(producer *platformkafka.Producer) *YjsMaintenanceResultProducer {
	return &YjsMaintenanceResultProducer{producer: producer}
}

func (p *YjsMaintenanceResultProducer) Produce(
	ctx context.Context,
	source eventcontract.EventEnvelope[json.RawMessage],
	result yjsworkereventscontract.YjsMaintenanceResultData,
) error {
	forwarded := eventcontract.EventEnvelope[durablejobeventscontract.YjsMaintenanceResultData]{
		SchemaVersion: eventcontract.Version,
		EventId:       uuid.New(),
		EventType:     durablejobeventscontract.EventType_YjsMaintenanceCompleted,
		AggregateType: durablejobeventscontract.AggregateType_BlockPack,
		AggregateId:   result.BlockPackId,
		KafkaKey:      result.BlockPackId.String(),
		OccurredAt:    time.Now().UTC(),
		CorrelationId: source.CorrelationId,
		CausationId:   &source.EventId,
		Trace:         source.Trace,
		Data: durablejobeventscontract.YjsMaintenanceResultData{
			RequestId:              result.RequestId,
			BlockPackId:            result.BlockPackId,
			DocumentId:             result.DocumentId,
			Operation:              result.Operation,
			TargetSequence:         result.TargetSequence,
			Success:                result.Success,
			CompactedUntilSequence: result.CompactedUntilSequence,
			ProjectedUntilSequence: result.ProjectedUntilSequence,
			Error:                  result.Error,
		},
	}
	payload, err := json.Marshal(forwarded)
	if err != nil {
		return err
	}

	return p.producer.Produce(
		ctx,
		durablejobeventscontract.DurableJobCoreYjsMaintenanceResultTopic.String(),
		result.BlockPackId.String(),
		payload,
	)
}
