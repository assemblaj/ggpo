package ggpo

import (
	"errors"
	"log"

	"github.com/assemblaj/GGPO-Go/internal/input"
	"github.com/assemblaj/GGPO-Go/internal/transport"
)

type Sync struct {
	session    Session
	savedState savedState
	config     SyncConfig

	rollingBack         bool
	lastConfirmedFrame  int
	frameCount          int
	maxPredictionFrames int

	inputQueues []input.InputQueue

	localConnectStatus []transport.UdpConnectStatus
}

//const MaxPredictionFrames int = 8

type SyncConfig struct {
	session             Session
	numPredictionFrames int
	numPlayers          int
	inputSize           int
}

type SyncEvntType int

const SyncConfirmedInput SyncEvntType = 0

type savedFrame struct {
	frame    int
	checksum int
}

type savedState struct {
	frames []savedFrame
	head   int
}

func NewSync(status []transport.UdpConnectStatus, config *SyncConfig) Sync {
	s := Sync{
		config:              *config,
		session:             config.session,
		maxPredictionFrames: config.numPredictionFrames,
		localConnectStatus:  status,
		frameCount:          0,
		lastConfirmedFrame:  -1,
		rollingBack:         false,
		savedState: savedState{
			frames: make([]savedFrame, MaxPredictionFrames+2)},
	}
	s.CreateQueues(*config)
	return s
}
func NeweSyncConfig(session Session, numPredictionFrames int,
	numPlayers int, inputSize int) SyncConfig {
	return SyncConfig{
		session:             session,
		numPredictionFrames: numPredictionFrames,
		numPlayers:          numPlayers,
		inputSize:           inputSize,
	}
}
func (s *SyncConfig) InputSize() int {
	return s.inputSize
}

// using close to mean delete
func (s *Sync) Close() {
	// delete frames manually here rather than in a destructor of the sendFrame
	// structure so we cna efficiently copy frames via weak references
	// - pond3r
	s.inputQueues = nil
}

func (s *Sync) SetLastConfirmedFrame(frame int) {
	s.lastConfirmedFrame = frame
	if s.lastConfirmedFrame > 0 {
		for i := 0; i < s.config.numPlayers; i++ {
			err := s.inputQueues[i].DiscardConfirmedFrames(frame - 1)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (s *Sync) AddLocalInput(queue int, input *input.GameInput) bool {
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
	err := s.inputQueues[queue].AddInput(input)
	if err != nil {
		panic(err)
	}

	return true
}

func (s *Sync) AddRemoteInput(queue int, input *input.GameInput) {

	err := s.inputQueues[queue].AddInput(input)
	if err != nil {
		panic(err)
	}
}

// originally took in a void ptr buffer and filled it with input
// maybe i should return that the filled buffer instead idk
// used by p2pbackend
func (s *Sync) GetConfirmedInputs(frame int) ([][]byte, int) {
	disconnectFlags := 0
	//Assert(size >= s.config.numPlayers*s.config.inputSize)

	//values := make([]byte, size)
	var values [][]byte
	for i := 0; i < s.config.numPlayers; i++ {
		var input input.GameInput
		if (s.localConnectStatus[i].Disconnected) && (int32(frame) > s.localConnectStatus[i].LastFrame) {
			disconnectFlags |= (1 << i)
			input.Erase()
		} else {
			_, err := s.inputQueues[i].GetConfirmedInput(frame, &input)
			if err != nil {
				panic(err)
			}
		}
		// this was originally a memcpy
		values = append(values, input.Bits)
	}
	return values, disconnectFlags
}

// used by p2pbackend
func (s *Sync) SynchronizeInputs() ([][]byte, int) {
	disconnectFlags := 0
	//Assert(size >= s.config.numPlayers*s.config.inputSize)

	//values := make([]byte, size)
	var values [][]byte
	for i := 0; i < s.config.numPlayers; i++ {
		var input input.GameInput
		if s.localConnectStatus[i].Disconnected && int32(s.frameCount) > s.localConnectStatus[i].LastFrame {
			disconnectFlags |= (1 << i)
			input.Erase()
		} else {
			_, err := s.inputQueues[i].GetInput(s.frameCount, &input)
			if err != nil {
				panic(err)
			}
		}
		values = append(values, input.Bits)
	}
	return values, disconnectFlags
}

func (s *Sync) CheckSimulation(timeout int) {

	var seekTo int
	if !s.CheckSimulationConsistency(&seekTo) {
		err := s.AdjustSimulation(seekTo)
		if err != nil {
			panic(err)
		}
	}
}

func (s *Sync) AdvanceFrame() {
	s.frameCount++

	s.SaveCurrentFrame()
}

func (s *Sync) AdjustSimulation(seekTo int) error {
	frameCount := s.frameCount
	count := s.frameCount - seekTo

	log.Printf("Catching up\n")
	s.rollingBack = true

	// flush our input queue and load the last frame
	err := s.LoadFrame(seekTo)
	if err != nil {
		panic(err)
	}

	if s.frameCount != seekTo {
		return errors.New("ggpo Sync AdjustSimulation: s.frameCount != seekTo")
	}

	// Advance frame by frame (stuffing notifications back to
	// the master).
	s.ResetPrediction(s.frameCount)
	for i := 0; i < count; i++ {
		s.session.AdvanceFrame(0)
	}

	if s.frameCount != frameCount {
		return errors.New("ggpo Sync AdjustSimulation: s.frameCount != frameCount")
	}
	s.rollingBack = false

	log.Printf("---\n")
	return nil
}

func (s *Sync) LoadFrame(frame int) error {
	if frame == s.frameCount {
		log.Printf("Skipping NOP.\n")
		return nil
	}
	// Move the head pointer back and load it up
	var err error
	s.savedState.head, err = s.FindSavedFrameIndex(frame)
	if err != nil {
		panic(err)
	}
	state := s.savedState.frames[s.savedState.head]

	log.Printf("=== Loading frame info %d (checksum: %08x).\n",
		state.frame, state.checksum)
	//s.callbacks.LoadGameState(state.buf, state.cbuf)
	s.session.LoadGameState(s.savedState.head)

	// Reset framecount and the head of the state ring-buffer to point in
	// advance of the current frame (as if we had just finished executing it).
	s.frameCount = state.frame

	s.savedState.head = (s.savedState.head + 1) % len(s.savedState.frames)
	return nil
}

func (s *Sync) SaveCurrentFrame() {

	// originally was
	// SavedFrame *state = _savedstate.frames + _savedstate.head;
	state := s.savedState.frames[s.savedState.head]
	state.frame = s.frameCount
	checksum := s.session.SaveGameState(s.savedState.head)
	state.checksum = checksum

	s.savedState.frames[s.savedState.head] = state
	log.Printf("=== Saved frame info %d (checksum: %08x).\n", state.frame, state.checksum)
	s.savedState.head = (s.savedState.head + 1) % len(s.savedState.frames)
}

func (s *Sync) GetLastSavedFrame() savedFrame {
	i := s.savedState.head - 1
	for i < 0 {
		i = len(s.savedState.frames) - 1
	}
	return s.savedState.frames[i]
}

// Trying to load a frame when it hasn't been saved causes an error
// that will panic up the chair.
func (s *Sync) FindSavedFrameIndex(frame int) (int, error) {

	count := len(s.savedState.frames)
	var i int
	for i = 0; i < count; i++ {
		if s.savedState.frames[i].frame == frame {
			break
		}
	}
	if i == count {
		return 0, errors.New("ggpo Sync FindSavedFrameIndex: i == count")
	}
	return i, nil
}

func (s *Sync) CreateQueues(config SyncConfig) bool {

	s.inputQueues = make([]input.InputQueue, s.config.numPlayers)
	for i := 0; i < s.config.numPlayers; i++ {
		s.inputQueues[i] = input.NewInputQueue(i, s.config.inputSize)
	}
	return true
}

func (s *Sync) CheckSimulationConsistency(seekTo *int) bool {

	firstInorrect := input.NullFrame
	for i := 0; i < s.config.numPlayers; i++ {
		incorrect := s.inputQueues[i].FirstIncorrectFrame()
		log.Printf("considering incorrect frame %d reported by queue %d.\n", incorrect, i)

		if incorrect != input.NullFrame && (firstInorrect == input.NullFrame || incorrect < firstInorrect) {
			firstInorrect = incorrect
		}
	}

	if firstInorrect == input.NullFrame {
		log.Printf("prediction ok.  proceeding.\n")
		return true
	}
	*seekTo = firstInorrect
	return false
}

func (s *Sync) SetFrameDelay(queue int, delay int) {
	s.inputQueues[queue].SetFrameDelay(delay)
}

func (s *Sync) ResetPrediction(frameNumber int) {
	for i := 0; i < s.config.numPlayers; i++ {
		err := s.inputQueues[i].ResetPrediction(frameNumber)
		if err != nil {
			panic(err)
		}
	}
}

func (s *Sync) FrameCount() int {
	return s.frameCount
}

func (s *Sync) InRollback() bool {
	return s.rollingBack
}
