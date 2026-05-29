package inputs

import (
	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
)

type CreateRoutineTagInput struct {
	Id    *uuid.UUID           `json:"id" gorm:"column:id;"`
	Name  string               `json:"name" gorm:"column:name;"`
	Color string               `json:"color" gorm:"column:color;"`
	Icon  *enums.SupportedIcon `json:"icon" gorm:"column:icon;"`
}

type BulkCreateRoutineTagInput struct {
	Id        *uuid.UUID           `json:"id" gorm:"column:id;"`
	StationId uuid.UUID            `json:"stationId" gorm:"column:station_id;"`
	Name      string               `json:"name" gorm:"column:name;"`
	Color     string               `json:"color" gorm:"column:color;"`
	Icon      *enums.SupportedIcon `json:"icon" gorm:"column:icon;"`
}

type UpdateRoutineTagInput struct {
	StationId *uuid.UUID           `json:"stationId" gorm:"column:station_id;"`
	Name      *string              `json:"name" gorm:"column:name;"`
	Color     *string              `json:"color" gorm:"column:color;"`
	Icon      *enums.SupportedIcon `json:"icon" gorm:"column:icon;"`
}

type PartialUpdateRoutineTagInput = PartialUpdateInput[UpdateRoutineTagInput]

type BulkUpdateRoutineTagInput struct {
	Id                 uuid.UUID                                 `json:"id" gorm:"column:id;"`
	PartialUpdateInput PartialUpdateInput[UpdateRoutineTagInput] `json:"partialUpdateInput"`
}
