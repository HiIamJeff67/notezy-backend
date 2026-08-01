package concurrency

type Job[T any] struct {
	Index int
	Data  T
}

type Result[R any] struct {
	Index int
	Data  R
	Err   error
}

type ProcessorFunc[T any, R any] func(data T) (R, error)
