package ggthx

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

func (r RingBuffer[T]) Front() T {
	Assert(r.size != r.capacity)
	return r.elements[r.tail]
}

func (r RingBuffer[T]) Item(i int) T {
	Assert(i < r.size)
	return r.elements[(r.tail+i)%r.capacity]
}

// hmm this fails if its at its max size
// seems to not make sense
func (r *RingBuffer[T]) Pop() {
	Assert(r.size != r.capacity)
	r.tail = (r.tail + 1) % r.capacity
	r.size--
}

func (r *RingBuffer[T]) Push(element T) {
	Assert(r.size != r.capacity-1)
	r.elements[r.head] = element
	r.head = (r.head + 1) % r.capacity
	r.size++
}

func (r RingBuffer[T]) Size() int {
	return r.size
}

func (r RingBuffer[T]) Empty() bool {
	return r.size == 0
}
