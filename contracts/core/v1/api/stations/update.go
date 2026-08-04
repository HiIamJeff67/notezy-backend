package stationsdto

import (
	"time"

	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
	stationstypes "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/types/stations"
)

type UpdateMyStationByIdRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			Values struct {
				Name                *string `json:"name,omitempty" validate:"omitnil,min=1,max=128"`
				Description         *string `json:"description,omitempty" validate:"omitnil,max=1024"`
				Icon                *string `json:"icon,omitempty" validate:"omitnil,issupportedicon"`
				HeaderBackgroundURL *string `json:"headerBackgroundURL,omitempty" validate:"omitnil,isimageurl"`
			} `json:"values"`
			SetNull *map[string]bool `json:"setNull,omitempty"`
		},
		struct {
			StationId uuid.UUID `json:"stationId" validate:"required"`
		},
		struct{},
	]
}

type UpdateMyStationByIdResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateMyStationsByIdsRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UpdatedStations []stationstypes.UpdatableStation `json:"updatedStations" validate:"required,min=1,max=1024,dive"`
		},
		struct{},
		struct{},
	]
}

type UpdateMyStationsByIdsResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
