package durablejobproducers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	durablejobeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1/events"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/events"
	platformkafka "github.com/HiIamJeff67/notezy-backend/internal/platform/kafka"
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
	result durablejobeventscontract.YjsMaintenanceResultData,
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
		Data:          result,
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
