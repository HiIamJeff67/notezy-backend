package durablejobeventscontract

import (
	"github.com/google/uuid"

	yjsworkereventscontract "github.com/HiIamJeff67/notezy-backend/contracts/yjs-worker/v1/events"
)

type YjsMaintenanceRequestData struct {
	RequestId      uuid.UUID                                       `json:"requestId"`
	BlockPackId    uuid.UUID                                       `json:"blockPackId"`
	DocumentId     uuid.UUID                                       `json:"documentId"`
	Operation      yjsworkereventscontract.YjsMaintenanceOperation `json:"operation"`
	TargetSequence int64                                           `json:"targetSequence"`
	CorrelationId  string                                          `json:"correlationId"`
}

type YjsMaintenanceResultData struct {
	RequestId              uuid.UUID                                       `json:"requestId"`
	BlockPackId            uuid.UUID                                       `json:"blockPackId"`
	DocumentId             uuid.UUID                                       `json:"documentId"`
	Operation              yjsworkereventscontract.YjsMaintenanceOperation `json:"operation"`
	TargetSequence         int64                                           `json:"targetSequence"`
	Success                bool                                            `json:"success"`
	CompactedUntilSequence int64                                           `json:"compactedUntilSequence"`
	ProjectedUntilSequence int64                                           `json:"projectedUntilSequence"`
	Error                  string                                          `json:"error,omitempty"`
}
