package stationstypes

import "github.com/google/uuid"

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
