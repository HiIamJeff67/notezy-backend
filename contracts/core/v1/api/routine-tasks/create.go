package routinetasksdto

import (
	"time"

	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
	routinetasktypes "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/types/routine-tasks"
)

type CreateRoutineTaskByRoutineIdRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		routinetasktypes.CreatableRoutineTask,
		struct{},
		struct{},
	]
}
type CreateRoutineTaskByRoutineIdResponseDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}
