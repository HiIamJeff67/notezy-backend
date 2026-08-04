package subshelvestypes

import "github.com/google/uuid"

type UpdatableSubShelf struct {
	SubShelfId uuid.UUID `json:"subShelfId" validate:"required"`
	Values     struct {
		Name *string `json:"name" validate:"omitnil,min=1,max=128,isshelfname"`
	} `json:"values"`
	SetNull *map[string]bool `json:"setNull,omitempty"`
}
