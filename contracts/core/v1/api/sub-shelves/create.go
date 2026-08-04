package subshelvesdto

import (
	"time"

	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
	subshelvestypes "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/types/sub-shelves"
)

type CreateSubShelfByRootShelfIdRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		subshelvestypes.CreatableSubShelf,
		struct{},
		struct{},
	]
}

type CreateSubShelfByRootShelfIdResponseDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateSubShelvesByRootShelfIdsRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			CreatedSubShelves []subshelvestypes.CreatableSubShelf `json:"createdSubShelves" validate:"required,min=1,max=1024,dive"`
		},
		struct{},
		struct{},
	]
}

type CreateSubShelvesByRootShelfIdsResponseDto struct {
	Ids       []uuid.UUID `json:"ids"`
	CreatedAt time.Time   `json:"createdAt"`
}
