package stack

import "errors"

type node[T any] struct {
	element T
	prev    *node[T]
}

type Stack[T any] struct {
	top      *node[T]
	size     int
	capacity int
}

func NewStack[T any](capacity int) Stack[T] {
	return Stack[T]{
		top:      nil,
		size:     0,
		capacity: max(0, capacity),
	}
}

// push an element to the top of the stack
func (st *Stack[T]) Push(element T) error {
	if st.IsFull() {
		return errors.New("stack is full")
	}

	newNode := &node[T]{
		element: element,
		prev:    nil,
	}

	if st.top == nil {
		st.top = newNode
	} else {
		newNode.prev = st.top
		st.top = newNode
	}

	st.size++
	return nil
}

// pop an element (include returning the popped element) from the top of the stack
func (st *Stack[T]) Pop() (T, error) {
	if st.IsEmpty() {
		var zero T
		return zero, errors.New("stack is empty")
	}

	result := st.top.element
	st.top = st.top.prev
	st.size--

	return result, nil
}

// get the top of the stack
func (st *Stack[T]) Top() (T, error) {
	if st.top == nil || st.IsEmpty() {
		var zero T
		return zero, errors.New("stack is empty")
	}
	return st.top.element, nil
}

// get the current size of the stack
func (st *Stack[T]) Size() int {
	return st.size
}

// get the capacity, the maximum elements that the stack is affordable
func (st *Stack[T]) Capacity() int {
	return st.capacity
}

// set the capacity to the given value
func (st *Stack[T]) SetCapacity(capacity int) error {
	if capacity < st.size {
		return errors.New("the new capacity is less than the current size of the stack, this operation will cause data loss which is unacceptable")
	}
	st.capacity = capacity
	return nil
}

func (st *Stack[T]) IsFull() bool {
	return st.size >= st.capacity
}

func (st *Stack[T]) IsEmpty() bool {
	return st.size == 0
}
