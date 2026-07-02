package dtos

import (
	"time"

	"github.com/google/uuid"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
)

/* ============================== Request DTO ============================== */

type GetAllMyRoutineTaskRecordsByRoutineTaskIdReqDto struct {
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
			Limit         int       `form:"limit" validate:"omitempty,min=1,max=500"`
		},
	]
}

type VisualizeMyRoutineTaskRecordStatusCountReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		any,
		struct {
			Permission     enums.AccessControlPermission `json:"permission" validate:"isaccesscontrolpermission,required"`
			RoutineTaskIds []uuid.UUID                   `form:"routineTaskIds" validate:"omitempty,max=1024"`
		},
	]
}

type VisualizeMyRoutineTaskRecordPurposeCountReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID
		},
		any,
		struct {
			Permission     enums.AccessControlPermission `json:"permission" validate:"isaccesscontrolpermission,required"`
			RoutineTaskIds []uuid.UUID                   `form:"routineTaskIds" validate:"omitempty,max=1024"`
		},
	]
}

type VisualizeMyRoutineTaskRecordScheduledAtCountReqDto struct {
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
			RoutineTaskIds      []uuid.UUID                   `form:"routineTaskIds" validate:"omitempty,max=1024"`
			TimeHourUnit        int                           `json:"timeHourUnit" validate:"required,min=1"`
			QueryRangeStartedAt time.Time                     `json:"queryRangeStartedAt" validate:"required"`
			QueryRangeEndedAt   time.Time                     `json:"queryRangeEndedAt" validate:"required"`
		},
	]
}

type VisualizeMyRoutineTaskRecordActualStartedAtCountReqDto struct {
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
			RoutineTaskIds      []uuid.UUID                   `form:"routineTaskIds" validate:"omitempty,max=1024"`
			TimeHourUnit        int                           `json:"timeHourUnit" validate:"required,min=1"`
			QueryRangeStartedAt time.Time                     `json:"queryRangeStartedAt" validate:"required"`
			QueryRangeEndedAt   time.Time                     `json:"queryRangeEndedAt" validate:"required"`
		},
	]
}

type VisualizeMyRoutineTaskRecordActualEndedAtCountReqDto struct {
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
			RoutineTaskIds      []uuid.UUID                   `form:"routineTaskIds" validate:"omitempty,max=1024"`
			TimeHourUnit        int                           `json:"timeHourUnit" validate:"required,min=1"`
			QueryRangeStartedAt time.Time                     `json:"queryRangeStartedAt" validate:"required"`
			QueryRangeEndedAt   time.Time                     `json:"queryRangeEndedAt" validate:"required"`
		},
	]
}

/* ============================== Response DTO ============================== */

type GetAllMyRoutineTaskRecordsByRoutineTaskIdResDto = []struct {
	Id              uuid.UUID                         `json:"id"`
	RoutineTaskId   uuid.UUID                         `json:"routineTaskId"`
	Purpose         enums.RoutineTaskPurpose          `json:"purpose"`
	Status          enums.RoutineTaskRecordStatus     `json:"status"`
	ErrorCode       *enums.RoutineTaskRecordErrorCode `json:"errorCode"`
	ErrorReason     *string                           `json:"errorReason"`
	CostUnit        int64                             `json:"costUnit"`
	TotalAttempts   int64                             `json:"totalAttempts"`
	ScheduledAt     time.Time                         `json:"scheduledAt"`
	ActualStartedAt *time.Time                        `json:"actualStartedAt"`
	ActualEndedAt   *time.Time                        `json:"actualEndedAt"`
	UpdatedAt       time.Time                         `json:"updatedAt"`
	CreatedAt       time.Time                         `json:"createdAt"`
}

type VisualizeMyRoutineTaskRecordStatusCountResDto = TwoDimensionalData[int64]

type VisualizeMyRoutineTaskRecordPurposeCountResDto = TwoDimensionalData[int64]

type VisualizeMyRoutineTaskRecordScheduledAtCountResDto = TwoDimensionalData[int64]

type VisualizeMyRoutineTaskRecordActualStartedAtCountResDto = TwoDimensionalData[int64]

type VisualizeMyRoutineTaskRecordActualEndedAtCountResDto = TwoDimensionalData[int64]
