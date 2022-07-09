package ggthx

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"strconv"
	"strings"
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

func getPeerAddress(address string) peerAddress {
	peerIPSlice := strings.Split(address, ":")
	if len(peerIPSlice) < 2 {
		panic("Please enter IP as ip:port")
	}
	peerPort, err := strconv.Atoi(peerIPSlice[1])
	if err != nil {
		panic("Please enter integer port")
	}
	return peerAddress{
		ip:   peerIPSlice[0],
		port: peerPort,
	}
}

type MessageHandler interface {
	HandleMessage(ipAddress string, port int, msg *UdpMsg, len int)
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
func (u *Udp) SendTo(msg *UdpMsg, remoteIp string, remotePort int) {
	if msg == nil || remoteIp == "" {
		return
	}

	RemoteEP := net.UDPAddr{IP: net.ParseIP(remoteIp), Port: remotePort}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	encErr := enc.Encode(msg)
	if encErr != nil {
		log.Fatal("encode error ", encErr)
	}

	u.listener.WriteTo(buf.Bytes(), &RemoteEP)
}

func (u *Udp) Read() {
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
			peer := getPeerAddress(addr.String())

			buf := bytes.NewBuffer(recvBuf)
			dec := gob.NewDecoder(buf)
			msg := UdpMsg{}
			if err = dec.Decode(&msg); err != nil {
				log.Fatal(err)
			}
			u.messageHandler.HandleMessage(peer.ip, peer.port, &msg, len)
		}

	}
}

func (u *Udp) IsInitialized() bool {
	return u.listener != nil
}
