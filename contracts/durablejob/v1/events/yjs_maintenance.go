package durablejobeventscontract

import (
	"time"

	"github.com/google/uuid"
)

type YjsMaintenanceOperation string

const (
	YjsMaintenanceOperation_Compact YjsMaintenanceOperation = "compact"
	YjsMaintenanceOperation_Project YjsMaintenanceOperation = "project"
)

type YjsMaintenanceHintData struct {
	BlockPackId            uuid.UUID  `json:"blockPackId"`
	DocumentId             uuid.UUID  `json:"documentId"`
	LatestUpdateSequence   int64      `json:"latestUpdateSequence"`
	CompactedUntilSequence int64      `json:"compactedUntilSequence"`
	ProjectedUntilSequence int64      `json:"projectedUntilSequence"`
	LastCompactedAt        *time.Time `json:"lastCompactedAt,omitempty"`
	UncompactedUpdateCount int64      `json:"uncompactedUpdateCount"`
	SnapshotBytes          int        `json:"snapshotBytes"`
	StateVectorBytes       int        `json:"stateVectorBytes"`
	Reason                 string     `json:"reason"`
}

type YjsMaintenanceRequestData struct {
	RequestId      uuid.UUID               `json:"requestId"`
	BlockPackId    uuid.UUID               `json:"blockPackId"`
	DocumentId     uuid.UUID               `json:"documentId"`
	Operation      YjsMaintenanceOperation `json:"operation"`
	TargetSequence int64                   `json:"targetSequence"`
	CorrelationId  string                  `json:"correlationId"`
}

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
