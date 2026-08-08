package apicontract

import (
	"time"

	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
	coretypes "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/types/root-shelves"
)

type UpdateMyRootShelfByIdRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			Values struct {
				Name *string `json:"name" validate:"omitnil,min=1,max=128,isshelfname"`
			} `json:"values"`
			SetNull *map[string]bool `json:"setNull,omitempty"`
		},
		struct {
			RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
		},
		struct{},
	]
}

type UpdateMyRootShelfByIdResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateMyRootShelvesByIdsRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UpdatedRootShelves []coretypes.UpdatableRootShelf `json:"updatedRootShelves" validate:"required,min=1,max=1024,dive"`
		},
		struct{},
		struct{},
	]
}

type UpdateMyRootShelvesByIdsResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
