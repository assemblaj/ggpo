package ggthx_test

import (
	"testing"

	ggthx "github.com/assemblaj/ggthx/src"
)

func TestRingBufferSizeZero(t *testing.T) {
	rb := ggthx.NewRingBuffer[int](0)
	rbSize := rb.Size()
	if rbSize != 0 {
		t.Errorf("Expected Size 0, got %d", rbSize)
	}

}
func TestRingBufferIsEmpty(t *testing.T) {
	rb := ggthx.NewRingBuffer[int](1)

	if !rb.Empty() {
		t.Errorf("Expected to be empty, instead not empty")
	}

}

func TestRingBufferPush(t *testing.T) {
	rb := ggthx.NewRingBuffer[int](2)
	rb.Push(1)
	if rb.Empty() {
		t.Errorf("Expected to have value, instead was empty")
	}

}

func TestRingBufferPop(t *testing.T) {
	rb := ggthx.NewRingBuffer[int](2)
	rb.Push(1)
	rb.Pop()
	if !rb.Empty() {
		t.Errorf("Expected to be empty, instead not empty")
	}

}

func TestRingBufferItem(t *testing.T) {
	rb := ggthx.NewRingBuffer[int](3)
	rb.Push(1)
	rb.Push(2)
	var val int

	val, _ = rb.Item(0)
	if val != 1 {
		t.Errorf("Expected to have value: 1, instead was %d", val)
	}

	val, _ = rb.Item(1)
	if val != 2 {
		t.Errorf("Expected to have value: 2, instead was %d", val)
	}
}

func TestRingBufferFront(t *testing.T) {
	rb := ggthx.NewRingBuffer[int](3)
	rb.Push(1)
	rb.Push(2)
	var val int

	val, _ = rb.Front()
	if val != 1 {
		t.Errorf("Expected to have value: 1, instead was %d", val)
	}

	rb.Pop()
	val, _ = rb.Front()
	if val != 2 {
		t.Errorf("Expected to have value: 2, instead was %d", val)
	}
}

func TestRingBufferPushOverMaxError(t *testing.T) {
	rb := ggthx.NewRingBuffer[int](3)
	rb.Push(1)
	rb.Push(1)
	rb.Push(1)
	err := rb.Push(1)
	if err == nil {
		t.Errorf("Pushing past capcity should create an arror.")
	}
}

func TestRingBufferItemError(t *testing.T) {
	rb := ggthx.NewRingBuffer[int](3)
	_, err := rb.Item(1232)
	if err == nil {
		t.Errorf("Trying to get an item larger than the size of the buffer should be an error.")
	}
}
