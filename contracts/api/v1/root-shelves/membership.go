package rootshelvesdto

import (
	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/api/v1"
)

type LeaveMyRootShelfRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct {
			RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
		},
		struct{},
	]
}

type LeaveMyRootShelfResponseDto struct{}

type LeaveMyRootShelvesRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			RootShelves []struct {
				RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
			} `json:"rootShelves" validate:"required,min=1,max=1024,dive"`
		},
		struct{},
		struct{},
	]
}

type LeaveMyRootShelvesResponseDto struct{}
