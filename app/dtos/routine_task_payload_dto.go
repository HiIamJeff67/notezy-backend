package dtos

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
)

/* ============================== Root Shelf Routine Task Payload ============================== */

type CreateRootShelfRoutineTaskPayload struct {
	Id   *uuid.UUID `json:"id" validate:"omitnil"`
	Name string     `json:"name" validate:"required,min=1,max=128,isshelfname"`
}

type UpdateRootShelfRoutineTaskPayload struct {
	RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
	Name        *string   `json:"name" validate:"omitnil,min=1,max=128,isshelfname"`
}

type ResetRootShelfRoutineTaskPayload struct {
	RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
}

/* ============================== Sub Shelf Routine Task Payload ============================== */

type CreateSubShelfRoutineTaskPayload struct {
	Id             *uuid.UUID `json:"id" validate:"omitnil"`
	RootShelfId    uuid.UUID  `json:"rootShelfId" validate:"required"`
	PrevSubShelfId *uuid.UUID `json:"prevSubShelfId" validate:"omitnil"`
	Name           string     `json:"name" validate:"required,min=1,max=128,isshelfname"`
}

type UpdateSubShelfRoutineTaskPayload struct {
	SubShelfId uuid.UUID `json:"subShelfId" validate:"required"`
	Name       *string   `json:"name" validate:"omitnil,min=1,max=128,isshelfname"`
}

type ResetSubShelfRoutineTaskPayload struct {
	SubShelfId uuid.UUID `json:"subShelfId" validate:"required"`
}

/* ============================== Block Pack Routine Task Payload ============================== */

type CreateBlockPackRoutineTaskTemplate struct {
	Name                    string               `json:"name" validate:"required,min=1,max=128"`
	Icon                    *enums.SupportedIcon `json:"icon" validate:"omitnil,issupportedicon"`
	HeaderBackgroundURL     *string              `json:"headerBackgroundURL" validate:"omitnil"`
	FinalBlockGroupClientId *string              `json:"finalBlockGroupClientId" validate:"omitnil"`
	BlockGroups             []struct {
		ClientId               string                 `json:"clientId" validate:"required"`
		PrevClientId           *string                `json:"prevClientId" validate:"omitnil"`
		ArborizedEditableBlock ArborizedEditableBlock `json:"arborizedEditableBlock" validate:"required"`
	} `json:"blockGroups" validate:"required,min=1"`
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

type UpdateRoutineRoutineTaskPayload struct {
	RoutineId        uuid.UUID            `json:"routineId" validate:"required"`
	Title            *string              `json:"title" validate:"omitnil,min=1,max=128"`
	Description      *string              `json:"description" validate:"omitnil,max=1024"`
	Status           *enums.RoutineStatus `json:"status" validate:"omitnil,isroutinestatus"`
	IsPinned         *bool                `json:"isPinned" validate:"omitnil"`
	ScheduledStartAt *time.Time           `json:"scheduledStartAt" validate:"omitnil"`
	ScheduledEndAt   *time.Time           `json:"scheduledEndAt" validate:"omitnil"`
	Period           *enums.RoutinePeriod `json:"period" validate:"omitnil,isroutineperiod"`
	Timezone         *string              `json:"timezone" validate:"omitnil,max=64,istimezone"`
}
