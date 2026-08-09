package routinetasktypes

import "github.com/google/uuid"

type CreateRootShelfRoutineTaskPayload struct {
	Id      *uuid.UUID         `json:"id" validate:"omitnil"`
	Name    string             `json:"name" validate:"required,min=1,max=128,isshelfname"`
	Pattern RoutineTaskPattern `json:"pattern" validate:"omitempty,dive"`
}

type UpdateRootShelfRoutineTaskPayload struct {
	RootShelfId uuid.UUID          `json:"rootShelfId" validate:"required"`
	Name        *string            `json:"name" validate:"omitnil,min=1,max=128,isshelfname"`
	Pattern     RoutineTaskPattern `json:"pattern" validate:"omitempty,dive"`
}

type ResetRootShelfRoutineTaskPayload struct {
	RootShelfId uuid.UUID `json:"rootShelfId" validate:"required"`
}
