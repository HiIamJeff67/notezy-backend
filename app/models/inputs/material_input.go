package inputs

import (
	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
)

type CreateMaterialInput struct {
	Id            uuid.UUID                 `json:"id" gorm:"column:id;"` // we allowed the API to generate its id, since it's faster
	RootShelfId   uuid.UUID                 `json:"rootShelfId" gorm:"column:root_shelf_id;"`
	ParentShelfId uuid.UUID                 `json:"parentShelfId" gorm:"column:parent_shelf_id;"`
	Name          string                    `json:"name" gorm:"column:name;"`
	Type          enums.MaterialType        `json:"type" gorm:"column:type;"`
	Size          int64                     `json:"size" gorm:"column:size;"`
	ContentKey    string                    `json:"contentKey" gorm:"column:content_key;"`
	ContentType   enums.MaterialContentType `json:"contentType" gorm:"column:content_type;"`
}

type UpdateMaterialInput struct {
	ParentShelfId *uuid.UUID                 `json:"parentShelfId" gorm:"column:parent_shelf_id"`
	Name          *string                    `json:"name" gorm:"column:name;"`
	Type          *enums.MaterialType        `json:"type" gorm:"column:type;"`
	Size          *int64                     `json:"size" gorm:"column:size;"`
	ContentKey    *string                    `json:"contentKey" gorm:"column:content_key;"`
	ContentType   *enums.MaterialContentType `json:"contentType" gorm:"column:content_type;"`
}

type PartialUpdateMaterialInput = PartialUpdateInput[UpdateMaterialInput]
