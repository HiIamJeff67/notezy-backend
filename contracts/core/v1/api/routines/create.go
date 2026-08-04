package routinesdto

import (
	"time"

	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
	routinestypes "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/types/routines"
)

type CreateRoutineByStationIdRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		routinestypes.CreatableRoutine,
		struct{},
		struct{},
	]
}
type CreateRoutineByStationIdResponseDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}
type CreateRoutinesByStationIdsRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			CreatedRoutines []routinestypes.CreatableRoutine `json:"createdRoutines" validate:"required,min=1,max=1024,dive"`
		},
		struct{},
		struct{},
	]
}
type CreateRoutinesByStationIdsResponseDto struct {
	Ids       []uuid.UUID `json:"ids"`
	CreatedAt time.Time   `json:"createdAt"`
}
