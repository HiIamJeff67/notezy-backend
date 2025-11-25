package dtos

import (
	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
	"time"

	"github.com/google/uuid"
)

/* ============================== Request DTO ============================== */

type GetMyBlockPackByIdReqDto struct {
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

type GetMyBlockPackAndItsParentByIdReqDto struct {
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

type GetAllMyBlockPacksByParentSubShelfIdReqDto struct {
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

type GetAllMyBlockPacksByRootShelfIdReqDto struct {
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

type CreateBlockPackReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			ParentSubShelfId    uuid.UUID                     `json:"parentSubShelfId" validate:"required"`
			Name                string                        `json:"name" validate:"required,min=1,max=128"`
			Icon                *enums.SupportedBlockPackIcon `json:"icon" validate:"omitnil,issupportedblockpackicon"`
			HeaderBackgroundURL *string                       `json:"headerBackgroundURL" validate:"omitnil"`
		},
		any,
	]
}

type UpdateMyBlockPackByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockPackId uuid.UUID `json:"blockPackId" validate:"required"`
			PartialUpdateDto[struct {
				Name                *string                       `json:"name" validate:"omitnil,min=1,max=128"`
				Icon                *enums.SupportedBlockPackIcon `json:"icon" validate:"omitnil,issupportedblockpackicon"`
				HeaderBackgroundURL *string                       `json:"headerBackgroundURL" validate:"omitnil"`
			}]
		},
		any,
	]
}

type MoveMyBlockPackByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockPackId                 uuid.UUID `json:"blockPackId" validate:"required"`
			DestinationParentSubShelfId uuid.UUID `json:"destinationParentSubShelfId" validate:"required"`
		},
		any,
	]
}

type MoveMyBlockPacksByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockPackIds                []uuid.UUID `json:"blockPackIds" validate:"required,min=1,max=100"`
			DestinationParentSubShelfId uuid.UUID   `json:"destinationParentSubShelfId" validate:"required"`
		},
		any,
	]
}

type RestoreMyBlockPackByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockPackId uuid.UUID `json:"blockPackId" validate:"required"`
		},
		any,
	]
}

type RestoreMyBlockPacksByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockPackIds []uuid.UUID `json:"blockPackIds" validate:"required,min=1,max=128"`
		},
		any,
	]
}

type DeleteMyBlockPackByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockPackId uuid.UUID `json:"blockPackId" validate:"required"`
		},
		any,
	]
}

type DeleteMyBlockPacksByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			BlockPackIds []uuid.UUID `json:"blockPackIds" validate:"required,min=1,max=128"`
		},
		any,
	]
}

/* ============================== Response DTO ============================== */

type GetMyBlockPackByIdResDto struct {
	Id                  uuid.UUID                     `json:"id" gorm:"column:id;"`
	ParentSubShelfId    uuid.UUID                     `json:"parentSubShelfId" gorm:"column:parent_sub_shelf_id;"`
	Name                string                        `json:"name" gorm:"column:name;"`
	Icon                *enums.SupportedBlockPackIcon `json:"icon" gorm:"column:icon;"`
	HeaderBackgroundURL *string                       `json:"headerBackgroundURL" gorm:"column:header_background_url;"`
	BlockCount          int32                         `json:"blockCount" gorm:"column:block_count;"`
	DeletedAt           *time.Time                    `json:"deletedAt" gorm:"column:deleted_at;"`
	UpdatedAt           time.Time                     `json:"updatedAt" gorm:"column:updated_at;"`
	CreatedAt           time.Time                     `json:"createdAt" gorm:"column:created_at;"`
}

type GetMyBlockPackAndItsParentByIdResDto struct {
	Id                           uuid.UUID                     `json:"id" gorm:"column:id;"`
	Name                         string                        `json:"name" gorm:"column:name;"`
	Icon                         *enums.SupportedBlockPackIcon `json:"icon" gorm:"column:icon;"`
	HeaderBackgroundURL          *string                       `json:"headerBackgroundURL" gorm:"column:header_background_url;"`
	BlockCount                   int32                         `json:"blockCount" gorm:"column:block_count;"`
	DeletedAt                    *time.Time                    `json:"deletedAt" gorm:"column:deleted_at;"`
	UpdatedAt                    time.Time                     `json:"updatedAt" gorm:"column:updated_at;"`
	CreatedAt                    time.Time                     `json:"createdAt" gorm:"column:created_at;"`
	RootShelfId                  uuid.UUID                     `json:"rootShelfId" gorm:"column:root_shelf_id;"`
	ParentSubShelfId             uuid.UUID                     `json:"parentSubShelfId" gorm:"column:parent_sub_shelf_id;"`
	ParentSubShelfPrevSubShelfId *uuid.UUID                    `json:"parentSubShelfPrevSubShelfId" gorm:"column:parent_sub_shelf_prev_sub_shelf_id;"`
	ParentSubShelfName           string                        `json:"parentSubShelfName" gorm:"column:parent_sub_shelf_name;"`
	ParentSubShelfPath           types.UUIDArray               `json:"parentSubShelfPath" gorm:"column:parent_sub_shelf_path;"`
	ParentSubShelfDeletedAt      time.Time                     `json:"parentSubShelfDeletedAt" gorm:"column:parent_sub_shelf_deleted_at;"`
	ParentSubShelfUpdatedAt      time.Time                     `json:"parentSubShelfUpdatedAt" gorm:"column:parent_sub_shelf_updated_at;"`
	ParentSubShelfCreatedAt      time.Time                     `json:"parentSubShelfCreatedAt" gorm:"column:parent_sub_shelf_created_at;"`
}

type GetAllMyBlockPacksByParentSubShelfIdResDto = []GetMyBlockPackByIdResDto

type GetAllMyBlockPacksByRootShelfIdResDto = []GetMyBlockPackByIdResDto

type CreateBlockPackResDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type UpdateMyBlockPackByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveMyBlockPackByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveMyBlockPacksByIdsResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type RestoreMyBlockPackByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type RestoreMyBlockPacksByIdsResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type DeleteMyBlockPackByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type DeleteMyBlockPacksByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
