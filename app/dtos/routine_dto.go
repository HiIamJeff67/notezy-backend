package dtos

import (
	"time"

	"github.com/google/uuid"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
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
			RoutineId   uuid.UUID      `form:"routineId" validate:"required"`
			OnlyDeleted *types.Ternary `form:"onlyDeleted" validate:"omitnil,min=0,max=2"`
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
