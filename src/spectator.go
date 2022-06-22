package ggthx

import (
	"log"
)

const SPECTATOR_FRAME_BUFFER_SIZE int = 64

type SpectatorBackend struct {
	callbacks       GGTHXSessionCallbacks
	poll            Poll
	udp             Udp
	host            UdpProtocol
	synchonizing    bool
	inputSize       int
	numPlayers      int
	nextInputToSend int
	inputs          []GameInput
}

func NewSpectatorBackend(cb *GGTHXSessionCallbacks,
	gameName string, lcaolPort int, numPlayers int, inputSize int, hostIp string, hostPort int) SpectatorBackend {
	s := SpectatorBackend{}
	s.numPlayers = numPlayers
	s.inputSize = inputSize
	s.nextInputToSend = 0

	s.callbacks = *cb
	s.synchonizing = true

	inputs := make([]GameInput, SPECTATOR_FRAME_BUFFER_SIZE)
	for _, i := range inputs {
		i.Frame = -1
	}
	s.inputs = inputs
	//port := strconv.Itoa(hostPort)
	s.udp = NewUdp(hostIp, hostPort, &s.poll, &s)
	s.host = NewUdpProtocol(&s.udp, &s.poll, 0, hostIp, nil)
	s.poll = NewPoll()
	s.callbacks.BeginGame(gameName)
	return s
}

func (s *SpectatorBackend) DoPoll(timeout int) GGTHXErrorCode {
	s.poll.Pump(0)

	s.PollUdpProtocolEvents()
	return GGTHX_OK
}

func (s *SpectatorBackend) SyncInput(disconnectFlags *int) ([]byte, GGTHXErrorCode) {
	// Wait until we've started to return inputs
	if s.synchonizing {
		return nil, GGTHX_ERRORCODE_NOT_SYNCHRONIZED
	}

	input := s.inputs[s.nextInputToSend%SPECTATOR_FRAME_BUFFER_SIZE]
	if input.Frame < s.nextInputToSend {
		// Haved recieved input from the host yet. Wait
		return nil, GGTHX_ERRORCODE_PREDICTION_THRESHOLD
	}
	if input.Frame > s.nextInputToSend {
		// The host is way way way far ahead of the spetator. How'd this
		// happen? Any, the input we need is gone forever.
		return nil, GGTHX_ERRORCODE_GENERAL_FAILURE
	}

	size := len(input.Bits)
	Assert(size >= s.inputSize*s.numPlayers)

	values := make([]byte, size)
	copy(values, input.Bits)

	if disconnectFlags != nil {
		*disconnectFlags = 0 // xxx: we should get them from the host! -pond3r
	}
	s.nextInputToSend++
	return values, GGTHX_OK
}

func (s *SpectatorBackend) IncrementFrame() GGTHXErrorCode {
	log.Printf("End of frame (%d)...\n", s.nextInputToSend-1)
	s.DoPoll(0)
	s.PollUdpProtocolEvents()

	return GGTHX_OK
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
	var info GGTHXEvent

	switch evt.eventType {
	case ConnectedEvent:
		info.Code = GGTHX_EVENTCODE_CONNECTED_TO_PEER
		info.player = 0
		s.callbacks.OnEvent(&info)

	case SynchronizingEvent:
		info.Code = GGTHX_EVENTCODE_SYNCHRONIZING_WITH_PEER
		info.player = 0
		info.count = evt.count
		info.total = evt.total
		s.callbacks.OnEvent(&info)

	case SynchronziedEvent:
		if s.synchonizing {
			info.Code = GGTHX_EVENTCODE_SYNCHRONIZED_WITH_PEER
			info.player = 0
			s.callbacks.OnEvent(&info)

			info.Code = GGTHX_EVENTCODE_RUNNING
			s.callbacks.OnEvent(&info)
			s.synchonizing = false
		}

	case NetworkInterruptedEvent:
		info.Code = GGTHX_EVENTCODE_CONNECTION_INTERRUPTED
		info.player = 0
		info.disconnectTimeout = evt.disconnectTimeout
		s.callbacks.OnEvent(&info)

	case NetworkResumedEvent:
		info.Code = GGTHX_EVENTCODE_CONNECTION_RESUMED
		info.player = 0
		s.callbacks.OnEvent(&info)

	case DisconnectedEvent:
		info.Code = GGTHX_EVENTCODE_DISCONNECTED_FROM_PEER
		info.player = 0
		s.callbacks.OnEvent(&info)

	case InputEvent:
		input := evt.input

		s.host.SetLocalFrameNumber(input.Frame)
		s.host.SendInputAck()
		s.inputs[input.Frame%SPECTATOR_FRAME_BUFFER_SIZE] = input
	}
}

func (s *SpectatorBackend) OnMsg(msg *UdpMsg, len int) {
	//if s.host.HandlesMsg()
	s.host.OnMsg(msg, len)
}
