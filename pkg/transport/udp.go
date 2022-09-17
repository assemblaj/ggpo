package transport

import (
	"log"
	"net"
	"strconv"

	"github.com/assemblaj/GGPO-Go/internal/messages"
)

const (
	MaxUDPEndpoints  = 16
	MaxUDPPacketSize = 4096
)

type Udp struct {
	Stats UdpStats // may not need this, may just be a service used by others

	socket         net.Conn
	messageHandler MessageHandler
	listener       net.PacketConn
	localPort      int
	ipAddress      string
}

type UdpStats struct {
	BytesSent   int
	PacketsSent int
	KbpsSent    float64
}
type peerAddress struct {
	ip   string
	port int
}

func getPeerAddress(address net.Addr) peerAddress {
	switch addr := address.(type) {
	case *net.UDPAddr:
		return peerAddress{
			ip:   addr.IP.String(),
			port: addr.Port,
		}
	}
	return peerAddress{}
}

func (u Udp) Close() {
	if u.listener != nil {
		u.listener.Close()
	}
}

func NewUdp(messageHandler MessageHandler, localPort int) Udp {
	u := Udp{}
	u.messageHandler = messageHandler

	portStr := strconv.Itoa(localPort)

	u.localPort = localPort
	log.Printf("binding udp socket to port %d.\n", localPort)
	//u.socket = u.CreateSocket(ipAdress, portStr, 0)
	u.listener, _ = net.ListenPacket("udp", "0.0.0.0:"+portStr)
	return u
}

// dst should be sockaddr
// maybe create Gob encoder and decoder members
// instead of creating them on each message send
func (u Udp) SendTo(msg messages.UDPMessage, remoteIp string, remotePort int) {
	if msg == nil || remoteIp == "" {
		return
	}

	RemoteEP := net.UDPAddr{IP: net.ParseIP(remoteIp), Port: remotePort}

	//buf, err := EncodeMessage(msg)
	//buf, err := EncodeMessageBinary(msg)
	buf := msg.ToBytes()
	/*
		if err != nil {
			log.Fatalf("encode error %s", err)
		}*/
	u.listener.WriteTo(buf, &RemoteEP)
}

func (u Udp) Read() {
	defer u.listener.Close()
	recvBuf := make([]byte, MaxUDPPacketSize*2)
	for {
		len, addr, err := u.listener.ReadFrom(recvBuf)

		if err != nil {
			log.Printf("conn.Read error returned: %s\n", err)
			break
		} else if len <= 0 {
			log.Printf("no data recieved\n")
		} else if len > 0 {
			log.Printf("recvfrom returned (len:%d  from:%s).\n", len, addr.String())
			peer := getPeerAddress(addr)

			//msg, err := DecodeMessage(recvBuf)
			msg, err := messages.DecodeMessageBinary(recvBuf)

			if err != nil {
				log.Printf("Error decoding message: %s", err)
				continue
			}
			u.messageHandler.HandleMessage(peer.ip, peer.port, msg, len)
		}

	}
}

func (u Udp) IsInitialized() bool {
	return u.listener != nil
}
