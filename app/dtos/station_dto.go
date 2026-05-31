package dtos

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

/* ============================== Request DTO ============================== */

type GetMyStationByIdReqDto struct {
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

type CreateStationReqDto struct {
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

type CreateStationsReqDto struct {
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

type UpdateMyStationByIdReqDto struct {
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

type UpdateMyStationsByIdsReqDto struct {
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

type RestoreMyStationByIdReqDto struct {
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

type RestoreMyStationsByIdsReqDto struct {
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

type DeleteMyStationByIdReqDto struct {
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

type DeleteMyStationsByIdsReqDto struct {
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

type HardDeleteMyStationByIdReqDto struct {
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

type HardDeleteMyStationsByIdsReqDto struct {
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

type GetMyStationByIdResDto struct {
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

type CreateStationResDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateStationsResDto struct {
	Ids       []uuid.UUID `json:"ids"`
	CreatedAt time.Time   `json:"createdAt"`
}

type UpdateMyStationByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateMyStationsByIdsResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type RestoreMyStationByIdResDto struct {
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

type RestoreMyStationsByIdsResDto = []RestoreMyStationByIdResDto

type DeleteMyStationByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type DeleteMyStationsByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type HardDeleteMyStationByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type HardDeleteMyStationsByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
