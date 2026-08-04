package rootshelvestypes

import "github.com/google/uuid"

type UpdatableRootShelf struct {
	RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
	Values      struct {
		Name *string `json:"name" validate:"omitnil,min=1,max=128,isshelfname"`
	} `json:"values"`
	SetNull *map[string]bool `json:"setNull,omitempty"`
}
