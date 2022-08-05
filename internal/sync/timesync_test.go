package sync_test

import (
	"testing"

	"github.com/assemblaj/ggthx/internal/input"
	"github.com/assemblaj/ggthx/internal/sync"
)

func TestTimeSyncRecommendFrameDuration(t *testing.T) {
	ts := sync.NewTimeSync()
	frame := 5
	bytes := []byte{1, 2, 3, 4}
	size := 4
	input, _ := input.NewGameInput(frame, bytes, size)
	ts.AdvanceFrames(&input, 8, 9)
	want := 0
	got := ts.ReccomendFrameWaitDuration(false)
	if want != got {
		t.Errorf("expected '%d' but got '%d'", want, got)
	}
}
func TestTimeSyncRecommendFrameDurationIdleInput(t *testing.T) {
	ts := sync.NewTimeSync()
	frame := 5
	bytes := []byte{1, 2, 3, 4}
	size := 4
	input, _ := input.NewGameInput(frame, bytes, size)
	ts.AdvanceFrames(&input, 8, 9)
	want := 0
	got := ts.ReccomendFrameWaitDuration(true)
	if want != got {
		t.Errorf("expected '%d' but got '%d'", want, got)
	}
}

func TestTimeSyncHighLocalFrameAdvantage(t *testing.T) {
	ts := sync.NewTimeSync()
	frame := 0
	bytes := []byte{1, 2, 3, 4}
	size := 4
	input, _ := input.NewGameInput(frame, bytes, size)
	ts.AdvanceFrames(&input, 9, 800)
	want := 9
	got := ts.ReccomendFrameWaitDuration(false)
	if want != got {
		t.Errorf("expected '%d' but got '%d'", want, got)
	}

}
func TestTimeSyncHighLocalFrameAdvantageRequireIdleInputPanic(t *testing.T) {
	ts := sync.NewTimeSync()
	frame := 0
	bytes := []byte{1, 2, 3, 4}
	size := 4
	input, _ := input.NewGameInput(frame, bytes, size)
	ts.AdvanceFrames(&input, 9, 800)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic when, within timeSync GameInput was compared to a nil pointer because there's only one input.")
		}
	}()
	ts.ReccomendFrameWaitDuration(true)
}
func TestTimeSyncHighLocalFrameAdvantageRequireIdleInput(t *testing.T) {
	ts := sync.NewTimeSync()
	frameCount := 20
	for i := 0; i < frameCount; i++ {
		frame := i
		bytes := []byte{1, 2, 3, 4}
		size := 4
		input, _ := input.NewGameInput(frame, bytes, size)
		ts.AdvanceFrames(&input, 9, 800)
	}
	want := 9
	got := ts.ReccomendFrameWaitDuration(true)
	if want != got {
		t.Errorf("expected '%d' but got '%d'", want, got)
	}
}

func TestTimeSyncHighRemoteFrameAdvantage(t *testing.T) {
	ts := sync.NewTimeSync()
	frame := 0
	bytes := []byte{1, 2, 3, 4}
	size := 4
	input, _ := input.NewGameInput(frame, bytes, size)
	ts.AdvanceFrames(&input, 800, 9)
	want := 0
	got := ts.ReccomendFrameWaitDuration(false)
	if want != got {
		t.Errorf("expected '%d' but got '%d'", want, got)
	}

}

func TestTimeSyncNoFrameAdvantage(t *testing.T) {
	ts := sync.NewTimeSync()
	frame := 0
	bytes := []byte{1, 2, 3, 4}
	size := 4
	input, _ := input.NewGameInput(frame, bytes, size)
	ts.AdvanceFrames(&input, 0, 0)
	want := 0
	got := ts.ReccomendFrameWaitDuration(false)
	if want != got {
		t.Errorf("expected '%d' but got '%d'", want, got)
	}
}

func TestTimeSyncNegativeLocalFrameAdvantage(t *testing.T) {
	ts := sync.NewTimeSync()
	frame := 0
	bytes := []byte{1, 2, 3, 4}
	size := 4
	input, _ := input.NewGameInput(frame, bytes, size)
	ts.AdvanceFrames(&input, -1, 9)
	want := 0
	got := ts.ReccomendFrameWaitDuration(false)
	if want != got {
		t.Errorf("expected '%d' but got '%d'", want, got)
	}
}
func TestTimeSyncBothNegativeFrameAdvantage(t *testing.T) {
	ts := sync.NewTimeSync()
	frame := 0
	bytes := []byte{1, 2, 3, 4}
	size := 4
	input, _ := input.NewGameInput(frame, bytes, size)
	ts.AdvanceFrames(&input, -2000, -2000)
	want := 0
	got := ts.ReccomendFrameWaitDuration(false)
	if want != got {
		t.Errorf("expected '%d' but got '%d'", want, got)
	}
}

func TestTimeSyncAdvanceFramesAndAdvantage(t *testing.T) {
	ts := sync.NewTimeSync()
	totalFrames := 20
	for i := 0; i < totalFrames; i++ {
		frame := i
		bytes := []byte{1, 2, 3, 4}
		size := 4
		input, _ := input.NewGameInput(frame, bytes, size)
		ts.AdvanceFrames(&input, 0, i)
	}
	want := 0
	got := ts.ReccomendFrameWaitDuration(false)
	if want != got {
		t.Errorf("expected '%d' but got '%d'", want, got)
	}
}
