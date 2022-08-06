package mocks

import (
	"fmt"
	"strconv"

	"github.com/assemblaj/GGPO-Go/internal/transport"
)

type FakeConnection struct {
	SendMap         map[string][]transport.UDPMessage
	LastSentMessage transport.UDPMessage
}

func NewFakeConnection() FakeConnection {
	return FakeConnection{
		SendMap: make(map[string][]transport.UDPMessage),
	}
}
func (f *FakeConnection) SendTo(msg transport.UDPMessage, remoteIp string, remotePort int) {
	portStr := strconv.Itoa(remotePort)
	addresssStr := remoteIp + ":" + portStr
	sendSlice, ok := f.SendMap[addresssStr]
	if !ok {
		sendSlice := make([]transport.UDPMessage, 2)
		f.SendMap[addresssStr] = sendSlice
	}
	sendSlice = append(sendSlice, msg)
	f.SendMap[addresssStr] = sendSlice
	f.LastSentMessage = msg
}

func (f *FakeConnection) Read() {

}

func (f *FakeConnection) Close() {

}

type FakeP2PConnection struct {
	remoteHandler   transport.MessageHandler
	localPort       int
	remotePort      int
	localIP         string
	printOutput     bool
	LastSentMessage transport.UDPMessage
	MessageHistory  []transport.UDPMessage
}

func (f *FakeP2PConnection) SendTo(msg transport.UDPMessage, remoteIp string, remotePort int) {
	if f.printOutput {
		fmt.Printf("f.localIP %s f.localPort %d msg %s size %d\n", f.localIP, f.localPort, msg, msg.PacketSize())
	}
	f.LastSentMessage = msg
	f.MessageHistory = append(f.MessageHistory, msg)
	f.remoteHandler.HandleMessage(f.localIP, f.localPort, msg, msg.PacketSize())
}

func (f *FakeP2PConnection) Read() {
}

func (f *FakeP2PConnection) Close() {

}

func NewFakeP2PConnection(remoteHandler transport.MessageHandler, localPort int, localIP string) FakeP2PConnection {
	f := FakeP2PConnection{}
	f.remoteHandler = remoteHandler
	f.localPort = localPort
	f.localIP = localIP
	f.MessageHistory = make([]transport.UDPMessage, 10)
	return f
}

type FakeMultiplePeerConnection struct {
	remoteHandler []transport.MessageHandler
	localPort     int
	localIP       string
	printOutput   bool
}

func (f *FakeMultiplePeerConnection) SendTo(msg transport.UDPMessage, remoteIp string, remotePort int) {
	if f.printOutput {
		fmt.Printf("f.localIP %s f.localPort %d msg %s size %d\n", f.localIP, f.localPort, msg, msg.PacketSize())
	}
	for _, r := range f.remoteHandler {
		r.HandleMessage(f.localIP, f.localPort, msg, msg.PacketSize())
	}
}

func (f *FakeMultiplePeerConnection) Read() {
}

func (f *FakeMultiplePeerConnection) Close() {

}

func NewFakeMultiplePeerConnection(remoteHandler []transport.MessageHandler, localPort int, localIP string) FakeMultiplePeerConnection {
	f := FakeMultiplePeerConnection{}
	f.remoteHandler = remoteHandler
	f.localPort = localPort
	f.localIP = localIP
	return f
}
