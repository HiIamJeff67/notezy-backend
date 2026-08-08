package apicontract

import (
	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
)

type GetMyRoutineTaskByIdRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct {
			RoutineTaskId uuid.UUID `json:"routineTaskId" validate:"required"`
			IsDeleted     *bool     `json:"isDeleted" validate:"omitnil"`
		},
		struct{},
	]
}
type GetMyRoutineTaskByIdResponseDto = RoutineTaskResponseDto
type GetAllMyRoutineTasksByRoutineIdsRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct {
			RoutineIds []uuid.UUID `json:"routineIds" validate:"required,min=1,max=1024"`
			AreDeleted *bool       `json:"areDeleted" validate:"omitnil"`
		},
		struct{},
	]
}
type GetAllMyRoutineTasksByRoutineIdsResponseDto []RoutineTaskResponseDto
type GetAllMyRoutineTasksRequestDto struct {
	coreapicontract.RequestDto[
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
type GetAllMyRoutineTasksResponseDto []RoutineTaskResponseDto
