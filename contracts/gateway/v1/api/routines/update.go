package routinesdto

import (
	"time"

	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api"
)

type UpdateMyRoutineByIdRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		UpdatableRoutine,
		struct {
			RoutineId uuid.UUID `json:"routineId" validate:"required"`
		},
		struct{},
	]
}
type UpdateMyRoutineByIdResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
type UpdateMyRoutinesByIdsRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UpdatedRoutines []UpdatableRoutine `json:"updatedRoutines" validate:"required,min=1,max=1024,dive"`
		},
		struct{},
		struct{},
	]
}
type UpdateMyRoutinesByIdsResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
