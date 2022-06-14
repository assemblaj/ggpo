package ggthx

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
)

const MAX_UDP_ENDPOINTS int = 16
const MAX_UDP_PACKET_SIZE int = 4096

type Udp struct {
	Stats UdpStats // may not need this, may just be a service used by others

	socket    net.Conn
	callbacks UdpCallbacks
	poll      Poll
}

type UdpStats struct {
	BytesSent   int
	PacketsSent int
	KbpsSent    float64
}

type UdpCallbacks interface {
	// from should be sockaddr_in
	//OnMsg(from string, msg *UdpMsg, len int)
	OnMsg(msg *UdpMsg, len int)
}

func (u *Udp) CreateSocket(ipAddress, port string, retries int) net.Conn {
	conn, err := net.Dial("udp", ipAddress+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func (u *Udp) Close() {
	u.socket.Close()
}

// Give one of these to every protocol object.
// In AddRemotePlayer / AddSpectator
// Or maybe even inside the protocol itself
// AddRemotePlayer / AddSpectator is where we get the IP
// Which is called by backend::AddPlayer
// which is called by Session:AddPlayer
// and obtained via the GGPOPlayer object
func (u *Udp) Init(ipAdress string, port string, p *Poll, callbacks UdpCallbacks) {
	u.callbacks = callbacks
	u.poll = *p
	u.poll.RegisterLoop(u, nil)

	log.Printf("binding udp socket to port %s.\n", port)
	u.socket = u.CreateSocket(ipAdress, port, 0)
}

// dst should be sockaddr
func (u *Udp) SendTo(buffer []byte) {
	u.socket.Write(buffer)
}

func (u *Udp) OnLoopPoll(cookie []byte) bool {
	recvBuf := make([]byte, MAX_UDP_PACKET_SIZE)

	for true {
		len, err := u.socket.Read(recvBuf)
		if err != nil {
			log.Printf("conn.Read error returned: %s\n", err)
			break
		} else if len > 0 {
			log.Printf("recvfrom returned (len:%d  from:%s).\n", len, u.socket.RemoteAddr())
			buf := bytes.NewBuffer(recvBuf)
			dec := gob.NewDecoder(buf)
			msg := UdpMsg{}
			if err = dec.Decode(&msg); err != nil {
				log.Fatal(err)
			}
			u.callbacks.OnMsg(&msg, len)
		}

	}
	return true
}
