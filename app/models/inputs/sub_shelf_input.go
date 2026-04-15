package inputs

import "github.com/google/uuid"

type CreateSubShelfInput struct {
	PrevSubShelfId *uuid.UUID `json:"prevSubShelfId" gorm:"column:prev_sub_shelf_id;"`
	Name           string     `json:"name" gorm:"column:name;"`
	// will be automatically set to the path of the prevSubShelf
	// Path []uuid.UUID `json:"path" gorm:"column:path;"`
}

type BulkCreateSubShelfInput struct {
	RootShelfId    uuid.UUID  `json:"rootShelfId" gorm:"column:root_shelf_id;"`
	PrevSubShelfId *uuid.UUID `json:"prevSubShelfId" gorm:"column:prev_sub_shelf_id;"`
	Name           string     `json:"name" gorm:"column:name;"`
}

type UpdateSubShelfInput struct {
	Name *string `json:"name" gorm:"column:name;"`
}

type PartialUpdateSubShelfInput = PartialUpdateInput[UpdateSubShelfInput]

type BulkUpdateSubShelfInput struct {
	Id                 uuid.UUID                               `json:"id" gorm:"column:id;"`
	PartialUpdateInput PartialUpdateInput[UpdateSubShelfInput] `json:"partialUpdateInput"`
}
