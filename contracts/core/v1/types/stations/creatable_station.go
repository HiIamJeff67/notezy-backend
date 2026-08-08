package coretypes

import "github.com/google/uuid"

type CreatableStation struct {
	Id                  *uuid.UUID `json:"id,omitempty" validate:"omitnil"`
	Name                string     `json:"name" validate:"required,min=1,max=128"`
	Description         string     `json:"description" validate:"max=1024"`
	Icon                *string    `json:"icon,omitempty" validate:"omitnil,issupportedicon"`
	HeaderBackgroundURL *string    `json:"headerBackgroundURL,omitempty" validate:"omitnil,isimageurl"`
}
