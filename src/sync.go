package main

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
	localConnectStatus *UdpConnectStatus
}

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

type saveFrame struct {
	buf      []byte
	cbuf     int
	fream    int
	checksum int
}

type savedState struct {
	frames []saveFrame
	head   int
}

type SyncEvent struct {
	syncType       SyncEvntType
	confirmedInput syncConfirmedInput
}

func NewSync(status UdpConnectStatus, config *SyncConfig) Sync {
	s := Sync{
		config:              *config,
		callbacks:           config.callbacks,
		maxPredictionFrames: config.numPredictionFrames,
		localConnectStatus:  &status,
		frameCount:          0,
		lastConfirmedFrame:  -1,
		rollingBack:         false,
	}
	return s
	//s.CreateQueues(config)
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
