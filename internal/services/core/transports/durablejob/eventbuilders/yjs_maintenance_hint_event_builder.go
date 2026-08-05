package eventbuilders

import (
	"time"

	durablejobeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1/events"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/events"
	"github.com/google/uuid"
)

type YjsMaintenanceHintEventBuilder struct{}

func NewYjsMaintenanceHintEventBuilder() *YjsMaintenanceHintEventBuilder {
	return &YjsMaintenanceHintEventBuilder{}
}

func (b *YjsMaintenanceHintEventBuilder) Build(
	hint durablejobeventscontract.YjsMaintenanceHintData,
	correlationId string,
	occurredAt time.Time,
) eventcontract.EventEnvelope[durablejobeventscontract.YjsMaintenanceHintData] {
	return eventcontract.EventEnvelope[durablejobeventscontract.YjsMaintenanceHintData]{
		SchemaVersion: eventcontract.Version,
		EventId:       uuid.New(),
		EventType:     durablejobeventscontract.EventType_YjsMaintenanceHint,
		AggregateType: durablejobeventscontract.AggregateType_BlockPack,
		AggregateId:   hint.BlockPackId,
		KafkaKey:      hint.BlockPackId.String(),
		OccurredAt:    occurredAt,
		CorrelationId: correlationId,
		Data:          hint,
	}
}
