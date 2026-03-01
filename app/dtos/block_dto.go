package dtos

import (
	"encoding/json"
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
			BlockIds []uuid.UUID `json:"blockIds" form:"blockIds" validate:"required"`
		},
	]
}

type GetMyBlocksByBlockGroupIdReqDto struct {
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

type GetMyBlocksByBlockGroupIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			BlockGroupIds []uuid.UUID `json:"blockGroupIds" form:"blockGroupIds" validate:"required"`
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
			BlockPackId uuid.UUID `json:"blockPackId" validate:"required"`
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
			ParentBlockId          *uuid.UUID             `json:"parentBlockId" validate:"required"`
			BlockGroupId           uuid.UUID              `json:"blockGroupId" validate:"required"`
			ArborizedEditableBlock ArborizedEditableBlock `json:"arborizedEditableBlock" validate:"required"`
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
				ParentBlockId          *uuid.UUID             `json:"parentBlockId" validate:"required"`
				BlockGroupId           uuid.UUID              `json:"blockGroupId" validate:"required"`
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
				ParentBlockId *uuid.UUID       `json:"parentBlockId" validate:"omitnil"`
				BlockGroupId  *uuid.UUID       `json:"blockGroupId" validate:"omitnil"`
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
					ParentBlockId *uuid.UUID       `json:"parentBlockId" validate:"omitnil"`
					BlockGroupId  *uuid.UUID       `json:"blockGroupId" validate:"omitnil"`
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
			BlockGroupId  uuid.UUID  `json:"blockGroupId" validate:"required"`
			ParentBlockId *uuid.UUID `json:"parentBlockId" validate:"omitnil"`
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
	ParentBlockId *uuid.UUID      `json:"parentBlockId"`
	BlockGroupId  uuid.UUID       `json:"blockGroupId"`
	Type          enums.BlockType `json:"type"`
	Props         datatypes.JSON  `json:"props"`
	Content       datatypes.JSON  `json:"content"`
	DeletedAt     *time.Time      `json:"deletedAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
	CreatedAt     time.Time       `json:"createdAt"`
}

type GetMyBlocksByIdsResDto = []GetMyBlockByIdResDto

type GetMyBlocksByBlockGroupIdResDto struct {
	RawArborizedEditableBlock RawArborizedEditableBlock `json:"rawArborizedEditableBlock"`
}

type GetMyBlocksByBlockGroupIdsResDto = []GetMyBlocksByBlockGroupIdResDto

type GetMyBlocksByBlockPackIdResDto = []GetMyBlockByIdResDto

type GetAllMyBlocksResDto = []GetMyBlockByIdResDto

type InsertBlockResDto struct {
	CreatedAt time.Time `json:"createdAt"`
}

type InsertBlocksResDto struct {
	IsAllSuccess                 bool  `json:"isAllSuccess"`
	FailedIndexes                []int `json:"failedIndexes"`
	SuccessIndexes               []int `json:"successIndexes"`
	SuccessBlockGroupAndBlockIds []struct {
		BlockGroupId uuid.UUID   `json:"blockGroupId"`
		BlockIds     []uuid.UUID `json:"blockIds"`
	} `json:"successBlockGroupAndBlockIds"`
	CreatedAt time.Time `json:"createdAt"`
}

type UpdateMyBlockByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateMyBlocksByIdsResDto struct {
	IsAllSuccess                 bool  `json:"isAllSuccess"`
	FailedIndexes                []int `json:"failedIndexes"`
	SuccessIndexes               []int `json:"successIndexes"`
	SuccessBlockGroupAndBlockIds []struct {
		BlockGroupId uuid.UUID   `json:"blockGroupId"`
		BlockIds     []uuid.UUID `json:"blockIds"`
	} `json:"successBlockGroupAndBlockIds"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveMyBlockByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type RestoreMyBlockByIdResDto struct {
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

type RestoreMyBlocksByIdsResDto = []RestoreMyBlockByIdResDto

type DeleteMyBlockByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type DeleteMyBlocksByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
