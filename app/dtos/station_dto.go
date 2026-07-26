package dtos

import (
	"time"

	"github.com/google/uuid"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
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
			StationId uuid.UUID `form:"stationId" validate:"required"`
			IsDeleted *bool     `form:"isDeleted" validate:"omitnil"`
		},
	]
}

type GetAllMyStationsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		any,
		struct {
			AreDeleted *bool `form:"areDeleted" validate:"omitnil"`
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
			} `json:"createdStations" validate:"required,min=1,max=200"`
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

type GetMyStationPermissionReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		any,
		struct {
			StationId    uuid.UUID `uri:"stationId" validate:"required"`
			UserPublicId uuid.UUID `uri:"userPublicId" validate:"required"`
		},
	]
}

type UpsertMyStationPermissionReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			Permission enums.AccessControlPermission `json:"permission" validate:"required,isaccesscontrolpermission"`
		},
		struct {
			StationId    uuid.UUID `uri:"stationId" validate:"required"`
			UserPublicId uuid.UUID `uri:"userPublicId" validate:"required"`
		},
	]
}

type CreateMyStationPermissionReqDto = UpsertMyStationPermissionReqDto
type UpdateMyStationPermissionReqDto = UpsertMyStationPermissionReqDto

type UpsertMyStationPermissionsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			Permissions []struct {
				UserPublicId uuid.UUID                     `json:"userPublicId" validate:"required"`
				Permission   enums.AccessControlPermission `json:"permission" validate:"required,isaccesscontrolpermission"`
			} `json:"permissions" validate:"required,min=1,max=1024,dive"`
		},
		struct {
			StationId uuid.UUID `uri:"stationId" validate:"required"`
		},
	]
}

type TransferMyStationOwnershipReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			TargetUserPublicId uuid.UUID `json:"targetUserPublicId" validate:"required"`
		},
		struct {
			StationId uuid.UUID `uri:"stationId" validate:"required"`
		},
	]
}

type DeleteMyStationPermissionReqDto = GetMyStationPermissionReqDto

type DeleteMyStationPermissionsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			UserPublicIds []uuid.UUID `json:"userPublicIds" validate:"required,min=1,max=1024,dive,required"`
		},
		struct {
			StationId uuid.UUID `uri:"stationId" validate:"required"`
		},
	]
}

type LeaveMyStationReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		any,
		struct {
			StationId uuid.UUID `uri:"stationId" validate:"required"`
		},
	]
}

type LeaveMyStationsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			Stations []struct {
				StationId uuid.UUID `json:"stationId" validate:"required"`
			} `json:"stations" validate:"required,min=1,max=1024,dive"`
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

type VisualizeMyTotalCountReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		any,
		struct {
			Permission enums.AccessControlPermission `json:"permission" validate:"isaccesscontrolpermission,required"`
		},
	]
}

/* ============================== Response DTO ============================== */

type GetMyStationByIdResDto struct {
	Id                  uuid.UUID                     `json:"id"`
	Name                string                        `json:"name"`
	Description         string                        `json:"description"`
	Icon                *enums.SupportedIcon          `json:"icon"`
	HeaderBackgroundURL *string                       `json:"headerBackgroundURL"`
	Permission          enums.AccessControlPermission `json:"permission"`
	RoutineCount        int64                         `json:"routineCount"`
	DeletedAt           *time.Time                    `json:"deletedAt"`
	UpdatedAt           time.Time                     `json:"updatedAt"`
	CreatedAt           time.Time                     `json:"createdAt"`
}

type GetAllMyStationsResDto = []struct {
	Id                  uuid.UUID                     `json:"id"`
	Name                string                        `json:"name"`
	Icon                *enums.SupportedIcon          `json:"icon"`
	HeaderBackgroundURL *string                       `json:"headerBackgroundURL"`
	Permission          enums.AccessControlPermission `json:"permission"`
	RoutineCount        int64                         `json:"routineCount"`
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

type UpsertMyStationPermissionResDto struct {
	UserPublicId uuid.UUID                     `json:"userPublicId"`
	Permission   enums.AccessControlPermission `json:"permission"`
	UpdatedAt    time.Time                     `json:"updatedAt"`
	CreatedAt    time.Time                     `json:"createdAt"`
}

type StationPermissionResDto = UpsertMyStationPermissionResDto
type GetMyStationPermissionResDto = UpsertMyStationPermissionResDto
type CreateMyStationPermissionResDto = UpsertMyStationPermissionResDto
type UpdateMyStationPermissionResDto = UpsertMyStationPermissionResDto

type UpsertMyStationPermissionsResDto struct {
	Permissions []UpsertMyStationPermissionResDto `json:"permissions"`
}

type TransferMyStationOwnershipResDto struct {
	StationId                 uuid.UUID `json:"stationId"`
	PreviousOwnerUserPublicId uuid.UUID `json:"previousOwnerUserPublicId"`
	NewOwnerUserPublicId      uuid.UUID `json:"newOwnerUserPublicId"`
	UpdatedAt                 time.Time `json:"updatedAt"`
}

type RestoreMyStationByIdResDto struct {
	Id                  uuid.UUID                     `json:"id"`
	Name                string                        `json:"name"`
	Description         string                        `json:"description"`
	Icon                *enums.SupportedIcon          `json:"icon"`
	HeaderBackgroundURL *string                       `json:"headerBackgroundURL"`
	Permission          enums.AccessControlPermission `json:"permission"`
	RoutineCount        int64                         `json:"routineCount"`
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

type VisualizeMyTotalCountResDto = TwoDimensionalData[int64]
