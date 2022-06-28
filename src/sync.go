package ggthx

import (
	"log"
)

type Sync struct {
	callbacks  GGTHXSessionCallbacks
	savedState savedState
	config     SyncConfig

	rollingBack         bool
	lastConfirmedFrame  int
	frameCount          int
	maxPredictionFrames int

	inputQueues []InputQueue

	eventQueue         RingBuffer[SyncEvent]
	localConnectStatus []UdpConnectStatus
}

const MAX_PREDICTION_FRAMES int = 8

type SyncConfig struct {
	callbacks           GGTHXSessionCallbacks
	numPredictionFrames int
	numPlayers          int
	inputSize           int
}

type SyncEvntType int

const SyncConfirmedInput SyncEvntType = 0

// have to get better about lowercase names for things I don't want to be public
type syncConfirmedInput struct {
	input GameInput
}

type savedFrame struct {
	buf      []byte
	cbuf     int
	frame    int
	checksum int
}

type savedState struct {
	frames []savedFrame
	head   int
}

type SyncEvent struct {
	syncType       SyncEvntType
	confirmedInput syncConfirmedInput
}

func NewSync(status []UdpConnectStatus, config *SyncConfig) Sync {
	s := Sync{
		config:              *config,
		callbacks:           config.callbacks,
		maxPredictionFrames: config.numPredictionFrames,
		localConnectStatus:  status,
		frameCount:          0,
		lastConfirmedFrame:  -1,
		rollingBack:         false,
		savedState: savedState{
			frames: make([]savedFrame, MAX_PREDICTION_FRAMES+2)},
		eventQueue: NewRingBuffer[SyncEvent](32),
	}
	s.CreateQueues(*config)
	return s
}

// using close to mean delete
func (s *Sync) Close() {
	// delete frames manually here rather than in a destructor of the sendFrame
	// structure so we cna efficiently copy frames via weak references
	// - pond3r
	for i := 0; i < len(s.savedState.frames); i++ {
		s.callbacks.FreeBuffer(s.savedState.frames[i].buf)
	}
	s.inputQueues = nil
}

func (s *Sync) SetLastConfirmedFrame(frame int) {
	s.lastConfirmedFrame = frame
	if s.lastConfirmedFrame > 0 {
		for i := 0; i < s.config.numPlayers; i++ {
			s.inputQueues[i].DiscardConfirmedFrames(frame - 1)
		}
	}
}

func (s *Sync) AddLocalInput(queue int, input *GameInput) bool {
	framesBehind := s.frameCount - s.lastConfirmedFrame

	if s.frameCount >= s.maxPredictionFrames && framesBehind >= s.maxPredictionFrames {
		log.Printf("Rejecting input from emulator: reached prediction barrier.\n")
		return false
	}

	if s.frameCount == 0 {
		s.SaveCurrentFrame()
	}

	log.Printf("Sending undelayed local frame %d to queue %d.\n", s.frameCount, queue)
	input.Frame = s.frameCount
	s.inputQueues[queue].AddInput(input)

	return true
}

func (s *Sync) AddRemoteInput(queue int, input *GameInput) {
	s.inputQueues[queue].AddInput(input)
}

// originally took in a void ptr buffer and filled it with input
// maybe i should return that the filled buffer instead idk
// used by p2pbackend
func (s Sync) GetConfirmedInputs(frame int) ([][]byte, int) {
	disconnectFlags := 0

	//Assert(size >= s.config.numPlayers*s.config.inputSize)

	//values := make([]byte, size)
	var values [][]byte
	for i := 0; i < s.config.numPlayers; i++ {
		var input GameInput
		if s.localConnectStatus[i].Disconnected > 0 && frame > s.localConnectStatus[i].LastFrame {
			disconnectFlags |= (1 << i)
			input.Erase()
		} else {
			s.inputQueues[i].GetConfirmedInput(frame, &input)
		}
		// this was originally a memcpy
		values = append(values, input.Bits)
	}
	return values, disconnectFlags
}

// used by p2pbackend
func (s Sync) SynchronizeInputs() ([][]byte, int) {
	disconnectFlags := 0

	//Assert(size >= s.config.numPlayers*s.config.inputSize)

	//values := make([]byte, size)
	var values [][]byte
	for i := 0; i < s.config.numPlayers; i++ {
		var input GameInput
		if s.localConnectStatus[i].Disconnected > 0 && s.frameCount > s.localConnectStatus[i].LastFrame {
			disconnectFlags |= (1 << i)
			input.Erase()
		} else {
			s.inputQueues[i].GetInput(s.frameCount, &input)
		}
		values = append(values, input.Bits)
	}
	return values, disconnectFlags
}

func (s Sync) CheckSimulation(timeout int) {
	var seekTo int
	if !s.CheckSimulationConsistency(&seekTo) {
		s.AdjustSimulation(seekTo)
	}
}

func (s *Sync) IncrementFrame() {
	s.frameCount++
	s.SaveCurrentFrame()
}

func (s *Sync) AdjustSimulation(seekTo int) {
	frameCount := s.frameCount
	count := s.frameCount - seekTo

	log.Printf("Catching up\n")
	s.rollingBack = true

	// flush our input queue and load the last frame
	s.LoadFrame(seekTo)
	Assert(s.frameCount == seekTo)

	// Advance frame by frame (stuffing notifications back to
	// the master).
	s.ResetPrediction(s.frameCount)
	for i := 0; i < count; i++ {
		s.callbacks.AdvanceFrame(0)
	}
	Assert(s.frameCount == frameCount)

	s.rollingBack = false

	log.Printf("---\n")
}

func (s *Sync) LoadFrame(frame int) {
	if frame == s.frameCount {
		log.Printf("Skipping NOP.\n")
		return
	}
	// Move the head pointer back and load it up
	s.savedState.head = s.FindSavedFrameIndex(frame)
	state := s.savedState.frames[s.savedState.head]

	log.Printf("=== Loading frame info %d (size: %d  checksum: %08x).\n",
		state.frame, state.cbuf, state.checksum)
	Assert(state.buf != nil && state.cbuf > 0)
	s.callbacks.LoadGameState(state.buf, state.cbuf)

	// Reset framecount and the head of the state ring-buffer to point in
	// advance of the current frame (as if we had just finished executing it).
	s.frameCount = state.frame
	s.savedState.head = (s.savedState.head + 1) % len(s.savedState.frames)
}

func (s *Sync) SaveCurrentFrame() {
	// originally was
	// SavedFrame *state = _savedstate.frames + _savedstate.head;
	state := s.savedState.frames[s.savedState.head]
	if state.buf != nil {
		s.callbacks.FreeBuffer(state.buf)
		state.buf = nil
	}
	state.frame = s.frameCount
	buf, result := s.callbacks.SaveGameState(&state.cbuf, &state.checksum, state.frame)
	if result {
		state.buf = buf
	}
	s.savedState.frames[s.savedState.head] = state
	log.Printf("=== Saved frame info %d (size: %d  checksum: %08x).\n", state.frame, state.cbuf, state.checksum)
	s.savedState.head = (s.savedState.head + 1) % len(s.savedState.frames)
}

func (s *Sync) GetLastSavedFrame() savedFrame {
	i := s.savedState.head - 1
	for i < 0 {
		i = len(s.savedState.frames) - 1
	}
	return s.savedState.frames[i]
}

func (s *Sync) FindSavedFrameIndex(frame int) int {
	count := len(s.savedState.frames)
	var i int
	for i = 0; i < count; i++ {
		if s.savedState.frames[i].frame == frame {
			break
		}
	}
	if i == count {
		Assert(false)
	}
	return i
}

func (s *Sync) CreateQueues(config SyncConfig) bool {
	s.inputQueues = make([]InputQueue, s.config.numPlayers)
	for i := 0; i < s.config.numPlayers; i++ {
		s.inputQueues[i] = NewInputQueue(i, s.config.inputSize)
	}
	return true
}

func (s Sync) CheckSimulationConsistency(seekTo *int) bool {
	firstInorrect := NullFrame
	for i := 0; i < s.config.numPlayers; i++ {
		incorrect := s.inputQueues[i].GetFirstIncorrectFrame()
		log.Printf("considering incorrect frame %d reported by queue %d.\n", incorrect, i)

		if incorrect != NullFrame && (firstInorrect == NullFrame || incorrect < firstInorrect) {
			firstInorrect = incorrect
		}
	}

	if firstInorrect == NullFrame {
		log.Printf("prediction ok.  proceeding.\n")
		return true
	}
	*seekTo = firstInorrect
	return false
}

func (s Sync) SetFrameDelay(queue int, delay int) {
	s.inputQueues[queue].SetFrameDelay(delay)
}

func (s Sync) ResetPrediction(frameNumber int) {
	for i := 0; i < s.config.numPlayers; i++ {
		s.inputQueues[i].ResetPrediction(frameNumber)
	}
}

func (s Sync) GetEvent(e *SyncEvent) bool {
	if s.eventQueue.Size() > 0 {
		*e = s.eventQueue.Front()
		s.eventQueue.Pop()
		return true
	}
	return false
}

func (s Sync) GetFrameCount() int {
	return s.frameCount
}

func (s Sync) InRollback() bool {
	return s.rollingBack
}
