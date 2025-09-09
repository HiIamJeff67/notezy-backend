package dtos

import (
	"io"
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
			MaterialId  uuid.UUID `json:"materialId" validate:"required"`
			RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
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
			RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
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
			UserId       uuid.UUID // extracted from the access token of AuthMiddleware()
			UserPublicId uuid.UUID // extracted from the AuthMiddleware()
		},
		struct {
			RootShelfId   uuid.UUID `json:"rootShelfId" validate:"required"`
			ParentShelfId uuid.UUID `json:"parentShelfId" validate:"required"`
			Name          string    `json:"name" validate:"required,min=1,max=128"`
		},
		any,
	]
}

type SaveMyMaterialByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId       uuid.UUID // extracted from the access token of AuthMiddleware()
			UserPublicId uuid.UUID // extracted from the AuthMiddleware()
		},
		struct {
			MaterialId    uuid.UUID `json:"materialId" validate:"required"`
			RootShelfId   uuid.UUID `json:"rootShelfId" validate:"required"`
			PartialUpdate PartialUpdateDto[struct {
				// we are not allowed the user to move the material to other shelves by this API
				Name *string `json:"name" validate:"omitnil,min=1,max=128"`
			}]
			ContentFile io.Reader `json:"contentFile" validate:"omitnil"` // extracted from the context of MultipartAdapter()
			Size        *int64    `json:"size" validate:"omitnil"`
			// Note that io.Reader is an interface, it can be nil although we declare the type of io.Reader instead of *io.Reader
		},
		any,
	]
}

type MoveMyMaterialByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			MaterialId               uuid.UUID `json:"materialId" validate:"required"`
			SourceRootShelfId        uuid.UUID `json:"sourceRootShelfId" validate:"required"`
			DestinationRootShelfId   uuid.UUID `json:"destinationRootShelfId" validate:"required"`
			DestinationParentShelfId uuid.UUID `json:"parentShelfId" validate:"required"`
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
			MaterialId  uuid.UUID `json:"materialId" validate:"required"`
			RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
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
			RootShelfId uuid.UUID   `json:"rootShelfId" validate:"required"`
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
			MaterialId  uuid.UUID `json:"materialId" validate:"required"`
			RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
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
			RootShelfId uuid.UUID   `json:"rootShelfId" validate:"required"`
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
	DownloadURL   string                    `json:"downloadURL"`
	ContentType   enums.MaterialContentType `json:"contentType"`
	DeletedAt     *time.Time                `json:"deletedAt"`
	UpdatedAt     time.Time                 `json:"updatedAt"`
	CreatedAt     time.Time                 `json:"createdAt"`
}

type SearchMyMaterialsByShelfIdResDto []GetMyMaterialByIdResDto

type CreateMaterialResDto struct {
	CreatedAt time.Time `json:"createdAt"`
}

type SaveMyMaterialByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveMyMaterialByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
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
