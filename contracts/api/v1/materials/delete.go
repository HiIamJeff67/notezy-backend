package materialsdto

import (
	"time"

	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/api/v1"
)

type DeleteMyMaterialByIdRequestDto struct {
	apiv1.RequestDto[
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
	apiv1.RequestDto[
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
