package dtos

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
)

/* ============================== Request DTO ============================== */

type GetMyBlockByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			BlockId uuid.UUID `form:"blockId" validate:"required"`
		},
	]
}

type GetMyBlocksByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			BlockIds []uuid.UUID `form:"blockIds" validate:"required"`
		},
	]
}

type GetMyBlocksByBlockPackIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			BlockPackId uuid.UUID `form:"blockPackId" validate:"required"`
		},
	]
}

/* ============================== Response DTO ============================== */

type GetMyBlockByIdResDto struct {
	Id            uuid.UUID       `json:"id"`
	BlockPackId   uuid.UUID       `json:"blockPackId"`
	ParentBlockId *uuid.UUID      `json:"parentBlockId"`
	PrevBlockId   *uuid.UUID      `json:"prevBlockId"`
	NextBlockId   *uuid.UUID      `json:"nextBlockId"`
	Type          enums.BlockType `json:"type"`
	Props         datatypes.JSON  `json:"props"`
	Content       datatypes.JSON  `json:"content"`
	UpdatedAt     time.Time       `json:"updatedAt"`
	CreatedAt     time.Time       `json:"createdAt"`
}

type GetMyBlocksByIdsResDto = []GetMyBlockByIdResDto

type GetMyBlocksByBlockPackIdResDto = []GetMyBlockByIdResDto
