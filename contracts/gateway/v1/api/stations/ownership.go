package stationsdto

import (
	"time"

	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api"
)

type TransferMyStationOwnershipRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			TargetUserPublicId uuid.UUID `json:"targetUserPublicId" validate:"required"`
		},
		struct {
			StationId uuid.UUID `json:"stationId" validate:"required"`
		},
		struct{},
	]
}

type TransferMyStationOwnershipResponseDto struct {
	StationId                 uuid.UUID `json:"stationId"`
	PreviousOwnerUserPublicId uuid.UUID `json:"previousOwnerUserPublicId"`
	NewOwnerUserPublicId      uuid.UUID `json:"newOwnerUserPublicId"`
	UpdatedAt                 time.Time `json:"updatedAt"`
}
