package ggthx_test

import (
	"strconv"
	"testing"
	"time"

	ggthx "github.com/assemblaj/ggthx/src"
)

type FakeConnection struct {
	sendMap         map[string][]ggthx.UDPMessage
	lastSentMessage ggthx.UDPMessage
}

func NewFakeConnection() FakeConnection {
	return FakeConnection{
		sendMap: make(map[string][]ggthx.UDPMessage),
	}
}
func (f *FakeConnection) SendTo(msg ggthx.UDPMessage, remoteIp string, remotePort int) {
	portStr := strconv.Itoa(remotePort)
	addresssStr := remoteIp + ":" + portStr
	sendSlice, ok := f.sendMap[addresssStr]
	if !ok {
		sendSlice := make([]ggthx.UDPMessage, 2)
		f.sendMap[addresssStr] = sendSlice
	}
	sendSlice = append(sendSlice, msg)
	f.sendMap[addresssStr] = sendSlice
	f.lastSentMessage = msg
}

func (f *FakeConnection) Read() {

}

func (f *FakeConnection) Close() {

}

type FakeMessageHandler struct {
	endpoint *ggthx.UdpProtocol
}

func (f *FakeMessageHandler) HandleMessage(ipAddress string, port int, msg ggthx.UDPMessage, length int) {
	if f.endpoint.HandlesMsg(ipAddress, port) {
		f.endpoint.OnMsg(msg, length)
	}
}

func NewFakeMessageHandler(endpoint *ggthx.UdpProtocol) FakeMessageHandler {
	f := FakeMessageHandler{}
	f.endpoint = endpoint
	return f
}

func TestMakeUDPProtocol(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
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
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	input := ggthx.GameInput{Bits: []byte{1, 2, 3, 4}}
	endpoint.SendInput(&input)
	portStr := strconv.Itoa(peerPort)
	messages, ok := connection.sendMap[peerAdress+":"+portStr]
	if ok != true {
		t.Errorf("The message was never sent. ")
	}
	inputPacket := messages[0].(*ggthx.InputPacket)
	got := inputPacket.Bits
	if got != nil {
		t.Errorf("expected '%#v' but got '%#v'", nil, got)
	}
}

func TestUDPProtocolSendMultipleInput(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	input := ggthx.GameInput{Size: 4, Bits: []byte{1, 2, 3, 4}}
	numInputs := 8
	for i := 0; i < numInputs; i++ {
		endpoint.SendInput(&input)
	}
	portStr := strconv.Itoa(peerPort)
	messages, ok := connection.sendMap[peerAdress+":"+portStr]
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
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	endpoint.Synchronize()
	portStr := strconv.Itoa(peerPort)
	messages, ok := connection.sendMap[peerAdress+":"+portStr]
	if ok != true {
		t.Errorf("The message was not sent. ")
	}

	syncPacket := messages[0].(*ggthx.SyncRequestPacket)
	if syncPacket.Header().HeaderType != uint8(ggthx.SyncRequestMsg) {
		t.Errorf("The message that was sent/recieved wsa not a SyncRequestMessage. ")
	}
}

func TestUDPProtocolSendInputAck(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	endpoint.SendInputAck()
	portStr := strconv.Itoa(peerPort)
	messages, ok := connection.sendMap[peerAdress+":"+portStr]
	if ok != true {
		t.Errorf("The message was not sent. ")
	}

	inputAckMessage := messages[0].(*ggthx.InputAckPacket)
	if inputAckMessage.Header().HeaderType != uint8(ggthx.InputAckMsg) {
		t.Errorf("The message that was sent/recieved wsa not a SyncRequestMessage. ")
	}
}

func TestUDPProtocolOnQualityReport(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	portStr := strconv.Itoa(peerPort)
	msg := ggthx.NewUDPMessage(ggthx.QualityReportMsg)
	qualityReportPacket := msg.(*ggthx.QualityReportPacket)
	qualityReportPacket.FrameAdvantage = 6
	qualityReportPacket.Ping = 50
	endpoint.OnQualityReport(qualityReportPacket, qualityReportPacket.PacketSize())
	messages, ok := connection.sendMap[peerAdress+":"+portStr]
	if ok != true {
		t.Errorf("The message was not sent. ")
	}

	qualityReplyPacket := messages[0].(*ggthx.QualityReplyPacket)
	if qualityReplyPacket.Header().HeaderType != uint8(ggthx.QualityReplyMsg) {
		t.Errorf("The message that was sent/recieved wsa not a SyncRequestMessage. ")
	}
}

func TestUDPProtocolOnSyncRequest(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	portStr := strconv.Itoa(peerPort)
	msg := ggthx.NewUDPMessage(ggthx.SyncRequestMsg)
	syncRequestPacket := msg.(*ggthx.SyncRequestPacket)

	endpoint.OnSyncRequest(syncRequestPacket, syncRequestPacket.PacketSize())

	messages, ok := connection.sendMap[peerAdress+":"+portStr]
	if ok != true {
		t.Errorf("The message was not sent. ")
	}

	syncReplyPacket := messages[0].(*ggthx.SyncReplyPacket)
	if syncReplyPacket.Header().HeaderType != uint8(ggthx.SyncReplyMsg) {
		t.Errorf("The message that was sent/recieved wsa not a SyncRequestMessage. ")
	}
}

func TestUDPProtocolGetPeerConnectStatus(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	var frame int32
	want := true
	got := endpoint.GetPeerConnectStatus(0, &frame)

	if want != got {
		t.Errorf("expected '%t' but got '%t'", want, got)
	}

	wantFrame := int32(ggthx.NullFrame)
	gotFrame := frame
	if wantFrame != gotFrame {
		t.Errorf("expected '%d' but got '%d'", wantFrame, gotFrame)
	}

}

func TestUDPProtocolHandlesMessage(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	want := true
	got := endpoint.HandlesMsg(peerAdress, peerPort)

	if want != got {
		t.Errorf("expected '%t' but got '%t'", want, got)
	}
}

func TestUDPProtocolHandlesMessageFalse(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	want := false
	got := endpoint.HandlesMsg("1.2.3.4", 0)

	if want != got {
		t.Errorf("expected '%t' but got '%t'", want, got)
	}
}

func TestUDPProtocolSetLocalFrameNumber(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	endpoint.SetLocalFrameNumber(8)
	var stats ggthx.NetworkStats
	endpoint.GetNetworkStats(&stats)
	want := -9 //
	got := stats.Timesync.LocalFramesBehind
	if want != got {
		t.Errorf("expected '%d' but got '%d'", want, got)
	}
}

func TestUDPProtocolOnQualityReply(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	msg := ggthx.NewUDPMessage(ggthx.QualityReplyMsg)
	qualityReplyPacket := msg.(*ggthx.QualityReplyPacket)
	qualityReplyPacket.Pong = 0
	checkInterval := 60
	endpoint.OnQualityReply(qualityReplyPacket, qualityReplyPacket.PacketSize())

	var stats ggthx.NetworkStats
	endpoint.GetNetworkStats(&stats)
	want := -9 //
	got := stats.Timesync.LocalFramesBehind
	now := int(time.Now().UnixMilli())
	if now-checkInterval > stats.Network.Ping {
		t.Errorf("expected '%d' but got '%d'", want, got)
	}
}

func TestUDPProtocolQueEventPanic(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	event := ggthx.UdpProtocolEvent{}
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
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	_, err := endpoint.GetEvent()
	if err == nil {
		t.Errorf("The program did not return an error when trying to get an event from an empty event queue.")
	}
}

func TestUDPProtocolGetEvent(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	want := ggthx.UdpProtocolEvent{}
	endpoint.QueueEvent(&want)
	got, _ := endpoint.GetEvent()
	if want.String() != got.String() {
		t.Errorf("expected '%s' but got '%s'", want, got)
	}
}

func TestUDPProtocolSyncchronize(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	endpoint.Synchronize()
	recvMessage := connection.lastSentMessage
	syncRequest := recvMessage.(*ggthx.SyncRequestPacket)

	msg := ggthx.NewUDPMessage(ggthx.SyncReplyMsg)
	syncReply := msg.(*ggthx.SyncReplyPacket)
	syncReply.RandomReply = syncRequest.RandomRequest
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		endpoint.OnSyncReply(syncReply, syncReply.PacketSize())
		syncRequest = connection.lastSentMessage.(*ggthx.SyncRequestPacket)
		syncReply.RandomReply = syncRequest.RandomRequest
	}

	evt, err := endpoint.GetEvent()
	if err != nil {
		t.Errorf("Got an error when should be recievving events")
	}

	if evt.Type() != ggthx.ConnectedEvent {
		t.Errorf("First popped event should be connected.")
	}
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		evt, err := endpoint.GetEvent()
		if err != nil {
			t.Errorf("Got an error when should be recievving events")
		}
		if i >= 0 && i < ggthx.NumSyncPackets-1 {
			if evt.Type() != ggthx.SynchronizingEvent {
				t.Errorf("These should be Synchronizing Events.")
			}
		} else if i == ggthx.NumSyncPackets-1 {
			if evt.Type() != ggthx.SynchronziedEvent {
				t.Errorf("This should be a Synchronized Event")
			}
		}
	}

}

func TestUDPProtocolOnLoopPoll(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	endpoint.Synchronize()
	recvMessage := connection.lastSentMessage
	syncRequest := recvMessage.(*ggthx.SyncRequestPacket)

	msg := ggthx.NewUDPMessage(ggthx.SyncReplyMsg)
	syncReply := msg.(*ggthx.SyncReplyPacket)
	syncReply.RandomReply = syncRequest.RandomRequest
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		endpoint.OnSyncReply(syncReply, syncReply.PacketSize())
		syncRequest = connection.lastSentMessage.(*ggthx.SyncRequestPacket)
		syncReply.RandomReply = syncRequest.RandomRequest
	}

	for i := 0; i < ggthx.NumSyncPackets+1; i++ {
		endpoint.GetEvent()
	}
	endpoint.OnLoopPoll(ggthx.DefaultTime)
	if connection.lastSentMessage.Type() != ggthx.QualityReportMsg {
		t.Errorf("This expected the OnLoopPoll to send a quality report message")
	}

}

func TestUDPProtocolHeartbeatGameInput(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	endpoint.Synchronize()
	recvMessage := connection.lastSentMessage
	syncRequest := recvMessage.(*ggthx.SyncRequestPacket)

	msg := ggthx.NewUDPMessage(ggthx.SyncReplyMsg)
	syncReply := msg.(*ggthx.SyncReplyPacket)
	syncReply.RandomReply = syncRequest.RandomRequest
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		endpoint.OnSyncReply(syncReply, syncReply.PacketSize())
		syncRequest = connection.lastSentMessage.(*ggthx.SyncRequestPacket)
		syncReply.RandomReply = syncRequest.RandomRequest
	}

	for i := 0; i < ggthx.NumSyncPackets+1; i++ {
		endpoint.GetEvent()
	}
	heartbeatTriggerInterval := 2
	for i := 0; i < heartbeatTriggerInterval; i++ {
		endpoint.OnLoopPoll(ggthx.DefaultTime)
	}

	if connection.lastSentMessage.Type() != ggthx.InputMsg {
		t.Errorf("This expected the OnLoopPoll to send a heartbeat game input")
	}

}

func TestUDPProtocolOnInputDefaultPanic(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	endpoint.Synchronize()
	recvMessage := connection.lastSentMessage
	syncRequest := recvMessage.(*ggthx.SyncRequestPacket)

	msg := ggthx.NewUDPMessage(ggthx.SyncReplyMsg)
	syncReply := msg.(*ggthx.SyncReplyPacket)
	syncReply.RandomReply = syncRequest.RandomRequest
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		endpoint.OnSyncReply(syncReply, syncReply.PacketSize())
		syncRequest = connection.lastSentMessage.(*ggthx.SyncRequestPacket)
		syncReply.RandomReply = syncRequest.RandomRequest
	}

	for i := 0; i < ggthx.NumSyncPackets+1; i++ {
		endpoint.GetEvent()
	}
	heartbeatTriggerInterval := 2
	for i := 0; i < heartbeatTriggerInterval; i++ {
		endpoint.OnLoopPoll(ggthx.DefaultTime)
	}

	if connection.lastSentMessage.Type() != ggthx.InputMsg {
		t.Errorf("This expected the OnLoopPoll to send a heartbeat game input")
	}
	msg = ggthx.NewUDPMessage(ggthx.InputMsg)
	inputPacket := msg.(*ggthx.InputPacket)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic when OnInput recieved a completely empty input packet.")
		}
	}()
	endpoint.OnInput(inputPacket, inputPacket.PacketSize())
}

func TestUDPProtocolOnInputPanicWithNonEqualConnectStatus(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	endpoint.Synchronize()
	recvMessage := connection.lastSentMessage
	syncRequest := recvMessage.(*ggthx.SyncRequestPacket)

	msg := ggthx.NewUDPMessage(ggthx.SyncReplyMsg)
	syncReply := msg.(*ggthx.SyncReplyPacket)
	syncReply.RandomReply = syncRequest.RandomRequest
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		endpoint.OnSyncReply(syncReply, syncReply.PacketSize())
		syncRequest = connection.lastSentMessage.(*ggthx.SyncRequestPacket)
		syncReply.RandomReply = syncRequest.RandomRequest
	}

	for i := 0; i < ggthx.NumSyncPackets+1; i++ {
		endpoint.GetEvent()
	}
	heartbeatTriggerInterval := 2
	for i := 0; i < heartbeatTriggerInterval; i++ {
		endpoint.OnLoopPoll(ggthx.DefaultTime)
	}

	if connection.lastSentMessage.Type() != ggthx.InputMsg {
		t.Errorf("This expected the OnLoopPoll to send a heartbeat game input")
	}
	msg = ggthx.NewUDPMessage(ggthx.InputMsg)
	inputPacket := msg.(*ggthx.InputPacket)
	inputPacket.PeerConnectStatus = connectStatus
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic when OnInput recieved a packet with number of connect statuses not equal to its own.")
		}
	}()
	endpoint.OnInput(inputPacket, inputPacket.PacketSize())
}

func TestUDPProtocolOnInputAfterSynchronizeCharacterization(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	endpoint.Synchronize()
	recvMessage := connection.lastSentMessage
	syncRequest := recvMessage.(*ggthx.SyncRequestPacket)

	msg := ggthx.NewUDPMessage(ggthx.SyncReplyMsg)
	syncReply := msg.(*ggthx.SyncReplyPacket)
	syncReply.RandomReply = syncRequest.RandomRequest
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		endpoint.OnSyncReply(syncReply, syncReply.PacketSize())
		syncRequest = connection.lastSentMessage.(*ggthx.SyncRequestPacket)
		syncReply.RandomReply = syncRequest.RandomRequest
	}

	for i := 0; i < ggthx.NumSyncPackets+1; i++ {
		endpoint.GetEvent()
	}
	heartbeatTriggerInterval := 2
	for i := 0; i < heartbeatTriggerInterval; i++ {
		endpoint.OnLoopPoll(ggthx.DefaultTime)
	}

	if connection.lastSentMessage.Type() != ggthx.InputMsg {
		t.Errorf("This expected the OnLoopPoll to send a heartbeat game input")
	}
	msg = ggthx.NewUDPMessage(ggthx.InputMsg)
	inputPacket := msg.(*ggthx.InputPacket)
	inputPacket.PeerConnectStatus = make([]ggthx.UdpConnectStatus, 4)
	inputPacket.Bits = []byte{1, 2, 3, 4}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic when OnInput recieved a packet without its imput size set")
		}
	}()
	endpoint.OnInput(inputPacket, inputPacket.PacketSize())
}

func TestUDPProtocolOnInputAfterSynchronize(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)

	endpoint.Synchronize()
	recvMessage := connection.lastSentMessage
	syncRequest := recvMessage.(*ggthx.SyncRequestPacket)

	msg := ggthx.NewUDPMessage(ggthx.SyncReplyMsg)
	syncReply := msg.(*ggthx.SyncReplyPacket)
	syncReply.RandomReply = syncRequest.RandomRequest
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		endpoint.OnSyncReply(syncReply, syncReply.PacketSize())
		syncRequest = connection.lastSentMessage.(*ggthx.SyncRequestPacket)
		syncReply.RandomReply = syncRequest.RandomRequest
	}

	for i := 0; i < ggthx.NumSyncPackets+1; i++ {
		endpoint.GetEvent()
	}
	heartbeatTriggerInterval := 2
	for i := 0; i < heartbeatTriggerInterval; i++ {
		endpoint.OnLoopPoll(ggthx.DefaultTime)
	}

	if connection.lastSentMessage.Type() != ggthx.InputMsg {
		t.Errorf("This expected the OnLoopPoll to send a heartbeat game input")
	}
	msg = ggthx.NewUDPMessage(ggthx.InputMsg)
	inputPacket := msg.(*ggthx.InputPacket)
	inputPacket.PeerConnectStatus = make([]ggthx.UdpConnectStatus, 4)
	inputPacket.Bits = []byte{1, 2, 3, 4}
	inputPacket.InputSize = 4
	endpoint.OnInput(inputPacket, inputPacket.PacketSize())
	evt, err := endpoint.GetEvent()
	if err != nil {
		t.Errorf("Expected there to be game input event, not error.")
	}
	if evt.Type() != ggthx.InputEvent {
		t.Errorf("Expected the event to be InputEvent, not %s", evt)
	}
}

func TestUDPProtocolFakeP2PandMessageHandler(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	peerAdress := "127.2.1.1"
	peerPort := 7001
	port2 := 7000
	f := FakeMessageHandler{}
	f2 := FakeMessageHandler{}
	connection := NewFakeP2PConnection(&f, peerPort, peerAdress)
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, port2, &connectStatus)

	connection2 := NewFakeP2PConnection(&f2, port2, peerAdress)
	endpoint2 := ggthx.NewUdpProtocol(&connection2, 0, peerAdress, peerPort, &connectStatus)
	f2.endpoint = &endpoint
	f.endpoint = &endpoint2

	//ggthx.EnableLogger()
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
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	endpoint.Disconnect()
	if endpoint.IsRunning() {
		t.Errorf("The endpoint should be disconnected after running the disconnect method.")
	}
}

func TestUDPProtocolIsInitalized(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(nil, 0, peerAdress, peerPort, &connectStatus)
	if endpoint.IsInitialized() {
		t.Errorf("The endpoint should not be initialized if connection is nil.")
	}
}
func TestUDPProtocolOnInvalid(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	invalidMessageType := 88
	msg := ggthx.NewUDPMessage(ggthx.UDPMessageType(invalidMessageType))
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
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	endpoint.SendPendingOutput()
	msg := connection.lastSentMessage
	inputPacket := msg.(*ggthx.InputPacket)
	if inputPacket.StartFrame != 0 {
		t.Errorf("Inputs sent when there's no pending output should have startframe 0 ")
	}
	if inputPacket.InputSize != 0 {
		t.Errorf("Inputs sent when there's no pending output should have Input size 0 ")
	}
}
func TestUDPProtocolMagicNumberReject(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	msg := ggthx.NewUDPMessage(ggthx.QualityReportMsg)
	msg.SetHeader(3, 0)
	endpoint.OnMsg(msg, msg.PacketSize())
	if connection.lastSentMessage != nil {
		t.Errorf("No messages should have been sent in response to the quality report message because the magic number doesn't match the default magic number (0)")
	}
}

func TestUDPProtocolSequenceNumberReject(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	connection := NewFakeConnection()
	peerAdress := "127.2.1.1"
	peerPort := 7001
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, peerPort, &connectStatus)
	msg := ggthx.NewUDPMessage(ggthx.QualityReportMsg)
	msg.SetHeader(0, ggthx.MaxSeqDistance+1)
	endpoint.OnMsg(msg, msg.PacketSize())
	if connection.lastSentMessage != nil {
		t.Errorf("No messages should have been sent in response to the quality report message because of the sequence number. ")
	}
}

func TestUDPProtocolKeepAlive(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	peerAdress := "127.2.1.1"
	peerPort := 7001
	port2 := 7000
	f := FakeMessageHandler{}
	f2 := FakeMessageHandler{}
	connection := NewFakeP2PConnection(&f, peerPort, peerAdress)
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, port2, &connectStatus)

	connection2 := NewFakeP2PConnection(&f2, port2, peerAdress)
	endpoint2 := ggthx.NewUdpProtocol(&connection2, 0, peerAdress, peerPort, &connectStatus)
	f2.endpoint = &endpoint
	f.endpoint = &endpoint2

	//ggthx.EnableLogger()
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
	if connection.lastSentMessage.Header().HeaderType != uint8(ggthx.KeepAliveMsg) {
		t.Errorf("Endpoint should've sent keep alive packet.")
	}
}

func TestUDPProtocolHeartBeatCharacterization(t *testing.T) {
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	peerAdress := "127.2.1.1"
	peerPort := 7001
	port2 := 7000
	f := FakeMessageHandler{}
	f2 := FakeMessageHandler{}
	connection := NewFakeP2PConnection(&f, peerPort, peerAdress)
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, port2, &connectStatus)

	connection2 := NewFakeP2PConnection(&f2, port2, peerAdress)
	endpoint2 := ggthx.NewUdpProtocol(&connection2, 0, peerAdress, peerPort, &connectStatus)
	f2.endpoint = &endpoint
	f.endpoint = &endpoint2

	//ggthx.EnableLogger()
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
	connectStatus := []ggthx.UdpConnectStatus{
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
		{Disconnected: false, LastFrame: 20},
		{Disconnected: false, LastFrame: 22},
	}
	peerAdress := "127.2.1.1"
	peerPort := 7001
	port2 := 7000
	f := FakeMessageHandler{}
	f2 := FakeMessageHandler{}
	connection := NewFakeP2PConnection(&f, peerPort, peerAdress)
	endpoint := ggthx.NewUdpProtocol(&connection, 0, peerAdress, port2, &connectStatus)

	connection2 := NewFakeP2PConnection(&f2, port2, peerAdress)
	endpoint2 := ggthx.NewUdpProtocol(&connection2, 0, peerAdress, peerPort, &connectStatus)
	f2.endpoint = &endpoint
	f.endpoint = &endpoint2

	//ggthx.EnableLogger()
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
	e1k := connection.messageHistory[len(connection.messageHistory)-1]
	e1i := connection.messageHistory[len(connection.messageHistory)-2]
	e2k := connection2.messageHistory[len(connection2.messageHistory)-1]
	e2i := connection2.messageHistory[len(connection2.messageHistory)-2]
	if e1k.Header().HeaderType != uint8(ggthx.KeepAliveMsg) {
		t.Errorf("Endpoint 1 should've sent a keep alive msg")
	}

	if e1i.Header().HeaderType != uint8(ggthx.InputMsg) {
		t.Errorf("Endpoint 1 should've sent a heartbeat input prior to the keep alive message")
	}

	if e2k.Header().HeaderType != uint8(ggthx.KeepAliveMsg) {
		t.Errorf("Endpoint 2 should've sent a keep alive msg")
	}

	if e2i.Header().HeaderType != uint8(ggthx.InputMsg) {
		t.Errorf("Endpoint 2 should've sent a heartbeat input prior to the keep alive message")
	}

}
