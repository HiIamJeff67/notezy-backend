package coretypes

import "github.com/google/uuid"

type CreatableSubShelf struct {
	Id             *uuid.UUID `json:"id" validate:"omitnil"`
	RootShelfId    uuid.UUID  `json:"rootShelfId" validate:"required"`
	PrevSubShelfId *uuid.UUID `json:"prevSubShelfId" validate:"omitnil"`
	Name           string     `json:"name" validate:"required,min=1,max=128,isshelfname"`
}
