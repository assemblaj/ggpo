package transport

type Connection interface {
	SendTo(msg UDPMessage, remoteIp string, remotePort int)
	Read()
	Close()
}
