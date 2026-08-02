package routinetasksdto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"
)

type RoutineTaskResponseDto struct {
	Id              uuid.UUID                `json:"id"`
	RoutineId       uuid.UUID                `json:"routineId"`
	Title           string                   `json:"title"`
	Purpose         enums.RoutineTaskPurpose `json:"purpose"`
	Payload         datatypes.JSON           `json:"payload"`
	CostUnit        int64                    `json:"costUnit"`
	Priority        int32                    `json:"priority"`
	Status          enums.RoutineTaskStatus  `json:"status"`
	Attempts        int32                    `json:"attempts"`
	MaxAttempts     int32                    `json:"maxAttempts"`
	Period          *enums.RoutinePeriod     `json:"period"`
	NextScheduledAt time.Time                `json:"nextScheduledAt"`
	ScheduledAt     time.Time                `json:"scheduledAt"`
	ActualStartedAt *time.Time               `json:"actualStartedAt"`
	ActualEndedAt   *time.Time               `json:"actualEndedAt"`
	UpdatedAt       time.Time                `json:"updatedAt"`
	CreatedAt       time.Time                `json:"createdAt"`
}

type CreatableRoutineTask struct {
	RoutineId       uuid.UUID                `json:"routineId" validate:"required"`
	Title           string                   `json:"title" validate:"required,min=1,max=128"`
	Purpose         enums.RoutineTaskPurpose `json:"purpose" validate:"required,isroutinetaskpurpose"`
	Payload         datatypes.JSON           `json:"payload" validate:"omitempty,max=16777216"`
	Priority        int32                    `json:"priority" validate:"omitempty,min=0,max=100"`
	MaxAttempts     int32                    `json:"maxAttempts" validate:"omitempty,min=1,max=20"`
	Period          *enums.RoutinePeriod     `json:"period" validate:"omitnil,isroutineperiod"`
	NextScheduledAt time.Time                `json:"nextScheduledAt" validate:"required"`
}

type UpdatableRoutineTask struct {
	RoutineTaskId uuid.UUID `json:"routineTaskId" validate:"required"`
	Values        struct {
		RoutineId       *uuid.UUID                `json:"routineId" validate:"omitnil"`
		Title           *string                   `json:"title" validate:"omitnil,min=1,max=128"`
		Purpose         *enums.RoutineTaskPurpose `json:"purpose" validate:"omitnil,isroutinetaskpurpose"`
		Payload         *datatypes.JSON           `json:"payload" validate:"omitnil,max=16777216"`
		Priority        *int32                    `json:"priority" validate:"omitnil,min=0,max=100"`
		MaxAttempts     *int32                    `json:"maxAttempts" validate:"omitnil,min=1,max=20"`
		Period          *enums.RoutinePeriod      `json:"period" validate:"omitnil,isroutineperiod"`
		NextScheduledAt *time.Time                `json:"nextScheduledAt" validate:"omitnil"`
	} `json:"values"`
	SetNull *map[string]bool `json:"setNull,omitempty"`
}
