package sync

import (
	"github.com/assemblaj/GGPO-Go/internal/input"
	"github.com/assemblaj/GGPO-Go/internal/util"
)

const (
	FrameWindowSize   = 120
	MinUniqueFrames   = 10
	MinFrameAdvantage = 3
	MaxFrameAdvantage = 9
)

type TimeSync struct {
	local            []float32
	remote           []float32
	lastInputs       []*input.GameInput
	nextPrediction   int
	RemoteFrameDelay int
	FrameDelay2      int
	numFrame         int
	avgLocal         float32
	avgRemote        float32
	clearedInitial   bool
}

func NewTimeSync() TimeSync {
	return TimeSync{
		local:          make([]float32, FrameWindowSize),
		remote:         make([]float32, FrameWindowSize),
		lastInputs:     make([]*input.GameInput, MinUniqueFrames),
		nextPrediction: FrameWindowSize * 3,
	}
}

func (t *TimeSync) SetFrameDelay(frame int) {
	t.FrameDelay2 = frame
}

func (t *TimeSync) AdvanceFrames(input *input.GameInput, advantage float32, radvantage float32) {
	// Remember the last frame and frame advantage
	t.lastInputs[input.Frame%len(t.lastInputs)] = input
	t.local[input.Frame%len(t.local)] = advantage
	t.remote[input.Frame%len(t.remote)] = radvantage

	t.avgLocal = ((float32(t.numFrame) * t.avgLocal) + advantage) / (float32(t.numFrame) + 1)
	t.avgRemote = ((float32(t.numFrame) * t.avgRemote) + radvantage) / (float32(t.numFrame) + 1)

	t.numFrame++
	// Clear after first 3 seconds, as this is a bit crazy
	if !t.clearedInitial && t.numFrame == 240 {
		t.clearedInitial = true
		t.numFrame = 0
	}
}

func (t *TimeSync) LocalAdvantage() float32 {
	var i int
	var sum float32
	var advantage float32
	for i = 0; i < len(t.local); i++ {
		sum += t.local[i]
	}
	advantage = float32(sum) / float32(len(t.local))
	return advantage
}

func (t *TimeSync) RemoteAdvantage() float32 {
	var i int
	var sum float32
	var radvantage float32
	for i = 0; i < len(t.remote); i++ {
		sum += t.remote[i]
	}
	radvantage = float32(sum) / float32(len(t.remote))
	return radvantage
}

func (t *TimeSync) AvgLocalAdvantageSinceStart() float32 {
	return t.avgLocal
}

func (t *TimeSync) AvgRemoteAdvantageSinceStart() float32 {
	return t.avgRemote
}

func (t *TimeSync) ReccomendFrameWaitDuration(requireIdleInput bool) float32 {
	// Average our local and remote frame advantages

	advantage := t.LocalAdvantage()

	radvantage := t.RemoteAdvantage()

	// var count int = 0
	// count++

	// See if someone should take action.  The person furthest ahead
	// needs to slow down so the other user can catch up.
	// Only do this if both clients agree on who's ahead!!
	// if advantage >= radvantage {
	// 	return 0
	// }

	// Both clients agree that we're the one ahead.  Split
	// the difference between the two to figure out how long to
	// sleep for.
	sleepFrames := ((radvantage - advantage) / 2.0)

	//log.Printf("iteration %d:  sleep frames is %d\n", count, sleepFrames)
	util.Log.Printf("In TimeSync: sleep frames is %f\n", sleepFrames)

	// Some things just aren't worth correcting for.  Make sure
	// the difference is relevant before proceeding.
	// if sleepFrames < MinFrameAdvantage {
	// 	return 0
	// }

	// Make sure our input had been "idle enough" before recommending
	// a sleep.  This tries to make the emulator sleep while the
	// user's input isn't sweeping in arcs (e.g. fireball motions in
	// Street Fighter), which could cause the player to miss moves.
	if requireIdleInput {
		for i = 1; i < len(t.lastInputs); i++ {
			equal, err := t.lastInputs[i].Equal(t.lastInputs[0], true)
			if err != nil {
				panic(err)
			}
			if !equal {
				util.Log.Printf("iteration %d:  rejecting due to input stuff at position %d...!!!\n", count, i)
				return 0
			}
		}
	}


	// Success!!! Recommend the number of frames to sleep and adjust
	if sleepFrames > 0 {
		return float32(util.Min(sleepFrames, MaxFrameAdvantage))
	} else {
		return float32(util.Max(sleepFrames, -MaxFrameAdvantage))
	}
}
