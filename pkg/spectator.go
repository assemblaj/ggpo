package ggthx

import (
	"log"

	"github.com/assemblaj/ggthx/internal/input"
	"github.com/assemblaj/ggthx/internal/polling"
	"github.com/assemblaj/ggthx/internal/protocol"
	"github.com/assemblaj/ggthx/internal/transport"
)

const SpectatorFrameBufferSize int = 32
const DefaultMaxFramesBehind int = 10
const DefaultCatchupSpeed int = 1

type SpectatorBackend struct {
	callbacks       SessionCallbacks
	poll            polling.Poller
	connection      transport.Connection
	host            protocol.UdpProtocol
	synchonizing    bool
	inputSize       int
	numPlayers      int
	nextInputToSend int
	inputs          []input.GameInput
	hostIp          string
	hostPort        int
	framesBehind    int
	localPort       int
	currentFrame    int
}

func NewSpectatorBackend(cb *SessionCallbacks,
	gameName string, localPort int, numPlayers int, inputSize int, hostIp string, hostPort int) SpectatorBackend {
	s := SpectatorBackend{}
	s.numPlayers = numPlayers
	s.inputSize = inputSize
	s.nextInputToSend = 0

	s.callbacks = *cb
	s.synchonizing = true

	inputs := make([]input.GameInput, SpectatorFrameBufferSize)
	for _, i := range inputs {
		i.Frame = -1
	}
	s.inputs = inputs
	//port := strconv.Itoa(hostPort)
	//s.udp = NewUdp(&s, localPort)
	s.hostIp = hostIp
	s.hostPort = hostPort
	s.localPort = localPort
	var poll polling.Poll = polling.NewPoll()
	s.poll = &poll
	s.callbacks.BeginGame(gameName)
	//go s.udp.Read()
	return s
}

func (s *SpectatorBackend) DoPoll(timeout int, timeFunc ...polling.FuncTimeType) error {
	if len(timeFunc) == 0 {
		s.poll.Pump()
	} else {
		s.poll.Pump(timeFunc[0])
	}
	s.PollUdpProtocolEvents()

	if s.framesBehind > 0 {
		for s.nextInputToSend < s.currentFrame {
			s.callbacks.AdvanceFrame(0)
			log.Printf("In Spectator: skipping frame %d\n", s.nextInputToSend)
			s.nextInputToSend++
		}
		s.framesBehind = 0
	}

	return nil
}

func (s *SpectatorBackend) SyncInput(disconnectFlags *int) ([][]byte, error) {
	// Wait until we've started to return inputs
	if s.synchonizing {
		return nil, Error{Code: ErrorCodeNotSynchronized, Name: "ErrorCodeNotSynchronized"}
	}

	input := s.inputs[s.nextInputToSend%SpectatorFrameBufferSize]
	s.currentFrame = input.Frame
	if input.Frame < s.nextInputToSend {
		// Haved recieved input from the host yet. Wait
		return nil, Error{Code: ErrorCodePredictionThreshod, Name: "ErrorCodePredictionThreshod"}

	}
	if input.Frame > s.nextInputToSend {
		s.framesBehind = input.Frame - s.nextInputToSend
		// The host is way way way far ahead of the spetator. How'd this
		// happen? Any, the input we need is gone forever.
		return nil, Error{Code: ErrorCodeGeneralFailure, Name: "ErrorCodeGeneralFailure"}
	}
	//s.framesBehind = 0

	//Assert(size >= s.inputSize*s.numPlayers)
	values := make([][]byte, s.numPlayers)
	offset := 0
	counter := 0
	for offset < len(input.Bits) {
		values[counter] = input.Bits[offset : s.inputSize+offset]
		offset += s.inputSize
		counter++
	}

	if disconnectFlags != nil {
		*disconnectFlags = 0 // xxx: we should get them from the host! -pond3r
	}
	s.nextInputToSend++
	return values, nil
}

func (s *SpectatorBackend) IncrementFrame() error {
	log.Printf("End of frame (%d)...\n", s.nextInputToSend-1)
	s.DoPoll(0)
	s.PollUdpProtocolEvents()

	return nil
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

func (s *SpectatorBackend) OnUdpProtocolEvent(evt *protocol.UdpProtocolEvent) {
	var info Event
	switch evt.Type() {
	case protocol.ConnectedEvent:
		info.Code = EventCodeConnectedToPeer
		info.player = 0
		s.callbacks.OnEvent(&info)

	case protocol.SynchronizingEvent:
		info.Code = EventCodeSynchronizingWithPeer
		info.player = 0
		info.count = evt.Count
		info.total = evt.Total
		s.callbacks.OnEvent(&info)

	case protocol.SynchronziedEvent:
		if s.synchonizing {
			info.Code = EventCodeSynchronizedWithPeer
			info.player = 0
			s.callbacks.OnEvent(&info)

			info.Code = EventCodeRunning
			s.callbacks.OnEvent(&info)
			s.synchonizing = false
		}

	case protocol.NetworkInterruptedEvent:
		info.Code = EventCodeConnectionInterrupted
		info.player = 0
		info.disconnectTimeout = evt.DisconnectTimeout
		s.callbacks.OnEvent(&info)

	case protocol.NetworkResumedEvent:
		info.Code = EventCodeConnectionResumed
		info.player = 0
		s.callbacks.OnEvent(&info)

	case protocol.DisconnectedEvent:
		info.Code = EventCodeDisconnectedFromPeer
		info.player = 0
		s.callbacks.OnEvent(&info)

	case protocol.InputEvent:
		input := evt.Input

		s.host.SetLocalFrameNumber(input.Frame)
		s.host.SendInputAck()
		s.inputs[input.Frame%SpectatorFrameBufferSize] = input
	}
}

func (s *SpectatorBackend) HandleMessage(ipAddress string, port int, msg transport.UDPMessage, len int) {
	if s.host.HandlesMsg(ipAddress, port) {
		s.host.OnMsg(msg, len)
	}
}

func (p *SpectatorBackend) AddLocalInput(player PlayerHandle, values []byte, size int) error {
	return nil
}

func (s *SpectatorBackend) AddPlayer(player *Player, handle *PlayerHandle) error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}

// We must 'impliment' these for this to be a true Session
func (s *SpectatorBackend) Chat(text string) error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}
func (s *SpectatorBackend) DisconnectPlayer(handle PlayerHandle) error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}
func (s *SpectatorBackend) GetNetworkStats(stats *protocol.NetworkStats, handle PlayerHandle) error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}
func (s *SpectatorBackend) Logv(format string, args ...int) error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}
func (s *SpectatorBackend) SetFrameDelay(player PlayerHandle, delay int) error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}
func (s *SpectatorBackend) SetDisconnectTimeout(timeout int) error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}
func (s *SpectatorBackend) SetDisconnectNotifyStart(timeout int) error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}
func (s *SpectatorBackend) Close() error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}
func (s *SpectatorBackend) InitalizeConnection(c ...transport.Connection) error {
	if len(c) == 0 {
		s.connection = transport.NewUdp(s, s.localPort)
		return nil
	}
	s.connection = c[0]
	return nil
}

func (s *SpectatorBackend) Start() {
	//s.udp.messageHandler = s
	go s.connection.Read()

	s.host = protocol.NewUdpProtocol(s.connection, 0, s.hostIp, s.hostPort, nil)
	s.poll.RegisterLoop(&s.host, nil)
	s.host.Synchronize()

}
