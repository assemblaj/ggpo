package mocks

import (
	"github.com/assemblaj/ggpo/internal/messages"
	"github.com/assemblaj/ggpo/internal/protocol"
)

type FakeMessageHandler struct {
	Endpoint *protocol.UdpProtocol
}

func (f *FakeMessageHandler) HandleMessage(ipAddress string, port int, msg messages.UDPMessage, length int) {
	if f.Endpoint.HandlesMsg(ipAddress, port) {
		f.Endpoint.OnMsg(msg, length)
	}
}

func NewFakeMessageHandler(endpoint *protocol.UdpProtocol) FakeMessageHandler {
	f := FakeMessageHandler{}
	f.Endpoint = endpoint
	return f
}
