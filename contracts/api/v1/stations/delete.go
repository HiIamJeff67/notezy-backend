package stationsdto

import (
	"time"

	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/api/v1"
)

type DeleteMyStationByIdRequestDto struct {
	apiv1.RequestDto[
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

type DeleteMyStationByIdResponseDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type DeleteMyStationsByIdsRequestDto struct {
	apiv1.RequestDto[
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

type DeleteMyStationsByIdsResponseDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type HardDeleteMyStationByIdRequestDto struct {
	apiv1.RequestDto[
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

type HardDeleteMyStationByIdResponseDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type HardDeleteMyStationsByIdsRequestDto struct {
	apiv1.RequestDto[
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

type HardDeleteMyStationsByIdsResponseDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
