package ggpo

import (
	"log"
	"os"

	"github.com/assemblaj/GGPO-Go/internal/buffer"
	"github.com/assemblaj/GGPO-Go/internal/input"
	"github.com/assemblaj/GGPO-Go/internal/polling"
	"github.com/assemblaj/GGPO-Go/internal/protocol"
	"github.com/assemblaj/GGPO-Go/pkg/transport"
)

type SyncTest struct {
	session       Session
	sync          Sync
	numPlayers    int
	checkDistance int
	lastVerified  int
	rollingBack   bool
	running       bool
	logFp         os.File
	game          string

	currentInput input.GameInput
	lastInput    input.GameInput
	savedFrames  buffer.RingBuffer[savedInfo]
}

type savedInfo struct {
	frame    int
	checksum int
	buf      []byte
	cbuf     int
	input    input.GameInput
}

func NewSyncTest(cb Session,
	numPlayers int,
	frames int, inputSize int) SyncTest {
	s := SyncTest{
		session:       cb,
		numPlayers:    numPlayers,
		checkDistance: frames,
		savedFrames:   buffer.NewRingBuffer[savedInfo](32)}
	s.currentInput.Erase()
	s.currentInput.Bits = make([]byte, input.GameInputMaxBytes*input.GameInputMaxPlayers)
	s.currentInput.Size = inputSize
	var config SyncConfig
	config.session = s.session
	config.numPredictionFrames = MaxPredictionFrames
	config.inputSize = inputSize
	s.sync = NewSync(nil, &config)

	return s
}

func (s *SyncTest) Idle(timeout int, timeFunc ...polling.FuncTimeType) error {
	if !s.running {
		var info Event
		info.Code = EventCodeRunning
		s.session.OnEvent(&info)
		s.running = true
	}
	return nil
}

func (s *SyncTest) AddPlayer(player *Player, handle *PlayerHandle) error {
	if player.PlayerNum < 1 || player.PlayerNum > s.numPlayers {
		return Error{Code: ErrorCodePlayerOutOfRange, Name: "ErrorCodePlayerOutOfRange"}
	}
	*handle = (PlayerHandle(player.PlayerNum - 1))
	return nil
}

func (s *SyncTest) AddLocalInput(player PlayerHandle, values []byte, size int) error {
	if !s.running {
		return Error{Code: ErrorCodeNotSynchronized, Name: "ErrorCodeNotSynchronized"}
	}

	//index := int(player)
	//for i := 0; i < size; i++ {
	//	s.currentInput.Bits[index*size+i] |= values[i]
	//}
	//copy(s.currentInput.Bits, values)
	index := int(player)
	start := index * size
	end := start + size
	copy(s.currentInput.Bits[start:end], values)
	return nil
}

func (s *SyncTest) SyncInput(discconectFlags *int) ([][]byte, error) {
	if s.rollingBack {
		var info savedInfo
		var err error
		info, err = s.savedFrames.Front()
		if err != nil {
			panic(err)
		}
		s.lastInput = *info.input.Clone()
	} else {
		if s.sync.FrameCount() == 0 {
			s.sync.SaveCurrentFrame()
		}
		s.lastInput = *s.currentInput.Clone()
	}
	//var values = make([]byte, len(s.lastInput.Bits))
	//copy(values, s.lastInput.Bits)

	values := make([][]byte, s.numPlayers)
	offset := 0
	counter := 0
	for offset < len(s.lastInput.Bits) && counter < s.numPlayers {
		values[counter] = s.lastInput.Bits[offset : s.lastInput.Size+offset]
		offset += s.lastInput.Size
		counter++
	}

	if *discconectFlags > 0 {
		*discconectFlags = 0
	}
	return values, nil
}

func (s *SyncTest) AdvanceFrame() error {
	s.sync.AdvanceFrame()
	s.currentInput.Erase()

	log.Printf("End of frame(%d)...\n", s.sync.FrameCount())

	if s.rollingBack {
		return nil
	}

	frame := s.sync.FrameCount()
	// Hold onto the current frame in our queue of saved states. We'll need
	// the checksum later to verify that our replay of the same frame got the
	// same results
	var err error
	var info savedInfo
	info.frame = frame
	info.input = s.lastInput
	info.checksum = s.sync.GetLastSavedFrame().checksum

	err = s.savedFrames.Push(info)
	if err != nil {
		panic(err)
	}

	if frame-s.lastVerified == s.checkDistance {
		// We've gone far enough ahead and should now now start replaying frames
		// Load the last verified frame and set the rollback flag to true.
		err = s.sync.LoadFrame(s.lastVerified)
		if err != nil {
			panic(err)
		}

		s.rollingBack = true
		for !s.savedFrames.Empty() {
			s.session.AdvanceFrame(0)

			// Verify that the checksum of this frame is the same as the one in our
			// list
			info, err = s.savedFrames.Front()
			if err != nil {
				panic(err)
			}
			err = s.savedFrames.Pop()
			if err != nil {
				panic(err)
			}
			if info.frame != s.sync.FrameCount() {
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
	return nil
}

func (s *SyncTest) LogGameStates(info savedInfo) {
	//s.session.LogGameState("saved:", info.buf, len(info.buf))
	//s.session.LogGameState("loaded:", s.sync.GetLastSavedFrame().buf, s.sync.GetLastSavedFrame().cbuf)
}

// We must 'impliment' these for this to be a true Session
func (s *SyncTest) DisconnectPlayer(handle PlayerHandle) error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}
func (s *SyncTest) GetNetworkStats(handle PlayerHandle) (protocol.NetworkStats, error) {
	return protocol.NetworkStats{}, Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}
func (s *SyncTest) SetFrameDelay(player PlayerHandle, delay int) error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}
func (s *SyncTest) SetDisconnectTimeout(timeout int) error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}
func (s *SyncTest) SetDisconnectNotifyStart(timeout int) error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}
func (s *SyncTest) Close() error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}

func (s *SyncTest) Start() {}

func (s *SyncTest) InitializeConnection(c ...transport.Connection) error {
	return nil
}
