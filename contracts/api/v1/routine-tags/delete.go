package routinetagsdto

import (
	"time"

	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/api/v1"
)

type HardDeleteMyRoutineTagByIdRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct {
			RoutineTagId uuid.UUID `json:"routineTagId" validate:"required"`
		},
		struct{},
	]
}

type HardDeleteMyRoutineTagByIdResponseDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type HardDeleteMyRoutineTagsByIdsRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			RoutineTagIds []uuid.UUID `json:"routineTagIds" validate:"required,min=1,max=1024"`
		},
		struct{},
		struct{},
	]
}

type HardDeleteMyRoutineTagsByIdsResponseDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
