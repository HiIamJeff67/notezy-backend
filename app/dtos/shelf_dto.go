package dtos

import (
	"time"

	"github.com/google/uuid"
)

/* ============================== Request DTO ============================== */

type SynchronizeShelvesReqDto struct {
	OwnerId        uuid.UUID   // extracted from the access token of AuthMiddleware()
	ShelfIds       []uuid.UUID `json:"shelfIds" validate:"required"`
	PartialUpdates []PartialUpdateDto[struct {
		Name             *string `json:"name" validate:"omitnil,max=128,isshelfname"`
		EncodedStructure *[]byte `json:"encodedStructure" validate:"omitnil"`
	}]
}

/* ============================== Response DTO ============================== */

type SynchronizeShelvesResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
