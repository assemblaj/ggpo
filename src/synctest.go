package ggthx

import (
	"log"
	"os"
)

type SyncTestBackend struct {
	callbacks     SessionCallbacks
	sync          Sync
	numPlayers    int
	checkDistance int
	lastVerified  int
	rollingBack   bool
	running       bool
	logFp         os.File
	game          string

	currentInput GameInput
	lastInput    GameInput
	savedFrames  RingBuffer[savedInfo]
}

type savedInfo struct {
	frame    int
	checksum int
	buf      []byte
	cbuf     int
	input    GameInput
}

func NewSyncTestBackend(cb *SessionCallbacks,
	gameName string,
	numPlayers int,
	frames int, inputSize int) SyncTestBackend {
	s := SyncTestBackend{
		callbacks:     *cb,
		numPlayers:    numPlayers,
		checkDistance: frames,
		savedFrames:   NewRingBuffer[savedInfo](32)}
	s.currentInput.Erase()
	s.currentInput.Bits = make([]byte, GameInputMaxBytes*GameInputMaxPlayers)
	s.currentInput.Size = inputSize
	var config SyncConfig
	config.callbacks = s.callbacks
	config.numPredictionFrames = MAX_PREDICTION_FRAMES
	config.inputSize = inputSize
	s.sync = NewSync(nil, &config)

	s.callbacks.BeginGame(gameName)
	return s
}

func (s *SyncTestBackend) DoPoll(timeout int) ErrorCode {
	if !s.running {
		var info Event
		info.Code = EventCodeRunning
		s.callbacks.OnEvent(&info)
		s.running = true
	}
	return Ok
}

func (s *SyncTestBackend) AddPlayer(player *Player, handle *PlayerHandle) ErrorCode {
	if player.PlayerNum < 1 || player.PlayerNum > s.numPlayers {
		return ErrorCodePlayerOutOfRange
	}
	*handle = (PlayerHandle(player.PlayerNum - 1))
	return Ok
}

func (s *SyncTestBackend) AddLocalInput(player PlayerHandle, values []byte, size int) ErrorCode {
	if !s.running {
		return ErrorCodeNotSynchronized
	}

	//index := int(player)
	//for i := 0; i < size; i++ {
	//	s.currentInput.Bits[index*size+i] |= values[i]
	//}
	s.currentInput.Bits = make([]byte, len(values))
	copy(s.currentInput.Bits, values)
	return Ok
}

func (s *SyncTestBackend) SyncInput(discconectFlags *int) ([][]byte, ErrorCode) {
	if s.rollingBack {
		var info savedInfo
		var err error
		info, err = s.savedFrames.Front()
		if err != nil {
			panic(err)
		}
		s.lastInput = info.input
	} else {
		if s.sync.GetFrameCount() == 0 {
			s.sync.SaveCurrentFrame()
		}
		s.lastInput = s.currentInput
	}
	var values = make([]byte, len(s.lastInput.Bits))
	copy(values, s.lastInput.Bits)

	if *discconectFlags > 0 {
		*discconectFlags = 0
	}
	return [][]byte{values}, Ok
}

func (s *SyncTestBackend) IncrementFrame() ErrorCode {
	s.sync.IncrementFrame()
	//s.currentInput.Erase()

	log.Printf("End of frame(%d)...\n", s.sync.GetFrameCount())

	if s.rollingBack {
		return Ok
	}

	frame := s.sync.GetFrameCount()
	// Hold onto the current frame in our queue of saved states. We'll need
	// the checksum later to verify that our replay of the same frame got the
	// same results
	var err error
	var info savedInfo
	info.frame = frame
	info.input = s.lastInput
	info.cbuf = s.sync.GetLastSavedFrame().cbuf
	info.buf = make([]byte, len(s.sync.GetLastSavedFrame().buf))
	copy(info.buf, s.sync.GetLastSavedFrame().buf)
	info.checksum = s.sync.GetLastSavedFrame().checksum

	s.savedFrames.Push(info)

	if frame-s.lastVerified == s.checkDistance {
		// We've gone far enough ahead and should now now start replaying frames
		// Load the last verified frame and set the rollback flag to true.
		s.sync.LoadFrame(s.lastVerified)

		s.rollingBack = true
		for !s.savedFrames.Empty() {
			s.callbacks.AdvanceFrame(0)

			// Verify that the checksum of this frame is the same as the one in our
			// list
			info, err = s.savedFrames.Front()
			if err != nil {
				panic(err)
			}
			s.savedFrames.Pop()

			if info.frame != s.sync.GetFrameCount() {
				log.Printf("Frame number %d does not match saved frame number %d", info.frame, frame)
				panic("RaiseSyncError")
			}
			checksum := s.sync.GetLastSavedFrame().checksum
			if info.checksum != checksum {
				s.LogGameStates(info)
				log.Printf("Checksum for frame %d does not match saved (%d != %d)", frame, checksum, info.checksum)
				panic("RaiseSyncError")
			}
			log.Printf("Checksum %08d for frame %d matches.\n", checksum, info.frame)
		}
		s.lastVerified = frame
		s.rollingBack = false
	}
	return Ok
}

func (s SyncTestBackend) LogGameStates(info savedInfo) {
	s.callbacks.LogGameState("saved:", info.buf, len(info.buf))
	s.callbacks.LogGameState("loaded:", s.sync.GetLastSavedFrame().buf, s.sync.GetLastSavedFrame().cbuf)
}

// We must 'impliment' these for this to be a true Session
func (s *SyncTestBackend) Chat(text string) ErrorCode { return ErrorCodeInvalidRequest }
func (s *SyncTestBackend) DisconnectPlayer(handle *PlayerHandle) ErrorCode {
	return ErrorCodeInvalidRequest
}
func (s *SyncTestBackend) GetNetworkStats(stats *NetworkStats, handle PlayerHandle) ErrorCode {
	return ErrorCodeInvalidRequest
}
func (s *SyncTestBackend) Logv(format string, args ...int) ErrorCode {
	return ErrorCodeInvalidRequest
}
func (s *SyncTestBackend) SetFrameDelay(player PlayerHandle, delay int) ErrorCode {
	return ErrorCodeInvalidRequest
}
func (s *SyncTestBackend) SetDisconnectTimeout(timeout int) ErrorCode {
	return ErrorCodeInvalidRequest
}
func (s *SyncTestBackend) SetDisconnectNotifyStart(timeout int) ErrorCode {
	return ErrorCodeInvalidRequest
}
func (s *SyncTestBackend) Close() ErrorCode { return ErrorCodeInvalidRequest }
