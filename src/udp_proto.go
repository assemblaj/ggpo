package ggthx

import (
	"errors"
	"log"
	"math/rand"
	"time"
)

const UDP_HEADER_SIZE int = 28 /* Size of IP + UDP headers */
const NUM_SYNC_PACKETS int = 5
const SYNC_RETRY_INTERVAL int = 2000
const SYNC_FIRST_RETRY_INTERVAL int = 500
const RUNNING_RETRY_INTERVAL int = 200
const KEEP_ALIVE_INTERVAL int = 200
const QUALITY_REPORT_INTERVAL int = 1000
const NETWORK_STATS_INTERVAL int = 1000
const UDP_SHUTDOWN_TIMER int = 5000
const MAX_SEQ_DISTANCE int = (1 << 15)

type UdpProtocol struct {
	stats UdpProtocolStats // may not need these
	event UdpProtocolEvent //

	// Network transmission information
	udp               Udp
	peerAddress       string
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
	localConnectStatus *UdpConnectStatus
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
}

func NewQueEntry(time int, destIp string, m *UdpMsg) QueueEntry {
	return QueueEntry{
		queueTime: time,
		destIp:    destIp,
		msg:       m,
	}
}

type OoPacket struct {
	sendTime int
	destIp   string
	msg      *UdpMsg
}

func NewUdpProtocol(udp *Udp, poll *Poll, queue int, ip string, status *UdpConnectStatus) UdpProtocol {
	var magicNumber uint16
	for {
		magicNumber = uint16(rand.Int())
		if magicNumber != 0 {
			break
		}
	}
	protocol := UdpProtocol{
		udp:                *udp,
		queue:              queue,
		localConnectStatus: status,
		peerAddress:        ip,
		magicNumber:        magicNumber,
	}
	poll.RegisterLoop(protocol, nil)
	return protocol
}

func (u UdpProtocol) OnLoopPoll(cookie []byte) bool {

	//if u.udp == nil  {
	//	return true
	//}

	now := uint(time.Now().UnixMilli())
	var nextInterval uint

	u.PumpSendQueue()
	switch u.currentState {
	case SyncingState:
		if int(u.state.roundTripRemaining) == NUM_SYNC_PACKETS {
			nextInterval = uint(SYNC_FIRST_RETRY_INTERVAL)
		} else {
			nextInterval = uint(SYNC_RETRY_INTERVAL)
		}
		if u.lastSendTime > 0 && u.lastSendTime+nextInterval < now {
			log.Printf("No luck syncing after %d ms... Re-queueing sync packet.\n", nextInterval)
			u.SendSyncRequest()
		}
		break
	case RunningState:
		if u.state.lastInputPacketRecvTime > 0 || u.state.lastInputPacketRecvTime+uint32(RUNNING_RETRY_INTERVAL) > uint32(now) {
			log.Printf("Haven't exchanged packets in a while (last received:%d  last sent:%d).  Resending.\n",
				u.lastRecievedInput.Frame, u.lastSentInput.Frame)
			u.SendPendingOutput()
			u.state.lastInputPacketRecvTime = uint32(now)
		}

		//if (!u.State.running.last_quality_report_time || _state.running.last_quality_report_time + QUALITY_REPORT_INTERVAL < now) {
		if u.state.lastQualityReportTime == 0 || uint32(u.state.lastQualityReportTime)+uint32(QUALITY_REPORT_INTERVAL) < uint32(now) {
			msg := NewUdpMsg(QualityReportMsg)
			msg.QualityReport.Ping = uint32(time.Now().UnixMilli())
			msg.QualityReport.FrameAdvantage = int8(u.localFrameAdvantage)
			u.SendMsg(&msg)
			u.state.lastQualityReportTime = uint32(now)
		}

		if u.state.lastNetworkStatsInterval == 0 || u.state.lastNetworkStatsInterval+uint32(NETWORK_STATS_INTERVAL) < uint32(now) {
			u.UpdateNetworkStats()
			u.state.lastNetworkStatsInterval = uint32(now)
		}

		if u.lastSendTime > 0 && u.lastSendTime+uint(KEEP_ALIVE_INTERVAL) < uint(now) {
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

		if u.disconnectTimeout > 0 && (u.lastRecvTime+u.disconnectTimeout > now) {
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
func (u UdpProtocol) SendPendingOutput() {
	msg := NewUdpMsg(InputMsg)
	var j, offset int

	if u.pendingOutput.Size() > 0 {
		last := u.lastAckedInput
		// bits = msg.Input.Bits

		msg.Input.StartFrame = uint32(u.pendingOutput.Front().Frame)
		msg.Input.InputSize = uint8(u.pendingOutput.Front().Size)

		Assert(last.Frame == -1 || last.Frame+1 == int(msg.Input.StartFrame))
		for j = 0; j < u.pendingOutput.Size(); j++ {
			current := u.pendingOutput.Item(j)
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
		for s := 0; s < UDP_MSG_MAX_PLAYERS; s++ {
			msg.Input.PeerConnectStatus = append(msg.Input.PeerConnectStatus, *u.localConnectStatus)
		}
	} else {
		// huh memset(msg->u.input.peer_connect_status, 0, sizeof(UdpMsg::connect_status) * UDP_MSG_MAX_PLAYERS);
		for s := 0; s < UDP_MSG_MAX_PLAYERS; s++ {
			msg.Input.PeerConnectStatus = append(msg.Input.PeerConnectStatus, UdpConnectStatus{})
		}

	}

	Assert(offset < MAX_COMPRESSED_BITS)
	u.SendMsg(&msg)
}

func (u UdpProtocol) SendInputAck() {
	msg := NewUdpMsg(InputAckMsg)
	msg.InputAck.AckFrame = u.lastRecievedInput.Frame
	u.SendMsg(&msg)
}

func (u UdpProtocol) GetEvent() (*UdpProtocolEvent, error) {
	if u.eventQueue.Size() == 0 {
		return nil, errors.New("No events")
	}
	e := u.eventQueue.Front()
	u.eventQueue.Pop()
	return &e, nil
}

func (u *UdpProtocol) QueueEvent(evt *UdpProtocolEvent) {
	log.Printf("Queueing event %s", *evt)
	u.eventQueue.Push(*evt)
}

func (u *UdpProtocol) Disconnect() {
	u.currentState = DisconnectedState
	u.shutdownTimeout = uint(time.Now().UnixMilli()) + uint(UDP_SHUTDOWN_TIMER)
}

func (u *UdpProtocol) SendSyncRequest() {
	u.state.random = uint32(rand.Int() & 0xFFFF)
	msg := NewUdpMsg(SyncRequestMsg)
	msg.SyncRequest.RandomRequest = u.state.random
	u.SendMsg(&msg)
}

func (u *UdpProtocol) SendMsg(msg *UdpMsg) {
	log.Printf("send %s", msg)

	u.packetsSent++
	u.lastSendTime = uint(time.Now().UnixMilli())
	u.bytesSent += msg.PacketSize()
	msg.Header.Magic = u.magicNumber
	msg.Header.SequenceNumber = u.nextSendSeq
	u.nextSendSeq++

	u.sendQueue.Push(NewQueEntry(
		int(time.Now().UnixMilli()), u.peerAddress, msg))
	u.PumpSendQueue()
}

// finally bitvector at work
func (u *UdpProtocol) OnInput(msg *UdpMsg, length int) bool {
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
			Assert(remoteStatus[i].LastFrame == u.peerConnectStatus[i].LastFrame)
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

	Assert(u.lastRecievedInput.Frame >= lastRecievedFrameNumber)
	// Get rid of our buffered input
	for u.pendingOutput.Size() > 0 && (u.pendingOutput.Front().Frame < msg.Input.AckFrame) {
		log.Printf("Throwing away pending output frame %d\n", u.pendingOutput.Front().Frame)
		u.lastAckedInput = u.pendingOutput.Front()
		u.pendingOutput.Pop()
	}
	return true
}

func (u *UdpProtocol) OnInputAck(msg *UdpMsg, len int) bool {
	// Get rid of our buffered input
	for u.pendingOutput.Size() > 0 && u.pendingOutput.Front().Frame > msg.Input.AckFrame {
		log.Printf("Throwing away pending output frame %d\n", u.pendingOutput.Front().Frame)
		u.lastAckedInput = u.pendingOutput.Front()
		u.pendingOutput.Pop()
	}
	return true
}

// why the builder
func (u *UdpProtocol) OnQualityReport(msg *UdpMsg, len int) bool {
	reply := NewUdpMsg(QualityReplyMsg)
	reply.QualityReply.Pong = msg.QualityReport.Ping
	//reply.BuildQualityReply(int(msg.qualityReport.ping))
	u.SendMsg(&reply)

	u.remoteFrameAdvantage = int(msg.QualityReport.FrameAdvantage)
	return true
}

func (u *UdpProtocol) OnQualityReply(msg *UdpMsg, len int) bool {
	u.roundTripTime = int(time.Now().UnixMilli()) - int(msg.QualityReply.Pong)
	return true
}

func (u UdpProtocol) OnKeepAlive(msg *UdpMsg, len int) bool {
	return true
}

// why the builder
// couldn't set those values for s without a builder
// inconsistent
func (u UdpProtocol) GetNetworkStats(s *GGTHXNetworkStats) {
	s.BuildGGTHXNetworkNetworkStats(u.pendingOutput.Size(), 0,
		u.roundTripTime, u.kbpsSent)
	s.BuildGGTHXNetworkTimeSyncStats(u.localFrameAdvantage, u.remoteFrameAdvantage)
}

func (u *UdpProtocol) SetLocalFrameNumber(localFrame int) {
	remoteFrame := u.lastRecievedInput.Frame + (u.roundTripTime * 60 / 1000)
	u.localFrameAdvantage = remoteFrame - localFrame
}

func (u UdpProtocol) RecommendFrameDelay() int {
	return u.timesync.ReccomendFrameWaitDuration(false)
}

func (u *UdpProtocol) SetDisconnectTimeout(timeout int) {
	u.disconnectTimeout = uint(timeout)
}

func (u *UdpProtocol) SetDisconnectNotifyStart(timeout int) {
	u.disconnectNotifyStart = uint(timeout)
}

func (u *UdpProtocol) PumpSendQueue() {
	for !u.sendQueue.Empty() {
		entry := u.sendQueue.Front()

		if u.sendLatency > 0 {
			jitter := (u.sendLatency * 2 / 3) + ((rand.Int() % u.sendLatency) / 3)
			if int(time.Now().UnixMilli()) < u.sendQueue.Front().queueTime+jitter {
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
			Assert(entry.destIp != "")

			u.udp.SendTo(entry.msg)

			// would delete the udpmsg here
		}
		u.sendQueue.Pop()
	}
	if u.ooPacket.msg != nil && u.ooPacket.sendTime < int(time.Now().UnixMilli()) {
		log.Printf("sending rogue oop!")
		u.udp.SendTo(u.ooPacket.msg)
		// u.ooPacket = nil
	}
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
}

/* don't really need this ?
func (u *UdpProtocol) HandlesMsg(from string, msg *UdpMsg) bool {
	return u.peerAddress == from
}
*/
func (u *UdpProtocol) SendInput(input *GameInput) {
	if u.currentState == RunningState {
		// check to see if this is a good time to adjust for the rift
		u.timesync.AdvanceFrames(input, u.localFrameAdvantage, u.remoteFrameAdvantage)

		// Save this input packet.
		u.pendingOutput.Push(*input)
	}
	u.SendPendingOutput()
}

func (u *UdpProtocol) UpdateNetworkStats() {
	now := int(time.Now().UnixMilli())
	if u.statsStartTime == 0 {
		u.statsStartTime = now
	}

	totalBytesSent := u.bytesSent + (UDP_HEADER_SIZE * u.packetsSent)
	seconds := float64(now-u.statsStartTime) / 1000.0
	bps := float64(totalBytesSent / int(seconds))
	udpOverhead := float64(100.0 * (float64(UDP_HEADER_SIZE * u.packetsSent)) / float64(u.bytesSent))
	u.kbpsSent = int(bps / 1024)

	log.Printf("Network Stats -- Bandwidth: %.2f KBps Packets Sent: %5d (%.2f pps) KB Sent: %.2f UDP Overhead: %.2f %%.\n",
		float64(u.kbpsSent),
		u.packetsSent,
		float64(u.packetsSent*1000/(now-u.statsStartTime)),
		float64(totalBytesSent/1024.0),
		udpOverhead)
}

func (u *UdpProtocol) Synchronize() {
	u.currentState = SyncingState
	u.state.roundTripRemaining = uint32(NUM_SYNC_PACKETS)
	u.SendSyncRequest()
}

func (u *UdpProtocol) GetPeerConnectStatus(id int, frame *int) bool {
	*frame = u.peerConnectStatus[id].LastFrame
	// !u.peerConnectStatus[id].Disconnected from the C/++ world
	return u.peerConnectStatus[id].Disconnected == 0
}

func (u *UdpProtocol) OnInvalid(msg *UdpMsg, len int) bool {
	//  Assert(false) // ? ASSERT(FALSE && "Invalid msg in UdpProtocol");
	// ah
	log.Printf("Invalid msg in UdpProtocol ")
	Assert(false)
	return false
}

func (u *UdpProtocol) OnSyncRequest(msg *UdpMsg, len int) bool {
	if u.remoteMagicNumber != 0 && msg.Header.Magic != u.remoteMagicNumber {
		log.Printf("Ignoring sync request from unknown endpoint (%d != %d).\n",
			msg.Header.Magic, u.remoteMagicNumber)
		return false
	}
	reply := NewUdpMsg(SyncReplyMsg)
	reply.SyncReply.RandomReply = msg.SyncRequest.RandomRequest
	u.SendMsg(&reply)
	return true
}

func (u *UdpProtocol) OnMsg(msg *UdpMsg, length int) {
	handled := false

	type UdpProtocolDispatchFunc func(msg *UdpMsg, length int) bool

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
		if skipped > uint16(MAX_SEQ_DISTANCE) {
			log.Printf("dropping out of order packet (seq: %d, last seq:%d)\n", seq, u.nextRecvSeq)
			return
		}
	}

	u.nextRecvSeq = seq
	log.Printf("recv %s", msg)
	if int(msg.Header.HeaderType) > len(table) {
		u.OnInvalid(msg, length)
	} else {
		handled = table[int(msg.Header.HeaderType)](msg, length)
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

func (u *UdpProtocol) OnSyncReply(msg *UdpMsg, length int) bool {
	if u.currentState != SyncingState {
		log.Println("Ignoring SyncReply while not synching.")
		return msg.Header.Magic == u.remoteMagicNumber
	}

	if msg.SyncReply.RandomReply != u.state.random {
		log.Printf("sync reply %d != %d.  Keep looking...\n",
			msg.SyncReply.RandomReply, u.state.random)
		return false
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
		evt.total = NUM_SYNC_PACKETS
		evt.count = NUM_SYNC_PACKETS - int(u.state.roundTripRemaining)
		u.QueueEvent(&evt)
		u.SendSyncRequest()
	}

	return true
}
