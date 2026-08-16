package apicontract

import (
	"time"

	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api"
)

type SaveMyMaterialByIdRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			ContentFile []byte `json:"contentFile" validate:"required"`
		},
		struct {
			MaterialId uuid.UUID `json:"materialId" validate:"required"`
		},
		struct{},
	]
}

type SaveMyMaterialByIdResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
