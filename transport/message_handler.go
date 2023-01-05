package transport

import "github.com/assemblaj/ggpo/internal/messages"

type MessageHandler interface {
	HandleMessage(ipAddress string, port int, msg messages.UDPMessage, len int)
}
