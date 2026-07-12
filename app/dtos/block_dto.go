package dtos

import (
	"encoding/json"
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

type InsertBlockReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockPackId            uuid.UUID              `json:"blockPackId" validate:"required"`
			ParentBlockId          *uuid.UUID             `json:"parentBlockId" validate:"omitnil"`
			PrevBlockId            *uuid.UUID             `json:"prevBlockId" validate:"omitnil"`
			ArborizedEditableBlock ArborizedEditableBlock `json:"arborizedEditableBlock" validate:"required"`
		},
		any,
	]
}

type AppendBlockReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockPackId            uuid.UUID              `json:"blockPackId" validate:"required"`
			ArborizedEditableBlock ArborizedEditableBlock `json:"arborizedEditableBlock" validate:"required"`
		},
		any,
	]
}

type AppendBlocksReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			AppendedBlocks []struct {
				BlockPackId            uuid.UUID              `json:"blockPackId" validate:"required"`
				ArborizedEditableBlock ArborizedEditableBlock `json:"arborizedEditableBlock" validate:"required"`
			} `json:"appendedBlocks" validate:"required,min=1"`
		},
		any,
	]
}

type InsertBlocksReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			InsertedBlocks []struct {
				BlockPackId            uuid.UUID              `json:"blockPackId" validate:"required"`
				ParentBlockId          *uuid.UUID             `json:"parentBlockId" validate:"omitnil"`
				PrevBlockId            *uuid.UUID             `json:"prevBlockId" validate:"omitnil"`
				ArborizedEditableBlock ArborizedEditableBlock `json:"arborizedEditableBlock" validate:"required"`
			} `json:"insertedBlocks" validate:"required"`
		},
		any,
	]
}

type UpdateMyBlockByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockId uuid.UUID `json:"blockId" validate:"required"`
			PartialUpdateDto[struct {
				BlockPackId   *uuid.UUID       `json:"blockPackId" validate:"omitnil"`
				ParentBlockId *uuid.UUID       `json:"parentBlockId" validate:"omitnil"`
				PrevBlockId   *uuid.UUID       `json:"prevBlockId" validate:"omitnil"`
				NextBlockId   *uuid.UUID       `json:"nextBlockId" validate:"omitnil"`
				Type          *enums.BlockType `json:"type" validate:"omitnil,isblocktype"`
				Props         *json.RawMessage `json:"props"`
				Content       *json.RawMessage `json:"content"`
			}]
		},
		any,
	]
}

type UpdateMyBlocksByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			UpdatedBlocks []struct {
				BlockId uuid.UUID `json:"blockId" validate:"required"`
				PartialUpdateDto[struct {
					BlockPackId   *uuid.UUID       `json:"blockPackId" validate:"omitnil"`
					ParentBlockId *uuid.UUID       `json:"parentBlockId" validate:"omitnil"`
					PrevBlockId   *uuid.UUID       `json:"prevBlockId" validate:"omitnil"`
					NextBlockId   *uuid.UUID       `json:"nextBlockId" validate:"omitnil"`
					Type          *enums.BlockType `json:"type" validate:"omitnil,isblocktype"`
					Props         *json.RawMessage `json:"props"`
					Content       *json.RawMessage `json:"content"`
				}]
			} `json:"updatedBlocks" validated:"required"`
		},
		any,
	]
}

type MoveMyBlockByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			Id            uuid.UUID  `json:"id" validate:"required"`
			BlockPackId   uuid.UUID  `json:"blockPackId" validate:"required"`
			ParentBlockId *uuid.UUID `json:"parentBlockId" validate:"omitnil"`
			PrevBlockId   *uuid.UUID `json:"prevBlockId" validate:"omitnil"`
		},
		any,
	]
}

type RestoreMyBlockByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockId uuid.UUID `json:"blockId" validate:"required"`
		},
		any,
	]
}

type RestoreMyBlocksByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockIds []uuid.UUID `json:"blockIds" validate:"required"`
		},
		any,
	]
}

type DeleteMyBlockByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockId uuid.UUID `json:"blockId" validate:"required"`
		},
		any,
	]
}

type DeleteMyBlocksByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockIds []uuid.UUID `json:"blockIds" validate:"required"`
		},
		any,
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

type GetAllMyBlocksResDto = []GetMyBlockByIdResDto

type InsertBlockResDto struct {
	CreatedAt time.Time `json:"createdAt"`
}

type AppendBlockResDto struct {
	BlockPackId uuid.UUID   `json:"blockPackId"`
	BlockIds    []uuid.UUID `json:"blockIds"`
	CreatedAt   time.Time   `json:"createdAt"`
}

type AppendBlocksResDto struct {
	IsAllSuccess                bool  `json:"isAllSuccess"`
	FailedIndexes               []int `json:"failedIndexes"`
	SuccessIndexes              []int `json:"successIndexes"`
	SuccessBlockPackAppendItems []struct {
		BlockPackId uuid.UUID   `json:"blockPackId"`
		BlockIds    []uuid.UUID `json:"blockIds"`
	} `json:"successBlockPackAppendItems"`
	CreatedAt time.Time `json:"createdAt"`
}

type InsertBlocksResDto struct {
	IsAllSuccess                bool  `json:"isAllSuccess"`
	FailedIndexes               []int `json:"failedIndexes"`
	SuccessIndexes              []int `json:"successIndexes"`
	SuccessBlockPackAndBlockIds []struct {
		BlockPackId uuid.UUID   `json:"blockPackId"`
		BlockIds    []uuid.UUID `json:"blockIds"`
	} `json:"successBlockPackAndBlockIds"`
	CreatedAt time.Time `json:"createdAt"`
}

type UpdateMyBlockByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateMyBlocksByIdsResDto struct {
	IsAllSuccess                bool  `json:"isAllSuccess"`
	FailedIndexes               []int `json:"failedIndexes"`
	SuccessIndexes              []int `json:"successIndexes"`
	SuccessBlockPackAndBlockIds []struct {
		BlockPackId uuid.UUID   `json:"blockPackId"`
		BlockIds    []uuid.UUID `json:"blockIds"`
	} `json:"successBlockPackAndBlockIds"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveMyBlockByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type RestoreMyBlockByIdResDto struct {
	Id            uuid.UUID       `json:"id"`
	BlockPackId   uuid.UUID       `json:"blockPackId"`
	ParentBlockId *uuid.UUID      `json:"parentBlockId"`
	PrevBlockId   *uuid.UUID      `json:"prevBlockId"`
	NextBlockId   *uuid.UUID      `json:"nextBlockId"`
	Type          enums.BlockType `json:"type"`
	Props         datatypes.JSON  `json:"props"`
	Content       datatypes.JSON  `json:"content"`
	DeletedAt     *time.Time      `json:"deletedAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
	CreatedAt     time.Time       `json:"createdAt"`
}

type RestoreMyBlocksByIdsResDto = []RestoreMyBlockByIdResDto

type DeleteMyBlockByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type DeleteMyBlocksByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
