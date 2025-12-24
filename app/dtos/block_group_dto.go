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

type InsertBlockGroupByBlockPackIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockPackId      uuid.UUID  `json:"blockPackId" validate:"required"`
			PrevBlockGroupId *uuid.UUID `json:"prevBlockGroupId" validate:"omitempty"`
		},
		any,
	]
}

type InsertBlockGroupAndItsBlocksByBlockPackIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockPackId            uuid.UUID              `json:"blockPackId" validate:"required"`
			PrevBlockGroupId       *uuid.UUID             `json:"prevBlockGroupId" validate:"omitempty"`
			ArborizedEditableBlock ArborizedEditableBlock `json:"arborizedEditableBlock" validate:"required"`
		},
		any,
	]
}

type InsertBlockGroupsAndTheirBlocksByBlockPackIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockPackId        uuid.UUID `json:"blockPackId" validate:"required"`
			BlockGroupContents []struct {
				PrevBlockGroupId       *uuid.UUID             `json:"prevBlockGroupId" validate:"omitempty"`
				ArborizedEditableBlock ArborizedEditableBlock `json:"arborizedEditableBlock" validate:"required"`
			} `json:"blockGroupContents" validate:"required"`
		},
		any,
	]
}

type InsertSequentialBlockGroupsAndTheirBlocksByBlockPackIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockPackId             uuid.UUID                `json:"blockPackId" validate:"required"`
			PrevBlockGroupId        *uuid.UUID               `json:"prevBlockGroupId" validate:"omitempty"`
			ArborizedEditableBlocks []ArborizedEditableBlock `json:"arborizedEditableBlocks" validate:"required"`
		},
		any,
	]
}

type MoveMyBlockGroupsByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockPackId              uuid.UUID    `json:"blockPackId" validate:"required"`
			MovableBlockGroupIds     []uuid.UUID  `json:"movableBlockGroupIds" validate:"required"`
			MovablePrevBlockGroupIds []*uuid.UUID `json:"movablePrevBlockGroupIds" validate:"required"`
			DestinationBlockGroupId  *uuid.UUID   `json:"destinationBlockGroupId" validate:"omitnil"` // expect result to place next to the destination block group
		},
		any,
	]
}

type RestoreMyBlockGroupByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockGroupId uuid.UUID `json:"blockGroupId" validate:"required"`
		},
		any,
	]
}

type RestoreMyBlockGroupsByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockGroupIds []uuid.UUID `json:"blockGroupIds" validate:"required"`
		},
		any,
	]
}

type DeleteMyBlockGroupByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockGroupId uuid.UUID `json:"blockGroupId" validate:"required"`
		},
		any,
	]
}

type DeleteMyBlockGroupsByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockGroupIds []uuid.UUID `json:"blockGroupIds" validate:"required"`
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
	Id                        uuid.UUID                 `json:"id"`
	BlockPackId               uuid.UUID                 `json:"blockPackId"`
	PrevBlockGroupId          *uuid.UUID                `json:"prevBlockGroupId"`
	SyncBlockGroupId          *uuid.UUID                `json:"syncBlockGroupId"`
	MegaByteSize              float64                   `json:"megaByteSize"`
	DeletedAt                 *time.Time                `json:"deletedAt"`
	UpdatedAt                 time.Time                 `json:"updatedAt"`
	CreatedAt                 time.Time                 `json:"createdAt"`
	RawArborizedEditableBlock RawArborizedEditableBlock `json:"rawArborizedEditableBlock"`
}

type GetMyBlockGroupsAndTheirBlocksByBlockPackIdResDto = []GetMyBlockGroupAndItsBlocksByIdResDto

type GetMyBlockGroupsByPrevBlockGroupIdResDto = []GetMyBlockGroupByIdResDto

type GetAllMyBlockGroupsByBlockPackIdResDto = []GetMyBlockGroupByIdResDto

type InsertBlockGroupByBlockPackIdResDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type InsertBlockGroupAndItsBlocksByBlockPackIdResDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type InsertBlockGroupsAndTheirBlocksByBlockPackIdResDto struct {
	IsAllSuccess                 bool  `json:"isAllSuccess"`
	FailedIndexes                []int `json:"failedIndexes"`
	SuccessIndexes               []int `json:"successIndexes"`
	SuccessBlockGroupAndBlockIds []struct {
		BlockGroupId uuid.UUID
		BlockIds     []uuid.UUID
	} `json:"successBlockGroupAndBlockIds"`
	CreatedAt time.Time `json:"createdAt"`
}

type InsertSequentialBlockGroupsAndTheirBlocksByBlockPackIdResDto struct {
	IsAllSuccess                 bool  `json:"isAllSuccess"`
	FailedIndexes                []int `json:"failedIndexes"`
	SuccessIndexes               []int `json:"successIndexes"`
	SuccessBlockGroupAndBlockIds []struct {
		BlockGroupId uuid.UUID
		BlockIds     []uuid.UUID
	} `json:"successBlockGroupAndBlockIds"`
	CreatedAt time.Time `json:"createdAt"`
}

type MoveMyBlockGroupsByIdsResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type RestoreMyBlockGroupByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type RestoreMyBlockGroupsByIdsResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type DeleteMyBlockGroupByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type DeleteMyBlockGroupsByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
