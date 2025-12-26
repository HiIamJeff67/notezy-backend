package dtos

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	enums "notezy-backend/app/models/schemas/enums"
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
			BlockId uuid.UUID `json:"blockId" validate:"required"`
		},
	]
}

type GetAllMyBlocksReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		any,
	]
}

/* ============================== Response DTO ============================== */

type GetMyBlockByIdResDto struct {
	Id            uuid.UUID       `json:"id"`
	ParentBlockId *uuid.UUID      `json:"parentBlockId"`
	BlockGroupId  uuid.UUID       `json:"blockGroupId"`
	Type          enums.BlockType `json:"type"`
	Props         datatypes.JSON  `json:"props"`
	Content       datatypes.JSON  `json:"content"`
	DeletedAt     *time.Time      `json:"deletedAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
	CreatedAt     time.Time       `json:"createdAt"`
}

type GetAllMyBlocksResDto = []GetMyBlockByIdResDto
