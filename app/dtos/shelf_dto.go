package dtos

import (
	"time"

	"github.com/google/uuid"
)

/* ============================== Request DTO ============================== */

type CreateShelfReqDto struct {
	OwnerId uuid.UUID // extracted from the access token of AuthMiddleware()
	Name    string    `json:"name" validate:"required,max=128"`
}

type SynchronizeShelvesReqDto struct {
	OwnerId        uuid.UUID   // extracted from the access token of AuthMiddleware()
	ShelfIds       []uuid.UUID `json:"shelfIds" validate:"required"`
	PartialUpdates []PartialUpdateDto[struct {
		Name             *string `json:"name" validate:"omitnil,max=128,isshelfname"`
		EncodedStructure *[]byte `json:"encodedStructure" validate:"omitnil"`
	}]
}

/* ============================== Response DTO ============================== */

type CreateShelfResDto struct {
	CreatedAt time.Time `json:"createdAt"`
}

type SynchronizeShelvesResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
