package inputs

import "github.com/google/uuid"

type CreateSubShelfInput struct {
	Name           string     `json:"name" gorm:"column:name;"`
	PrevSubShelfId *uuid.UUID `json:"prevSubShelfId" gorm:"column:prev_sub_shelf_id;"`
	// will be automatically set to the path of the prevSubShelf
	// Path []uuid.UUID `json:"path" gorm:"column:path;"`
}

type UpdateSubShelfInput struct {
	Name        *string    `json:"name" gorm:"column:name;"`
	RootShelfId *uuid.UUID `json:"rootShelfId" gorm:"column:root_shelf_id;"`
}

type PartialUpdateSubShelfInput = PartialUpdateInput[UpdateSubShelfInput]
