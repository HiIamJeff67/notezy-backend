package eventscontract

import (
	"time"

	"github.com/google/uuid"
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
