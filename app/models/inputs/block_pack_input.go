package inputs

import (
	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
)

type CreateBlockPackInput struct {
	Name                string                        `json:"name" gorm:"column:name;"`
	Icon                *enums.SupportedBlockPackIcon `json:"icon" gorm:"column:icon;"`
	HeaderBackgroundURL *string                       `json:"headerBackgroundURL" gorm:"header_background_url;"`
}

type UpdateBlockPackInput struct {
	ParentSubShelfId    *uuid.UUID                    `json:"parentSubShelfId" gorm:"column:parent_sub_shelf_id;"`
	Name                *string                       `json:"name" gorm:"column:name;"`
	Icon                *enums.SupportedBlockPackIcon `json:"icon" gorm:"column:icon;"`
	HeaderBackgroundURL *string                       `json:"headerBackgroundURL" gorm:"header_background_url;"`
}

type PartialUpdateBlockPackInput = PartialUpdateInput[UpdateBlockPackInput]
