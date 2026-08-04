package routinetasktypes

import (
	"github.com/google/uuid"

	typescontract "github.com/HiIamJeff67/notezy-backend/contracts/types"
)

type AppendBlockRoutineTaskPayload struct {
	BlockPackId            uuid.UUID                            `json:"blockPackId" validate:"required"`
	Pattern                RoutineTaskPattern                   `json:"pattern" validate:"omitempty,dive"`
	ArborizedEditableBlock typescontract.ArborizedEditableBlock `json:"arborizedEditableBlock" validate:"required"`
}

type UpdateBlockRoutineTaskPayload struct {
	BlockId                uuid.UUID                             `json:"blockId" validate:"required"`
	Pattern                RoutineTaskPattern                    `json:"pattern" validate:"omitempty,dive"`
	ArborizedEditableBlock *typescontract.ArborizedEditableBlock `json:"arborizedEditableBlock" validate:"required"`
}

type ResetBlockRoutineTaskPayload struct {
	BlockId uuid.UUID `json:"blockId" validate:"required"`
}
