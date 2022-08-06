package buffer_test

import (
	"testing"

	"github.com/assemblaj/GGPO-Go/internal/buffer"
)

func TestNewStaticBuffer(t *testing.T) {
	sb := buffer.NewStaticBuffer[int](16)
	pushedValue := 88
	err := sb.PushBack(pushedValue)
	if err != nil {
		t.Errorf("Got error when pushing back value onto static buffer")
	}
	got, err := sb.Get(0)
	if err != nil {
		t.Errorf("Got error when trying to get value from static buffer")
	}
	want := pushedValue
	if want != got {
		t.Errorf("expected '%d' but got '%d'", want, got)
	}

}

func TestStaticBufferGetNegativeValue(t *testing.T) {
	sb := buffer.NewStaticBuffer[int](16)
	pushedValue := 88
	err := sb.PushBack(pushedValue)
	if err != nil {
		t.Errorf("Got error when pushing back value onto static buffer")
	}
	_, err = sb.Get(-1)
	if err == nil {
		t.Errorf("Should have returned error when trying to retrieve negative index.")
	}

}

func TestStaticBufferGetEmptyBuffer(t *testing.T) {
	sb := buffer.NewStaticBuffer[int](16)
	_, err := sb.Get(0)
	if err == nil {
		t.Errorf("Should have returned error when trying to retrieve from an empty buffer.")
	}

}

func TestStaticBufferAddOverCapacity(t *testing.T) {
	capacity := 16
	sb := buffer.NewStaticBuffer[int](capacity)
	for i := 0; i < capacity-1; i++ {
		sb.PushBack(i)
	}
	err := sb.PushBack(capacity)
	if err == nil {
		t.Errorf("Should have returned error for trying to push back past the capacity.")
	}
}

func TestStaticBufferSize(t *testing.T) {
	capacity := 16
	toAdd := 5
	sb := buffer.NewStaticBuffer[int](capacity)
	for i := 0; i < toAdd; i++ {
		sb.PushBack(i)
	}
	want := toAdd
	got := sb.Size()
	if want != got {
		t.Errorf("expected '%d' but got '%d'", want, got)
	}
}
