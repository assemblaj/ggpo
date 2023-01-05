package sync_test

import (
	"math"
	"testing"

	"github.com/assemblaj/ggpo/internal/input"
	"github.com/assemblaj/ggpo/internal/sync"
)

/*
These are all charectization tests for the new rift synchronization now.
*/
var Delta float64 = 0.000001

func AlmostEqual(a, b float32) bool {
	return math.Abs(float64(a-b)) <= Delta
}

func TestTimeSyncRecommendFrameDuration(t *testing.T) {
	ts := sync.NewTimeSync()
	frame := 5
	bytes := []byte{1, 2, 3, 4}
	size := 4
	input, _ := input.NewGameInput(frame, bytes, size)
	ts.AdvanceFrames(&input, 8, 9)
	want := 0.004167
	got := ts.ReccomendFrameWaitDuration(false)
	if !AlmostEqual(float32(want), got) {
		t.Errorf("expected '%f' but got '%f'", want, got)
	}
}
func TestTimeSyncRecommendFrameDurationIdleInput(t *testing.T) {
	ts := sync.NewTimeSync()
	frame := 5
	bytes := []byte{1, 2, 3, 4}
	size := 4
	input, _ := input.NewGameInput(frame, bytes, size)
	ts.AdvanceFrames(&input, 8, 9)
	want := 0.004167
	got := ts.ReccomendFrameWaitDuration(true)
	if !AlmostEqual(float32(want), got) {
		t.Errorf("expected '%f' but got '%f'", want, got)
	}
}

func TestTimeSyncHighLocalFrameAdvantage(t *testing.T) {
	ts := sync.NewTimeSync()
	frame := 0
	bytes := []byte{1, 2, 3, 4}
	size := 4
	input, _ := input.NewGameInput(frame, bytes, size)
	ts.AdvanceFrames(&input, 9, 800)
	want := 3.295833
	got := ts.ReccomendFrameWaitDuration(false)
	if !AlmostEqual(float32(want), got) {
		t.Errorf("expected '%f' but got '%f'", want, got)
	}

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
	want := 9.0
	got := ts.ReccomendFrameWaitDuration(true)
	if !AlmostEqual(float32(want), got) {
		t.Errorf("expected '%f' but got '%f'", want, got)
	}
}

func TestTimeSyncHighRemoteFrameAdvantage(t *testing.T) {
	ts := sync.NewTimeSync()
	frame := 0
	bytes := []byte{1, 2, 3, 4}
	size := 4
	input, _ := input.NewGameInput(frame, bytes, size)
	ts.AdvanceFrames(&input, 800, 9)
	want := float32(-3.295833)
	got := ts.ReccomendFrameWaitDuration(false)
	if !AlmostEqual(float32(want), got) {
		t.Errorf("expected '%f' but got '%f'", want, got)
	}

}

func TestTimeSyncNoFrameAdvantage(t *testing.T) {
	ts := sync.NewTimeSync()
	frame := 0
	bytes := []byte{1, 2, 3, 4}
	size := 4
	input, _ := input.NewGameInput(frame, bytes, size)
	ts.AdvanceFrames(&input, 0, 0)
	want := 0.0
	got := ts.ReccomendFrameWaitDuration(false)
	if !AlmostEqual(float32(want), got) {
		t.Errorf("expected '%f' but got '%f'", want, got)
	}
}

func TestTimeSyncNegativeLocalFrameAdvantage(t *testing.T) {
	ts := sync.NewTimeSync()
	frame := 0
	bytes := []byte{1, 2, 3, 4}
	size := 4
	input, _ := input.NewGameInput(frame, bytes, size)
	ts.AdvanceFrames(&input, -1, 9)
	want := 0.041667
	got := ts.ReccomendFrameWaitDuration(false)
	if !AlmostEqual(float32(want), got) {
		t.Errorf("expected '%f' but got '%f'", want, got)
	}
}
func TestTimeSyncBothNegativeFrameAdvantage(t *testing.T) {
	ts := sync.NewTimeSync()
	frame := 0
	bytes := []byte{1, 2, 3, 4}
	size := 4
	input, _ := input.NewGameInput(frame, bytes, size)
	ts.AdvanceFrames(&input, -2000, -2000)
	want := 0.0
	got := ts.ReccomendFrameWaitDuration(false)
	if !AlmostEqual(float32(want), got) {
		t.Errorf("expected '%f' but got '%f'", want, got)
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
		ts.AdvanceFrames(&input, 0, float32(i))
	}
	want := 0.791667
	got := ts.ReccomendFrameWaitDuration(false)
	if !AlmostEqual(float32(want), got) {
		t.Errorf("expected '%f' but got '%f'", want, got)
	}
}
