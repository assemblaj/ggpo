package transport

import (
	"github.com/assemblaj/ggpo/internal/messages"
)

type Connection interface {
	SendTo(msg messages.UDPMessage, remoteIp string, remotePort int)
	Close()
	Read(messageChan chan MessageChannelItem)
}

type peerAddress struct {
	Ip   string
	Port int
}

type MessageChannelItem struct {
	Peer    peerAddress
	Message messages.UDPMessage
	Length  int
}
