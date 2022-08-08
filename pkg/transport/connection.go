package transport

import "github.com/assemblaj/GGPO-Go/internal/messages"

type Connection interface {
	SendTo(msg messages.UDPMessage, remoteIp string, remotePort int)
	Read()
	Close()
}
