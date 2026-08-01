package rootshelvesdto

import (
	"time"

	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/api/v1"
)

type TransferMyRootShelfOwnershipRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			TargetUserPublicId uuid.UUID `json:"targetUserPublicId" validate:"required"`
		},
		struct {
			RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
		},
		struct{},
	]
}

type TransferMyRootShelfOwnershipResponseDto struct {
	RootShelfId               uuid.UUID `json:"rootShelfId"`
	PreviousOwnerUserPublicId uuid.UUID `json:"previousOwnerUserPublicId"`
	NewOwnerUserPublicId      uuid.UUID `json:"newOwnerUserPublicId"`
	UpdatedAt                 time.Time `json:"updatedAt"`
}
