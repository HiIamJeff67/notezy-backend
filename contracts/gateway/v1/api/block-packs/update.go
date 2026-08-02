package blockpacksdto

import (
	"time"

	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"
)

type UpdateMyBlockPackByIdRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			Values struct {
				Name                *string              `json:"name" validate:"omitnil,min=1,max=128"`
				Icon                *enums.SupportedIcon `json:"icon" validate:"omitnil,issupportedicon"`
				HeaderBackgroundURL *string              `json:"headerBackgroundURL" validate:"omitnil"`
			} `json:"values"`
			SetNull *map[string]bool `json:"setNull,omitempty"`
		},
		struct {
			BlockPackId uuid.UUID `json:"blockPackId" validate:"required"`
		},
		struct{},
	]
}

type UpdateMyBlockPackByIdResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateMyBlockPacksByIdsRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UpdatedBlockPacks []UpdatableBlockPack `json:"updatedBlockPacks" validate:"required,min=1,max=1024,dive"`
		},
		struct{},
		struct{},
	]
}

type UpdateMyBlockPacksByIdsResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
