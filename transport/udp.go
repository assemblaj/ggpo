package transport

import (
	"net"
	"strconv"

	"github.com/assemblaj/ggpo/internal/messages"
	"github.com/assemblaj/ggpo/internal/util"
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

func getPeerAddress(address net.Addr) peerAddress {
	switch addr := address.(type) {
	case *net.UDPAddr:
		return peerAddress{
			Ip:   addr.IP.String(),
			Port: addr.Port,
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
	util.Log.Printf("binding udp socket to port %d.\n", localPort)
	u.listener, _ = net.ListenPacket("udp", ":"+portStr)
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
	buf := msg.ToBytes()
	u.listener.WriteTo(buf, &RemoteEP)
}

func (u Udp) Read(messageChan chan MessageChannelItem) {
	defer u.listener.Close()
	recvBuf := make([]byte, MaxUDPPacketSize*2)
	for {
		len, addr, err := u.listener.ReadFrom(recvBuf)
		if err != nil {
			util.Log.Printf("conn.Read error returned: %s\n", err)
			break
		} else if len <= 0 {
			util.Log.Printf("no data recieved\n")
		} else if len > 0 {
			util.Log.Printf("recvfrom returned (len:%d  from:%s).\n", len, addr.String())
			peer := getPeerAddress(addr)

			msg, err := messages.DecodeMessageBinary(recvBuf)
			if err != nil {
				util.Log.Printf("Error decoding message: %s", err)
				continue
			}
			messageChan <- MessageChannelItem{Peer: peer, Message: msg, Length: len}
		}

	}
}

func (u Udp) IsInitialized() bool {
	return u.listener != nil
}
