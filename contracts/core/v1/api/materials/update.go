package materialsdto

import (
	"time"

	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
	materialstypes "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/types/materials"
)

type UpdateMyMaterialByIdRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			Values  materialstypes.UpdatableMaterial `json:"values"`
			SetNull *map[string]bool                 `json:"setNull,omitempty"`
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
