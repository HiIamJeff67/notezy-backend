package dtos

type GetManyDto struct {
	Query  string `form:"query" validate:"omitempty,max=256"`
	Limit  int32  `form:"limit" validate:"omitempty,min=0"`
	Offset int32  `form:"offset" validate:"omitempty,min=1"`
}
