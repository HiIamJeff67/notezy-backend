package apicontract

import (
	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
)

type LeaveMyRootShelfRequestDto struct {
	coreapicontract.RequestDto[
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
	coreapicontract.RequestDto[
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
