package util

func DerefOrNil[T any](ptr *T) interface{} {
	if ptr != nil {
		return *ptr
	}
	return nil
}
