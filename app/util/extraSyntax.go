package util

func AssignIfNotNil[T any](
	updatedObject *T, 
	value *T, 
) bool {
	if value != nil {
		*updatedObject = *value
		return true
	}
	return false
}

