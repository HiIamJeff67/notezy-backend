package blockpacksdto

import (
	"time"

	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
)

type DeleteMyBlockPackByIdRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct {
			BlockPackId uuid.UUID `json:"blockPackId" validate:"required"`
		},
		struct{},
	]
}

type DeleteMyBlockPackByIdResponseDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type DeleteMyBlockPacksByIdsRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			BlockPackIds []uuid.UUID `json:"blockPackIds" validate:"required,min=1,max=1024"`
		},
		struct{},
		struct{},
	]
}

type DeleteMyBlockPacksByIdsResponseDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
