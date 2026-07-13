package dtos

import (
	"time"

	"github.com/google/uuid"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
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
			BlockPackId uuid.UUID `form:"blockPackId" validate:"required"`
			IsDeleted   *bool     `form:"isDeleted" validate:"omitnil"`
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
			BlockPackId uuid.UUID `form:"blockPackId" validate:"required"`
			IsDeleted   *bool     `form:"isDeleted" validate:"omitnil"`
		},
	]
}

type GetMyBlockPacksByParentSubShelfIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		struct {
			ParentSubShelfId uuid.UUID `form:"parentSubShelfId" validate:"required"`
			AreDeleted       *bool     `form:"areDeleted" validate:"omitnil"`
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
			RootShelfId uuid.UUID `form:"rootShelfId" validate:"required"`
			AreDeleted  *bool     `form:"areDeleted" validate:"omitnil"`
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
			Id                  *uuid.UUID           `json:"id" validate:"omitnil"`
			ParentSubShelfId    uuid.UUID            `json:"parentSubShelfId" validate:"required"`
			Name                string               `json:"name" validate:"required,min=1,max=128"`
			Icon                *enums.SupportedIcon `json:"icon" validate:"omitnil,issupportedicon"`
			HeaderBackgroundURL *string              `json:"headerBackgroundURL" validate:"omitnil"`
		},
		any,
	]
}

type CreateBlockPacksReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			CreatedBlockPacks []struct {
				Id                  *uuid.UUID           `json:"id" validate:"omitnil"`
				ParentSubShelfId    uuid.UUID            `json:"parentSubShelfId" validate:"required"`
				Name                string               `json:"name" validate:"required,min=1,max=128"`
				Icon                *enums.SupportedIcon `json:"icon" validate:"omitnil,issupportedicon"`
				HeaderBackgroundURL *string              `json:"headerBackgroundURL" validate:"omitnil"`
			} `json:"createdBlockPacks" validate:"required"`
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
				Name                *string              `json:"name" validate:"omitnil,min=1,max=128"`
				Icon                *enums.SupportedIcon `json:"icon" validate:"omitnil,issupportedicon"`
				HeaderBackgroundURL *string              `json:"headerBackgroundURL" validate:"omitnil"`
			}]
		},
		any,
	]
}

type UpdateMyBlockPacksByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			UpdatedBlockPacks []struct {
				BlockPackId uuid.UUID `json:"blockPackId" validate:"required"`
				PartialUpdateDto[struct {
					Name                *string              `json:"name" validate:"omitnil,min=1,max=128"`
					Icon                *enums.SupportedIcon `json:"icon" validate:"omitnil,issupportedicon"`
					HeaderBackgroundURL *string              `json:"headerBackgroundURL" validate:"omitnil"`
				}]
			} `json:"updatedBlockPacks" validate:"required"`
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

type MoveMyBlockPacksByParentSubShelfIdReqDto struct {
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

type MoveMyBlockPacksByParentSubShelfIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			MovedBlockPacks []struct {
				BlockPackIds                []uuid.UUID `json:"blockPackIds" validate:"required,min=1,max=100"`
				DestinationParentSubShelfId uuid.UUID   `json:"destinationParentSubShelfId" validate:"required"`
			} `json:"movedBlockPacks" validate:"required"`
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
			BlockPackIds []uuid.UUID `json:"blockPackIds" validate:"required,min=1,max=1024"`
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
			BlockPackIds []uuid.UUID `json:"blockPackIds" validate:"required,min=1,max=1024"`
		},
		any,
	]
}

/* ============================== Response DTO ============================== */

type GetMyBlockPackByIdResDto struct {
	Id                     uuid.UUID            `json:"id" gorm:"column:id;"`
	ParentSubShelfId       uuid.UUID            `json:"parentSubShelfId" gorm:"column:parent_sub_shelf_id;"`
	Name                   string               `json:"name" gorm:"column:name;"`
	Icon                   *enums.SupportedIcon `json:"icon" gorm:"column:icon;"`
	HeaderBackgroundURL    *string              `json:"headerBackgroundURL" gorm:"column:header_background_url;"`
	BlockCount             int64                `json:"blockCount" gorm:"column:block_count;"`
	LastUpdateSequence     int64                `json:"lastUpdateSequence" gorm:"column:last_update_sequence;"`
	CompactedUntilSequence int64                `json:"compactedUntilSequence" gorm:"column:compacted_until_sequence;"`
	ProjectedUntilSequence int64                `json:"projectedUntilSequence" gorm:"column:projected_until_sequence;"`
	IsProjectionCurrent    bool                 `json:"isProjectionCurrent" gorm:"column:is_projection_current;"`
	DeletedAt              *time.Time           `json:"deletedAt" gorm:"column:deleted_at;"`
	UpdatedAt              time.Time            `json:"updatedAt" gorm:"column:updated_at;"`
	CreatedAt              time.Time            `json:"createdAt" gorm:"column:created_at;"`
}

type GetMyBlockPackAndItsParentByIdResDto struct {
	Id                           uuid.UUID            `json:"id" gorm:"column:id;"`
	Name                         string               `json:"name" gorm:"column:name;"`
	Icon                         *enums.SupportedIcon `json:"icon" gorm:"column:icon;"`
	HeaderBackgroundURL          *string              `json:"headerBackgroundURL" gorm:"column:header_background_url;"`
	BlockCount                   int64                `json:"blockCount" gorm:"column:block_count;"`
	LastUpdateSequence           int64                `json:"lastUpdateSequence" gorm:"column:last_update_sequence;"`
	CompactedUntilSequence       int64                `json:"compactedUntilSequence" gorm:"column:compacted_until_sequence;"`
	ProjectedUntilSequence       int64                `json:"projectedUntilSequence" gorm:"column:projected_until_sequence;"`
	IsProjectionCurrent          bool                 `json:"isProjectionCurrent" gorm:"column:is_projection_current;"`
	DeletedAt                    *time.Time           `json:"deletedAt" gorm:"column:deleted_at;"`
	UpdatedAt                    time.Time            `json:"updatedAt" gorm:"column:updated_at;"`
	CreatedAt                    time.Time            `json:"createdAt" gorm:"column:created_at;"`
	RootShelfId                  uuid.UUID            `json:"rootShelfId" gorm:"column:root_shelf_id;"`
	ParentSubShelfId             uuid.UUID            `json:"parentSubShelfId" gorm:"column:parent_sub_shelf_id;"`
	ParentSubShelfPrevSubShelfId *uuid.UUID           `json:"parentSubShelfPrevSubShelfId" gorm:"column:parent_sub_shelf_prev_sub_shelf_id;"`
	ParentSubShelfName           string               `json:"parentSubShelfName" gorm:"column:parent_sub_shelf_name;"`
	ParentSubShelfPath           types.UUIDArray      `json:"parentSubShelfPath" gorm:"column:parent_sub_shelf_path;"`
	ParentSubShelfDeletedAt      *time.Time           `json:"parentSubShelfDeletedAt" gorm:"column:parent_sub_shelf_deleted_at;"`
	ParentSubShelfUpdatedAt      time.Time            `json:"parentSubShelfUpdatedAt" gorm:"column:parent_sub_shelf_updated_at;"`
	ParentSubShelfCreatedAt      time.Time            `json:"parentSubShelfCreatedAt" gorm:"column:parent_sub_shelf_created_at;"`
}

type GetMyBlockPacksByParentSubShelfIdResDto = []GetMyBlockPackByIdResDto

type GetAllMyBlockPacksByRootShelfIdResDto = []GetMyBlockPackByIdResDto

type CreateBlockPackResDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateBlockPacksResDto struct {
	Ids       []uuid.UUID `json:"ids"`
	CreatedAt time.Time   `json:"createdAt"`
}

type UpdateMyBlockPackByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateMyBlockPacksByIdsResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveMyBlockPackByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveMyBlockPacksByParentSubShelfIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveMyBlockPacksByParentSubShelfIdsResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type RestoreMyBlockPackByIdResDto struct {
	Id                  uuid.UUID            `json:"id" gorm:"column:id;"`
	ParentSubShelfId    uuid.UUID            `json:"parentSubShelfId" gorm:"column:parent_sub_shelf_id;"`
	Name                string               `json:"name" gorm:"column:name;"`
	Icon                *enums.SupportedIcon `json:"icon" gorm:"column:icon;"`
	HeaderBackgroundURL *string              `json:"headerBackgroundURL" gorm:"column:header_background_url;"`
	BlockCount          int64                `json:"blockCount" gorm:"column:block_count;"`
	DeletedAt           *time.Time           `json:"deletedAt"`
	UpdatedAt           time.Time            `json:"updatedAt"`
	CreatedAt           time.Time            `json:"createdAt"`
}

type RestoreMyBlockPacksByIdsResDto = []RestoreMyBlockPackByIdResDto

type DeleteMyBlockPackByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type DeleteMyBlockPacksByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
