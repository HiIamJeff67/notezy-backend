package routinestypes

import "github.com/google/uuid"

type LinkableRoutineAndTag struct {
	RoutineId    uuid.UUID `json:"routineId" validate:"required"`
	RoutineTagId uuid.UUID `json:"routineTagId" validate:"required"`
}
