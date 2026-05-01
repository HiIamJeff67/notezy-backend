# Stack Library

## Overview

`app/lib/stack` provides a generic LIFO stack implementation with fixed capacity.

It is designed as a lightweight in-memory data structure for internal use cases
that need predictable push/pop behavior and explicit capacity control.

## Files

- `stack.go`: stack node type, stack struct, and stack operations.

## Public API

```go
func NewStack[T any](capacity int) Stack[T]
func (st *Stack[T]) Push(element T) error
func (st *Stack[T]) Pop() (T, error)
func (st *Stack[T]) Top() (T, error)
func (st *Stack[T]) Size() int
func (st *Stack[T]) Capacity() int
func (st *Stack[T]) SetCapacity(capacity int) error
func (st *Stack[T]) IsFull() bool
func (st *Stack[T]) IsEmpty() bool
```

## Behavior Notes

- `Push` returns `stack is full` when size reaches capacity.
- `Pop` and `Top` return `stack is empty` on an empty stack.
- `NewStack` normalizes negative capacity to `0`.
- `SetCapacity` rejects values smaller than current stack size to avoid data loss.
- This implementation is not thread-safe; synchronize externally for concurrent access.

## Example

```go
package main

import (
	"fmt"

	stack "notezy-backend/shared/lib/stack"
)

func main() {
	st := stack.NewStack[int](3)

	_ = st.Push(10)
	_ = st.Push(20)

	top, _ := st.Top()
	fmt.Println("top:", top) // 20

	v, _ := st.Pop()
	fmt.Println("pop:", v) // 20

	fmt.Println("size:", st.Size()) // 1
}
```

## File Structure

```text
app/lib/stack/
├── README.md
├── LICENSE.md
└── stack.go
```
