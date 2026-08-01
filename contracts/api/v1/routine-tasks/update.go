package routinetasksdto

import (
	"time"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/api/v1"
)

type UpdateMyRoutineTaskByIdRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		UpdatableRoutineTask,
		struct{},
		struct{},
	]
}
type UpdateMyRoutineTaskByIdResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
