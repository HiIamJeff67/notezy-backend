package stationsdto

import "github.com/google/uuid"

type CreatableStation struct {
	Id                  *uuid.UUID `json:"id,omitempty" validate:"omitnil"`
	Name                string     `json:"name" validate:"required,min=1,max=128"`
	Description         string     `json:"description" validate:"max=1024"`
	Icon                *string    `json:"icon,omitempty" validate:"omitnil,issupportedicon"`
	HeaderBackgroundURL *string    `json:"headerBackgroundURL,omitempty" validate:"omitnil,isimageurl"`
}

type UpdatableStation struct {
	StationId uuid.UUID `json:"stationId" validate:"required"`
	Values    struct {
		Name                *string `json:"name,omitempty" validate:"omitnil,min=1,max=128"`
		Description         *string `json:"description,omitempty" validate:"omitnil,max=1024"`
		Icon                *string `json:"icon,omitempty" validate:"omitnil,issupportedicon"`
		HeaderBackgroundURL *string `json:"headerBackgroundURL,omitempty" validate:"omitnil,isimageurl"`
	} `json:"values"`
	SetNull *map[string]bool `json:"setNull,omitempty"`
}
