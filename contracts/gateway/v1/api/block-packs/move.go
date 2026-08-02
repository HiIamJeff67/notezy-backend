package blockpacksdto

import (
	"time"

	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api"
)

type MoveMyBlockPackByParentSubShelfIdRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			BlockPackId                 uuid.UUID `json:"blockPackId" validate:"required"`
			DestinationParentSubShelfId uuid.UUID `json:"destinationParentSubShelfId" validate:"required"`
		},
		struct{},
		struct{},
	]
}

type MoveMyBlockPackByParentSubShelfIdResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveMyBlockPacksByParentSubShelfIdRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		MovableBlockPack,
		struct{},
		struct{},
	]
}

type MoveMyBlockPacksByParentSubShelfIdResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveMyBlockPacksByParentSubShelfIdsRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			MovedBlockPacks []MovableBlockPack `json:"movedBlockPacks" validate:"required,min=1,max=1024,dive"`
		},
		struct{},
		struct{},
	]
}

type MoveMyBlockPacksByParentSubShelfIdsResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
