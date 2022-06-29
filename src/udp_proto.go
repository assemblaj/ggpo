package ggthx

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
)

const (
	// UDPHeaderSize is the size of the IP + UDP headers.
	UDPHeaderSize          = 28
	NumSyncPackets         = 5
	SyncRetryInterval      = 2000
	SyncFirstRetryInterval = 500
	RunningRetryInterval   = 200
	KeepAliveInterval      = 200
	QualityReportInterval  = 1000
	NetworkStatsInterval   = 1000
	UDOShutdownTimer       = 5000
	MaxSeqDistance         = 1 << 15
)

type UdpProtocol struct {
	stats UdpProtocolStats // may not need these
	event UdpProtocolEvent //

	// Network transmission information
	udp               Udp
	peerAddress       string
	peerPort          int
	magicNumber       uint16
	queue             int
	remoteMagicNumber uint16
	connected         bool
	sendLatency       int
	ooPercent         int
	ooPacket          OoPacket
	sendQueue         RingBuffer[QueueEntry]

	// Stats
	roundTripTime  int
	packetsSent    int
	bytesSent      int
	kbpsSent       int
	statsStartTime int

	// The State Machine
	localConnectStatus []UdpConnectStatus
	peerConnectStatus  []UdpConnectStatus
	currentState       UdpProtocolState
	state              UdpProtocolStateInfo

	// Fairness
	localFrameAdvantage  int
	remoteFrameAdvantage int

	// Packet Loss
	pendingOutput         RingBuffer[GameInput]
	lastRecievedInput     GameInput
	lastSentInput         GameInput
	lastAckedInput        GameInput
	lastSendTime          uint
	lastRecvTime          uint
	shutdownTimeout       uint
	disconnectEventSent   bool
	disconnectTimeout     uint
	disconnectNotifyStart uint
	disconnectNotifySent  bool

	nextSendSeq uint16
	nextRecvSeq uint16

	// Rift synchronization
	timesync TimeSync

	// Event Queue
	eventQueue RingBuffer[UdpProtocolEvent]
}

type UdpProtocolStats struct {
	ping                int
	remoteFrameAdvtange int
	localFrameAdvantage int
	sendQueueLen        int
	udp                 UdpStats
}

type UdpProtocolEvent struct {
	eventType         UdpProtocolEventType
	input             GameInput // for Input message
	total             int       // for synchronizing
	count             int       //
	disconnectTimeout int       // network interrupted
}

// for logging purposes only
func (upe UdpProtocolEvent) String() string {
	str := "(event:"
	switch upe.eventType {
	case UnknownEvent:
		str += "Unknown"
		break
	case ConnectedEvent:
		str += "Connected"
		break
	case SynchronizingEvent:
		str += "Synchronizing"
		break
	case SynchronziedEvent:
		str += "Synchronized"
		break
	case InputEvent:
		str += "Input"
		break
	case DisconnectedEvent:
		str += "Disconnected"
		break
	case NetworkInterruptedEvent:
		str += "NetworkInterrupted"
		break
	case NetworkResumedEvent:
		str += "NetworkResumed"
		break
	}
	str += ").\n"
	return str
}

type UdpProtocolEventType int

const (
	UnknownEvent UdpProtocolEventType = iota - 1
	ConnectedEvent
	SynchronizingEvent
	SynchronziedEvent
	InputEvent
	DisconnectedEvent
	NetworkInterruptedEvent
	NetworkResumedEvent
)

type UdpProtocolState int

const (
	SyncingState UdpProtocolState = iota
	SynchronziedState
	RunningState
	DisconnectedState
)

type UdpProtocolStateInfo struct {
	roundTripRemaining uint32 // sync
	random             uint32

	lastQualityReportTime    uint32 // running
	lastNetworkStatsInterval uint32
	lastInputPacketRecvTime  uint32
}

type QueueEntry struct {
	queueTime int
	destIp    string
	msg       *UdpMsg
	destPort  int
}

func (q QueueEntry) String() string {
	return fmt.Sprintf("Entry : queueTime %d destIp %s msg %s", q.queueTime, q.destIp, *q.msg)
}

func NewQueEntry(time int, destIp string, destPort int, m *UdpMsg) QueueEntry {
	return QueueEntry{
		queueTime: time,
		destIp:    destIp,
		destPort:  destPort,
		msg:       m,
	}
}

type OoPacket struct {
	sendTime int
	destIp   string
	msg      *UdpMsg
}

func NewUdpProtocol(udp *Udp, poll *Poll, queue int, ip string, port int, status []UdpConnectStatus) UdpProtocol {
	var magicNumber uint16
	for {
		magicNumber = uint16(rand.Int())
		if magicNumber != 0 {
			break
		}
	}
	peerConnectStatus := make([]UdpConnectStatus, UDPMsgMaxPlayers)
	for i := 0; i < len(peerConnectStatus); i++ {
		peerConnectStatus[i].LastFrame = -1
	}
	lastSentInput, _ := NewGameInput(-1, nil, 1)
	lastRecievedInput, _ := NewGameInput(-1, nil, 1)
	lastAckedInput, _ := NewGameInput(-1, nil, 1)

	protocol := UdpProtocol{
		udp:                *udp,
		queue:              queue,
		localConnectStatus: status,
		peerConnectStatus:  peerConnectStatus,
		peerAddress:        ip,
		peerPort:           port,
		magicNumber:        magicNumber,
		pendingOutput:      NewRingBuffer[GameInput](64),
		sendQueue:          NewRingBuffer[QueueEntry](64),
		eventQueue:         NewRingBuffer[UdpProtocolEvent](64),
		timesync:           NewTimeSync(),
		lastSentInput:      lastSentInput,
		lastRecievedInput:  lastRecievedInput,
		lastAckedInput:     lastAckedInput}
	//poll.RegisterLoop(&protocol, nil)
	return protocol
}

func (u *UdpProtocol) OnLoopPoll(cookie []byte) bool {

	// originally was if !udp
	if !u.udp.IsInitialized() {
		return true
	}

	now := uint(time.Now().UnixMilli())
	var nextInterval uint

	u.PumpSendQueue()
	switch u.currentState {
	case SyncingState:
		if int(u.state.roundTripRemaining) == NumSyncPackets {
			nextInterval = uint(SyncFirstRetryInterval)
		} else {
			nextInterval = uint(SyncRetryInterval)
		}
		if u.lastSendTime > 0 && u.lastSendTime+nextInterval < now {
			log.Printf("No luck syncing after %d ms... Re-queueing sync packet.\n", nextInterval)
			u.SendSyncRequest()
		}
		break
	case RunningState:
		if u.state.lastInputPacketRecvTime > 0 || u.state.lastInputPacketRecvTime+uint32(RunningRetryInterval) > uint32(now) {
			log.Printf("Haven't exchanged packets in a while (last received:%d  last sent:%d).  Resending.\n",
				u.lastRecievedInput.Frame, u.lastSentInput.Frame)
			u.SendPendingOutput()
			u.state.lastInputPacketRecvTime = uint32(now)
		}

		//if (!u.State.running.last_quality_report_time || _state.running.last_quality_report_time + QUALITY_REPORT_INTERVAL < now) {
		if u.state.lastQualityReportTime == 0 || uint32(u.state.lastQualityReportTime)+uint32(QualityReportInterval) < uint32(now) {
			msg := NewUdpMsg(QualityReportMsg)
			msg.QualityReport.Ping = uint32(time.Now().UnixMilli())
			msg.QualityReport.FrameAdvantage = int8(u.localFrameAdvantage)
			u.SendMsg(&msg)
			u.state.lastQualityReportTime = uint32(now)
		}

		if u.state.lastNetworkStatsInterval == 0 || u.state.lastNetworkStatsInterval+uint32(NetworkStatsInterval) < uint32(now) {
			u.UpdateNetworkStats()
			u.state.lastNetworkStatsInterval = uint32(now)
		}

		if u.lastSendTime > 0 && u.lastSendTime+uint(KeepAliveInterval) < uint(now) {
			log.Println("Sending keep alive packet")
			msg := NewUdpMsg(KeepAliveMsg)
			u.SendMsg(&msg)
		}

		if u.disconnectTimeout > 0 && u.disconnectNotifyStart > 0 &&
			!u.disconnectNotifySent && (u.lastRecvTime+u.disconnectNotifyStart < now) {
			log.Printf("Endpoint has stopped receiving packets for %d ms.  Sending notification.\n", u.disconnectNotifyStart)
			e := UdpProtocolEvent{
				eventType: NetworkInterruptedEvent}
			e.disconnectTimeout = int(u.disconnectTimeout) - int(u.disconnectNotifyStart)
			u.QueueEvent(&e)
			u.disconnectNotifySent = true
		}

		if u.disconnectTimeout > 0 && (u.lastRecvTime+u.disconnectTimeout < now) {
			if !u.disconnectEventSent {
				log.Printf("Endpoint has stopped receiving packets for %d ms.  Disconnecting.\n",
					u.disconnectTimeout)
				u.QueueEvent(&UdpProtocolEvent{
					eventType: DisconnectedEvent})
				u.disconnectEventSent = true
			}
		}
		break
	case DisconnectedState:
		if u.shutdownTimeout < now {
			log.Printf("Shutting down udp connection.\n")
			u.udp = Udp{}
			u.shutdownTimeout = 0
		}
	}
	return true
}

// all this method did, and all bitvector did, in the original
// was incode input as bits, etc
// go globs can do a lot of that for us, so i've forgone much of that logic
// https://github.com/pond3r/ggpo/blob/7ddadef8546a7d99ff0b3530c6056bc8ee4b9c0a/src/lib/ggpo/network/udp_proto.cpp#L111
func (u *UdpProtocol) SendPendingOutput() error {
	msg := NewUdpMsg(InputMsg)
	var j, offset int

	if u.pendingOutput.Size() > 0 {
		last := u.lastAckedInput
		// bits = msg.Input.Bits
		input, err := u.pendingOutput.Front()
		if err != nil {
			panic(err)
		}
		msg.Input.StartFrame = uint32(input.Frame)
		msg.Input.InputSize = uint8(input.Size)

		if !(last.Frame == -1 || last.Frame+1 == int(msg.Input.StartFrame)) {
			return errors.New("!((last.Frame == -1 || last.Frame+1 == int(msg.Input.StartFrame))) ")
		}

		for j = 0; j < u.pendingOutput.Size(); j++ {
			current, err := u.pendingOutput.Item(j)
			if err != nil {
				panic(err)
			}
			msg.Input.Bits = append(msg.Input.Bits, current.Bits...)
			last = current // might get rid of this
			u.lastSentInput = current
		}
	} else {
		msg.Input.StartFrame = 0
		msg.Input.InputSize = 0
	}

	msg.Input.AckFrame = u.lastRecievedInput.Frame
	msg.Input.NumBits = uint16(offset)

	// to represent : msg->u.input.disconnect_requested = _current_state == Disconnected;
	if u.currentState == DisconnectedState {
		msg.Input.DisconectRequested = 1
	} else {
		msg.Input.DisconectRequested = 0
	}

	if u.localConnectStatus != nil {
		for s := 0; s < UDPMsgMaxPlayers; s++ {
			msg.Input.PeerConnectStatus = append(msg.Input.PeerConnectStatus, u.localConnectStatus...)
		}
	} else {
		// huh memset(msg->u.input.peer_connect_status, 0, sizeof(UdpMsg::connect_status) * UDP_MSG_MAX_PLAYERS);
		for s := 0; s < UDPMsgMaxPlayers; s++ {
			msg.Input.PeerConnectStatus = append(msg.Input.PeerConnectStatus, UdpConnectStatus{})
		}

	}

	if offset >= MaxCompressedBits {
		return errors.New("offset >= MaxCompressedBits")
	}
	u.SendMsg(&msg)
	return nil
}

func (u *UdpProtocol) SendInputAck() {
	msg := NewUdpMsg(InputAckMsg)
	msg.InputAck.AckFrame = u.lastRecievedInput.Frame
	u.SendMsg(&msg)
}

func (u *UdpProtocol) GetEvent() (*UdpProtocolEvent, error) {
	if u.eventQueue.Size() == 0 {
		return nil, errors.New("No events")
	}
	e, err := u.eventQueue.Front()
	if err != nil {
		panic(err)
	}
	u.eventQueue.Pop()
	return &e, nil
}

func (u *UdpProtocol) QueueEvent(evt *UdpProtocolEvent) {
	log.Printf("Queueing event %s", *evt)
	u.eventQueue.Push(*evt)
}

func (u *UdpProtocol) Disconnect() {
	u.currentState = DisconnectedState
	u.shutdownTimeout = uint(time.Now().UnixMilli()) + uint(UDOShutdownTimer)
}

func (u *UdpProtocol) SendSyncRequest() {
	u.state.random = uint32(rand.Int() & 0xFFFF)
	msg := NewUdpMsg(SyncRequestMsg)
	msg.SyncRequest.RandomRequest = u.state.random
	u.SendMsg(&msg)
}

func (u *UdpProtocol) SendMsg(msg *UdpMsg) {
	log.Printf("In UdpProtocol send %s", msg)
	u.packetsSent++
	u.lastSendTime = uint(time.Now().UnixMilli())
	u.bytesSent += msg.PacketSize()
	msg.Header.Magic = u.magicNumber
	msg.Header.SequenceNumber = u.nextSendSeq
	u.nextSendSeq++
	if u.peerAddress == "" {
		panic("peerAdress empty, why?")
	}
	u.sendQueue.Push(NewQueEntry(
		int(time.Now().UnixMilli()), u.peerAddress, u.peerPort, msg))
	u.PumpSendQueue()
}

// finally bitvector at work
func (u *UdpProtocol) OnInput(msg *UdpMsg, length int) (bool, error) {
	// If a disconnect is requested, go ahead and disconnect now.
	disconnectRequested := msg.Input.DisconectRequested
	if disconnectRequested > 0 {
		if u.currentState != DisconnectedState && !u.disconnectEventSent {
			log.Printf("Disconnecting endpoint on remote request.\n")
			u.QueueEvent(&UdpProtocolEvent{
				eventType: DisconnectedEvent,
			})
			u.disconnectEventSent = true
		}
	} else {
		// update the peer connection status if this peer is still considered to be part
		// of the network
		remoteStatus := msg.Input.PeerConnectStatus
		for i := 0; i < len(u.peerConnectStatus); i++ {
			if remoteStatus[i].LastFrame < u.peerConnectStatus[i].LastFrame {
				return false, errors.New("remoteStatus[i].LastFrame < u.peerConnectStatus[i].LastFrame")
			}
			if u.peerConnectStatus[i].Disconnected > 0 || remoteStatus[i].Disconnected > 0 {
				u.peerConnectStatus[i].Disconnected = 1
			}
			u.peerConnectStatus[i].LastFrame = Max(u.peerConnectStatus[i].LastFrame, remoteStatus[i].LastFrame)
		}
	}

	// Decompress the input (gob should already have done this)
	lastRecievedFrameNumber := u.lastRecievedInput.Frame
	if msg.Input.NumBits > 0 {
		//offset := 0
		//numBits := msg.Input.NumBits
		currentFrame := msg.Input.StartFrame

		u.lastRecievedInput.Size = int(msg.Input.InputSize)
		if u.lastRecievedInput.Frame < 0 {
			u.lastRecievedInput.Frame = int(msg.Input.StartFrame - 1)
		}
		// my logic, not GGPO

		// send input as one big chunk

		u.lastRecievedInput.Frame = int(currentFrame)
		evt := UdpProtocolEvent{
			eventType: InputEvent,
			input:     u.lastRecievedInput,
		}

		u.state.lastInputPacketRecvTime = uint32(time.Now().UnixMilli())
		log.Printf("Sending frame %d to emu queue %d.\n", u.lastRecievedInput.Frame, u.queue)
		u.QueueEvent(&evt)
	}

	if u.lastRecievedInput.Frame < lastRecievedFrameNumber {
		return false, errors.New("u.lastRecievedInput.Frame < lastRecievedFrameNumber")
	}
	// Get rid of our buffered input
	input, err := u.pendingOutput.Front()
	if err != nil {
		panic(err)
	}

	for u.pendingOutput.Size() > 0 && (input.Frame < msg.Input.AckFrame) {
		log.Printf("Throwing away pending output frame %d\n", input.Frame)
		u.lastAckedInput = input
		u.pendingOutput.Pop()
	}
	return true, nil
}

func (u *UdpProtocol) OnInputAck(msg *UdpMsg, len int) (bool, error) {
	// Get rid of our buffered input
	input, err := u.pendingOutput.Front()
	if err != nil {
		panic(err)
	}

	for u.pendingOutput.Size() > 0 && input.Frame > msg.Input.AckFrame {
		log.Printf("Throwing away pending output frame %d\n", input.Frame)
		u.lastAckedInput = input
		u.pendingOutput.Pop()
	}
	return true, nil
}

func (u *UdpProtocol) OnQualityReport(msg *UdpMsg, len int) (bool, error) {
	reply := NewUdpMsg(QualityReplyMsg)
	reply.QualityReply.Pong = msg.QualityReport.Ping
	u.SendMsg(&reply)

	u.remoteFrameAdvantage = int(msg.QualityReport.FrameAdvantage)
	return true, nil
}

func (u *UdpProtocol) OnQualityReply(msg *UdpMsg, len int) (bool, error) {
	u.roundTripTime = int(time.Now().UnixMilli()) - int(msg.QualityReply.Pong)
	return true, nil
}

func (u *UdpProtocol) OnKeepAlive(msg *UdpMsg, len int) (bool, error) {
	return true, nil
}

func (u *UdpProtocol) GetNetworkStats(s *NetworkStats) {
	s.network.ping = u.roundTripTime
	s.network.sendQueueLen = u.pendingOutput.Size()
	s.network.kbpsSent = u.kbpsSent
	s.timesync.remoteFramesBehind = u.remoteFrameAdvantage
	s.timesync.localFramesBehind = u.localFrameAdvantage
}

func (u *UdpProtocol) SetLocalFrameNumber(localFrame int) {
	remoteFrame := u.lastRecievedInput.Frame + (u.roundTripTime * 60 / 1000)
	u.localFrameAdvantage = remoteFrame - localFrame
}

func (u *UdpProtocol) RecommendFrameDelay() int {
	return u.timesync.ReccomendFrameWaitDuration(false)
}

func (u *UdpProtocol) SetDisconnectTimeout(timeout int) {
	u.disconnectTimeout = uint(timeout)
}

func (u *UdpProtocol) SetDisconnectNotifyStart(timeout int) {
	u.disconnectNotifyStart = uint(timeout)
}

func (u *UdpProtocol) PumpSendQueue() error {
	var entry QueueEntry
	var err error

	for !u.sendQueue.Empty() {
		entry, err = u.sendQueue.Front()
		if err != nil {
			panic(err)
		}

		if u.sendLatency > 0 {
			jitter := (u.sendLatency * 2 / 3) + ((rand.Int() % u.sendLatency) / 3)
			if int(time.Now().UnixMilli()) < entry.queueTime+jitter {
				break
			}
		}
		if u.ooPercent > 0 && u.ooPacket.msg != nil && ((rand.Int() % 100) < u.ooPercent) {
			delay := rand.Int() % (u.sendLatency*10 + 1000)
			log.Printf("creating rogue oop (seq: %d  delay: %d)\n",
				entry.msg.Header.SequenceNumber, delay)
			u.ooPacket.sendTime = int(time.Now().UnixMilli()) + delay
			u.ooPacket.msg = entry.msg
			u.ooPacket.destIp = entry.destIp
		} else {
			return errors.New("entry.destIp == \"\"")

			u.udp.SendTo(entry.msg, entry.destIp, entry.destPort)

			// would delete the udpmsg here
		}
		u.sendQueue.Pop()
	}
	if u.ooPacket.msg != nil && u.ooPacket.sendTime < int(time.Now().UnixMilli()) {
		log.Printf("sending rogue oop!")
		u.udp.SendTo(u.ooPacket.msg, u.peerAddress, u.peerPort)
		u.ooPacket.msg = nil
	}
	return nil
}

func (u *UdpProtocol) ClearSendQueue() {
	for !u.sendQueue.Empty() {
		// i'd manually delete the QueueEntry in a language where I could
		u.sendQueue.Pop()
	}
}

// going to call deletes close
func (u *UdpProtocol) Close() {
	u.ClearSendQueue()
	u.udp.Close()
}

// this originally included something that looked like this
//    return _peer_addr.sin_addr.S_un.S_addr == from.sin_addr.S_un.S_addr &&
//   _peer_addr.sin_port == from.sin_port;
// however since in this version of GGPO we don't listen in for all packets
// and instead sendto/read from prespecified endpoints, I don't *think*
// I need to do this. Will come back around to it
func (u *UdpProtocol) HandlesMsg() bool {
	return u.udp.IsInitialized()
	//return u.peerAddress == from
}

func (u *UdpProtocol) SendInput(input *GameInput) {
	if u.udp.IsInitialized() {
		if u.currentState == RunningState {
			// check to see if this is a good time to adjust for the rift
			u.timesync.AdvanceFrames(input, u.localFrameAdvantage, u.remoteFrameAdvantage)

			// Save this input packet.
			u.pendingOutput.Push(*input)
		}
	}
	u.SendPendingOutput()
}

func (u *UdpProtocol) UpdateNetworkStats() {
	now := int(time.Now().UnixMilli())
	if u.statsStartTime == 0 {
		u.statsStartTime = now
	}

	totalBytesSent := u.bytesSent + (UDPHeaderSize * u.packetsSent)
	seconds := float64(now-u.statsStartTime) / 1000.0
	bps := float64(totalBytesSent) / seconds
	udpOverhead := float64(100.0 * (float64(UDPHeaderSize * u.packetsSent)) / float64(u.bytesSent))
	u.kbpsSent = int(bps / 1024)

	log.Printf("Network Stats -- Bandwidth: %.2f KBps Packets Sent: %5d (%.2f pps) KB Sent: %.2f UDP Overhead: %.2f %%.\n",
		float64(u.kbpsSent),
		u.packetsSent,
		float64(u.packetsSent*1000)/float64(now-u.statsStartTime),
		float64(totalBytesSent/1024.0),
		udpOverhead)
}

func (u *UdpProtocol) Synchronize() {
	u.currentState = SyncingState
	u.state.roundTripRemaining = uint32(NumSyncPackets)
	u.SendSyncRequest()
}

func (u *UdpProtocol) GetPeerConnectStatus(id int, frame *int) bool {
	*frame = u.peerConnectStatus[id].LastFrame
	// !u.peerConnectStatus[id].Disconnected from the C/++ world
	return u.peerConnectStatus[id].Disconnected == 0
}

func (u *UdpProtocol) OnInvalid(msg *UdpMsg, len int) (bool, error) {
	//  Assert(false) // ? ASSERT(FALSE && "Invalid msg in UdpProtocol");
	// ah
	log.Printf("Invalid msg in UdpProtocol ")
	return false, errors.New("Invalid msg in UdpProtocol")
}

func (u *UdpProtocol) OnSyncRequest(msg *UdpMsg, len int) (bool, error) {
	if u.remoteMagicNumber != 0 && msg.Header.Magic != u.remoteMagicNumber {
		log.Printf("Ignoring sync request from unknown endpoint (%d != %d).\n",
			msg.Header.Magic, u.remoteMagicNumber)
		return false, nil
	}
	reply := NewUdpMsg(SyncReplyMsg)
	reply.SyncReply.RandomReply = msg.SyncRequest.RandomRequest
	u.SendMsg(&reply)
	return true, nil
}

func (u *UdpProtocol) OnMsg(msg *UdpMsg, length int) {
	handled := false
	var err error

	type UdpProtocolDispatchFunc func(msg *UdpMsg, length int) (bool, error)

	table := []UdpProtocolDispatchFunc{
		u.OnInvalid,
		u.OnSyncRequest,
		u.OnSyncReply,
		u.OnInput,
		u.OnQualityReport,
		u.OnQualityReply,
		u.OnKeepAlive,
		u.OnInputAck}

	// filter out messages that don't match what we expect
	seq := msg.Header.SequenceNumber
	if msg.Header.HeaderType != uint8(SyncRequestMsg) && msg.Header.HeaderType != uint8(SyncReplyMsg) {
		if msg.Header.Magic != u.remoteMagicNumber {
			log.Printf("recv rejecting %s", msg)
			return
		}

		// filer out out-of-order packets
		skipped := seq - u.nextRecvSeq
		if skipped > uint16(MaxSeqDistance) {
			log.Printf("dropping out of order packet (seq: %d, last seq:%d)\n", seq, u.nextRecvSeq)
			return
		}
	}

	u.nextRecvSeq = seq
	log.Printf("recv %s", msg)
	if int(msg.Header.HeaderType) > len(table) {
		u.OnInvalid(msg, length)
	} else {
		handled, err = table[int(msg.Header.HeaderType)](msg, length)
	}
	if err != nil {
		panic(err)
	}

	if handled {
		u.lastRecvTime = uint(time.Now().UnixMilli())
		if u.disconnectNotifySent && u.currentState == RunningState {
			u.QueueEvent(
				&UdpProtocolEvent{
					eventType: NetworkResumedEvent,
				})
			u.disconnectNotifySent = false
		}
	}
}

func (u *UdpProtocol) OnSyncReply(msg *UdpMsg, length int) (bool, error) {
	if u.currentState != SyncingState {
		log.Println("Ignoring SyncReply while not synching.")
		return msg.Header.Magic == u.remoteMagicNumber, nil
	}

	if msg.SyncReply.RandomReply != u.state.random {
		log.Printf("sync reply %d != %d.  Keep looking...\n",
			msg.SyncReply.RandomReply, u.state.random)
		return false, nil
	}

	if !u.connected {
		u.QueueEvent(&UdpProtocolEvent{
			eventType: ConnectedEvent})
		u.connected = true
	}

	log.Printf("Checking sync state (%d round trips remaining).\n", u.state.roundTripRemaining)
	u.state.roundTripRemaining--
	if u.state.roundTripRemaining == 0 {
		log.Printf("Synchronized\n")
		u.QueueEvent(&UdpProtocolEvent{
			eventType: SynchronziedEvent,
		})
		u.currentState = RunningState
		u.lastRecievedInput.Frame = -1
		u.remoteMagicNumber = msg.Header.Magic
	} else {
		evt := UdpProtocolEvent{
			eventType: SynchronizingEvent,
		}
		evt.total = NumSyncPackets
		evt.count = NumSyncPackets - int(u.state.roundTripRemaining)
		u.QueueEvent(&evt)
		u.SendSyncRequest()
	}

	return true, nil
}

func (u *UdpProtocol) IsInitialized() bool {
	// Originially checked if udp object was nill
	return u.udp.IsInitialized()
}

func (u *UdpProtocol) IsSynchronized() bool {
	return u.currentState == RunningState
}

func (u *UdpProtocol) IsRunning() bool {
	return u.currentState == RunningState
}
