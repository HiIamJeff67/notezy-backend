package subshelvesdto

import (
	"time"

	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
	subshelvestypes "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/types/sub-shelves"
)

type MoveMySubShelfByRootShelfIdRequestDto struct {
	coreapicontract.RequestDto[
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
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		subshelvestypes.MovableSubShelf,
		struct{},
		struct{},
	]
}

type MoveMySubShelvesByRootShelfIdResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveMySubShelvesByRootShelfIdsRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			MovedSubShelves []subshelvestypes.MovableSubShelf `json:"movedSubShelves" validate:"required,min=1,max=1024,dive"`
		},
		struct{},
		struct{},
	]
}

type MoveMySubShelvesByRootShelfIdsResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
