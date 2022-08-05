package ggthx

import "errors"

type RingBuffer[T any] struct {
	elements []T
	head     int
	tail     int
	size     int
	capacity int
}

func NewRingBuffer[T any](capacity int) RingBuffer[T] {
	return RingBuffer[T]{
		elements: make([]T, capacity),
		capacity: capacity,
	}
}

func (r *RingBuffer[T]) Front() (T, error) {
	var element T
	if r.size == r.capacity {
		return element, errors.New("ggthx RingBuffer Front : r.size == capacity")
	}
	element = r.elements[r.tail]
	return element, nil
}

func (r *RingBuffer[T]) Item(i int) (T, error) {
	var element T
	if i >= r.size {
		return element, errors.New("ggthx RingBuffer Item: i >= r.size")
	}
	element = r.elements[(r.tail+i)%r.capacity]
	return element, nil
}

// hmm this fails if its at its max size
// seems to not make sense
func (r *RingBuffer[T]) Pop() error {
	if r.size == r.capacity {
		return errors.New("ggthx RingBuffer Pop : r.size == r.capacity")
	}
	r.tail = (r.tail + 1) % r.capacity
	r.size--
	return nil
}

func (r *RingBuffer[T]) Push(element T) error {
	if r.size == r.capacity-1 {
		return errors.New("ggthx RingBuffer Push : r.size == r.capacity -1 ")
	}
	r.elements[r.head] = element
	r.head = (r.head + 1) % r.capacity
	r.size++
	return nil
}

func (r *RingBuffer[T]) Size() int {
	return r.size
}

// It's supposed to be size == 0,  I have no idea why it's been making it to <1 but it's been
// breaking me code
// - me considering checking for r <= 0
func (r *RingBuffer[T]) Empty() bool {
	return r.size <= 0
}
