package ggthx

import (
	"log"
)

const SpectatorFrameBufferSize int = 64

type SpectatorBackend struct {
	callbacks       SessionCallbacks
	poll            Poll
	udp             Udp
	host            UdpProtocol
	synchonizing    bool
	inputSize       int
	numPlayers      int
	nextInputToSend int
	inputs          []GameInput
}

func NewSpectatorBackend(cb *SessionCallbacks,
	gameName string, localPort int, numPlayers int, inputSize int, hostIp string, hostPort int) SpectatorBackend {
	s := SpectatorBackend{}
	s.numPlayers = numPlayers
	s.inputSize = inputSize
	s.nextInputToSend = 0

	s.callbacks = *cb
	s.synchonizing = true

	inputs := make([]GameInput, SpectatorFrameBufferSize)
	for _, i := range inputs {
		i.Frame = -1
	}
	s.inputs = inputs
	//port := strconv.Itoa(hostPort)
	s.udp = NewUdp("127.2.1.1", localPort, &s.poll, &s)
	s.host = NewUdpProtocol(&s.udp, &s.poll, 0, hostIp, hostPort, nil)
	s.poll = NewPoll()
	s.callbacks.BeginGame(gameName)
	return s
}

func (s *SpectatorBackend) DoPoll(timeout int) ErrorCode {
	s.poll.Pump(0)

	s.PollUdpProtocolEvents()
	return Ok
}

func (s *SpectatorBackend) SyncInput(disconnectFlags *int) ([]byte, ErrorCode) {
	// Wait until we've started to return inputs
	if s.synchonizing {
		return nil, ErrorCodeNotSynchronized
	}

	input := s.inputs[s.nextInputToSend%SpectatorFrameBufferSize]
	if input.Frame < s.nextInputToSend {
		// Haved recieved input from the host yet. Wait
		return nil, ErrorCodePredictionThreshod
	}
	if input.Frame > s.nextInputToSend {
		// The host is way way way far ahead of the spetator. How'd this
		// happen? Any, the input we need is gone forever.
		return nil, ErrorCodeGeneralFailure
	}

	size := len(input.Bits)
	//Assert(size >= s.inputSize*s.numPlayers)

	values := make([]byte, size)
	copy(values, input.Bits)

	if disconnectFlags != nil {
		*disconnectFlags = 0 // xxx: we should get them from the host! -pond3r
	}
	s.nextInputToSend++
	return values, Ok
}

func (s *SpectatorBackend) IncrementFrame() ErrorCode {
	log.Printf("End of frame (%d)...\n", s.nextInputToSend-1)
	s.DoPoll(0)
	s.PollUdpProtocolEvents()

	return Ok
}

func (s *SpectatorBackend) PollUdpProtocolEvents() {
	for {
		evt, ok := s.host.GetEvent()
		if ok != nil {
			break
		} else {
			s.OnUdpProtocolEvent(evt)
		}
	}
}

func (s *SpectatorBackend) OnUdpProtocolEvent(evt *UdpProtocolEvent) {
	var info Event

	switch evt.eventType {
	case ConnectedEvent:
		info.Code = EventCodeConnectedToPeer
		info.player = 0
		s.callbacks.OnEvent(&info)

	case SynchronizingEvent:
		info.Code = EventCodeSynchronizingWithPeer
		info.player = 0
		info.count = evt.count
		info.total = evt.total
		s.callbacks.OnEvent(&info)

	case SynchronziedEvent:
		if s.synchonizing {
			info.Code = EventCodeSynchronizedWithPeer
			info.player = 0
			s.callbacks.OnEvent(&info)

			info.Code = EventCodeRunning
			s.callbacks.OnEvent(&info)
			s.synchonizing = false
		}

	case NetworkInterruptedEvent:
		info.Code = EventCodeConnectionInterrupted
		info.player = 0
		info.disconnectTimeout = evt.disconnectTimeout
		s.callbacks.OnEvent(&info)

	case NetworkResumedEvent:
		info.Code = EventCodeConnectionResumed
		info.player = 0
		s.callbacks.OnEvent(&info)

	case DisconnectedEvent:
		info.Code = EventCodeDisconnectedFromPeer
		info.player = 0
		s.callbacks.OnEvent(&info)

	case InputEvent:
		input := evt.input

		s.host.SetLocalFrameNumber(input.Frame)
		s.host.SendInputAck()
		s.inputs[input.Frame%SpectatorFrameBufferSize] = input
	}
}

func (s *SpectatorBackend) OnMsg(msg *UdpMsg, len int) {
	//if s.host.HandlesMsg()
	s.host.OnMsg(msg, len)
}
