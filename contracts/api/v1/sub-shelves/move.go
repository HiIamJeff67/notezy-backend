package subshelvesdto

import (
	"time"

	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/api/v1"
)

type MoveMySubShelfByRootShelfIdRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			SourceRootShelfId      uuid.UUID  `json:"sourceRootShelfId" validate:"required"`
			SourceSubShelfId       uuid.UUID  `json:"sourceSubShelfId" validate:"required"`
			DestinationRootShelfId uuid.UUID  `json:"destinationRootShelfId" validate:"required"`
			DestinationSubShelfId  *uuid.UUID `json:"destinationSubShelfId" validate:"omitnil"`
		},
		struct{},
		struct{},
	]
}

type MoveMySubShelfByRootShelfIdResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveMySubShelvesByRootShelfIdRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		MovableSubShelf,
		struct{},
		struct{},
	]
}

type MoveMySubShelvesByRootShelfIdResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveMySubShelvesByRootShelfIdsRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			MovedSubShelves []MovableSubShelf `json:"movedSubShelves" validate:"required,min=1,max=1024,dive"`
		},
		struct{},
		struct{},
	]
}

type MoveMySubShelvesByRootShelfIdsResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
