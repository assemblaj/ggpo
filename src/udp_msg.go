package ggthx

import (
	"fmt"
	"unsafe"
)

const MaxCompressedBits int = 4096
const UDPMsgMaxPlayers int = 4

// Bad bad bad. Want to find way around this. Giga packet.
// Original used Union
type UdpMsg struct {
	Header        UdpHeader
	MessageType   UdpMsgType
	SyncRequest   UdpSyncRequestPacket
	SyncReply     UdpSyncReplyPacket
	QualityReport UdpQualityReportPacket
	QualityReply  UdpQualityReplyPacket
	Input         UdpInputPacket
	InputAck      UdpInputAckPacket
}

type UdpMsgType int

const (
	InvalidMsg UdpMsgType = iota
	SyncRequestMsg
	SyncReplyMsg
	InputMsg
	QualityReportMsg
	QualityReplyMsg
	KeepAliveMsg
	InputAckMsg
)

type UdpConnectStatus struct {
	Disconnected uint
	LastFrame    int
}

type UdpHeader struct {
	Magic          uint16
	SequenceNumber uint16
	HeaderType     uint8
}

type UdpPacket interface{}

type UdpSyncRequestPacket struct {
	RandomRequest  uint32
	RemoteMagic    uint16
	RemoteEndpoint uint8
}

type UdpSyncReplyPacket struct {
	RandomReply uint32
}

type UdpQualityReportPacket struct {
	FrameAdvantage int8
	Ping           uint32
}

type UdpQualityReplyPacket struct {
	Pong uint32
}

type UdpInputPacket struct {
	PeerConnectStatus []UdpConnectStatus
	StartFrame        uint32

	DisconectRequested int
	AckFrame           int

	NumBits   uint16
	InputSize uint8
	Bits      []byte
}

type UdpInputAckPacket struct {
	AckFrame int
}

func NewUdpMsg(t UdpMsgType) UdpMsg {
	header := UdpHeader{HeaderType: uint8(t)}
	messageType := t
	var msg UdpMsg
	switch t {
	case SyncRequestMsg:
		msg = UdpMsg{
			Header:      header,
			MessageType: messageType,
			SyncRequest: UdpSyncRequestPacket{}}
	case SyncReplyMsg:
		msg = UdpMsg{
			Header:      header,
			MessageType: messageType,
			SyncReply:   UdpSyncReplyPacket{}}
	case QualityReportMsg:
		msg = UdpMsg{
			Header:        header,
			MessageType:   messageType,
			QualityReport: UdpQualityReportPacket{}}
	case QualityReplyMsg:
		msg = UdpMsg{
			Header:       header,
			MessageType:  messageType,
			QualityReply: UdpQualityReplyPacket{}}
	case InputAckMsg:
		msg = UdpMsg{
			Header:      header,
			MessageType: messageType,
			InputAck:    UdpInputAckPacket{}}
	case InputMsg:
		msg = UdpMsg{
			Header:      header,
			MessageType: messageType,
			Input:       UdpInputPacket{}}
	case KeepAliveMsg:
		fallthrough
	default:
		msg = UdpMsg{
			Header:      header,
			MessageType: messageType}

	}
	return msg
}

//func (u *UdpMsg) BuildQualityReply(pong int) {
//	u.qualityReply.pong = uint32(pong)
//}

func (u *UdpMsg) PacketSize() int {
	return int(unsafe.Sizeof(u.Header)) + u.PaylaodSize()
}

func (u *UdpMsg) PaylaodSize() int {
	var size int

	switch UdpMsgType(u.Header.HeaderType) {
	case SyncRequestMsg:
		return int(unsafe.Sizeof(u.SyncRequest))
	case SyncReplyMsg:
		return int(unsafe.Sizeof(u.SyncReply))
	case QualityReportMsg:
		return int(unsafe.Sizeof(u.QualityReport))
	case QualityReplyMsg:
		return int(unsafe.Sizeof(u.QualityReply))
	case InputAckMsg:
		return int(unsafe.Sizeof(u.InputAck))
	case KeepAliveMsg:
		return 0
	case InputMsg:
		for _, s := range u.Input.PeerConnectStatus {
			size += int(unsafe.Sizeof(s))
		}
		size += int(unsafe.Sizeof(u.Input.StartFrame))
		size += int(unsafe.Sizeof(u.Input.DisconectRequested))
		size += int(unsafe.Sizeof(u.Input.AckFrame))
		size += int(unsafe.Sizeof(u.Input.NumBits))
		size += int(unsafe.Sizeof(u.Input.InputSize))
		size += int(u.Input.NumBits+7) / 8
		return size
	}
	Assert(false)
	return 0
}

// might just wanna make this a log function specifically, but this'll do for not
func (u UdpMsg) String() string {
	str := ""
	switch UdpMsgType(u.Header.HeaderType) {
	case SyncRequestMsg:
		str = fmt.Sprintf("sync-request (%d).\n", u.SyncRequest.RandomRequest)
		break
	case SyncReplyMsg:
		str = fmt.Sprintf("sync-reply (%d).\n", u.SyncReply.RandomReply)
		break
	case QualityReportMsg:
		str = fmt.Sprintf("quality report.\n")
		break
	case QualityReplyMsg:
		str = fmt.Sprintf("quality reply.\n")
		break
	case InputAckMsg:
		str = fmt.Sprintf("input ack.\n")
		break
	case KeepAliveMsg:
		str = fmt.Sprintf("keep alive.\n")
		break
	case InputMsg:
		str = fmt.Sprintf("game-compressed-input %d (+ %d bits).\n",
			u.Input.StartFrame, u.Input.NumBits)
		break
	}
	return str
}
