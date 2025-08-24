package dtos

import (
	"time"

	"github.com/google/uuid"
)

/* ============================== Request DTO ============================== */

type GetRecentShelvesReqDto struct {
	NotezyRequest[
		any,
		struct {
			OwnerId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			GetManyDto
		},
	]
}

type CreateShelfReqDto struct {
	NotezyRequest[
		any,
		struct {
			OwnerId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			Name string `json:"name" validate:"required,max=128"`
		},
	]
}

type SynchronizeShelvesReqDto struct {
	NotezyRequest[
		any,
		struct {
			OwnerId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			ShelfIds       []uuid.UUID `json:"shelfIds" validate:"required"`
			PartialUpdates []PartialUpdateDto[struct {
				Name             *string `json:"name" validate:"omitnil,max=128,isshelfname"`
				EncodedStructure *[]byte `json:"encodedStructure" validate:"omitnil"`
			}] `json:"partialUpdates" validate:"required"`
		},
	]
}

/* ============================== Response DTO ============================== */

type GetRecentShelvesResDto struct {
	Name                     string    `json:"name"`
	EncodedStructure         []byte    `json:"encodedStructure"`
	EncodedStructureByteSize int64     `json:"encodedStructureByteSize"`
	TotalShelfNodes          int32     `json:"totalShelfNodes"`
	TotalMaterials           int32     `json:"totalMaterials"`
	MaxWidth                 int32     `json:"maxWidth"`
	MaxDepth                 int32     `json:"maxDepth"`
	UpdatedAt                time.Time `json:"updatedAt"`
	CreatedAt                time.Time `json:"createdAt"`
}

type CreateShelfResDto struct {
	EncodedStructure []byte    `json:"encodedStructure"`
	CreatedAt        time.Time `json:"createdAt"`
}

type SynchronizeShelvesResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
