package apicontract

import (
	"time"

	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
)

type RestoreMyRoutineByIdRequestDto struct {
	coreapicontract.RequestDto[
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
	coreapicontract.RequestDto[
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
	coreapicontract.RequestDto[
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
	coreapicontract.RequestDto[
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
