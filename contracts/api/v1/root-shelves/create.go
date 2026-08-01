package rootshelvesdto

import (
	"time"

	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/api/v1"
)

type CreateRootShelfRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		CreatableRootShelf,
		struct{},
		struct{},
	]
}

type CreateRootShelfResponseDto struct {
	Id             uuid.UUID `json:"id"`
	LastAnalyzedAt time.Time `json:"lastAnalyzedAt"`
	CreatedAt      time.Time `json:"createdAt"`
}

type CreateRootShelvesRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			RootShelves []CreatableRootShelf `json:"insertedRootShelves" validate:"required,min=1,max=1024,dive"`
		},
		struct{},
		struct{},
	]
}

type CreateRootShelvesResponseDto struct {
	Ids            []uuid.UUID `json:"ids"`
	LastAnalyzedAt time.Time   `json:"lastAnalyzedAt"`
	CreatedAt      time.Time   `json:"createdAt"`
}
