package ggthx

type StaticBuffer[T any] struct {
	elements []T
	size     int
	capacity int
}

func NewStaticBuffer[T any](capacity int) StaticBuffer[T] {
	return StaticBuffer[T]{
		elements: make([]T, capacity),
		capacity: capacity,
	}
}

func (s StaticBuffer[T]) Get(i int) T {
	Assert(i >= 0 && i < s.size)
	return s.elements[i]
}

func (s *StaticBuffer[T]) PushBack(t T) {
	Assert(s.size != s.capacity-1)
	s.elements[s.size] = t
	s.size++
}

func (s StaticBuffer[T]) Size() int {
	return s.size
}
