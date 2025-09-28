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
			Name           string     `json:"name" validate:"required,min=1,max=128"`
			PrevSubShelfId *uuid.UUID `json:"prevSubShelfId" validate:"omitnil"`
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
				Name *string `json:"name" validate:"omitnil,min=1,max=128"`
			}]
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

type CreateSubShelfByRootShelfIdResDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type UpdateMySubShelfByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveMySubShelfResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveMySubShelvesResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type RestoreMySubShelfByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type RestoreMySubShelvesByIdsResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type DeleteMySubShelfByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type DeleteMySubShelvesByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
