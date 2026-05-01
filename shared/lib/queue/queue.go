package queue

import "errors"

type node[T any] struct {
	element T
	next    *node[T]
}

type Queue[T any] struct {
	head     *node[T]
	tail     *node[T]
	size     int
	capacity int
}

func NewQueue[T any](capacity int) Queue[T] {
	return Queue[T]{
		head:     nil,
		tail:     nil,
		size:     0,
		capacity: max(0, capacity),
	}
}

// enqueue an element to the tail of the queue
func (q *Queue[T]) Enqueue(element T) error {
	if q.IsFull() {
		return errors.New("queue is full")
	}

	newNode := &node[T]{
		element: element,
		next:    nil,
	}

	if q.tail == nil {
		q.head = newNode
		q.tail = newNode
	} else {
		q.tail.next = newNode
		q.tail = q.tail.next
	}

	q.size++
	return nil
}

// dequeue an element from the head of the queue
func (q *Queue[T]) Dequeue() (T, error) {
	if q.head == nil || q.IsEmpty() {
		var zero T
		return zero, errors.New("queue is empty")
	}

	result := q.head.element
	q.head = q.head.next
	q.size--

	if q.head == nil {
		q.tail = nil
	}

	return result, nil
}

// get the head of the queue
func (q *Queue[T]) Front() (T, error) {
	if q.head == nil || q.IsEmpty() {
		var zero T
		return zero, errors.New("queue is empty")
	}
	return q.head.element, nil
}

// get the current size of the queue
func (q *Queue[T]) Size() int {
	return q.size
}

// get the capacity, the maximum elements that the queue is affordable
func (q *Queue[T]) Capacity() int {
	return q.capacity
}

// set the capacity to the given value
func (q *Queue[T]) SetCapacity(capacity int) error {
	if capacity < q.size {
		return errors.New("the new capacity is less than the current size of the queue, this operation will cause data loss which is unacceptable")
	}
	q.capacity = capacity
	return nil
}

func (q *Queue[T]) IsFull() bool {
	return q.size == q.capacity
}

func (q *Queue[T]) IsEmpty() bool {
	return q.size == 0
}
