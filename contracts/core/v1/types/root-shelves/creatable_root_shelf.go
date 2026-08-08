package coretypes

import "github.com/google/uuid"

type CreatableRootShelf struct {
	Id   *uuid.UUID `json:"id,omitempty" validate:"omitnil"`
	Name string     `json:"name" validate:"required,min=1,max=128,isshelfname"`
}
