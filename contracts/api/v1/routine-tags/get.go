package routinetagsdto

import (
	"time"

	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/api/v1"
)

type GetMyRoutineTagByIdRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct {
			RoutineTagId uuid.UUID `json:"routineTagId" validate:"required"`
			IsDeleted    *bool     `json:"isDeleted" validate:"omitnil"`
		},
		struct{},
	]
}

type RoutineTagResponseDto struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	Icon      *string   `json:"icon"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type GetMyRoutineTagByIdResponseDto = RoutineTagResponseDto

type GetAllMyRoutineTagsRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct {
			AreDeleted *bool `json:"areDeleted" validate:"omitnil"`
		},
		struct{},
	]
}

type GetAllMyRoutineTagsResponseDto []RoutineTagResponseDto
