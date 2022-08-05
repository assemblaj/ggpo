package ggthx_test

import (
	"testing"

	ggthx "github.com/assemblaj/ggthx/src"
)

/*
	Charecterization Tests
*/
func TestFirstLastConfirmedFrame(t *testing.T) {
	queue := ggthx.NewInputQueue(0, 50)
	want := ggthx.NullFrame
	got := queue.LastConfirmedFrame()
	if want != got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestFirstFirstIncorrectFrame(t *testing.T) {
	queue := ggthx.NewInputQueue(0, 50)
	want := ggthx.NullFrame
	got := queue.FirstIncorrectFrame()
	if want != got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

// AddInput
func TestAddFirstInput(t *testing.T) {
	queue := ggthx.NewInputQueue(0, 50)
	input, _ := ggthx.NewGameInput(0, nil, 50)
	queue.AddInput(&input)
	want := 0
	got := queue.LastConfirmedFrame()
	if want != got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestAddFirstInputLength(t *testing.T) {
	queue := ggthx.NewInputQueue(0, 50)
	input, _ := ggthx.NewGameInput(0, nil, 50)
	queue.AddInput(&input)
	want := 1
	got := queue.Length()
	if want != got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}

}

func TestAddTenInputsDifferentFrame(t *testing.T) {
	queue := ggthx.NewInputQueue(0, 50)
	for i := 0; i < 10; i++ {
		input, _ := ggthx.NewGameInput(i, nil, 50)
		queue.AddInput(&input)
	}
	want := 10
	got := queue.Length()
	if want != got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}

}

func TestAddTenInputsSameFrame(t *testing.T) {
	queue := ggthx.NewInputQueue(0, 50)
	for i := 0; i < 10; i++ {
		input, _ := ggthx.NewGameInput(0, nil, 50)
		queue.AddInput(&input)
	}
	want := 1
	got := queue.Length()
	if want != got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

// GetInput
func TestGetInputEmptyQueue(t *testing.T) {
	queue := ggthx.NewInputQueue(0, 50)
	var input ggthx.GameInput
	_, err := queue.GetInput(0, &input)
	if err != nil {
		t.Errorf("expected nil, got an error %s", err)
	}
}

func TestInputQueueNegativeInputSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic due to negative input size passed.")
		}
	}()
	ggthx.NewInputQueue(0, -80)
}
func TestInputQueueDiscardConfirmedFramesNegative(t *testing.T) {
	queue := ggthx.NewInputQueue(0, 4)
	err := queue.DiscardConfirmedFrames(-1)
	if err == nil {
		t.Errorf("DiscardConfirmedFrames should throw an error when the frame number passed is negative.")
	}
}
