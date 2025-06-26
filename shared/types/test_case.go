package types

type TestCase[ArgType any, ReturnType any] struct {
	Args    ArgType
	Returns ReturnType
}
