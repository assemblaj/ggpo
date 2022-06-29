package ggthx

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"strconv"
)

const (
	MaxUDPEndpoints  = 16
	MaxUDPPacketSize = 4096
)

type Udp struct {
	Stats UdpStats // may not need this, may just be a service used by others

	socket    net.Conn
	callbacks UdpCallbacks
	poll      Poll
	listener  net.PacketConn
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
func NewUdp(ipAdress string, localPort int, p *Poll, callbacks UdpCallbacks) Udp {
	u := Udp{}
	u.callbacks = callbacks
	u.poll = *p
	u.poll.RegisterLoop(&u, nil)

	portStr := strconv.Itoa(localPort)

	log.Printf("binding udp socket to port %d.\n", localPort)
	//u.socket = u.CreateSocket(ipAdress, portStr, 0)
	u.listener, _ = net.ListenPacket("udp", ipAdress+":"+portStr)

	return u
}

// dst should be sockaddr
// maybe create Gob encoder and decoder members
// instead of creating them on each message send
func (u *Udp) SendTo(msg *UdpMsg, remoteIp string, remotePort int) {
	//
	if msg == nil || remoteIp == "" {
		return
	}

	remotePortStr := strconv.Itoa(remotePort)
	conn, conErr := net.Dial("udp", remoteIp+":"+remotePortStr)
	if conErr != nil {
		log.Fatal("Udp dial error ", conErr)
	}
	defer conn.Close()

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	encErr := enc.Encode(msg)
	if encErr != nil {
		log.Fatal("encode error ", encErr)
	}
	conn.Write(buf.Bytes())
}

func (u *Udp) Read() {
	defer u.listener.Close()
	recvBuf := make([]byte, MaxUDPPacketSize)
	for {
		len, addr, err := u.listener.ReadFrom(recvBuf)

		if err != nil {
			log.Printf("conn.Read error returned: %s\n", err)
			break
		} else if len <= 0 {
			log.Printf("no data recieved\n")
		} else if len > 0 {
			log.Printf("recvfrom returned (len:%d  from:%s).\n", len, addr.String())
			buf := bytes.NewBuffer(recvBuf)
			dec := gob.NewDecoder(buf)
			msg := UdpMsg{}
			if err = dec.Decode(&msg); err != nil {
				log.Fatal(err)
			}
			u.callbacks.OnMsg(&msg, len)
		}

	}
}

func (u *Udp) OnLoopPoll(cookie []byte) bool {
	return true
}

func (u *Udp) IsInitialized() bool {
	return u.listener != nil
}
