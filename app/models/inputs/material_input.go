package inputs

import (
	"notezy-backend/app/models/schemas/enums"

	"github.com/google/uuid"
)

type CreateMaterialInput struct {
	RootShelfId   uuid.UUID                 `json:"rootShelfId" gorm:"column:root_shelf_id;"`
	ParentShelfId uuid.UUID                 `json:"parentShelfId" gorm:"column:parent_shelf_id;"`
	Name          string                    `json:"name" gorm:"column:name;"`
	Type          enums.MaterialType        `json:"type" gorm:"column:type;"`
	ContentURL    string                    `json:"contentURL" gorm:"column:content_url;"`
	ContentType   enums.MaterialContentType `json:"contentType" gorm:"column:content_type;"`
}

type UpdateMaterialInput struct {
	ParentShelfId *uuid.UUID                 `json:"parentShelfId" gorm:"column:parent_shelf_id"`
	Name          *string                    `json:"name" gorm:"column:name;"`
	Type          *enums.MaterialType        `json:"type" gorm:"column:type;"`
	ContentURL    *[]string                  `json:"contentURL" gorm:"column:content_url;"`
	ContentType   *enums.MaterialContentType `json:"contentType" gorm:"column:content_type;"`
}

type PartialUpdateMaterialInput = PartialUpdateInput[UpdateMaterialInput]
