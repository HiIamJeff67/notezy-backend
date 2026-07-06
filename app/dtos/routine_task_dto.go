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
			IsDeleted     *bool     `form:"isDeleted" validate:"omitnil"`
		},
	]
}

type GetAllMyRoutineTasksByRoutineIdsReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		any,
		struct {
			RoutineIds []uuid.UUID `form:"routineIds" validate:"required,min=1,max=1024"`
			AreDeleted *bool       `form:"areDeleted" validate:"omitnil"`
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
		struct {
			AreDeleted *bool `form:"areDeleted" validate:"omitnil"`
		},
	]
}

type CreateRoutineTaskByRoutineIdReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		struct {
			RoutineId       uuid.UUID                `json:"routineId" validate:"required"`
			Title           string                   `json:"title" validate:"required,min=1,max=128"`
			Purpose         enums.RoutineTaskPurpose `json:"purpose" validate:"required,isroutinetaskpurpose"`
			Payload         datatypes.JSON           `json:"payload" validate:"omitempty,max=16777216"`
			Priority        int32                    `json:"priority" validate:"omitempty,min=0,max=100"`
			MaxAttempts     int32                    `json:"maxAttempts" validate:"omitempty,min=1,max=20"`
			Period          *enums.RoutinePeriod     `json:"period" validate:"omitnil,isroutineperiod"`
			NextScheduledAt time.Time                `json:"nextScheduledAt" validate:"required"`
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
				RoutineId       *uuid.UUID                `json:"routineId" validate:"omitnil"`
				Title           *string                   `json:"title" validate:"omitnil,min=1,max=128"`
				Purpose         *enums.RoutineTaskPurpose `json:"purpose" validate:"omitnil,isroutinetaskpurpose"`
				Payload         *datatypes.JSON           `json:"payload" validate:"omitnil,max=16777216"`
				Priority        *int32                    `json:"priority" validate:"omitnil,min=0,max=100"`
				MaxAttempts     *int32                    `json:"maxAttempts" validate:"omitnil,min=1,max=20"`
				Period          *enums.RoutinePeriod      `json:"period" validate:"omitnil,isroutineperiod"`
				NextScheduledAt *time.Time                `json:"nextScheduledAt" validate:"omitnil"`
			}]
		},
		any,
	]
}

type PauseMyRoutineTaskByIdReqDto struct {
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

type ResumeMyRoutineTaskByIdReqDto struct {
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

type VisualizeMyRoutineTaskStatusCountReqDto struct {
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

type VisualizeMyRoutineTaskPurposeCountReqDto struct {
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

type VisualizeMyRoutineTaskScheduledAtCountReqDto struct {
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

type VisualizeMyRoutineTaskActualStartedAtCountReqDto struct {
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

type VisualizeMyRoutineTaskActualEndedAtCountReqDto struct {
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

type GetMyRoutineTaskByIdResDto struct {
	Id              uuid.UUID                `json:"id"`
	RoutineId       uuid.UUID                `json:"routineId"`
	Title           string                   `json:"title"`
	Purpose         enums.RoutineTaskPurpose `json:"purpose"`
	Payload         datatypes.JSON           `json:"payload"`
	CostUnit        int64                    `json:"costUnit"`
	Priority        int32                    `json:"priority"`
	Status          enums.RoutineTaskStatus  `json:"status"`
	Attempts        int32                    `json:"attempts"`
	MaxAttempts     int32                    `json:"maxAttempts"`
	Period          *enums.RoutinePeriod     `json:"period"`
	NextScheduledAt time.Time                `json:"nextScheduledAt"`
	ScheduledAt     time.Time                `json:"scheduledAt"`
	ActualStartedAt *time.Time               `json:"actualStartedAt"`
	ActualEndedAt   *time.Time               `json:"actualEndedAt"`
	UpdatedAt       time.Time                `json:"updatedAt"`
	CreatedAt       time.Time                `json:"createdAt"`
}

type GetAllMyRoutineTasksByRoutineIdsResDto = []struct {
	Id              uuid.UUID                `json:"id"`
	RoutineId       uuid.UUID                `json:"routineId"`
	Title           string                   `json:"title"`
	Purpose         enums.RoutineTaskPurpose `json:"purpose"`
	CostUnit        int64                    `json:"costUnit"`
	Priority        int32                    `json:"priority"`
	Status          enums.RoutineTaskStatus  `json:"status"`
	Attempts        int32                    `json:"attempts"`
	MaxAttempts     int32                    `json:"maxAttempts"`
	Period          *enums.RoutinePeriod     `json:"period"`
	NextScheduledAt time.Time                `json:"nextScheduledAt"`
	ScheduledAt     time.Time                `json:"scheduledAt"`
	ActualStartedAt *time.Time               `json:"actualStartedAt"`
	ActualEndedAt   *time.Time               `json:"actualEndedAt"`
	UpdatedAt       time.Time                `json:"updatedAt"`
	CreatedAt       time.Time                `json:"createdAt"`
}

type GetAllMyRoutineTasksResDto = []GetMyRoutineTaskByIdResDto

type CreateRoutineTaskByRoutineIdResDto struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type UpdateMyRoutineTaskByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type PauseMyRoutineTaskByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type ResumeMyRoutineTaskByIdResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type HardDeleteMyRoutineTaskByIdResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type HardDeleteMyRoutineTasksByIdsResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type VisualizeMyRoutineTaskStatusCountResDto = TwoDimensionalData[int64]

type VisualizeMyRoutineTaskPurposeCountResDto = TwoDimensionalData[int64]

type VisualizeMyRoutineTaskScheduledAtCountResDto = TwoDimensionalData[int64]

type VisualizeMyRoutineTaskActualStartedAtCountResDto = TwoDimensionalData[int64]

type VisualizeMyRoutineTaskActualEndedAtCountResDto = TwoDimensionalData[int64]
