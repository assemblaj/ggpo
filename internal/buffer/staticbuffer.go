package buffer

import "errors"

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

func (s *StaticBuffer[T]) Get(i int) (T, error) {
	var element T
	if i < 0 || i >= s.size {
		return element, errors.New("ggpo StaticBuffer Get: i < 0 || i >= s.size")
	}
	return s.elements[i], nil
}

func (s *StaticBuffer[T]) PushBack(t T) error {
	if s.size == s.capacity-1 {
		return errors.New("ggpo StaticBuffer PushBack :s.size == s.capacity-1")
	}
	s.elements[s.size] = t
	s.size++
	return nil
}

func (s *StaticBuffer[T]) Size() int {
	return s.size
}
