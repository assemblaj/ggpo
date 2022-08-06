package mocks

import (
	"github.com/assemblaj/GGPO-Go/internal/protocol"
	"github.com/assemblaj/GGPO-Go/internal/transport"
)

type FakeMessageHandler struct {
	Endpoint *protocol.UdpProtocol
}

func (f *FakeMessageHandler) HandleMessage(ipAddress string, port int, msg transport.UDPMessage, length int) {
	if f.Endpoint.HandlesMsg(ipAddress, port) {
		f.Endpoint.OnMsg(msg, length)
	}
}

func NewFakeMessageHandler(endpoint *protocol.UdpProtocol) FakeMessageHandler {
	f := FakeMessageHandler{}
	f.Endpoint = endpoint
	return f
}
