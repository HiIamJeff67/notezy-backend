package routinesdto

import (
	"time"

	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/api/v1"
)

type RestoreMyRoutineByIdRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			RoutineId uuid.UUID `json:"routineId" validate:"required"`
		},
		struct{},
		struct{},
	]
}
type RestoreMyRoutineByIdResponseDto = RoutineResponseDto
type RestoreMyRoutinesByIdsRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			RoutineIds []uuid.UUID `json:"routineIds" validate:"required,min=1,max=1024"`
		},
		struct{},
		struct{},
	]
}
type RestoreMyRoutinesByIdsResponseDto []RoutineResponseDto
type DeleteMyRoutineByIdRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			RoutineId uuid.UUID `json:"routineId" validate:"required"`
		},
		struct{},
		struct{},
	]
}
type DeleteMyRoutineByIdResponseDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
type DeleteMyRoutinesByIdsRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			RoutineIds []uuid.UUID `json:"routineIds" validate:"required,min=1,max=1024"`
		},
		struct{},
		struct{},
	]
}
type DeleteMyRoutinesByIdsResponseDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
type HardDeleteMyRoutineByIdRequestDto = DeleteMyRoutineByIdRequestDto
type HardDeleteMyRoutineByIdResponseDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
type HardDeleteMyRoutinesByIdsRequestDto = DeleteMyRoutinesByIdsRequestDto
type HardDeleteMyRoutinesByIdsResponseDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
