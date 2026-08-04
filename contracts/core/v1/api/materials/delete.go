package materialsdto

import (
	"time"

	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
)

type DeleteMyMaterialByIdRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct {
			MaterialId uuid.UUID `json:"materialId" validate:"required"`
		},
		struct{},
	]
}

type DeleteMyMaterialByIdResponseDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type DeleteMyMaterialsByIdsRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			MaterialIds []uuid.UUID `json:"materialIds" validate:"required,min=1,max=1024"`
		},
		struct{},
		struct{},
	]
}

type DeleteMyMaterialsByIdsResponseDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
