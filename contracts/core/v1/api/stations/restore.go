package apicontract

import (
	"time"

	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
)

type RestoreMyStationByIdRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			StationId uuid.UUID `json:"stationId" validate:"required"`
		},
		struct{},
		struct{},
	]
}

type RestoreMyStationByIdResponseDto struct {
	Id                  uuid.UUID  `json:"id"`
	Name                string     `json:"name"`
	Description         string     `json:"description"`
	Icon                *string    `json:"icon"`
	HeaderBackgroundURL *string    `json:"headerBackgroundURL"`
	RoutineCount        int64      `json:"routineCount"`
	DeletedAt           *time.Time `json:"deletedAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
	CreatedAt           time.Time  `json:"createdAt"`
}

type RestoreMyStationsByIdsRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			StationIds []uuid.UUID `json:"stationIds" validate:"required,min=1,max=1024"`
		},
		struct{},
		struct{},
	]
}

type RestoreMyStationsByIdsResponseDto []RestoreMyStationByIdResponseDto
