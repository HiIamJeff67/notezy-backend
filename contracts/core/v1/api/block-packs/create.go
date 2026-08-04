package blockpacksdto

import (
	"time"

	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
	blockpackstypes "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/types/block-packs"
)

type CreateBlockPackRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		blockpackstypes.CreatableBlockPack,
		struct{},
		struct{},
	]
}

type CreateBlockPackResponseDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateBlockPacksRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			CreatedBlockPacks []blockpackstypes.CreatableBlockPack `json:"createdBlockPacks" validate:"required,min=1,max=1024,dive"`
		},
		struct{},
		struct{},
	]
}

type CreateBlockPacksResponseDto struct {
	Ids       []uuid.UUID `json:"ids"`
	CreatedAt time.Time   `json:"createdAt"`
}
