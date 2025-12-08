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
		capacity: capacity,
	}
}

func (q *Queue[T]) maintain() {
	if q.tail == nil {
		return
	}

	for q.tail.next != nil {
		q.tail = q.tail.next
	}
}

func (q *Queue[T]) Enqueue(element T) error {
	if q.IsFull() {
		return errors.New("queue is full")
	}

	if q.tail == nil {
		q.tail = &node[T]{
			element: element,
			next:    nil,
		}
		q.head = q.tail
		q.size = 1
		return nil
	}

	q.maintain()
	q.tail.next = &node[T]{
		element: element,
		next:    nil,
	}
	q.tail = q.tail.next
	q.size++
	return nil
}

func (q *Queue[T]) Dequeue() (T, error) {
	if q.head == nil || q.IsEmpty() {
		var zero T
		return zero, errors.New("queue is empty")
	}

	result := q.head.element
	q.head = q.head.next
	q.size--
	return result, nil
}

func (q *Queue[T]) Top() (T, error) {
	if q.head == nil || q.IsEmpty() {
		var zero T
		return zero, errors.New("queue is empty")
	}
	return q.head.element, nil
}

func (q *Queue[T]) Size() int {
	return q.size
}

func (q *Queue[T]) Capacity() int {
	return q.capacity
}

func (q *Queue[T]) SetCapacity(capacity int) error {
	if capacity < q.size {
		return errors.New("the new capacity is less than the current size of the queue, may cause data loss")
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
