package dtos

import (
	"time"

	"github.com/google/uuid"

	types "notezy-backend/shared/types"
)

/* ============================== Request DTO ============================== */

type GetMySubShelfByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			SubShelfId uuid.UUID `json:"subShelfId" validate:"required"`
		},
	]
}

type GetMySubShelvesByPrevSubShelfIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			PrevSubShelfId uuid.UUID `json:"prevSubShelfId" validate:"required"`
		},
	]
}

type GetAllMySubShelvesByRootShelfIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
		},
	]
}

type GetMySubShelvesAndItemsByPrevSubShelfIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			PrevSubShelfId uuid.UUID `json:"prevSubShelfId" validate:"required"`
		},
	]
}

type CreateSubShelfByRootShelfIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			RootShelfId    uuid.UUID  `json:"rootShelfId" validate:"required"`
			PrevSubShelfId *uuid.UUID `json:"prevSubShelfId" validate:"omitnil"`
			Name           string     `json:"name" validate:"required,min=1,max=128,isshelfname"`
		},
		any,
	]
}

type CreateSubShelvesByRootShelfIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			CreatedSubShelves []struct {
				RootShelfId    uuid.UUID  `json:"rootShelfId" validate:"required"`
				PrevSubShelfId *uuid.UUID `json:"prevSubShelfId" validate:"omitnil"`
				Name           string     `json:"name" validate:"required,min=1,max=128,isshelfname"`
			} `json:"createdSubShelves" validate:"required"`
		},
		any,
	]
}

type UpdateMySubShelfByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			SubShelfId uuid.UUID `json:"subShelfId" validate:"required"`
			PartialUpdateDto[struct {
				Name *string `json:"name" validate:"omitnil,min=1,max=128,isshelfname"`
			}]
		},
		any,
	]
}

type UpdateMySubShelvesByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			UpdatedSubShelves []struct {
				SubShelfId uuid.UUID `json:"subShelfId" validate:"required"`
				PartialUpdateDto[struct {
					Name *string `json:"name" validate:"omitnil,min=1,max=128,isshelfname"`
				}]
			} `json:"updatedSubShelves" validate:"required"`
		},
		any,
	]
}

type MoveMySubShelfReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			SourceRootShelfId      uuid.UUID  `json:"sourceRootShelfId" validate:"required"`
			SourceSubShelfId       uuid.UUID  `json:"sourceSubShelfId" validate:"required"`
			DestinationRootShelfId uuid.UUID  `json:"destinationRootShelfId" validate:"required"`
			DestinationSubShelfId  *uuid.UUID `json:"destinationSubShelfId" validate:"omitnil"`
		},
		any,
	]
}

type MoveMySubShelvesReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			SourceRootShelfId      uuid.UUID   `json:"sourceRootShelfId" validate:"required"`
			SourceSubShelfIds      []uuid.UUID `json:"sourceSubShelfIds" validate:"required"`
			DestinationRootShelfId uuid.UUID   `json:"destinationRootShelfId" validate:"required"`
			DestinationSubShelfId  *uuid.UUID  `json:"destinationSubShelfId" validate:"omitnil"`
		},
		any,
	]
}

type BatchMoveMySubShelvesReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			MovedSubShelves []struct {
				SourceRootShelfId      uuid.UUID   `json:"sourceRootShelfId" validate:"required"`
				SourceSubShelfIds      []uuid.UUID `json:"sourceSubShelfIds" validate:"required"`
				DestinationRootShelfId uuid.UUID   `json:"destinationRootShelfId" validate:"required"`
				DestinationSubShelfId  *uuid.UUID  `json:"destinationSubShelfId" validate:"omitnil"`
			} `json:"moveSubShelves" validate:"required"`
		},
		any,
	]
}

type RestoreMySubShelfByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			SubShelfId uuid.UUID `json:"subShelfId" validate:"required"`
		},
		any,
	]
}

type RestoreMySubShelvesByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			SubShelfIds []uuid.UUID `json:"subShelfIds" validate:"required,min=1,max=128"`
		},
		any,
	]
}

type DeleteMySubShelfByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			SubShelfId uuid.UUID `json:"subShelfId" validate:"required"`
		},
		any,
	]
}

type DeleteMySubShelvesByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			SubShelfIds []uuid.UUID `json:"subShelfIds" validate:"required,min=1,max=128"`
		},
		any,
	]
}

/* ============================== Response DTO ============================== */

type GetMySubShelfByIdResDto struct {
	Id             uuid.UUID       `json:"id"`
	Name           string          `json:"name"`
	RootShelfId    uuid.UUID       `json:"rootShelfId"`
	PrevSubShelfId *uuid.UUID      `json:"prevSubShelfId"`
	Path           types.UUIDArray `json:"path"`
	DeletedAt      *time.Time      `json:"deletedAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	CreatedAt      time.Time       `json:"createdAt"`
}

type GetMySubShelvesByPrevSubShelfIdResDto = []GetMySubShelfByIdResDto

type GetAllMySubShelvesByRootShelfIdResDto = []GetMySubShelfByIdResDto

type GetMySubShelvesAndItemsByPrevSubShelfIdResDto struct {
	SubShelves []GetMySubShelfByIdResDto  `json:"subShelves"`
	Materials  []GetMyMaterialByIdResDto  `json:"materials"`
	BlockPacks []GetMyBlockPackByIdResDto `json:"blockPacks"`
}

type CreateSubShelfByRootShelfIdResDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateSubShelvesByRootShelfIdsResDto struct {
	Ids       []uuid.UUID `json:"id"`
	CreatedAt time.Time   `json:"createdAt"`
}

type UpdateMySubShelfByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateMySubShelvesByIdsResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveMySubShelfResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveMySubShelvesResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type BatchMoveMySubShelvesResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type RestoreMySubShelfByIdResDto struct {
	Id             uuid.UUID       `json:"id"`
	Name           string          `json:"name"`
	RootShelfId    uuid.UUID       `json:"rootShelfId"`
	PrevSubShelfId *uuid.UUID      `json:"prevSubShelfId"`
	Path           types.UUIDArray `json:"path"`
	DeletedAt      *time.Time      `json:"deletedAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	CreatedAt      time.Time       `json:"createdAt"`
}

type RestoreMySubShelvesByIdsResDto = []RestoreMySubShelfByIdResDto

type DeleteMySubShelfByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type DeleteMySubShelvesByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
