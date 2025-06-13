package dtos

type PartialUpdateDto[T any] struct {
	Values  T                `json:"values"`
	SetNull *map[string]bool `json:"setNull" validate:"omitempty"`
}
