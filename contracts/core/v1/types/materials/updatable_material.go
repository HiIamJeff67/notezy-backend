package coretypes

type UpdatableMaterial struct {
	Name *string `json:"name" validate:"omitnil,min=1,max=128"`
}
