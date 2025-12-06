package dtos

import (
	"time"

	"github.com/google/uuid"
)

/* ============================== Request DTO ============================== */

type GetMyBlockGroupByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			BlockGroupId uuid.UUID `json:"blockGroupId" validate:"required"`
		},
	]
}

type GetMyBlockGroupAndItsBlocksByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			BlockGroupId uuid.UUID `json:"blockGroupId" validate:"required"`
		},
	]
}

type GetMyBlockGroupsByPrevBlockGroupIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			PrevBlockGroupId uuid.UUID `json:"prevBlockGroupId" validate:"required"`
		},
	]
}

type GetAllMyBlockGroupsByBlockPackIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			BlockPackId uuid.UUID `json:"blockPackId" validate:"required"`
		},
	]
}

type GetMyBlockGroupsAndTheirBlocksByBlockPackIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			BlockPackId uuid.UUID `json:"blockPackId" validate:"required"`
		},
	]
}

type CreateBlockGroupByBlockPackIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockPackId      uuid.UUID  `json:"blockPackId" validate:"required"`
			PrevBlockGroupId *uuid.UUID `json:"prevBlockGroupId" validate:"required"`
		},
		any,
	]
}

/* ============================== Response DTO ============================== */

type GetMyBlockGroupByIdResDto struct {
	Id               uuid.UUID  `json:"id"`
	BlockPackId      uuid.UUID  `json:"blockPackId"`
	PrevBlockGroupId *uuid.UUID `json:"prevBlockGroupId"`
	SyncBlockGroupId *uuid.UUID `json:"syncBlockGroupId"`
	MegaByteSize     float64    `json:"megaByteSize"`
	DeletedAt        *time.Time `json:"deletedAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	CreatedAt        time.Time  `json:"createdAt"`
}

type GetMyBlockGroupAndItsBlocksByIdResDto struct {
	Id                   uuid.UUID            `json:"id"`
	BlockPackId          uuid.UUID            `json:"blockPackId"`
	PrevBlockGroupId     *uuid.UUID           `json:"prevBlockGroupId"`
	SyncBlockGroupId     *uuid.UUID           `json:"syncBlockGroupId"`
	MegaByteSize         float64              `json:"megaByteSize"`
	DeletedAt            *time.Time           `json:"deletedAt"`
	UpdatedAt            time.Time            `json:"updatedAt"`
	CreatedAt            time.Time            `json:"createdAt"`
	EditableBlockContent EditableBlockContent `json:"editableBlockContent"`
}

type GetMyBlockGroupsByPrevBlockGroupIdResDto = []GetMyBlockGroupByIdResDto

type GetAllMyBlockGroupsByBlockPackIdResDto = []GetMyBlockGroupByIdResDto

type GetMyBlockGroupsAndTheirBlocksByBlockPackIdResDto = []GetMyBlockGroupAndItsBlocksByIdResDto

type CreateBlockGroupByBlockPackIdResDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}
