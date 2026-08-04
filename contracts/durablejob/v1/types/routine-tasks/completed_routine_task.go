package routinetasktypes

import (
	"time"

	"github.com/google/uuid"
)

type CompletedRoutineTask struct {
	RoutineTaskId       uuid.UUID `json:"routineTaskId" validate:"required"`
	RoutineTaskRecordId uuid.UUID `json:"routineTaskRecordId" validate:"required"`
	CompletedAt         time.Time `json:"completedAt" validate:"required"`
}
