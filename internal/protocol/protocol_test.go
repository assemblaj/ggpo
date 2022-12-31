package protocol_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/assemblaj/GGPO-Go/internal/input"
	"github.com/assemblaj/GGPO-Go/internal/messages"
	"github.com/assemblaj/GGPO-Go/internal/mocks"
	"github.com/assemblaj/GGPO-Go/internal/polling"
	"github.com/assemblaj/GGPO-Go/internal/protocol"
)

func TestMakeUDPProtocol(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	if !endpoint.IsInitialized() {
		t.Errorf("The fake connection wasn't properly saved.")
	}
}

/*
	Sends.
*/
/*
	Characterization dunno why it works this way
*/
func TestUDPProtocolSendInput(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	input := input.GameInput{Bits: []byte{1, 2, 3, 4}}
	endpoint.SendInput(&input)
	portStr := strconv.Itoa(peerPort)
	msgs, ok := connection.SendMap[peerAdress+":"+portStr]
	if ok != true {
		t.Errorf("The message was never sent. ")
	}
	inputPacket := msgs[0].(*messages.InputPacket)
	got := inputPacket.Bits
	if got != nil {
		t.Errorf("expected '%#v' but got '%#v'", nil, got)
	}
}

func TestUDPProtocolSendMultipleInput(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	input := input.GameInput{Size: 4, Bits: []byte{1, 2, 3, 4}}
	numInputs := 8
	for i := 0; i < numInputs; i++ {
		endpoint.SendInput(&input)
	}
	portStr := strconv.Itoa(peerPort)
	messages, ok := connection.SendMap[peerAdress+":"+portStr]
	if ok != true {
		t.Errorf("The messages were never sent. ")
	}
	want := numInputs
	got := len(messages)
	if len(messages) != numInputs {
		t.Errorf("expected '%d' but got '%d'", want, got)
	}
}

func TestUDPProtocolSynchronize(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	endpoint.Synchronize()
	portStr := strconv.Itoa(peerPort)
	msgs, ok := connection.SendMap[peerAdress+":"+portStr]
	if ok != true {
		t.Errorf("The message was not sent. ")
	}

	syncPacket := msgs[0].(*messages.SyncRequestPacket)
	if syncPacket.Header().HeaderType != uint8(messages.SyncRequestMsg) {
		t.Errorf("The message that was sent/recieved wsa not a SyncRequestMessage. ")
	}
}

func TestUDPProtocolSendInputAck(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	endpoint.SendInputAck()
	portStr := strconv.Itoa(peerPort)
	msgs, ok := connection.SendMap[peerAdress+":"+portStr]
	if ok != true {
		t.Errorf("The message was not sent. ")
	}

	inputAckMessage := msgs[0].(*messages.InputAckPacket)
	if inputAckMessage.Header().HeaderType != uint8(messages.InputAckMsg) {
		t.Errorf("The message that was sent/recieved wsa not a SyncRequestMessage. ")
	}
}

func TestUDPProtocolOnQualityReport(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	portStr := strconv.Itoa(peerPort)
	msg := messages.NewUDPMessage(messages.QualityReportMsg)
	qualityReportPacket := msg.(*messages.QualityReportPacket)
	qualityReportPacket.FrameAdvantage = 6
	qualityReportPacket.Ping = 50
	endpoint.OnQualityReport(qualityReportPacket, qualityReportPacket.PacketSize())
	msgs, ok := connection.SendMap[peerAdress+":"+portStr]
	if ok != true {
		t.Errorf("The message was not sent. ")
	}

	qualityReplyPacket := msgs[0].(*messages.QualityReplyPacket)
	if qualityReplyPacket.Header().HeaderType != uint8(messages.QualityReplyMsg) {
		t.Errorf("The message that was sent/recieved wsa not a SyncRequestMessage. ")
	}
}

func TestUDPProtocolOnSyncRequest(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	portStr := strconv.Itoa(peerPort)
	msg := messages.NewUDPMessage(messages.SyncRequestMsg)
	syncRequestPacket := msg.(*messages.SyncRequestPacket)

	endpoint.OnSyncRequest(syncRequestPacket, syncRequestPacket.PacketSize())

	msgs, ok := connection.SendMap[peerAdress+":"+portStr]
	if ok != true {
		t.Errorf("The message was not sent. ")
	}

	syncReplyPacket := msgs[0].(*messages.SyncReplyPacket)
	if syncReplyPacket.Header().HeaderType != uint8(messages.SyncReplyMsg) {
		t.Errorf("The message that was sent/recieved wsa not a SyncRequestMessage. ")
	}
}

func TestUDPProtocolGetPeerConnectStatus(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	var frame int32
	want := true
	got := endpoint.GetPeerConnectStatus(0, &frame)

	if want != got {
		t.Errorf("expected '%t' but got '%t'", want, got)
	}

	wantFrame := int32(input.NullFrame)
	gotFrame := frame
	if wantFrame != gotFrame {
		t.Errorf("expected '%d' but got '%d'", wantFrame, gotFrame)
	}

}

func TestUDPProtocolHandlesMessage(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	want := true
	got := endpoint.HandlesMsg(peerAdress, peerPort)

	if want != got {
		t.Errorf("expected '%t' but got '%t'", want, got)
	}
}

func TestUDPProtocolHandlesMessageFalse(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	want := false
	got := endpoint.HandlesMsg("1.2.3.4", 0)

	if want != got {
		t.Errorf("expected '%t' but got '%t'", want, got)
	}
}

func TestUDPProtocolSetLocalFrameNumber(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	endpoint.SetLocalFrameNumber(8)
	stats := endpoint.GetNetworkStats()
	want := float32(-9.0) //
	got := stats.Timesync.LocalFramesBehind
	if want != got {
		t.Errorf("expected '%f' but got '%f'", want, got)
	}
}

func TestUDPProtocolOnQualityReply(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	msg := messages.NewUDPMessage(messages.QualityReplyMsg)
	qualityReplyPacket := msg.(*messages.QualityReplyPacket)
	qualityReplyPacket.Pong = 0
	checkInterval := 60
	endpoint.OnQualityReply(qualityReplyPacket, qualityReplyPacket.PacketSize())

	var stats protocol.NetworkStats
	stats = endpoint.GetNetworkStats()
	want := -9 //
	got := stats.Timesync.LocalFramesBehind
	now := int(time.Now().UnixMilli())
	if now-checkInterval > stats.Network.Ping {
		t.Errorf("expected '%d' but got '%d'", want, got)
	}
}

func TestUDPProtocolQueEventPanic(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	event := protocol.UdpProtocolEvent{}
	capcity := 64
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic when QueEvent attempted to add an event higher than the capacity.")
		}
	}()
	for i := 0; i < capcity+1; i++ {
		endpoint.QueueEvent(&event)
	}
}

func TestUDPProtocolGetEventError(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	_, err := endpoint.GetEvent()
	if err == nil {
		t.Errorf("The program did not return an error when trying to get an event from an empty event queue.")
	}
}

func TestUDPProtocolGetEvent(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	want := protocol.UdpProtocolEvent{}
	endpoint.QueueEvent(&want)
	got, _ := endpoint.GetEvent()
	if want.String() != got.String() {
		t.Errorf("expected '%s' but got '%s'", want, got)
	}
}

func TestUDPProtocolSyncchronize(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	endpoint.Synchronize()
	recvMessage := connection.LastSentMessage
	syncRequest := recvMessage.(*messages.SyncRequestPacket)

	msg := messages.NewUDPMessage(messages.SyncReplyMsg)
	syncReply := msg.(*messages.SyncReplyPacket)
	syncReply.RandomReply = syncRequest.RandomRequest
	for i := 0; i < protocol.NumSyncPackets; i++ {
		endpoint.OnSyncReply(syncReply, syncReply.PacketSize())
		syncRequest = connection.LastSentMessage.(*messages.SyncRequestPacket)
		syncReply.RandomReply = syncRequest.RandomRequest
	}

	evt, err := endpoint.GetEvent()
	if err != nil {
		t.Errorf("Got an error when should be recievving events")
	}

	if evt.Type() != protocol.ConnectedEvent {
		t.Errorf("First popped event should be connected.")
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		evt, err := endpoint.GetEvent()
		if err != nil {
			t.Errorf("Got an error when should be recievving events")
		}
		if i >= 0 && i < protocol.NumSyncPackets-1 {
			if evt.Type() != protocol.SynchronizingEvent {
				t.Errorf("These should be Synchronizing Events.")
			}
		} else if i == protocol.NumSyncPackets-1 {
			if evt.Type() != protocol.SynchronziedEvent {
				t.Errorf("This should be a Synchronized Event")
			}
		}
	}

}

func TestUDPProtocolOnLoopPoll(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	endpoint.Synchronize()
	recvMessage := connection.LastSentMessage
	syncRequest := recvMessage.(*messages.SyncRequestPacket)

	msg := messages.NewUDPMessage(messages.SyncReplyMsg)
	syncReply := msg.(*messages.SyncReplyPacket)
	syncReply.RandomReply = syncRequest.RandomRequest
	for i := 0; i < protocol.NumSyncPackets; i++ {
		endpoint.OnSyncReply(syncReply, syncReply.PacketSize())
		syncRequest = connection.LastSentMessage.(*messages.SyncRequestPacket)
		syncReply.RandomReply = syncRequest.RandomRequest
	}

	for i := 0; i < protocol.NumSyncPackets+1; i++ {
		endpoint.GetEvent()
	}
	endpoint.OnLoopPoll(polling.DefaultTime)
	if connection.LastSentMessage.Type() != messages.QualityReportMsg {
		t.Errorf("This expected the OnLoopPoll to send a quality report message")
	}

}

func TestUDPProtocolHeartbeatGameInput(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	endpoint.Synchronize()
	recvMessage := connection.LastSentMessage
	syncRequest := recvMessage.(*messages.SyncRequestPacket)

	msg := messages.NewUDPMessage(messages.SyncReplyMsg)
	syncReply := msg.(*messages.SyncReplyPacket)
	syncReply.RandomReply = syncRequest.RandomRequest
	for i := 0; i < protocol.NumSyncPackets; i++ {
		endpoint.OnSyncReply(syncReply, syncReply.PacketSize())
		syncRequest = connection.LastSentMessage.(*messages.SyncRequestPacket)
		syncReply.RandomReply = syncRequest.RandomRequest
	}

	for i := 0; i < protocol.NumSyncPackets+1; i++ {
		endpoint.GetEvent()
	}
	heartbeatTriggerInterval := 2
	for i := 0; i < heartbeatTriggerInterval; i++ {
		endpoint.OnLoopPoll(polling.DefaultTime)
	}

	if connection.LastSentMessage.Type() != messages.InputMsg {
		t.Errorf("This expected the OnLoopPoll to send a heartbeat game input")
	}

}

func TestUDPProtocolOnInputDefaultPanic(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	endpoint.Synchronize()
	recvMessage := connection.LastSentMessage
	syncRequest := recvMessage.(*messages.SyncRequestPacket)

	msg := messages.NewUDPMessage(messages.SyncReplyMsg)
	syncReply := msg.(*messages.SyncReplyPacket)
	syncReply.RandomReply = syncRequest.RandomRequest
	for i := 0; i < protocol.NumSyncPackets; i++ {
		endpoint.OnSyncReply(syncReply, syncReply.PacketSize())
		syncRequest = connection.LastSentMessage.(*messages.SyncRequestPacket)
		syncReply.RandomReply = syncRequest.RandomRequest
	}

	for i := 0; i < protocol.NumSyncPackets+1; i++ {
		endpoint.GetEvent()
	}
	heartbeatTriggerInterval := 2
	for i := 0; i < heartbeatTriggerInterval; i++ {
		endpoint.OnLoopPoll(polling.DefaultTime)
	}

	if connection.LastSentMessage.Type() != messages.InputMsg {
		t.Errorf("This expected the OnLoopPoll to send a heartbeat game input")
	}
	msg = messages.NewUDPMessage(messages.InputMsg)
	inputPacket := msg.(*messages.InputPacket)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic when OnInput recieved a completely empty input packet.")
		}
	}()
	endpoint.OnInput(inputPacket, inputPacket.PacketSize())
}

func TestUDPProtocolOnInputPanicWithNonEqualConnectStatus(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	endpoint.Synchronize()
	recvMessage := connection.LastSentMessage
	syncRequest := recvMessage.(*messages.SyncRequestPacket)

	msg := messages.NewUDPMessage(messages.SyncReplyMsg)
	syncReply := msg.(*messages.SyncReplyPacket)
	syncReply.RandomReply = syncRequest.RandomRequest
	for i := 0; i < protocol.NumSyncPackets; i++ {
		endpoint.OnSyncReply(syncReply, syncReply.PacketSize())
		syncRequest = connection.LastSentMessage.(*messages.SyncRequestPacket)
		syncReply.RandomReply = syncRequest.RandomRequest
	}

	for i := 0; i < protocol.NumSyncPackets+1; i++ {
		endpoint.GetEvent()
	}
	heartbeatTriggerInterval := 2
	for i := 0; i < heartbeatTriggerInterval; i++ {
		endpoint.OnLoopPoll(polling.DefaultTime)
	}

	if connection.LastSentMessage.Type() != messages.InputMsg {
		t.Errorf("This expected the OnLoopPoll to send a heartbeat game input")
	}
	msg = messages.NewUDPMessage(messages.InputMsg)
	inputPacket := msg.(*messages.InputPacket)
	inputPacket.PeerConnectStatus = connectStatus
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic when OnInput recieved a packet with number of connect statuses not equal to its own.")
		}
	}()
	endpoint.OnInput(inputPacket, inputPacket.PacketSize())
}

func TestUDPProtocolOnInputAfterSynchronizeCharacterization(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	endpoint.Synchronize()
	recvMessage := connection.LastSentMessage
	syncRequest := recvMessage.(*messages.SyncRequestPacket)

	msg := messages.NewUDPMessage(messages.SyncReplyMsg)
	syncReply := msg.(*messages.SyncReplyPacket)
	syncReply.RandomReply = syncRequest.RandomRequest
	for i := 0; i < protocol.NumSyncPackets; i++ {
		endpoint.OnSyncReply(syncReply, syncReply.PacketSize())
		syncRequest = connection.LastSentMessage.(*messages.SyncRequestPacket)
		syncReply.RandomReply = syncRequest.RandomRequest
	}

	for i := 0; i < protocol.NumSyncPackets+1; i++ {
		endpoint.GetEvent()
	}
	heartbeatTriggerInterval := 2
	for i := 0; i < heartbeatTriggerInterval; i++ {
		endpoint.OnLoopPoll(polling.DefaultTime)
	}

	if connection.LastSentMessage.Type() != messages.InputMsg {
		t.Errorf("This expected the OnLoopPoll to send a heartbeat game input")
	}
	msg = messages.NewUDPMessage(messages.InputMsg)
	inputPacket := msg.(*messages.InputPacket)
	inputPacket.PeerConnectStatus = make([]messages.UdpConnectStatus, 4)
	inputPacket.Bits = []byte{1, 2, 3, 4}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic when OnInput recieved a packet without its imput size set")
		}
	}()
	endpoint.OnInput(inputPacket, inputPacket.PacketSize())
}

func TestUDPProtocolOnInputAfterSynchronize(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	endpoint.Synchronize()
	recvMessage := connection.LastSentMessage
	syncRequest := recvMessage.(*messages.SyncRequestPacket)

	msg := messages.NewUDPMessage(messages.SyncReplyMsg)
	syncReply := msg.(*messages.SyncReplyPacket)
	syncReply.RandomReply = syncRequest.RandomRequest
	for i := 0; i < protocol.NumSyncPackets; i++ {
		endpoint.OnSyncReply(syncReply, syncReply.PacketSize())
		syncRequest = connection.LastSentMessage.(*messages.SyncRequestPacket)
		syncReply.RandomReply = syncRequest.RandomRequest
	}

	for i := 0; i < protocol.NumSyncPackets+1; i++ {
		endpoint.GetEvent()
	}
	heartbeatTriggerInterval := 2
	for i := 0; i < heartbeatTriggerInterval; i++ {
		endpoint.OnLoopPoll(polling.DefaultTime)
	}

	if connection.LastSentMessage.Type() != messages.InputMsg {
		t.Errorf("This expected the OnLoopPoll to send a heartbeat game input")
	}
	msg = messages.NewUDPMessage(messages.InputMsg)
	inputPacket := msg.(*messages.InputPacket)
	inputPacket.PeerConnectStatus = make([]messages.UdpConnectStatus, 4)
	inputPacket.Bits = []byte{1, 2, 3, 4}
	inputPacket.InputSize = 4
	endpoint.OnInput(inputPacket, inputPacket.PacketSize())
	evt, err := endpoint.GetEvent()
	if err != nil {
		t.Errorf("Expected there to be game input event, not error.")
	}
	if evt.Type() != protocol.InputEvent {
		t.Errorf("Expected the event to be InputEvent, not %s", evt)
	}
}

func TestUDPProtocolFakeP2PandMessageHandler(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	peerAdress := "127.2.1.1"
	peerPort := 7001
	port2 := 7000
	f := mocks.FakeMessageHandler{}
	f2 := mocks.FakeMessageHandler{}
	connection := mocks.NewFakeP2PConnection(&f, peerPort, peerAdress)
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, port2, &connectStatus)

	connection2 := mocks.NewFakeP2PConnection(&f2, port2, peerAdress)
	endpoint2 := protocol.NewUdpProtocol(&connection2, 0, peerAdress, peerPort, &connectStatus)
	f2.Endpoint = &endpoint
	f.Endpoint = &endpoint2

	//ggpo.EnableLogger()
	endpoint.Synchronize()
	endpoint2.Synchronize()

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 11000).UnixMilli()
	}
	for !endpoint2.IsSynchronized() {
		endpoint2.OnLoopPoll(advance)
	}

	if !endpoint.IsSynchronized() {
		t.Errorf("First endpoint never synchronized.")
	}

	if !endpoint2.IsSynchronized() {
		t.Errorf("Second endpoint never synchronized.")
	}
}

func TestUDPProtocolDiscconect(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	endpoint.Disconnect()
	if endpoint.IsRunning() {
		t.Errorf("The endpoint should be disconnected after running the disconnect method.")
	}
}

func TestUDPProtocolDiscconectOnLoopPoll(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	endpoint.Disconnect()
	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 8000).UnixMilli()
	}

	endpoint.OnLoopPoll(advance)
	if endpoint.IsInitialized() {
		t.Errorf("Disconnected endpoints should not still be 'initalized' after they've been polled. ")
	}
}

func TestUDPProtocolOnInputDisconnectedRequest(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	/*
		advance := func() int64 {
			return time.Now().Add(time.Millisecond * 8000).UnixMilli()
		}*/
	msg := messages.NewUDPMessage(messages.InputMsg)
	inputPacket := msg.(*messages.InputPacket)
	inputPacket.DisconectRequested = true
	endpoint.OnInput(inputPacket, inputPacket.PacketSize())
	evt, _ := endpoint.GetEvent()
	if evt.Type() != protocol.DisconnectedEvent {
		t.Errorf("Recieving an input with DisconnectRequested = true should create a DisconnectedEvent")
	}
}

func TestUDPProtocolIsInitalized(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(nil, 0, peerAdress, peerPort, &connectStatus)
	if endpoint.IsInitialized() {
		t.Errorf("The endpoint should not be initialized if connection is nil.")
	}
}
func TestUDPProtocolOnInvalid(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	invalidMessageType := 88
	msg := messages.NewUDPMessage(messages.UDPMessageType(invalidMessageType))
	endpoint.OnMsg(msg, msg.PacketSize())
	handled, err := endpoint.OnInvalid(msg, msg.PacketSize())
	if handled == true {
		t.Errorf("Invalid messasge shouldn't be handled")
	}
	if err == nil {
		t.Errorf("Invalid message should create an error")
	}
}

func TestUDPProtocolSendPendingOutputDefault(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	endpoint.SendPendingOutput()
	msg := connection.LastSentMessage
	inputPacket := msg.(*messages.InputPacket)
	if inputPacket.StartFrame != 0 {
		t.Errorf("Inputs sent when there's no pending output should have startframe 0 ")
	}
	if inputPacket.InputSize != 0 {
		t.Errorf("Inputs sent when there's no pending output should have Input size 0 ")
	}
}
func TestUDPProtocolMagicNumberReject(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	msg := messages.NewUDPMessage(messages.QualityReportMsg)
	msg.SetHeader(3, 0)
	endpoint.OnMsg(msg, msg.PacketSize())
	if connection.LastSentMessage != nil {
		t.Errorf("No messages should have been sent in response to the quality report message because the magic number doesn't match the default magic number (0)")
	}
}

func TestUDPProtocolSequenceNumberReject(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := mocks.NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	msg := messages.NewUDPMessage(messages.QualityReportMsg)
	msg.SetHeader(0, protocol.MaxSeqDistance+1)
	endpoint.OnMsg(msg, msg.PacketSize())
	if connection.LastSentMessage != nil {
		t.Errorf("No messages should have been sent in response to the quality report message because of the sequence number. ")
	}
}

func TestUDPProtocolKeepAlive(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	peerAdress := "127.2.1.1"
	peerPort := 7001
	port2 := 7000
	f := mocks.FakeMessageHandler{}
	f2 := mocks.FakeMessageHandler{}
	connection := mocks.NewFakeP2PConnection(&f, peerPort, peerAdress)
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, port2, &connectStatus)

	connection2 := mocks.NewFakeP2PConnection(&f2, port2, peerAdress)
	endpoint2 := protocol.NewUdpProtocol(&connection2, 0, peerAdress, peerPort, &connectStatus)
	f2.Endpoint = &endpoint
	f.Endpoint = &endpoint2

	//ggpo.EnableLogger()
	endpoint.Synchronize()
	endpoint2.Synchronize()

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 11000).UnixMilli()
	}
	for !endpoint2.IsSynchronized() {
		endpoint2.OnLoopPoll(advance)
	}

	endpoint.OnLoopPoll(advance)
	endpoint.OnLoopPoll(advance)
	if connection.LastSentMessage.Header().HeaderType != uint8(messages.KeepAliveMsg) {
		t.Errorf("Endpoint should've sent keep alive packet.")
	}
}

func TestUDPProtocolHeartBeatCharacterization(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	peerAdress := "127.2.1.1"
	peerPort := 7001
	port2 := 7000
	f := mocks.FakeMessageHandler{}
	f2 := mocks.FakeMessageHandler{}
	connection := mocks.NewFakeP2PConnection(&f, peerPort, peerAdress)
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, port2, &connectStatus)

	connection2 := mocks.NewFakeP2PConnection(&f2, port2, peerAdress)
	endpoint2 := protocol.NewUdpProtocol(&connection2, 0, peerAdress, peerPort, &connectStatus)
	f2.Endpoint = &endpoint
	f.Endpoint = &endpoint2

	//ggpo.EnableLogger()
	endpoint.Synchronize()
	endpoint2.Synchronize()

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 1000).UnixMilli()
	}
	for !endpoint2.IsSynchronized() {
		endpoint2.OnLoopPoll(advance)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic when OnInput recieved a connection status with length < 4")
		}
	}()

	for i := 0; i < 10; i++ {
		endpoint2.OnLoopPoll(advance)
	}
}

func TestUDPProtocolHeartBeat(t *testing.T) {
	connectStatus := []messages.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	peerAdress := "127.2.1.1"
	peerPort := 7001
	port2 := 7000
	f := mocks.FakeMessageHandler{}
	f2 := mocks.FakeMessageHandler{}
	connection := mocks.NewFakeP2PConnection(&f, peerPort, peerAdress)
	endpoint := protocol.NewUdpProtocol(&connection, 0, peerAdress, port2, &connectStatus)

	connection2 := mocks.NewFakeP2PConnection(&f2, port2, peerAdress)
	endpoint2 := protocol.NewUdpProtocol(&connection2, 0, peerAdress, peerPort, &connectStatus)
	f2.Endpoint = &endpoint
	f.Endpoint = &endpoint2

	//ggpo.EnableLogger()
	endpoint.Synchronize()
	endpoint2.Synchronize()

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 1000).UnixMilli()
	}
	for !endpoint2.IsSynchronized() {
		endpoint2.OnLoopPoll(advance)
	}

	for i := 0; i < 20; i++ {
		endpoint2.OnLoopPoll(advance)
		endpoint.OnLoopPoll(advance)
	}
	e1k := connection.MessageHistory[len(connection.MessageHistory)-1]
	e1i := connection.MessageHistory[len(connection.MessageHistory)-2]
	e2k := connection2.MessageHistory[len(connection2.MessageHistory)-1]
	e2i := connection2.MessageHistory[len(connection2.MessageHistory)-2]
	if e1k.Header().HeaderType != uint8(messages.KeepAliveMsg) {
		t.Errorf("Endpoint 1 should've sent a keep alive msg")
	}

	if e1i.Header().HeaderType != uint8(messages.InputMsg) {
		t.Errorf("Endpoint 1 should've sent a heartbeat input prior to the keep alive message")
	}

	if e2k.Header().HeaderType != uint8(messages.KeepAliveMsg) {
		t.Errorf("Endpoint 2 should've sent a keep alive msg")
	}

	if e2i.Header().HeaderType != uint8(messages.InputMsg) {
		t.Errorf("Endpoint 2 should've sent a heartbeat input prior to the keep alive message")
	}

}
