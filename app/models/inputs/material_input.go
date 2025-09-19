package inputs

import (
	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
)

type CreateMaterialInput struct {
	Id               uuid.UUID                 `json:"id" gorm:"column:id;"` // we allowed the API to generate its id, since it's faster
	ParentSubShelfId uuid.UUID                 `json:"parentSubShelfId" gorm:"column:parent_sub_shelf_id;"`
	Name             string                    `json:"name" gorm:"column:name;"`
	Type             enums.MaterialType        `json:"type" gorm:"column:type;"`
	Size             int64                     `json:"size" gorm:"column:size;"`
	ContentKey       string                    `json:"contentKey" gorm:"column:content_key;"`
	ContentType      enums.MaterialContentType `json:"contentType" gorm:"column:content_type;"`
	ParseMediaType   string                    `json:"parseMediaType" gorm:"column:parse_media_type;"`
}

type UpdateMaterialInput struct {
	ParentSubShelfId *uuid.UUID                 `json:"parentSubShelfId" gorm:"column:parent_sub_shelf_id;"`
	Name             *string                    `json:"name" gorm:"column:name;"`
	Type             *enums.MaterialType        `json:"type" gorm:"column:type;"`
	Size             *int64                     `json:"size" gorm:"column:size;"`
	ContentKey       *string                    `json:"contentKey" gorm:"column:content_key;"`
	ContentType      *enums.MaterialContentType `json:"contentType" gorm:"column:content_type;"`
	ParseMediaType   string                     `json:"parseMediaType" gorm:"column:parse_media_type;"`
}

type PartialUpdateMaterialInput = PartialUpdateInput[UpdateMaterialInput]
