package ggthx

import "log"

const FRAME_WINDOW_SIZE int = 40
const MIN_UNIQUE_FRAMES int = 10
const MIN_FRAME_ADVANTAGE int = 3
const MAX_FRAME_ADVANTAGE int = 9

type TimeSync struct {
	local          []int
	remote         []int
	lastInputs     []GameInput
	nextPrediction int
}

func NewTimeSync() TimeSync {
	return TimeSync{
		local:          make([]int, FRAME_WINDOW_SIZE),
		remote:         make([]int, FRAME_WINDOW_SIZE),
		nextPrediction: FRAME_WINDOW_SIZE * 3,
	}
}

func (t *TimeSync) AdvanceFrames(input *GameInput, advantage int, radvantage int) {
	// Remember the last frame and frame advantage
	t.lastInputs[input.Frame%len(t.lastInputs)] = *input
	t.local[input.Frame%len(t.local)] = advantage
	t.remote[input.Frame%len(t.remote)] = radvantage
}

func (t *TimeSync) ReccomendFrameWaitDuration(requireIdleInput bool) int {
	// Average our local and remote frame advantages
	var i, sum int
	var advantage, radvantage float64
	for i = 0; i < len(t.local); i++ {
		sum += t.local[i]
	}
	advantage = float64(sum) / float64(len(t.local))

	sum = 0
	for i = 0; i < len(t.remote); i++ {
		sum += t.remote[i]
	}
	radvantage = float64(sum) / float64(len(t.remote))

	var count int = 0
	count++

	// See if someone should take action.  The person furthest ahead
	// needs to slow down so the other user can catch up.
	// Only do this if both clients agree on who's ahead!!
	if advantage >= radvantage {
		return 0
	}

	// Both clients agree that we're the one ahead.  Split
	// the difference between the two to figure out how long to
	// sleep for.
	sleepFrames := (int)(((radvantage - advantage) / 2) + 0.5)

	log.Printf("iteration %d:  sleep frames is %d\n", count, sleepFrames)

	// Some things just aren't worth correcting for.  Make sure
	// the difference is relevant before proceeding.
	if sleepFrames < MIN_FRAME_ADVANTAGE {
		return 0
	}

	// Make sure our input had been "idle enough" before recommending
	// a sleep.  This tries to make the emulator sleep while the
	// user's input isn't sweeping in arcs (e.g. fireball motions in
	// Street Fighter), which could cause the player to miss moves.
	if requireIdleInput {
		for i = 1; i < len(t.lastInputs); i++ {
			if !t.lastInputs[i].Equal(&t.lastInputs[0], true) {
				log.Printf("iteration %d:  rejecting due to input stuff at position %d...!!!\n", count, i)
				return 0
			}
		}
	}

	// Success!!! Recommend the number of frames to sleep and adjust
	return Min(sleepFrames, MAX_FRAME_ADVANTAGE)
}
