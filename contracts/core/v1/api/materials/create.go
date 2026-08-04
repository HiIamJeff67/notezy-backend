package materialsdto

import (
	"time"

	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
	materialstypes "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/types/materials"
)

type CreateMyMaterialRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		materialstypes.CreatableMaterial,
		struct{},
		struct{},
	]
}

type CreateMyMaterialResponseDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}
