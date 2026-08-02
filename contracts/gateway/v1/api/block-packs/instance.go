package blockpacksdto

import (
	"github.com/google/uuid"

	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"
)

type CreatableBlockPack struct {
	Id                  *uuid.UUID           `json:"id" validate:"omitnil"`
	ParentSubShelfId    uuid.UUID            `json:"parentSubShelfId" validate:"required"`
	Name                string               `json:"name" validate:"required,min=1,max=128"`
	Icon                *enums.SupportedIcon `json:"icon" validate:"omitnil,issupportedicon"`
	HeaderBackgroundURL *string              `json:"headerBackgroundURL" validate:"omitnil"`
}

type UpdatableBlockPack struct {
	BlockPackId uuid.UUID `json:"blockPackId" validate:"required"`
	Values      struct {
		Name                *string              `json:"name" validate:"omitnil,min=1,max=128"`
		Icon                *enums.SupportedIcon `json:"icon" validate:"omitnil,issupportedicon"`
		HeaderBackgroundURL *string              `json:"headerBackgroundURL" validate:"omitnil"`
	} `json:"values"`
	SetNull *map[string]bool `json:"setNull,omitempty"`
}

type MovableBlockPack struct {
	BlockPackIds                []uuid.UUID `json:"blockPackIds" validate:"required,min=1,max=100"`
	DestinationParentSubShelfId uuid.UUID   `json:"destinationParentSubShelfId" validate:"required"`
}
