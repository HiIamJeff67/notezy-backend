package yjsworkereventscontract

import (
	"github.com/google/uuid"

	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/types"
)

const (
	YjsWorkerCoreCommandTopic            eventcontract.Topic = "notezy.yjsworker.core.command.v1"
	CoreYjsWorkerReplyTopic              eventcontract.Topic = "notezy.core.yjsworker.reply.v1"
	YjsWorkerCoreMaintenanceCommandTopic eventcontract.Topic = "notezy.core.yjsworker.maintenance-command.v1"
	CoreYjsWorkerMaintenanceResultTopic  eventcontract.Topic = "notezy.yjsworker.core.maintenance-result.v1"
)

const (
	EventType_YjsWorkerCommand          eventcontract.EventType = "YjsWorkerCommand"
	EventType_YjsWorkerCommandCompleted eventcontract.EventType = "YjsWorkerCommandCompleted"
	EventType_YjsMaintenanceCommand     eventcontract.EventType = "YjsMaintenanceCommand"
	EventType_YjsMaintenanceCompleted   eventcontract.EventType = "YjsMaintenanceCompleted"
)

const AggregateType_BlockPack eventcontract.AggregateType = "BlockPack"

type YjsMaintenanceOperation string

const (
	YjsMaintenanceOperation_Compact YjsMaintenanceOperation = "compact"
	YjsMaintenanceOperation_Project YjsMaintenanceOperation = "project"
)

type YjsMaintenanceCommandData struct {
	RequestId      uuid.UUID               `json:"requestId"`
	BlockPackId    uuid.UUID               `json:"blockPackId"`
	DocumentId     uuid.UUID               `json:"documentId"`
	Operation      YjsMaintenanceOperation `json:"operation"`
	TargetSequence int64                   `json:"targetSequence"`
	CorrelationId  string                  `json:"correlationId"`
}

type YjsMaintenanceResultData struct {
	RequestId              uuid.UUID               `json:"requestId"`
	BlockPackId            uuid.UUID               `json:"blockPackId"`
	DocumentId             uuid.UUID               `json:"documentId"`
	Operation              YjsMaintenanceOperation `json:"operation"`
	TargetSequence         int64                   `json:"targetSequence"`
	Success                bool                    `json:"success"`
	CompactedUntilSequence int64                   `json:"compactedUntilSequence"`
	ProjectedUntilSequence int64                   `json:"projectedUntilSequence"`
	Error                  string                  `json:"error,omitempty"`
}
