package rootshelvesdto

import "github.com/google/uuid"

type CreatableRootShelf struct {
	Id   *uuid.UUID `json:"id,omitempty" validate:"omitnil"`
	Name string     `json:"name" validate:"required,min=1,max=128,isshelfname"`
}

type UpdatableRootShelf struct {
	RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
	Values      struct {
		Name *string `json:"name" validate:"omitnil,min=1,max=128,isshelfname"`
	} `json:"values"`
	SetNull *map[string]bool `json:"setNull,omitempty"`
}
