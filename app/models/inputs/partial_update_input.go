package inputs

type PartialUpdateInput[T any] struct {
	Values  T
	SetNull *map[string]bool
}
