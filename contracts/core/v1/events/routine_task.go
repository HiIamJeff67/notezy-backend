package eventscontract

import (
	"time"

	"github.com/google/uuid"
)

type RoutineTaskCompletedData struct {
	RoutineTaskId       uuid.UUID `json:"routineTaskId"`
	RoutineTaskRecordId uuid.UUID `json:"routineTaskRecordId"`
	WorkerId            uuid.UUID `json:"workerId"`
	Attempt             int32     `json:"attempt"`
	CompletedAt         time.Time `json:"completedAt"`
}
