package materialsdto

import (
	"time"

	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/api/v1"
)

type UpdateMyMaterialByIdRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			Values  UpdatableMaterial `json:"values"`
			SetNull *map[string]bool  `json:"setNull,omitempty"`
		},
		struct {
			MaterialId uuid.UUID `json:"materialId" validate:"required"`
		},
		struct{},
	]
}

type UpdateMyMaterialByIdResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
