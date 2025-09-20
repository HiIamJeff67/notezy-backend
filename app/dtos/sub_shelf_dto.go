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

type GetAllSubShelvesByRootShelfIdReqDto struct {
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
			SourceSubShelfId      uuid.UUID `json:"sourceSubShelfId" validate:"required"`
			DestinationSubShelfId uuid.UUID `json:"destinationSubShelfId" validate:"required"`
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
			SourceSubShelfIds     []uuid.UUID `json:"sourceSubShelfIds" validate:"required"`
			DestinationSubShelfId uuid.UUID   `json:"destinationSubShelfId" validate:"required"`
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

type GetAllSubShelvesByRootShelfIdResDto = []GetMySubShelfByIdResDto

type CreateSubShelfByRootShelfIdResDto struct {
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
