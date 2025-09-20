package dtos

/* ============================== DTOs for the Highest Layer Wrapper ============================== */

type NotezyRequest[H, C, B, P any] struct {
	Header        H
	ContextFields C
	Body          B
	Param         P
}

// // Since we response it directly, so we don't have to define this
// type NotezyResponse[D any] struct {
// 	Success   bool
// 	Data      D
// 	Exception *exceptions.Exception
// }

/* ============================== Higher Layer DTO for Get Many Operations ============================== */

type SimpleSearchDto struct {
	Query  string `form:"query" validate:"omitempty,max=256"`
	Limit  int32  `form:"limit" validate:"omitempty,min=1"`
	Offset int32  `form:"offset" validate:"omitempty,min=0"`
}

/* ============================== Higher Layer DTO for Partial Update Operations ============================== */

type PartialUpdateDto[T any] struct {
	Values  T                `json:"values"`
	SetNull *map[string]bool `json:"setNull" validate:"omitempty"`
}
