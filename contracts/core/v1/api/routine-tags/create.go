package routinetagsdto

import (
	"time"

	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
	routinetagstypes "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/types/routine-tags"
)

type CreateRoutineTagRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		routinetagstypes.CreatableRoutineTag,
		struct{},
		struct{},
	]
}

type CreateRoutineTagResponseDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateRoutineTagsRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			CreatedRoutineTags []routinetagstypes.CreatableRoutineTag `json:"createdRoutineTags" validate:"required,min=1,max=1024,dive"`
		},
		struct{},
		struct{},
	]
}

type CreateRoutineTagsResponseDto struct {
	Ids       []uuid.UUID `json:"ids"`
	CreatedAt time.Time   `json:"createdAt"`
}
