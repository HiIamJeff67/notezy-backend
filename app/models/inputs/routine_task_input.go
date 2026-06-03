package inputs

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
)

type CreateRoutineTaskInput struct {
	Title       string                   `json:"title" gorm:"column:title;"`
	Purpose     enums.RoutineTaskPurpose `json:"purpose" gorm:"column:purpose;"`
	Payload     datatypes.JSON           `json:"payload" gorm:"column:payload;"`
	Priority    int32                    `json:"priority" gorm:"column:priority;"`
	MaxAttempts int32                    `json:"maxAttempts" gorm:"column:max_attempts;"`
}

type BulkCreateRoutineTaskInput struct {
	StationId   uuid.UUID                `json:"stationId" gorm:"column:station_id;"`
	Title       string                   `json:"title" gorm:"column:title;"`
	Purpose     enums.RoutineTaskPurpose `json:"purpose" gorm:"column:purpose;"`
	Payload     datatypes.JSON           `json:"payload" gorm:"column:payload;"`
	Priority    int32                    `json:"priority" gorm:"column:priority;"`
	MaxAttempts int32                    `json:"maxAttempts" gorm:"column:max_attempts;"`
}

type UpdateRoutineTaskInput struct {
	StationId   *uuid.UUID                `json:"stationId" gorm:"column:station_id;"`
	Title       *string                   `json:"title" gorm:"column:title;"`
	Purpose     *enums.RoutineTaskPurpose `json:"purpose" gorm:"column:purpose;"`
	Payload     *datatypes.JSON           `json:"payload" gorm:"column:payload;"`
	Priority    *int32                    `json:"priority" gorm:"column:priority;"`
	MaxAttempts *int32                    `json:"maxAttempts" gorm:"column:max_attempts;"`
}

type PartialUpdateRoutineTaskInput = PartialUpdateInput[UpdateRoutineTaskInput]

type BulkUpdateRoutineTaskInput struct {
	Id                 uuid.UUID                                  `json:"id" gorm:"column:id;"`
	PartialUpdateInput PartialUpdateInput[UpdateRoutineTaskInput] `json:"partialUpdateInput"`
}
