package coretypes

import "github.com/google/uuid"

type MovableSubShelf struct {
	SourceRootShelfId      uuid.UUID   `json:"sourceRootShelfId" validate:"required"`
	SourceSubShelfIds      []uuid.UUID `json:"sourceSubShelfIds" validate:"required,min=1,max=1024"`
	DestinationRootShelfId uuid.UUID   `json:"destinationRootShelfId" validate:"required"`
	DestinationSubShelfId  *uuid.UUID  `json:"destinationSubShelfId" validate:"omitnil"`
}
