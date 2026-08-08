package coretypes

import "github.com/google/uuid"

type CreatableRoutineTag struct {
	Id    *uuid.UUID `json:"id" validate:"omitnil"`
	Name  string     `json:"name" validate:"required,min=1,max=128"`
	Color string     `json:"color" validate:"omitempty,ishexcodecolor"`
	Icon  *string    `json:"icon" validate:"omitnil,issupportedicon"`
}
