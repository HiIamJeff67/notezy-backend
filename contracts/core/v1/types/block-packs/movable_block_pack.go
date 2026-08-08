package coretypes

import "github.com/google/uuid"

type MovableBlockPack struct {
	BlockPackIds                []uuid.UUID `json:"blockPackIds" validate:"required,min=1,max=100"`
	DestinationParentSubShelfId uuid.UUID   `json:"destinationParentSubShelfId" validate:"required"`
}
