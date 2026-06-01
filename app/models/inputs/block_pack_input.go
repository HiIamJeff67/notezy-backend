package inputs

import (
	"github.com/google/uuid"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
)

type CreateBlockPackInput struct {
	Id                  *uuid.UUID           `json:"id" gorm:"column:id;"`
	Name                string               `json:"name" gorm:"column:name;"`
	Icon                *enums.SupportedIcon `json:"icon" gorm:"column:icon;"`
	HeaderBackgroundURL *string              `json:"headerBackgroundURL" gorm:"header_background_url;"`
}

type BulkCreateBlockPackInput struct {
	Id                  *uuid.UUID           `json:"id" gorm:"column:id;"`
	ParentSubShelfId    uuid.UUID            `json:"parentSubShelfId" gorm:"column:parent_sub_shelf_id;"`
	Name                string               `json:"name" gorm:"column:name;"`
	Icon                *enums.SupportedIcon `json:"icon" gorm:"column:icon;"`
	HeaderBackgroundURL *string              `json:"headerBackgroundURL" gorm:"header_background_url;"`
}

type UpdateBlockPackInput struct {
	ParentSubShelfId    *uuid.UUID           `json:"parentSubShelfId" gorm:"column:parent_sub_shelf_id;"`
	Name                *string              `json:"name" gorm:"column:name;"`
	Icon                *enums.SupportedIcon `json:"icon" gorm:"column:icon;"`
	HeaderBackgroundURL *string              `json:"headerBackgroundURL" gorm:"header_background_url;"`
}

type PartialUpdateBlockPackInput = PartialUpdateInput[UpdateBlockPackInput]

type BulkUpdateBlockPackInput struct {
	Id                 uuid.UUID                                `json:"id" gorm:"column:id;"`
	PartialUpdateInput PartialUpdateInput[UpdateBlockPackInput] `json:"partialUpdateInput"`
}
