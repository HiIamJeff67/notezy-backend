package dtos

import (
	"time"

	"github.com/google/uuid"
)

/* ============================== Request DTO ============================== */

type GetMyShelfByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			OwnerId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			ShelfId uuid.UUID `json:"shelfId" validate:"required"`
		},
		any,
	]
}

type SearchRecentShelvesReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			OwnerId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			SimpleSearchDto
		},
	]
}

type CreateShelfReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			OwnerId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			Name string `json:"name" validate:"required,max=128"`
		},
		any,
	]
}

type SynchronizeShelvesReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			OwnerId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			ShelfIds       []uuid.UUID `json:"shelfIds" validate:"required"`
			PartialUpdates []PartialUpdateDto[struct {
				Name             *string    `json:"name" validate:"omitnil,max=128,isshelfname"`
				EncodedStructure *[]byte    `json:"encodedStructure" validate:"omitnil"`
				LastAnalyzedAt   *time.Time `json:"lastAnalyzedAt" validate:"omitnil,notfuture"`
			}] `json:"partialUpdates" validate:"required"`
		},
		any,
	]
}

type RestoreMyShelfReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			OwnerId uuid.UUID
		},
		struct {
			ShelfId uuid.UUID `json:"shelfId" validate:"required"`
		},
		any,
	]
}

type RestoreMyShelvesReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			OwnerId uuid.UUID
		},
		struct {
			ShelfIds []uuid.UUID `json:"shelfIds" validate:"required,min=1,max=32"`
		},
		any,
	]
}

type DeleteMyShelfReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			OwnerId uuid.UUID
		},
		struct {
			ShelfId uuid.UUID `json:"shelfId" validate:"required"`
		},
		any,
	]
}

type DeleteMyShelvesReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			OwnerId uuid.UUID
		},
		struct {
			ShelfIds []uuid.UUID `json:"shelfIds" validate:"required,min=1,max=32"`
		},
		any,
	]
}

/* ============================== Response DTO ============================== */

type GetMyShelfByIdResDto struct {
	Id                       uuid.UUID  `json:"id"`
	Name                     string     `json:"name"`
	EncodedStructure         []byte     `json:"encodedStructure"`
	EncodedStructureByteSize int64      `json:"encodedStructureByteSize"`
	TotalShelfNodes          int32      `json:"totalShelfNodes"`
	TotalMaterials           int32      `json:"totalMaterials"`
	MaxWidth                 int32      `json:"maxWidth"`
	MaxDepth                 int32      `json:"maxDepth"`
	LastAnalyzedAt           time.Time  `json:"lastAnalyzedAt"`
	DeletedAt                *time.Time `json:"deletedAt"`
	UpdatedAt                time.Time  `json:"updatedAt"`
	CreatedAt                time.Time  `json:"createdAt"`
}

type SearchRecentShelvesResDto []struct {
	Id                       uuid.UUID  `json:"id"`
	Name                     string     `json:"name"`
	EncodedStructure         []byte     `json:"encodedStructure"`
	EncodedStructureByteSize int64      `json:"encodedStructureByteSize"`
	TotalShelfNodes          int32      `json:"totalShelfNodes"`
	TotalMaterials           int32      `json:"totalMaterials"`
	MaxWidth                 int32      `json:"maxWidth"`
	MaxDepth                 int32      `json:"maxDepth"`
	LastAnalyzedAt           time.Time  `json:"lastAnalyzedAt"`
	DeletedAt                *time.Time `json:"deletedAt"`
	UpdatedAt                time.Time  `json:"updatedAt"`
	CreatedAt                time.Time  `json:"createdAt"`
}

type CreateShelfResDto struct {
	Id               uuid.UUID `json:"id"`
	EncodedStructure []byte    `json:"encodedStructure"`
	LastAnalyzedAt   time.Time `json:"lastAnalyzedAt"`
	CreatedAt        time.Time `json:"createdAt"`
}

type SynchronizeShelvesResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type RestoreMyShelfResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type RestoreMyShelvesResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type DeleteMyShelfResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type DeleteMyShelvesResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
