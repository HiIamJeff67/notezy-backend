// Package pointers provides small helpers for working with pointers.
package pointers

// ToPtr returns a pointer to value.
func ToPtr[T any](value T) *T {
	return &value
}
