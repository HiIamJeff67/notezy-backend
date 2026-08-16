package apicontract

import (
	"time"

	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api"
)

type DeleteMySubShelfByIdRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct {
			SubShelfId uuid.UUID `json:"subShelfId" validate:"required"`
		},
		struct{},
	]
}

type DeleteMySubShelfByIdResponseDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type DeleteMySubShelvesByIdsRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			SubShelfIds []uuid.UUID `json:"subShelfIds" validate:"required,min=1,max=1024"`
		},
		struct{},
		struct{},
	]
}

type DeleteMySubShelvesByIdsResponseDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
