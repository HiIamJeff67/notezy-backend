package routinetasksdto

import (
	"time"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
	routinetasktypes "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/types/routine-tasks"
)

type UpdateMyRoutineTaskByIdRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		routinetasktypes.UpdatableRoutineTask,
		struct{},
		struct{},
	]
}
type UpdateMyRoutineTaskByIdResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
