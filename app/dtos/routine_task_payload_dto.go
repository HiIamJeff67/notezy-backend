package dtos

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
)

type RoutineTaskNamePatternValue struct {
	Source   string  `json:"source" validate:"required,oneof=scheduledAt recordId shortRecordId routineTaskId"`
	Format   *string `json:"format" validate:"omitnil,max=64"`
	Timezone *string `json:"timezone" validate:"omitnil,max=64,istimezone"`
}

type RoutineTaskNamePattern map[string]RoutineTaskNamePatternValue

/* ============================== Root Shelf Routine Task Payload ============================== */

type CreateRootShelfRoutineTaskPayload struct {
	Id          *uuid.UUID             `json:"id" validate:"omitnil"`
	Name        string                 `json:"name" validate:"required,min=1,max=128,isshelfname"`
	NamePattern RoutineTaskNamePattern `json:"namePattern" validate:"omitempty,dive"`
}

type UpdateRootShelfRoutineTaskPayload struct {
	RootShelfId uuid.UUID               `json:"rootShelfId" validate:"required"`
	Name        *string                 `json:"name" validate:"omitnil,min=1,max=128,isshelfname"`
	NamePattern *RoutineTaskNamePattern `json:"namePattern" validate:"omitnil,dive"`
}

type ResetRootShelfRoutineTaskPayload struct {
	RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
}

/* ============================== Sub Shelf Routine Task Payload ============================== */

type CreateSubShelfRoutineTaskPayload struct {
	Id             *uuid.UUID             `json:"id" validate:"omitnil"`
	RootShelfId    uuid.UUID              `json:"rootShelfId" validate:"required"`
	PrevSubShelfId *uuid.UUID             `json:"prevSubShelfId" validate:"omitnil"`
	Name           string                 `json:"name" validate:"required,min=1,max=128,isshelfname"`
	NamePattern    RoutineTaskNamePattern `json:"namePattern" validate:"omitempty,dive"`
}

type UpdateSubShelfRoutineTaskPayload struct {
	SubShelfId  uuid.UUID               `json:"subShelfId" validate:"required"`
	Name        *string                 `json:"name" validate:"omitnil,min=1,max=128,isshelfname"`
	NamePattern *RoutineTaskNamePattern `json:"namePattern" validate:"omitnil,dive"`
}

type ResetSubShelfRoutineTaskPayload struct {
	SubShelfId uuid.UUID `json:"subShelfId" validate:"required"`
}

/* ============================== Block Pack Routine Task Payload ============================== */

type CreateBlockPackRoutineTaskTemplate struct {
	Name                string                 `json:"name" validate:"required,min=1,max=128"`
	NamePattern         RoutineTaskNamePattern `json:"namePattern" validate:"omitempty,dive"`
	Icon                *enums.SupportedIcon   `json:"icon" validate:"omitnil,issupportedicon"`
	HeaderBackgroundURL *string                `json:"headerBackgroundURL" validate:"omitnil"`
	Blocks              []struct {
		ClientId               string                 `json:"clientId" validate:"required"`
		PrevClientId           *string                `json:"prevClientId" validate:"omitnil"`
		ArborizedEditableBlock ArborizedEditableBlock `json:"arborizedEditableBlock" validate:"required"`
	} `json:"blocks" validate:"required,min=1"`
}

type CreateBlockPackRoutineTaskPayload struct {
	TargetSubShelfId uuid.UUID                          `json:"targetSubShelfId" validate:"required"`
	Template         CreateBlockPackRoutineTaskTemplate `json:"template" validate:"required"`
	Pattern          map[string]json.RawMessage         `json:"pattern" validate:"required"`
}

type UpdateBlockPackRoutineTaskPayload struct {
	BlockPackId   uuid.UUID `json:"blockPackId" validate:"required"`
	UpdatedBlocks []struct {
		BlockId                uuid.UUID               `json:"blockId" validate:"required"`
		ArborizedEditableBlock *ArborizedEditableBlock `json:"arborizedEditableBlock" validate:"required"`
	} `json:"updatedBlocks" validate:"required,min=1"`
}

type ResetBlockPackRoutineTaskPayload struct {
	BlockPackId uuid.UUID `json:"blockPackId" validate:"required"`
}

/* ============================== Block Routine Task Payload ============================== */

type AppendBlockRoutineTaskPayload struct {
	BlockPackId            uuid.UUID              `json:"blockPackId" validate:"required"`
	ArborizedEditableBlock ArborizedEditableBlock `json:"arborizedEditableBlock" validate:"required"`
}

type UpdateBlockRoutineTaskPayload struct {
	BlockId                uuid.UUID               `json:"blockId" validate:"required"`
	ArborizedEditableBlock *ArborizedEditableBlock `json:"arborizedEditableBlock" validate:"required"`
}

type ResetBlockRoutineTaskPayload struct {
	BlockId uuid.UUID `json:"blockId" validate:"required"`
}

/* ============================== Routine Routine Task Payload ============================== */

type CreateRoutineRoutineTaskPayload struct {
	Id               *uuid.UUID             `json:"id" validate:"omitnil"`
	StationId        uuid.UUID              `json:"stationId" validate:"required"`
	Title            string                 `json:"title" validate:"required,min=1,max=128"`
	TitlePattern     RoutineTaskNamePattern `json:"titlePattern" validate:"omitempty,dive"`
	Description      string                 `json:"description" validate:"max=1024"`
	Status           *enums.RoutineStatus   `json:"status" validate:"omitnil,isroutinestatus"`
	IsPinned         *bool                  `json:"isPinned" validate:"omitnil"`
	ScheduledStartAt *time.Time             `json:"scheduledStartAt" validate:"omitnil"`
	ScheduledEndAt   *time.Time             `json:"scheduledEndAt" validate:"omitnil"`
	Period           *enums.RoutinePeriod   `json:"period" validate:"omitnil,isroutineperiod"`
	Timezone         *string                `json:"timezone" validate:"omitnil,max=64,istimezone"`
}

type UpdateRoutineRoutineTaskPayload struct {
	RoutineId        uuid.UUID               `json:"routineId" validate:"required"`
	Title            *string                 `json:"title" validate:"omitnil,min=1,max=128"`
	TitlePattern     *RoutineTaskNamePattern `json:"titlePattern" validate:"omitnil,dive"`
	Description      *string                 `json:"description" validate:"omitnil,max=1024"`
	Status           *enums.RoutineStatus    `json:"status" validate:"omitnil,isroutinestatus"`
	IsPinned         *bool                   `json:"isPinned" validate:"omitnil"`
	ScheduledStartAt *time.Time              `json:"scheduledStartAt" validate:"omitnil"`
	ScheduledEndAt   *time.Time              `json:"scheduledEndAt" validate:"omitnil"`
	Period           *enums.RoutinePeriod    `json:"period" validate:"omitnil,isroutineperiod"`
	Timezone         *string                 `json:"timezone" validate:"omitnil,max=64,istimezone"`
}
