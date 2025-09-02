package inputs

import "notezy-backend/app/models/schemas/enums"

type CreateMaterialInput struct {
	Name        string                    `json:"name" gorm:"column:name;"`
	ContentURL  string                    `json:"contentURL" gorm:"column:content_url;"`
	Type        enums.MaterialType        `json:"type" gorm:"column:type;"`
	ContentType enums.MaterialContentType `json:"contentType" gorm:"column:content_type;"`
}

type UpdateMaterialInput struct {
	Name        *string                    `json:"name" gorm:"column:name;"`
	ContentURL  *[]string                  `json:"contentURL" gorm:"column:content_url;"`
	Type        *enums.MaterialType        `json:"type" gorm:"column:type;"`
	ContentType *enums.MaterialContentType `json:"contentType" gorm:"column:content_type;"`
}

type PartialUpdateMaterialInput = PartialUpdateInput[UpdateMaterialInput]
