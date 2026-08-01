package routinesdto

import (
	"time"

	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/schemas/enums"
	"github.com/google/uuid"
)

type RoutineResponseDto struct {
	Id               uuid.UUID            `json:"id"`
	StationId        uuid.UUID            `json:"stationId"`
	Title            string               `json:"title"`
	Description      string               `json:"description"`
	Status           enums.RoutineStatus  `json:"status"`
	IsPinned         bool                 `json:"isPinned"`
	ScheduledStartAt time.Time            `json:"scheduledStartAt"`
	ScheduledEndAt   time.Time            `json:"scheduledEndAt"`
	Period           *enums.RoutinePeriod `json:"period"`
	Timezone         string               `json:"timezone"`
	DeletedAt        *time.Time           `json:"deletedAt"`
	UpdatedAt        time.Time            `json:"updatedAt"`
	CreatedAt        time.Time            `json:"createdAt"`
	TagIds           []uuid.UUID          `json:"tagIds"`
	TaskIds          []uuid.UUID          `json:"taskIds"`
	ItemIds          []uuid.UUID          `json:"itemIds"`
}

type CreatableRoutine struct {
	Id               *uuid.UUID           `json:"id" validate:"omitnil"`
	StationId        uuid.UUID            `json:"stationId" validate:"required"`
	Title            string               `json:"title" validate:"required,min=1,max=128"`
	Description      string               `json:"description" validate:"max=1024"`
	Status           *enums.RoutineStatus `json:"status" validate:"omitnil,isroutinestatus"`
	IsPinned         *bool                `json:"isPinned" validate:"omitnil"`
	ScheduledStartAt *time.Time           `json:"scheduledStartAt" validate:"omitnil"`
	ScheduledEndAt   *time.Time           `json:"scheduledEndAt" validate:"omitnil"`
	Period           *enums.RoutinePeriod `json:"period" validate:"omitnil,isroutineperiod"`
	Timezone         *string              `json:"timezone" validate:"omitnil,max=64,istimezone"`
}

type UpdatableRoutine struct {
	RoutineId uuid.UUID `json:"routineId" validate:"required"`
	Values    struct {
		StationId        *uuid.UUID           `json:"stationId" validate:"omitnil"`
		Title            *string              `json:"title" validate:"omitnil,min=1,max=128"`
		Description      *string              `json:"description" validate:"omitnil,max=1024"`
		Status           *enums.RoutineStatus `json:"status" validate:"omitnil,isroutinestatus"`
		IsPinned         *bool                `json:"isPinned" validate:"omitnil"`
		ScheduledStartAt *time.Time           `json:"scheduledStartAt" validate:"omitnil"`
		ScheduledEndAt   *time.Time           `json:"scheduledEndAt" validate:"omitnil"`
		Period           *enums.RoutinePeriod `json:"period" validate:"omitnil,isroutineperiod"`
		Timezone         *string              `json:"timezone" validate:"omitnil,max=64,istimezone"`
	} `json:"values"`
	SetNull *map[string]bool `json:"setNull,omitempty"`
}

type LinkedRoutineAndTag struct {
	RoutineId    uuid.UUID `json:"routineId" validate:"required"`
	RoutineTagId uuid.UUID `json:"routineTagId" validate:"required"`
}

type LinkedRoutineAndItem struct {
	RoutineId uuid.UUID      `json:"routineId" validate:"required"`
	ItemId    uuid.UUID      `json:"itemId" validate:"required"`
	ItemType  enums.ItemType `json:"itemType" validate:"required,isitemtype"`
}
