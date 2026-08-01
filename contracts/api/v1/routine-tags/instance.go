package routinetagsdto

import "github.com/google/uuid"

type CreatableRoutineTag struct {
	Id    *uuid.UUID `json:"id" validate:"omitnil"`
	Name  string     `json:"name" validate:"required,min=1,max=128"`
	Color string     `json:"color" validate:"omitempty,ishexcodecolor"`
	Icon  *string    `json:"icon" validate:"omitnil,issupportedicon"`
}

type UpdatableRoutineTag struct {
	RoutineTagId uuid.UUID `json:"routineTagId" validate:"required"`
	Values       struct {
		Name  *string `json:"name" validate:"omitnil,min=1,max=128"`
		Color *string `json:"color" validate:"omitnil,ishexcodecolor"`
		Icon  *string `json:"icon" validate:"omitnil,issupportedicon"`
	} `json:"values"`
	SetNull *map[string]bool `json:"setNull,omitempty"`
}
