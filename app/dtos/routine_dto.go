package dtos

import (
	"time"

	"github.com/google/uuid"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
)

/* ============================== Request DTO ============================== */

type GetMyRoutineByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		any,
		struct {
			RoutineId uuid.UUID `form:"routineId" validate:"required"`
			IsDeleted *bool     `form:"isDeleted" validate:"omitnil"`
		},
	]
}

type GetMyRoutinesByStationIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		any,
		struct {
			StationId  uuid.UUID `form:"stationId" validate:"required"`
			AreDeleted *bool     `form:"areDeleted" validate:"omitnil"`
		},
	]
}

type GetAllMyRoutinesByTimeRangeReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		any,
		struct {
			From       time.Time   `form:"from" validate:"required"`
			To         time.Time   `form:"to" validate:"required"`
			StationIds []uuid.UUID `form:"stationIds" validate:"required,min=1,max=1024"`
			AreDeleted *bool       `form:"areDeleted" validate:"omitnil"`
		},
	]
}

type CreateRoutineByStationIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			Id               *uuid.UUID           `json:"id" validate:"omitnil"`
			StationId        uuid.UUID            `json:"stationId" validate:"required"`
			Title            string               `json:"title" validate:"required,min=1,max=128"`
			Description      string               `json:"description" validate:"max=1024"`
			Status           *enums.RoutineStatus `json:"status" validate:"omitnil,isroutinestatus"`
			IsPinned         *bool                `json:"isPinned" validate:"omitnil"`
			ScheduledStartAt *time.Time           `json:"scheduledStartAt" validate:"omitnil"`
			ScheduledEndAt   *time.Time           `json:"scheduledEndAt" validate:"omitnil"`
			Period           *enums.RoutinePeriod `json:"period" validate:"omitnil,isroutineperiod"`
			Timezone         *string              `json:"timezone" validate:"omitnil,max=64,istimezone"`
		},
		any,
	]
}

type CreateRoutinesByStationIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			CreatedRoutines []struct {
				Id               *uuid.UUID           `json:"id" validate:"omitnil"`
				StationId        uuid.UUID            `json:"stationId" validate:"required"`
				Title            string               `json:"title" validate:"required,min=1,max=128"`
				Description      string               `json:"description" validate:"max=1024"`
				Status           *enums.RoutineStatus `json:"status" validate:"omitnil,isroutinestatus"`
				IsPinned         *bool                `json:"isPinned" validate:"omitnil"`
				ScheduledStartAt *time.Time           `json:"scheduledStartAt" validate:"omitnil"`
				ScheduledEndAt   *time.Time           `json:"scheduledEndAt" validate:"omitnil"`
				Period           *enums.RoutinePeriod `json:"period" validate:"omitnil,isroutineperiod"`
				Timezone         *string              `json:"timezone" validate:"omitnil,max=64,istimezone"`
			} `json:"createdRoutines" validate:"required,min=1,max=1024"`
		},
		any,
	]
}

type UpdateMyRoutineByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			RoutineId uuid.UUID `json:"routineId" validate:"required"`
			PartialUpdateDto[struct {
				StationId        *uuid.UUID           `json:"stationId" validate:"omitnil"`
				Title            *string              `json:"title" validate:"omitnil,min=1,max=128"`
				Description      *string              `json:"description" validate:"omitnil,max=1024"`
				Status           *enums.RoutineStatus `json:"status" validate:"omitnil,isroutinestatus"`
				IsPinned         *bool                `json:"isPinned" validate:"omitnil"`
				ScheduledStartAt *time.Time           `json:"scheduledStartAt" validate:"omitnil"`
				ScheduledEndAt   *time.Time           `json:"scheduledEndAt" validate:"omitnil"`
				Period           *enums.RoutinePeriod `json:"period" validate:"omitnil,isroutineperiod"`
				Timezone         *string              `json:"timezone" validate:"omitnil,max=64,istimezone"`
			}]
		},
		any,
	]
}

type UpdateMyRoutinesByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			UpdatedRoutines []struct {
				RoutineId uuid.UUID `json:"routineId" validate:"required"`
				PartialUpdateDto[struct {
					StationId        *uuid.UUID           `json:"stationId" validate:"omitnil"`
					Title            *string              `json:"title" validate:"omitnil,min=1,max=128"`
					Description      *string              `json:"description" validate:"omitnil,max=1024"`
					Status           *enums.RoutineStatus `json:"status" validate:"omitnil,isroutinestatus"`
					IsPinned         *bool                `json:"isPinned" validate:"omitnil"`
					ScheduledStartAt *time.Time           `json:"scheduledStartAt" validate:"omitnil"`
					ScheduledEndAt   *time.Time           `json:"scheduledEndAt" validate:"omitnil"`
					Period           *enums.RoutinePeriod `json:"period" validate:"omitnil,isroutineperiod"`
					Timezone         *string              `json:"timezone" validate:"omitnil,max=64,istimezone"`
				}]
			} `json:"updatedRoutines" validate:"required,min=1,max=1024"`
		},
		any,
	]
}

type LinkRoutineTagByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			RoutineId    uuid.UUID `json:"routineId" validate:"required"`
			RoutineTagId uuid.UUID `json:"routineTagId" validate:"required"`
			IsUnlink     bool      `json:"isUnlink"`
		},
		any,
	]
}

type LinkRoutineTagsByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			LinkedRoutinesAndTags []struct {
				RoutineId    uuid.UUID `json:"routineId" validate:"required"`
				RoutineTagId uuid.UUID `json:"routineTagId" validate:"required"`
			} `json:"linkedRoutinesAndTags" validate:"required,min=1,max=1024"`
			IsUnlink bool `json:"isUnlink"`
		},
		any,
	]
}

type LinkRoutineTaskByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			RoutineId     uuid.UUID `json:"routineId" validate:"required"`
			RoutineTaskId uuid.UUID `json:"routineTaskId" validate:"required"`
			IsUnlink      bool      `json:"isUnlink"`
		},
		any,
	]
}

type LinkRoutineTasksByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			LinkedRoutinesAndTasks []struct {
				RoutineId     uuid.UUID `json:"routineId" validate:"required"`
				RoutineTaskId uuid.UUID `json:"routineTaskId" validate:"required"`
			} `json:"linkedRoutinesAndTasks" validate:"required,min=1,max=1024"`
			IsUnlink bool `json:"isUnlink"`
		},
		any,
	]
}

type LinkRoutineItemByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			RoutineId uuid.UUID      `json:"routineId" validate:"required"`
			ItemId    uuid.UUID      `json:"itemId" validate:"required"`
			ItemType  enums.ItemType `json:"itemType" validate:"required,isitemtype"`
			IsUnlink  bool           `json:"isUnlink"`
		},
		any,
	]
}

type LinkRoutineItemsByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			LinkedRoutinesAndItems []struct {
				RoutineId uuid.UUID      `json:"routineId" validate:"required"`
				ItemId    uuid.UUID      `json:"itemId" validate:"required"`
				ItemType  enums.ItemType `json:"itemType" validate:"required,isitemtype"`
			} `json:"linkedRoutinesAndItems" validate:"required,min=1,max=1024"`
			IsUnlink bool `json:"isUnlink"`
		},
		any,
	]
}

type RestoreMyRoutineByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			RoutineId uuid.UUID `json:"routineId" validate:"required"`
		},
		any,
	]
}

type RestoreMyRoutinesByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			RoutineIds []uuid.UUID `json:"routineIds" validate:"required,min=1,max=1024"`
		},
		any,
	]
}

type DeleteMyRoutineByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			RoutineId uuid.UUID `json:"routineId" validate:"required"`
		},
		any,
	]
}

type DeleteMyRoutinesByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			RoutineIds []uuid.UUID `json:"routineIds" validate:"required,min=1,max=1024"`
		},
		any,
	]
}

type HardDeleteMyRoutineByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			RoutineId uuid.UUID `json:"routineId" validate:"required"`
		},
		any,
	]
}

type HardDeleteMyRoutinesByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			RoutineIds []uuid.UUID `json:"routineIds" validate:"required,min=1,max=1024"`
		},
		any,
	]
}

type VisualizeMyRoutineStatusCountReqDto struct {
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

type VisualizeMyRoutinePeriodCountReqDto struct {
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

type VisualizeMyRoutineScheduledStartAtCountReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		any,
		struct {
			Permission          enums.AccessControlPermission `json:"permission" validate:"isaccesscontrolpermission,required"`
			TimeHourUnit        int                           `json:"timeHourUnit" validate:"required,min=1"`
			QueryRangeStartedAt time.Time                     `json:"queryRangeStartedAt" validate:"required"`
			QueryRangeEndedAt   time.Time                     `json:"queryRangeEndedAt" validate:"required"`
		},
	]
}

type VisualizeMyRoutineScheduledEndAtCountReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		any,
		struct {
			Permission          enums.AccessControlPermission `json:"permission" validate:"isaccesscontrolpermission,required"`
			TimeHourUnit        int                           `json:"timeHourUnit" validate:"required,min=1"`
			QueryRangeStartedAt time.Time                     `json:"queryRangeStartedAt" validate:"required"`
			QueryRangeEndedAt   time.Time                     `json:"queryRangeEndedAt" validate:"required"`
		},
	]
}

/* ============================== Response DTO ============================== */

type GetMyRoutineByIdResDto struct {
	Id               uuid.UUID            `json:"id"`
	StationId        uuid.UUID            `json:"stationId"`
	Title            string               `json:"title"`
	Description      string               `json:"description"`
	Status           enums.RoutineStatus  `json:"status"`
	IsPinned         bool                 `json:"isPinned"`
	ScheduledStartAt time.Time            `json:"scheduledStartAt"`
	ScheduledEndAt   time.Time            `json:"scheduledEndAt"`
	Period           *enums.RoutinePeriod `json:"period"`
	Timezone         string               `json:"timezone"`
	DeletedAt        *time.Time           `json:"deletedAt"`
	UpdatedAt        time.Time            `json:"updatedAt"`
	CreatedAt        time.Time            `json:"createdAt"`
	TagIds           []uuid.UUID          `json:"tagIds"`
	TaskIds          []uuid.UUID          `json:"taskIds"`
	ItemIds          []uuid.UUID          `json:"itemIds"`
}

type GetMyRoutinesByStationIdResDto []struct {
	Id               uuid.UUID            `json:"id"`
	StationId        uuid.UUID            `json:"stationId"`
	Title            string               `json:"title"`
	Status           enums.RoutineStatus  `json:"status"`
	IsPinned         bool                 `json:"isPinned"`
	ScheduledStartAt time.Time            `json:"scheduledStartAt"`
	ScheduledEndAt   time.Time            `json:"scheduledEndAt"`
	Period           *enums.RoutinePeriod `json:"period"`
	Timezone         string               `json:"timezone"`
	DeletedAt        *time.Time           `json:"deletedAt"`
	UpdatedAt        time.Time            `json:"updatedAt"`
	CreatedAt        time.Time            `json:"createdAt"`
	TagIds           []uuid.UUID          `json:"tagIds"`
	TaskIds          []uuid.UUID          `json:"taskIds"`
	ItemIds          []uuid.UUID          `json:"itemIds"`
}

type GetAllMyRoutinesByTimeRangeResDto = []struct {
	Id               uuid.UUID            `json:"id"`
	StationId        uuid.UUID            `json:"stationId"`
	Title            string               `json:"title"`
	Status           enums.RoutineStatus  `json:"status"`
	IsPinned         bool                 `json:"isPinned"`
	ScheduledStartAt time.Time            `json:"scheduledStartAt"`
	ScheduledEndAt   time.Time            `json:"scheduledEndAt"`
	Period           *enums.RoutinePeriod `json:"period"`
	Timezone         string               `json:"timezone"`
	DeletedAt        *time.Time           `json:"deletedAt"`
	UpdatedAt        time.Time            `json:"updatedAt"`
	CreatedAt        time.Time            `json:"createdAt"`
	TagIds           []uuid.UUID          `json:"tagIds"`
	TaskIds          []uuid.UUID          `json:"taskIds"`
	ItemIds          []uuid.UUID          `json:"itemIds"`
}

type CreateRoutineByStationIdResDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateRoutinesByStationIdsResDto struct {
	Ids       []uuid.UUID `json:"ids"`
	CreatedAt time.Time   `json:"createdAt"`
}

type UpdateMyRoutineByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateMyRoutinesByIdsResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type LinkRoutineTagByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type LinkRoutineTagsByIdsResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type LinkRoutineTaskByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type LinkRoutineTasksByIdsResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type LinkRoutineItemByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type LinkRoutineItemsByIdsResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type RestoreMyRoutineByIdResDto struct {
	Id               uuid.UUID            `json:"id"`
	StationId        uuid.UUID            `json:"stationId"`
	Title            string               `json:"title"`
	Description      string               `json:"description"`
	Status           enums.RoutineStatus  `json:"status"`
	IsPinned         bool                 `json:"isPinned"`
	ScheduledStartAt time.Time            `json:"scheduledStartAt"`
	ScheduledEndAt   time.Time            `json:"scheduledEndAt"`
	Period           *enums.RoutinePeriod `json:"period"`
	Timezone         string               `json:"timezone"`
	DeletedAt        *time.Time           `json:"deletedAt"`
	UpdatedAt        time.Time            `json:"updatedAt"`
	CreatedAt        time.Time            `json:"createdAt"`
}

type RestoreMyRoutinesByIdsResDto = []RestoreMyRoutineByIdResDto

type DeleteMyRoutineByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type DeleteMyRoutinesByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type HardDeleteMyRoutineByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type HardDeleteMyRoutinesByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type VisualizeMyRoutineStatusCountResDto = TwoDimensionalData[int64]

type VisualizeMyRoutinePeriodCountResDto = TwoDimensionalData[int64]

type VisualizeMyRoutineScheduledStartAtCountResDto = TwoDimensionalData[int64]

type VisualizeMyRoutineScheduledEndAtCountResDto = TwoDimensionalData[int64]
