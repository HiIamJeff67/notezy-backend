package dtos

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

/* ============================== Request DTO ============================== */

type GetOneStationByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		any,
		struct {
			StationId   uuid.UUID      `form:"stationId" validate:"required"`
			OnlyDeleted *types.Ternary `form:"onlyDeleted" validate:"omitnil,min=0,max=2"`
		},
	]
}

type CreateOneStationByOwnerIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			Id                  *uuid.UUID           `json:"id" validate:"omitnil"`
			Name                string               `json:"name" validate:"required,min=1,max=128"`
			Description         string               `json:"description" validate:"max=1024"`
			Icon                *enums.SupportedIcon `json:"icon" validate:"omitnil,issupportedicon"`
			HeaderBackgroundURL *string              `json:"headerBackgroundURL" validate:"omitnil,isimageurl"`
		},
		any,
	]
}

type CreateManyStationsByOwnerIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			CreatedStations []struct {
				Id                  *uuid.UUID           `json:"id" validate:"omitnil"`
				Name                string               `json:"name" validate:"required,min=1,max=128"`
				Description         string               `json:"description" validate:"max=1024"`
				Icon                *enums.SupportedIcon `json:"icon" validate:"omitnil,issupportedicon"`
				HeaderBackgroundURL *string              `json:"headerBackgroundURL" validate:"omitnil,isimageurl"`
			} `json:"createdStations" validate:"required,min=1,max=1024"`
		},
		any,
	]
}

type UpdateOneStationByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			StationId uuid.UUID `json:"stationId" validate:"required"`
			PartialUpdateDto[struct {
				Name                *string              `json:"name" validate:"omitnil,min=1,max=128"`
				Description         *string              `json:"description" validate:"omitnil,max=1024"`
				Icon                *enums.SupportedIcon `json:"icon" validate:"omitnil,issupportedicon"`
				HeaderBackgroundURL *string              `json:"headerBackgroundURL" validate:"omitnil,isimageurl"`
			}]
		},
		any,
	]
}

type BulkUpdateManyStationsByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			UpdatedStations []struct {
				StationId uuid.UUID `json:"stationId" validate:"required"`
				PartialUpdateDto[struct {
					Name                *string              `json:"name" validate:"omitnil,min=1,max=128"`
					Description         *string              `json:"description" validate:"omitnil,max=1024"`
					Icon                *enums.SupportedIcon `json:"icon" validate:"omitnil,issupportedicon"`
					HeaderBackgroundURL *string              `json:"headerBackgroundURL" validate:"omitnil,isimageurl"`
				}]
			} `json:"updatedStations" validate:"required,min=1,max=1024"`
		},
		any,
	]
}

type RestoreSoftDeletedOneStationByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			StationId uuid.UUID `json:"stationId" validate:"required"`
		},
		any,
	]
}

type RestoreSoftDeletedManyStationsByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			StationIds []uuid.UUID `json:"stationIds" validate:"required,min=1,max=1024"`
		},
		any,
	]
}

type SoftDeleteOneStationByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			StationId uuid.UUID `json:"stationId" validate:"required"`
		},
		any,
	]
}

type SoftDeleteManyStationsByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			StationIds []uuid.UUID `json:"stationIds" validate:"required,min=1,max=1024"`
		},
		any,
	]
}

type HardDeleteOneStationByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			StationId uuid.UUID `json:"stationId" validate:"required"`
		},
		any,
	]
}

type HardDeleteManyStationsByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			StationIds []uuid.UUID `json:"stationIds" validate:"required,min=1,max=1024"`
		},
		any,
	]
}

/* ============================== Response DTO ============================== */

type GetOneStationByIdResDto struct {
	Id                  uuid.UUID                     `json:"id"`
	OwnerId             uuid.UUID                     `json:"ownerId"`
	Name                string                        `json:"name"`
	Description         string                        `json:"description"`
	Icon                *enums.SupportedIcon          `json:"icon"`
	HeaderBackgroundURL *string                       `json:"headerBackgroundURL"`
	Permission          enums.AccessControlPermission `json:"permission"`
	RoutineCount        int32                         `json:"routineCount"`
	DeletedAt           *time.Time                    `json:"deletedAt"`
	UpdatedAt           time.Time                     `json:"updatedAt"`
	CreatedAt           time.Time                     `json:"createdAt"`
}

type CreateOneStationByOwnerIdResDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateManyStationsByOwnerIdResDto struct {
	Ids       []uuid.UUID `json:"ids"`
	CreatedAt time.Time   `json:"createdAt"`
}

type UpdateOneStationByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type BulkUpdateManyStationsByIdsResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type RestoreSoftDeletedOneStationByIdResDto struct {
	Id                  uuid.UUID                     `json:"id"`
	OwnerId             uuid.UUID                     `json:"ownerId"`
	Name                string                        `json:"name"`
	Description         string                        `json:"description"`
	Icon                *enums.SupportedIcon          `json:"icon"`
	HeaderBackgroundURL *string                       `json:"headerBackgroundURL"`
	Permission          enums.AccessControlPermission `json:"permission"`
	RoutineCount        int32                         `json:"routineCount"`
	DeletedAt           *time.Time                    `json:"deletedAt"`
	UpdatedAt           time.Time                     `json:"updatedAt"`
	CreatedAt           time.Time                     `json:"createdAt"`
}

type RestoreSoftDeletedManyStationsByIdsResDto = []RestoreSoftDeletedOneStationByIdResDto

type SoftDeleteOneStationByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type SoftDeleteManyStationsByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type HardDeleteOneStationByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type HardDeleteManyStationsByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
