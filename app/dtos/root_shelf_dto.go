package dtos

import (
	"time"

	"github.com/google/uuid"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
)

/* ============================== Request DTO ============================== */

type GetMyRootShelfByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			RootShelfId uuid.UUID `form:"rootShelfId" validate:"required"`
			IsDeleted   *bool     `form:"isDeleted" validate:"omitnil"`
		},
	]
}

type SearchRecentRootShelvesReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			SimpleSearchDto
		},
	]
}

type CreateRootShelfReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			Id   *uuid.UUID `json:"id" validate:"omitnil"`
			Name string     `json:"name" validate:"required,min=1,max=128,isshelfname"`
		},
		any,
	]
}

type CreateRootShelvesReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			CreatedRootShelves []struct {
				Id   *uuid.UUID `json:"id" validate:"omitnil"`
				Name string     `json:"name" validate:"required,min=1,max=128,isshelfname"`
			} `json:"insertedRootShelves" validate:"required"`
		},
		any,
	]
}

type UpdateMyRootShelfByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
			PartialUpdateDto[struct {
				Name *string `json:"name" validate:"omitnil,min=1,max=128,isshelfname"`
			}]
		},
		any,
	]
}

type UpdateMyRootShelvesByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			UpdatedRootShelves []struct {
				RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
				PartialUpdateDto[struct {
					Name *string `json:"name" validate:"omitnil,min=1,max=128,isshelfname"`
				}]
			} `json:"updatedRootShelves" validate:"required"`
		},
		any,
	]
}

type UpsertMyRootShelfPermissionReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			Permission enums.AccessControlPermission `json:"permission" validate:"required,isaccesscontrolpermission"`
		},
		struct {
			RootShelfId  uuid.UUID `uri:"rootShelfId" validate:"required"`
			UserPublicId uuid.UUID `uri:"userPublicId" validate:"required"`
		},
	]
}

type UpsertMyRootShelfPermissionsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			Permissions []struct {
				UserPublicId uuid.UUID                     `json:"userPublicId" validate:"required"`
				Permission   enums.AccessControlPermission `json:"permission" validate:"required,isaccesscontrolpermission"`
			} `json:"permissions" validate:"required,min=1,max=1024,dive"`
		},
		struct {
			RootShelfId uuid.UUID `uri:"rootShelfId" validate:"required"`
		},
	]
}

type RestoreMyRootShelfByIdReqDto struct {
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
		any,
	]
}

type RestoreMyRootShelvesByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			RootShelfIds []uuid.UUID `json:"rootShelfIds" validate:"required,min=1,max=1024"`
		},
		any,
	]
}

type DeleteMyRootShelfByIdReqDto struct {
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
		any,
	]
}

type DeleteMyRootShelvesByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			RootShelfIds []uuid.UUID `json:"rootShelfIds" validate:"required,min=1,max=1024"`
		},
		any,
	]
}

type DeleteMyRootShelfPermissionReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			RootShelfId  uuid.UUID `uri:"rootShelfId" validate:"required"`
			UserPublicId uuid.UUID `uri:"userPublicId" validate:"required"`
		},
	]
}

type DeleteMyRootShelfPermissionsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			UserPublicIds []uuid.UUID `json:"userPublicIds" validate:"required,min=1,max=1024,dive,required"`
		},
		struct {
			RootShelfId uuid.UUID `uri:"rootShelfId" validate:"required"`
		},
	]
}

/* ============================== Response DTO ============================== */

type GetMyRootShelfByIdResDto struct {
	Id             uuid.UUID                     `json:"id"`
	Name           string                        `json:"name"`
	Permission     enums.AccessControlPermission `json:"permission"`
	SubShelfCount  int64                         `json:"subShelfCount"`
	ItemCount      int64                         `json:"itemCount"`
	LastAnalyzedAt time.Time                     `json:"lastAnalyzedAt"`
	DeletedAt      *time.Time                    `json:"deletedAt"`
	UpdatedAt      time.Time                     `json:"updatedAt"`
	CreatedAt      time.Time                     `json:"createdAt"`
}

type SearchRecentRootShelvesResDto = []GetMyRootShelfByIdResDto

type CreateRootShelfResDto struct {
	Id             uuid.UUID `json:"id"`
	LastAnalyzedAt time.Time `json:"lastAnalyzedAt"`
	CreatedAt      time.Time `json:"createdAt"`
}

type CreateRootShelvesResDto struct {
	Ids            []uuid.UUID `json:"ids"`
	LastAnalyzedAt time.Time   `json:"lastAnalyzedAt"`
	CreatedAt      time.Time   `json:"createdAt"`
}

type UpdateMyRootShelfByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateMyRootShelvesByIdsResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpsertMyRootShelfPermissionResDto struct {
	UserPublicId uuid.UUID                     `json:"userPublicId"`
	Permission   enums.AccessControlPermission `json:"permission"`
	UpdatedAt    time.Time                     `json:"updatedAt"`
	CreatedAt    time.Time                     `json:"createdAt"`
}

type UpsertMyRootShelfPermissionsResDto struct {
	Permissions []UpsertMyRootShelfPermissionResDto `json:"permissions"`
}

type RestoreMyRootShelfByIdResDto struct {
	Id             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	SubShelfCount  int64      `json:"subShelfCount"`
	ItemCount      int64      `json:"itemCount"`
	LastAnalyzedAt time.Time  `json:"lastAnalyzedAt"`
	DeletedAt      *time.Time `json:"deletedAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	CreatedAt      time.Time  `json:"createdAt"`
}

type RestoreMyRootShelvesByIdsResDto = []RestoreMyRootShelfByIdResDto

type DeleteMyRootShelfByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type DeleteMyRootShelvesByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
