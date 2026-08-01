package routinesdto

import (
	"time"

	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/api/v1"
)

type GetMyRoutineByIdRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct {
			RoutineId uuid.UUID `json:"routineId" validate:"required"`
			IsDeleted *bool     `json:"isDeleted" validate:"omitnil"`
		},
		struct{},
	]
}
type GetMyRoutineByIdResponseDto = RoutineResponseDto

type GetMyRoutinesByStationIdRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct {
			StationId  uuid.UUID `json:"stationId" validate:"required"`
			AreDeleted *bool     `json:"areDeleted" validate:"omitnil"`
		},
		struct{},
	]
}
type GetMyRoutinesByStationIdResponseDto []RoutineResponseDto

type GetAllMyRoutinesByTimeRangeRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct {
			From       time.Time   `json:"from" validate:"required"`
			To         time.Time   `json:"to" validate:"required"`
			StationIds []uuid.UUID `json:"stationIds" validate:"required,min=1,max=1024"`
			AreDeleted *bool       `json:"areDeleted" validate:"omitnil"`
		},
		struct{},
	]
}
type GetAllMyRoutinesByTimeRangeResponseDto []RoutineResponseDto
