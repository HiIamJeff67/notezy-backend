package routinetasktypes

import (
	"github.com/google/uuid"

	blocknote "github.com/HiIamJeff67/notezy-backend/contracts/types/blocknote"
)

type AppendBlockRoutineTaskPayload struct {
	BlockPackId            uuid.UUID                        `json:"blockPackId" validate:"required"`
	Pattern                RoutineTaskPattern               `json:"pattern" validate:"omitempty,dive"`
	ArborizedEditableBlock blocknote.ArborizedEditableBlock `json:"arborizedEditableBlock" validate:"required"`
}

type UpdateBlockRoutineTaskPayload struct {
	BlockId                uuid.UUID                         `json:"blockId" validate:"required"`
	Pattern                RoutineTaskPattern                `json:"pattern" validate:"omitempty,dive"`
	ArborizedEditableBlock *blocknote.ArborizedEditableBlock `json:"arborizedEditableBlock" validate:"required"`
}

type ResetBlockRoutineTaskPayload struct {
	BlockId uuid.UUID `json:"blockId" validate:"required"`
}
