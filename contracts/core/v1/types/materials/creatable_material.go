package coretypes

import "github.com/google/uuid"

type CreatableMaterial struct {
	ParentSubShelfId uuid.UUID `json:"parentSubShelfId" validate:"required"`
	Name             string    `json:"name" validate:"required,min=1,max=128"`
}
