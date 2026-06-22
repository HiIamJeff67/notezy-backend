package dtos

import "encoding/json"

/* ============================== One Dimensional Data ============================== */

type OneDimensionalData[T ~int | ~int32 | ~int64 | ~float32 | float64] struct {
	Data []T `json:"data"`
}

/* ============================== Two Dimensional Data ============================== */

type TwoDimensionalDatum[T ~int | ~int32 | ~int64 | ~float32 | ~float64] struct {
	Id    string          `json:"id"`
	X     string          `json:"x"`
	Value T               `json:"value"`
	Meta  json.RawMessage `json:"meta"`
}

type TwoDimensionalData[T ~int | ~int32 | ~int64 | ~float32 | ~float64] struct {
	Data []TwoDimensionalDatum[T] `json:"data"`
}

/* ============================== Three Dimensional Data ============================== */

type ThreeDimensionalDatum[T ~int | ~int32 | ~int64 | ~float32 | ~float64] struct {
	Id    string          `json:"id"`
	X     string          `json:"x"`
	Y     string          `json:"y"`
	Value T               `json:"value"`
	Meta  json.RawMessage `json:"meta"`
}

type ThreeDimensionalData[T ~int | ~int32 | ~int64 | ~float32 | ~float64] struct {
	Data []ThreeDimensionalDatum[T] `json:"data"`
}
