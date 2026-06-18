package dtos

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
)

/* ============================== Request DTO ============================== */

type GetMyRoutineTaskByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		any,
		struct {
			RoutineTaskId uuid.UUID `form:"routineTaskId" validate:"required"`
		},
	]
}

type GetAllMyRoutineTasksByStationIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		any,
		struct {
			StationIds []uuid.UUID `form:"stationIds" validate:"required,min=1,max=1024"`
		},
	]
}

type GetAllMyRoutineTasksReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		any,
		any,
	]
}

type CreateRoutineTaskByStationIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			StationId   uuid.UUID                `json:"stationId" validate:"required"`
			Title       string                   `json:"title" validate:"required,min=1,max=128"`
			Purpose     enums.RoutineTaskPurpose `json:"purpose" validate:"required,isroutinetaskpurpose"`
			Payload     datatypes.JSON           `json:"payload" validate:"omitempty,max=2048"`
			Priority    int32                    `json:"priority" validate:"omitempty,min=0"`
			MaxAttempts int32                    `json:"maxAttempts" validate:"omitempty,min=1,max=20"`
		},
		any,
	]
}

type UpdateMyRoutineTaskByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			RoutineTaskId uuid.UUID `json:"routineTaskId" validate:"required"`
			PartialUpdateDto[struct {
				StationId   *uuid.UUID                `json:"stationId" validate:"omitnil"`
				Title       *string                   `json:"title" validate:"omitnil,min=1,max=128"`
				Purpose     *enums.RoutineTaskPurpose `json:"purpose" validate:"omitnil,isroutinetaskpurpose"`
				Payload     *datatypes.JSON           `json:"payload" validate:"omitnil,max=2048"`
				Priority    *int32                    `json:"priority" validate:"omitnil,min=0"`
				MaxAttempts *int32                    `json:"maxAttempts" validate:"omitnil,min=1,max=20"`
			}]
		},
		any,
	]
}

type HardDeleteMyRoutineTaskByIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			RoutineTaskId uuid.UUID `json:"routineTaskId" validate:"required"`
		},
		any,
	]
}

type HardDeleteMyRoutineTasksByIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			RoutineTaskIds []uuid.UUID `json:"routineTaskIds" validate:"required,min=1,max=1024"`
		},
		any,
	]
}

/* ============================== Response DTO ============================== */

type GetMyRoutineTaskByIdResDto struct {
	Id              uuid.UUID                `json:"id"`
	StationId       uuid.UUID                `json:"stationId"`
	Title           string                   `json:"title"`
	Purpose         enums.RoutineTaskPurpose `json:"purpose"`
	Payload         datatypes.JSON           `json:"payload"`
	Priority        int32                    `json:"priority"`
	Status          enums.RoutineTaskStatus  `json:"status"`
	Attempts        int32                    `json:"attempts"`
	MaxAttempts     int32                    `json:"maxAttempts"`
	ScheduledAt     time.Time                `json:"scheduledAt"`
	ActualStartedAt *time.Time               `json:"actualStartedAt"`
	ActualEndedAt   *time.Time               `json:"actualEndedAt"`
	UpdatedAt       time.Time                `json:"updatedAt"`
	CreatedAt       time.Time                `json:"createdAt"`
}

type GetAllMyRoutineTasksByStationIdsResDto = []struct {
	Id              uuid.UUID                `json:"id"`
	StationId       uuid.UUID                `json:"stationId"`
	Title           string                   `json:"title"`
	Purpose         enums.RoutineTaskPurpose `json:"purpose"`
	Priority        int32                    `json:"priority"`
	Status          enums.RoutineTaskStatus  `json:"status"`
	Attempts        int32                    `json:"attempts"`
	MaxAttempts     int32                    `json:"maxAttempts"`
	ScheduledAt     time.Time                `json:"scheduledAt"`
	ActualStartedAt *time.Time               `json:"actualStartedAt"`
	ActualEndedAt   *time.Time               `json:"actualEndedAt"`
	UpdatedAt       time.Time                `json:"updatedAt"`
	CreatedAt       time.Time                `json:"createdAt"`
}

type GetAllMyRoutineTasksResDto = []GetMyRoutineTaskByIdResDto

type CreateRoutineTaskByStationIdResDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type UpdateMyRoutineTaskByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type HardDeleteMyRoutineTaskByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type HardDeleteMyRoutineTasksByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
