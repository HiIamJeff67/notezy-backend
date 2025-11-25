package dtos

import (
	"io"
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
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
		any,
		struct {
			MaterialId uuid.UUID `json:"materialId" validate:"required"`
		},
	]
}

type GetMyMaterialAndItsParentByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			MaterialId uuid.UUID `json:"materialId" validate:"required"`
		},
	]
}

type GetAllMyMaterialsByParentSubShelfIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			ParentSubShelfId uuid.UUID `json:"parentSubShelfId" validate:"required"`
		},
	]
}

type GetAllMyMaterialsByRootShelfIdReqDto struct {
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

type CreateTextbookMaterialReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId       uuid.UUID // extracted from the access token of AuthMiddleware()
			UserPublicId uuid.UUID // extracted from the AuthMiddleware()
		},
		struct {
			ParentSubShelfId uuid.UUID `json:"parentSubShelfId" validate:"required"`
			Name             string    `json:"name" validate:"required,min=1,max=128"`
		},
		any,
	]
}

type CreateNotebookMaterialReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId       uuid.UUID // extracted from the access token of AuthMiddleware()
			UserPublicId uuid.UUID // extracted from the AuthMiddleware()
		},
		struct {
			ParentSubShelfId uuid.UUID `json:"parentSubShelfId" validate:"required"`
			Name             string    `json:"name" validate:"required,min=1,max=128"`
		},
		any,
	]
}

type UpdateMyMaterialByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			MaterialId uuid.UUID `json:"materialId" validate:"required"`
			PartialUpdateDto[struct {
				Name *string `json:"name" validate:"omitnil,min=1,max=128"`
			}]
			MaterialType enums.MaterialType `json:"type" validate:"required,ismaterialtype"` // for extra validation on the type
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
			Size         *int64    // extracted from the opened contentFile
		},
		struct {
			MaterialId  uuid.UUID `json:"materialId" validate:"required"`
			ContentFile io.Reader `json:"contentFile" validate:"required"` // from the context of MultipartAdapter()
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
			MaterialId                  uuid.UUID `json:"materialId" validate:"required"`
			DestinationParentSubShelfId uuid.UUID `json:"destinationParentSubShelfId" validate:"required"`
		},
		any,
	]
}

type MoveMyMaterialsByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			MaterialIds                 []uuid.UUID `json:"materialIds" validate:"required,min=1,max=100"`
			DestinationParentSubShelfId uuid.UUID   `json:"destinationParentSubShelfId" validate:"required"`
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
			MaterialId uuid.UUID `json:"materialId" validate:"required"`
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
			MaterialIds []uuid.UUID `json:"materialIds" validate:"required,min=1,max=128"`
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
			MaterialId uuid.UUID `json:"materialId" validate:"required"`
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
			MaterialIds []uuid.UUID `json:"materialIds" validate:"required,min=1,max=128"`
		},
		any,
	]
}

/* ============================== Response DTO ============================== */

type GetMyMaterialByIdResDto struct {
	Id               uuid.UUID          `json:"id"`
	ParentSubShelfId uuid.UUID          `json:"parentSubShelfId"`
	Name             string             `json:"name"`
	Type             enums.MaterialType `json:"type"`
	Size             int64              `json:"size"`
	DownloadURL      string             `json:"downloadURL"`
	DeletedAt        *time.Time         `json:"deletedAt"`
	UpdatedAt        time.Time          `json:"updatedAt"`
	CreatedAt        time.Time          `json:"createdAt"`
}

type GetMyMaterialAndItsParentByIdResDto struct {
	Id                           uuid.UUID          `json:"id"`
	Name                         string             `json:"name"`
	Type                         enums.MaterialType `json:"type"`
	Size                         int64              `json:"size"`
	DownloadURL                  string             `json:"downloadURL"`
	DeletedAt                    *time.Time         `json:"deletedAt"`
	UpdatedAt                    time.Time          `json:"updatedAt"`
	CreatedAt                    time.Time          `json:"createdAt"`
	RootShelfId                  uuid.UUID          `json:"rootShelfId"`
	ParentSubShelfId             uuid.UUID          `json:"parentSubShelfId"`
	ParentSubShelfPrevSubShelfId *uuid.UUID         `json:"parentSubShelfPrevSubShelfId"`
	ParentSubShelfName           string             `json:"parentSubShelfName"`
	ParentSubShelfPath           types.UUIDArray    `json:"parentSubShelfPath"`
	ParentSubShelfDeletedAt      time.Time          `json:"parentSubShelfDeletedAt"`
	ParentSubShelfUpdatedAt      time.Time          `json:"parentSubShelfUpdatedAt"`
	ParentSubShelfCreatedAt      time.Time          `json:"parentSubShelfCreatedAt"`
}

type GetAllMyMaterialsByParentSubShelfIdResDto []GetMyMaterialByIdResDto

type GetAllMyMaterialsByRootShelfIdResDto []GetMyMaterialByIdResDto

type CreateTextbookMaterialResDto struct {
	Id          uuid.UUID `json:"id"`
	DownloadURL string    `json:"downloadURL"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateNotebookMaterialResDto struct {
	Id          uuid.UUID `json:"id"`
	DownloadURL string    `json:"downloadURL"`
	CreatedAt   time.Time `json:"createdAt"`
}

type UpdateMyMaterialByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type SaveMyMaterialByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveMyMaterialByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveMyMaterialsByIdsResDto struct {
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
