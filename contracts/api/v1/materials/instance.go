package materialsdto

import "github.com/google/uuid"

type CreatableMaterial struct {
	ParentSubShelfId uuid.UUID `json:"parentSubShelfId" validate:"required"`
	Name             string    `json:"name" validate:"required,min=1,max=128"`
}

type UpdatableMaterial struct {
	Name *string `json:"name" validate:"omitnil,min=1,max=128"`
}
