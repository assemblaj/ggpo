package ggpo

import (
	"log"

	"github.com/assemblaj/GGPO-Go/internal/input"
	"github.com/assemblaj/GGPO-Go/internal/messages"
	"github.com/assemblaj/GGPO-Go/internal/polling"
	"github.com/assemblaj/GGPO-Go/internal/protocol"
	"github.com/assemblaj/GGPO-Go/pkg/transport"
)

const SpectatorFrameBufferSize int = 32
const DefaultMaxFramesBehind int = 10
const DefaultCatchupSpeed int = 1

type Spectator struct {
	session         Session
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
	messageChannel  chan transport.MessageChannelItem
}

func NewSpectator(cb Session, localPort int, numPlayers int, inputSize int, hostIp string, hostPort int) Spectator {
	s := Spectator{}
	s.numPlayers = numPlayers
	s.inputSize = inputSize
	s.nextInputToSend = 0

	s.session = cb
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
	s.messageChannel = make(chan transport.MessageChannelItem, 200)
	//go s.udp.Read()
	return s
}

func (s *Spectator) Idle(timeout int, timeFunc ...polling.FuncTimeType) error {
	s.HandleMessages()
	if len(timeFunc) == 0 {
		s.poll.Pump()
	} else {
		s.poll.Pump(timeFunc[0])
	}
	s.PollUdpProtocolEvents()

	if s.framesBehind > 0 {
		for s.nextInputToSend < s.currentFrame {
			s.session.AdvanceFrame(0)
			log.Printf("In Spectator: skipping frame %d\n", s.nextInputToSend)
			s.nextInputToSend++
		}
		s.framesBehind = 0
	}

	return nil
}

func (s *Spectator) SyncInput(disconnectFlags *int) ([][]byte, error) {
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

func (s *Spectator) AdvanceFrame(checksum uint32) error {
	log.Printf("End of frame (%d)...\n", s.nextInputToSend-1)
	s.Idle(0)
	s.PollUdpProtocolEvents()

	return nil
}

func (s *Spectator) PollUdpProtocolEvents() {
	for {
		evt, ok := s.host.GetEvent()
		if ok != nil {
			break
		} else {
			s.OnUdpProtocolEvent(evt)
		}
	}
}

func (s *Spectator) OnUdpProtocolEvent(evt *protocol.UdpProtocolEvent) {
	var info Event
	switch evt.Type() {
	case protocol.ConnectedEvent:
		info.Code = EventCodeConnectedToPeer
		info.Player = 0
		s.session.OnEvent(&info)

	case protocol.SynchronizingEvent:
		info.Code = EventCodeSynchronizingWithPeer
		info.Player = 0
		info.Count = evt.Count
		info.Total = evt.Total
		s.session.OnEvent(&info)

	case protocol.SynchronziedEvent:
		if s.synchonizing {
			info.Code = EventCodeSynchronizedWithPeer
			info.Player = 0
			s.session.OnEvent(&info)

			info.Code = EventCodeRunning
			s.session.OnEvent(&info)
			s.synchonizing = false
		}

	case protocol.NetworkInterruptedEvent:
		info.Code = EventCodeConnectionInterrupted
		info.Player = 0
		info.DisconnectTimeout = evt.DisconnectTimeout
		s.session.OnEvent(&info)

	case protocol.NetworkResumedEvent:
		info.Code = EventCodeConnectionResumed
		info.Player = 0
		s.session.OnEvent(&info)

	case protocol.DisconnectedEvent:
		info.Code = EventCodeDisconnectedFromPeer
		info.Player = 0
		s.session.OnEvent(&info)

	case protocol.InputEvent:
		input := evt.Input

		s.host.SetLocalFrameNumber(input.Frame)
		s.host.SendInputAck()
		s.inputs[input.Frame%SpectatorFrameBufferSize] = input
	}
}

func (s *Spectator) HandleMessage(ipAddress string, port int, msg messages.UDPMessage, len int) {
	if s.host.HandlesMsg(ipAddress, port) {
		s.host.OnMsg(msg, len)
	}
}

func (p *Spectator) AddLocalInput(player PlayerHandle, values []byte, size int) error {
	return nil
}

func (s *Spectator) AddPlayer(player *Player, handle *PlayerHandle) error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}

// We must 'impliment' these for this to be a true Session
func (s *Spectator) DisconnectPlayer(handle PlayerHandle) error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}
func (s *Spectator) GetNetworkStats(handle PlayerHandle) (protocol.NetworkStats, error) {
	return protocol.NetworkStats{}, Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}
func (s *Spectator) SetFrameDelay(player PlayerHandle, delay int) error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}
func (s *Spectator) SetDisconnectTimeout(timeout int) error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}
func (s *Spectator) SetDisconnectNotifyStart(timeout int) error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}
func (s *Spectator) Close() error {
	return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
}
func (s *Spectator) InitializeConnection(c ...transport.Connection) error {
	if len(c) == 0 {
		s.connection = transport.NewUdp(s, s.localPort)
		return nil
	}
	s.connection = c[0]
	return nil
}

func (s *Spectator) HandleMessages() {
	for i := 0; i < len(s.messageChannel); i++ {
		mi := <-s.messageChannel
		s.HandleMessage(mi.Peer.Ip, mi.Peer.Port, mi.Message, mi.Length)
	}
}

func (s *Spectator) Start() {
	go s.connection.Read(s.messageChannel)

	s.host = protocol.NewUdpProtocol(s.connection, 0, s.hostIp, s.hostPort, nil)
	s.poll.RegisterLoop(&s.host, nil)
	s.host.Synchronize()

}
