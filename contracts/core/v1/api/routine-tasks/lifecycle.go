package routinetasksdto

import (
	"time"

	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
)

type PauseMyRoutineTaskByIdRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			RoutineTaskId uuid.UUID `json:"routineTaskId" validate:"required"`
		},
		struct{},
		struct{},
	]
}
type PauseMyRoutineTaskByIdResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
type ResumeMyRoutineTaskByIdRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			RoutineTaskId uuid.UUID `json:"routineTaskId" validate:"required"`
		},
		struct{},
		struct{},
	]
}
type ResumeMyRoutineTaskByIdResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
type HardDeleteMyRoutineTaskByIdRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			RoutineTaskId uuid.UUID `json:"routineTaskId" validate:"required"`
		},
		struct{},
		struct{},
	]
}
type HardDeleteMyRoutineTaskByIdResponseDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
type HardDeleteMyRoutineTasksByIdsRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			RoutineTaskIds []uuid.UUID `json:"routineTaskIds" validate:"required,min=1,max=1024"`
		},
		struct{},
		struct{},
	]
}
type HardDeleteMyRoutineTasksByIdsResponseDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
