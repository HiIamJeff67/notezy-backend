package routinesdto

import (
	"time"

	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"
)

type LinkRoutineTagByIdRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			RoutineId    uuid.UUID `json:"routineId" validate:"required"`
			RoutineTagId uuid.UUID `json:"routineTagId" validate:"required"`
			IsUnlink     bool      `json:"isUnlink"`
		},
		struct{},
		struct{},
	]
}
type LinkRoutineTagByIdResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
type LinkRoutineTagsByIdsRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			LinkedRoutinesAndTags []LinkedRoutineAndTag `json:"linkedRoutinesAndTags" validate:"required,min=1,max=1024,dive"`
			IsUnlink              bool                  `json:"isUnlink"`
		},
		struct{},
		struct{},
	]
}
type LinkRoutineTagsByIdsResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
type LinkRoutineItemByIdRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			RoutineId uuid.UUID      `json:"routineId" validate:"required"`
			ItemId    uuid.UUID      `json:"itemId" validate:"required"`
			ItemType  enums.ItemType `json:"itemType" validate:"required,isitemtype"`
			IsUnlink  bool           `json:"isUnlink"`
		},
		struct{},
		struct{},
	]
}
type LinkRoutineItemByIdResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
type LinkRoutineItemsByIdsRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			LinkedRoutinesAndItems []LinkedRoutineAndItem `json:"linkedRoutinesAndItems" validate:"required,min=1,max=1024,dive"`
			IsUnlink               bool                   `json:"isUnlink"`
		},
		struct{},
		struct{},
	]
}
type LinkRoutineItemsByIdsResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
