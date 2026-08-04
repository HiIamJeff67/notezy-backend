package subshelvesdto

import (
	"time"

	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
	subshelvestypes "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/types/sub-shelves"
)

type UpdateMySubShelfByIdRequestDto struct {
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
			SubShelfId uuid.UUID `json:"subShelfId" validate:"required"`
		},
		struct{},
	]
}

type UpdateMySubShelfByIdResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateMySubShelvesByIdsRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UpdatedSubShelves []subshelvestypes.UpdatableSubShelf `json:"updatedSubShelves" validate:"required,min=1,max=1024,dive"`
		},
		struct{},
		struct{},
	]
}

type UpdateMySubShelvesByIdsResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
