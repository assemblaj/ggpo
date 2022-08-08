package transport

type MessageHandler interface {
	HandleMessage(ipAddress string, port int, msg UDPMessage, len int)
}
