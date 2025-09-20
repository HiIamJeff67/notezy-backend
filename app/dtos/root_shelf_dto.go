package dtos

import (
	"time"

	"github.com/google/uuid"
)

/* ============================== Request DTO ============================== */

type GetMyRootShelfByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			RootShelfId uuid.UUID `form:"rootShelfId" validate:"required"`
		},
	]
}

type SearchRecentRootShelvesReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			SimpleSearchDto
		},
	]
}

type CreateRootShelfReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			OwnerId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			Name string `json:"name" validate:"required,min=1,max=128"`
		},
		any,
	]
}

type UpdateMyRootShelfByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
			PartialUpdateDto[struct {
				Name *string `json:"name" validate:"omitnil,min=1,max=128"`
			}]
		},
		any,
	]
}

type RestoreMyRootShelfByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			OwnerId uuid.UUID
		},
		struct {
			RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
		},
		any,
	]
}

type RestoreMyRootShelvesByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			OwnerId uuid.UUID
		},
		struct {
			RootShelfIds []uuid.UUID `json:"rootShelfIds" validate:"required,min=1,max=128"`
		},
		any,
	]
}

type DeleteMyRootShelfByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			OwnerId uuid.UUID
		},
		struct {
			RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
		},
		any,
	]
}

type DeleteMyRootShelvesByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			OwnerId uuid.UUID
		},
		struct {
			RootShelfIds []uuid.UUID `json:"rootShelfIds" validate:"required,min=1,max=128"`
		},
		any,
	]
}

/* ============================== Response DTO ============================== */

type GetMyRootShelfByIdResDto struct {
	Id              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	TotalShelfNodes int32      `json:"totalShelfNodes"`
	TotalMaterials  int32      `json:"totalMaterials"`
	LastAnalyzedAt  time.Time  `json:"lastAnalyzedAt"`
	DeletedAt       *time.Time `json:"deletedAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	CreatedAt       time.Time  `json:"createdAt"`
}

type SearchRecentRootShelvesResDto []GetMyRootShelfByIdResDto

type CreateRootShelfResDto struct {
	Id             uuid.UUID `json:"id"`
	LastAnalyzedAt time.Time `json:"lastAnalyzedAt"`
	CreatedAt      time.Time `json:"createdAt"`
}

type UpdateMyRootShelfByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type RestoreMyRootShelfByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type RestoreMyRootShelvesByIdsResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type DeleteMyRootShelfByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type DeleteMyRootShelvesByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
