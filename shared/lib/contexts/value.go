package contexts

import (
	"context"
	"errors"
)

func WithValue[T any](ctx context.Context, key any, value T) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetValue[T any](ctx context.Context, key any) (T, error) {
	value := ctx.Value(key)
	if value == nil {
		var zero T
		return zero, errors.New("context value not found")
	}

	converted, ok := value.(T)
	if !ok {
		var zero T
		return zero, errors.New("context value has an unexpected type")
	}

	return converted, nil
}
