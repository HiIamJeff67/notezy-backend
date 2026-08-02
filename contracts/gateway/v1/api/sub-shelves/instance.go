package subshelvesdto

import "github.com/google/uuid"

type CreatableSubShelf struct {
	Id             *uuid.UUID `json:"id" validate:"omitnil"`
	RootShelfId    uuid.UUID  `json:"rootShelfId" validate:"required"`
	PrevSubShelfId *uuid.UUID `json:"prevSubShelfId" validate:"omitnil"`
	Name           string     `json:"name" validate:"required,min=1,max=128,isshelfname"`
}

type UpdatableSubShelf struct {
	SubShelfId uuid.UUID `json:"subShelfId" validate:"required"`
	Values     struct {
		Name *string `json:"name" validate:"omitnil,min=1,max=128,isshelfname"`
	} `json:"values"`
	SetNull *map[string]bool `json:"setNull,omitempty"`
}

type MovableSubShelf struct {
	SourceRootShelfId      uuid.UUID   `json:"sourceRootShelfId" validate:"required"`
	SourceSubShelfIds      []uuid.UUID `json:"sourceSubShelfIds" validate:"required,min=1,max=1024"`
	DestinationRootShelfId uuid.UUID   `json:"destinationRootShelfId" validate:"required"`
	DestinationSubShelfId  *uuid.UUID  `json:"destinationSubShelfId" validate:"omitnil"`
}
