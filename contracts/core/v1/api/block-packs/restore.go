package blockpacksdto

import (
	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
)

type RestoreMyBlockPackByIdRequestDto struct {
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

type RestoreMyBlockPackByIdResponseDto = BlockPackResponseDto

type RestoreMyBlockPacksByIdsRequestDto struct {
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

type RestoreMyBlockPacksByIdsResponseDto []BlockPackResponseDto
