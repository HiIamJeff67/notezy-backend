package dtos

import (
	"encoding/json"

	"github.com/google/uuid"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
)

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

type DeleteBlockPackRoutineTaskPayload struct {
	BlockPackId uuid.UUID `json:"blockPackId" validate:"required"`
}

type CreateBlockRoutineTaskPayload struct {
	BlockGroupId           uuid.UUID              `json:"blockGroupId" validate:"required"`
	ParentBlockId          *uuid.UUID             `json:"parentBlockId" validate:"omitnil"`
	ArborizedEditableBlock ArborizedEditableBlock `json:"arborizedEditableBlock" validate:"required"`
}

type UpdateBlockRoutineTaskPayload struct {
	BlockId                uuid.UUID               `json:"blockId" validate:"required"`
	ArborizedEditableBlock *ArborizedEditableBlock `json:"arborizedEditableBlock" validate:"required"`
}

type DeleteBlockRoutineTaskPayload struct {
	BlockId uuid.UUID `json:"blockId" validate:"required"`
}
