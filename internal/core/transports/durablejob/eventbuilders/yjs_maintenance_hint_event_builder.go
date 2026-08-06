package eventbuilders

import (
	"time"

	"github.com/google/uuid"

	coreeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/events"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/types"
)

type YjsMaintenanceHintEventBuilder struct{}

func NewYjsMaintenanceHintEventBuilder() *YjsMaintenanceHintEventBuilder {
	return &YjsMaintenanceHintEventBuilder{}
}

func (b *YjsMaintenanceHintEventBuilder) Build(
	hint coreeventscontract.YjsMaintenanceHintData,
	correlationId string,
	occurredAt time.Time,
) eventcontract.EventEnvelope[coreeventscontract.YjsMaintenanceHintData] {
	return eventcontract.EventEnvelope[coreeventscontract.YjsMaintenanceHintData]{
		SchemaVersion: eventcontract.Version,
		EventId:       uuid.New(),
		EventType:     coreeventscontract.EventType_YjsMaintenanceHint,
		AggregateType: coreeventscontract.AggregateType_BlockPack,
		AggregateId:   hint.BlockPackId,
		KafkaKey:      hint.BlockPackId.String(),
		OccurredAt:    occurredAt,
		CorrelationId: correlationId,
		Data:          hint,
	}
}
