package dtos

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
)

/* ============================== Request DTO ============================== */

type GetMyMaterialByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			MaterialId uuid.UUID `json:"materialId" validate:"required"`
		},
		any,
	]
}

type SearchMyMaterialsByShelfIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			ShelfId uuid.UUID `json:"shelfId" validate:"required"`
		},
		struct {
			SimpleSearchDto
		},
	]
}

type CreateMaterialReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			RootShelfId   uuid.UUID `json:"rootShelfId" validate:"required"`
			ParentShelfId uuid.UUID `json:"parentShelfId" validate:"required"`
			Name          string    `json:"name" validate:"required,min=1,max=128"`
		},
		any,
	]
}

type RestoreMyMaterialByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			MaterialId uuid.UUID `json:"materialId" validate:"required"`
		},
		any,
	]
}

type RestoreMyMaterialsByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			MaterialIds []uuid.UUID `json:"materialIds" validate:"required,min=1,max=32"`
		},
		any,
	]
}

type DeleteMyMaterialByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			MaterialId uuid.UUID `json:"materialId" validate:"required"`
		},
		any,
	]
}

type DeleteMyMaterialsByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			MaterialIds []uuid.UUID `json:"materialIds" validate:"required,min=1,max=32"`
		},
		any,
	]
}

/* ============================== Response DTO ============================== */

type GetMyMaterialByIdResDto struct {
	Id            uuid.UUID                 `json:"id"`
	RootShelfId   uuid.UUID                 `json:"rootShelfId"`
	ParentShelfId uuid.UUID                 `json:"parentShelfId"`
	Name          string                    `json:"name"`
	Type          enums.MaterialType        `json:"type"`
	ContentURL    string                    `json:"contentURL"`
	ContentType   enums.MaterialContentType `json:"contentType"`
	DeletedAt     *time.Time                `json:"deletedAt"`
	UpdatedAt     time.Time                 `json:"updatedAt"`
	CreatedAt     time.Time                 `json:"createdAt"`
}

type SearchMyMaterialsByShelfIdResDto []struct {
	Id            uuid.UUID                 `json:"id"`
	RootShelfId   uuid.UUID                 `json:"rootShelfId"`
	ParentShelfId uuid.UUID                 `json:"parentShelfId"`
	Name          string                    `json:"name"`
	Type          enums.MaterialContentType `json:"type"`
	ContentURL    string                    `json:"contentURL"`
	ContentType   enums.MaterialContentType `json:"contentType"`
	DeletedAt     *time.Time                `json:"deletedAt"`
	UpdatedAt     time.Time                 `json:"updatedAt"`
	CreatedAt     time.Time                 `json:"createdAt"`
}

type CreateMaterialResDto struct {
	CreatedAt time.Time `json:"createdAt"`
}

type RestoreMyMaterialByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type RestoreMyMaterialsByIdsResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type DeleteMyMaterialByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type DeleteMyMaterialsByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
