package ggthx

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"unsafe"
)

const (
	MaxCompressedBits = 4096
	UDPMsgMaxPlayers  = 4
)

func init() {
	gob.Register(&SyncRequestPacket{})
	gob.Register(&SyncReplyPacket{})
	gob.Register(&QualityReportPacket{})
	gob.Register(&QualityReplyPacket{})
	gob.Register(&InputPacket{})
	gob.Register(&InputAckPacket{})
	gob.Register(&KeepAlivePacket{})
}

// Bad bad bad. Want to find way around this. Giga packet.
// Original used Union
type UdpMsg struct {
	Header        UDPHeader
	MessageType   UDPMessageType
	SyncRequest   SyncRequestPacket
	SyncReply     SyncReplyPacket
	QualityReport QualityReportPacket
	QualityReply  QualityReplyPacket
	Input         InputPacket
	InputAck      InputAckPacket
}

type UDPMessage interface {
	Type() UDPMessageType
	Header() UDPHeader
	SetHeader(magicNumber uint16, sequenceNumber uint16)
	PacketSize() int
}

type UDPMessageType int

const (
	InvalidMsg UDPMessageType = iota
	SyncRequestMsg
	SyncReplyMsg
	InputMsg
	QualityReportMsg
	QualityReplyMsg
	KeepAliveMsg
	InputAckMsg
)

type UdpConnectStatus struct {
	Disconnected bool
	LastFrame    int
}

type UDPHeader struct {
	Magic          uint16
	SequenceNumber uint16
	HeaderType     uint8
}

type UdpPacket struct {
	msg *UdpMsg
	len int
}

type SyncRequestPacket struct {
	MessageHeader  UDPHeader
	RandomRequest  uint32
	RemoteMagic    uint16
	RemoteEndpoint uint8
}

func (s *SyncRequestPacket) Type() UDPMessageType { return SyncRequestMsg }
func (s *SyncRequestPacket) Header() UDPHeader    { return s.MessageHeader }
func (s *SyncRequestPacket) SetHeader(magicNumber uint16, sequenceNumber uint16) {
	s.MessageHeader.Magic = magicNumber
	s.MessageHeader.SequenceNumber = sequenceNumber
}
func (s *SyncRequestPacket) PacketSize() int { return int(unsafe.Sizeof(s)) }
func (s *SyncRequestPacket) String() string {
	return fmt.Sprintf("sync-request (%d).\n", s.RandomRequest)
}

type SyncReplyPacket struct {
	MessageHeader UDPHeader
	RandomReply   uint32
}

func (s *SyncReplyPacket) Type() UDPMessageType { return SyncReplyMsg }
func (s *SyncReplyPacket) Header() UDPHeader    { return s.MessageHeader }
func (s *SyncReplyPacket) SetHeader(magicNumber uint16, sequenceNumber uint16) {
	s.MessageHeader.Magic = magicNumber
	s.MessageHeader.SequenceNumber = sequenceNumber
}
func (s *SyncReplyPacket) PacketSize() int { return int(unsafe.Sizeof(s)) }
func (s *SyncReplyPacket) String() string  { return fmt.Sprintf("sync-reply (%d).\n", s.RandomReply) }

type QualityReportPacket struct {
	MessageHeader  UDPHeader
	FrameAdvantage int8
	Ping           uint32
}

func (q *QualityReportPacket) Type() UDPMessageType { return QualityReportMsg }
func (q *QualityReportPacket) Header() UDPHeader    { return q.MessageHeader }
func (q *QualityReportPacket) SetHeader(magicNumber uint16, sequenceNumber uint16) {
	q.MessageHeader.Magic = magicNumber
	q.MessageHeader.SequenceNumber = sequenceNumber
}
func (q *QualityReportPacket) PacketSize() int { return int(unsafe.Sizeof(q)) }
func (s *QualityReportPacket) String() string  { return "quality report.\n" }

type QualityReplyPacket struct {
	MessageHeader UDPHeader
	Pong          uint32
}

func (q *QualityReplyPacket) Type() UDPMessageType { return QualityReplyMsg }
func (q *QualityReplyPacket) Header() UDPHeader    { return q.MessageHeader }
func (q *QualityReplyPacket) SetHeader(magicNumber uint16, sequenceNumber uint16) {
	q.MessageHeader.Magic = magicNumber
	q.MessageHeader.SequenceNumber = sequenceNumber
}
func (q *QualityReplyPacket) PacketSize() int { return int(unsafe.Sizeof(q)) }
func (q *QualityReplyPacket) String() string  { return "quality reply.\n" }

type InputPacket struct {
	MessageHeader     UDPHeader
	PeerConnectStatus []UdpConnectStatus
	StartFrame        uint32

	DisconectRequested bool
	AckFrame           int

	NumBits     uint16
	InputSize   uint8
	Bits        [][]byte
	Sizes       []int
	IsSpectator bool
}

func (i *InputPacket) Type() UDPMessageType { return InputMsg }
func (i *InputPacket) Header() UDPHeader    { return i.MessageHeader }
func (i *InputPacket) SetHeader(magicNumber uint16, sequenceNumber uint16) {
	i.MessageHeader.Magic = magicNumber
	i.MessageHeader.SequenceNumber = sequenceNumber
}

// may go back and make this calculate the Bits, Sizes and IsSpectator better
func (i *InputPacket) PacketSize() int {
	size := 0
	size += int(unsafe.Sizeof(i.MessageHeader))
	for _, s := range i.PeerConnectStatus {
		size += int(unsafe.Sizeof(s))
	}
	size += int(unsafe.Sizeof(i.StartFrame))
	size += int(unsafe.Sizeof(i.DisconectRequested))
	size += int(unsafe.Sizeof(i.AckFrame))
	size += int(unsafe.Sizeof(i.NumBits))
	size += int(unsafe.Sizeof(i.InputSize))
	size += int(i.NumBits+7) / 8
	return size
}
func (i InputPacket) String() string {
	return fmt.Sprintf("game-compressed-input %d (+ %d bits).\n",
		i.StartFrame, i.NumBits)
}

type InputAckPacket struct {
	MessageHeader UDPHeader
	AckFrame      int
}

func (i *InputAckPacket) Type() UDPMessageType { return InputAckMsg }
func (i *InputAckPacket) Header() UDPHeader    { return i.MessageHeader }
func (i *InputAckPacket) SetHeader(magicNumber uint16, sequenceNumber uint16) {
	i.MessageHeader.Magic = magicNumber
	i.MessageHeader.SequenceNumber = sequenceNumber
}
func (i *InputAckPacket) PacketSize() int { return int(unsafe.Sizeof(i)) }
func (i *InputAckPacket) String() string  { return "input ack.\n" }

type KeepAlivePacket struct {
	MessageHeader UDPHeader
}

func (k *KeepAlivePacket) Type() UDPMessageType { return KeepAliveMsg }
func (k *KeepAlivePacket) Header() UDPHeader    { return k.MessageHeader }
func (k *KeepAlivePacket) SetHeader(magicNumber uint16, sequenceNumber uint16) {
	k.MessageHeader.Magic = magicNumber
	k.MessageHeader.SequenceNumber = sequenceNumber
}
func (k *KeepAlivePacket) PacketSize() int { return int(unsafe.Sizeof(k)) }
func (k *KeepAlivePacket) String() string  { return "keep alive.\n" }

func NewUdpMsg(t UDPMessageType) UdpMsg {
	header := UDPHeader{HeaderType: uint8(t)}
	messageType := t
	var msg UdpMsg
	switch t {
	case SyncRequestMsg:
		msg = UdpMsg{
			Header:      header,
			MessageType: messageType,
			SyncRequest: SyncRequestPacket{}}
	case SyncReplyMsg:
		msg = UdpMsg{
			Header:      header,
			MessageType: messageType,
			SyncReply:   SyncReplyPacket{}}
	case QualityReportMsg:
		msg = UdpMsg{
			Header:        header,
			MessageType:   messageType,
			QualityReport: QualityReportPacket{}}
	case QualityReplyMsg:
		msg = UdpMsg{
			Header:       header,
			MessageType:  messageType,
			QualityReply: QualityReplyPacket{}}
	case InputAckMsg:
		msg = UdpMsg{
			Header:      header,
			MessageType: messageType,
			InputAck:    InputAckPacket{}}
	case InputMsg:
		msg = UdpMsg{
			Header:      header,
			MessageType: messageType,
			Input:       InputPacket{}}
	case KeepAliveMsg:
		fallthrough
	default:
		msg = UdpMsg{
			Header:      header,
			MessageType: messageType}

	}
	return msg
}

func NewUDPMessage(t UDPMessageType) UDPMessage {
	header := UDPHeader{HeaderType: uint8(t)}
	var msg UDPMessage
	switch t {
	case SyncRequestMsg:
		msg = &SyncRequestPacket{
			MessageHeader: header}
	case SyncReplyMsg:
		msg = &SyncReplyPacket{
			MessageHeader: header}
	case QualityReportMsg:
		msg = &QualityReportPacket{
			MessageHeader: header}
	case QualityReplyMsg:
		msg = &QualityReplyPacket{
			MessageHeader: header}
	case InputAckMsg:
		msg = &InputAckPacket{
			MessageHeader: header}
	case InputMsg:
		msg = &InputPacket{
			MessageHeader: header}
	case KeepAliveMsg:
		fallthrough
	default:
		msg = &KeepAlivePacket{
			MessageHeader: header}
	}
	return msg
}

func (u *UdpMsg) PacketSize() int {
	size, err := u.PaylaodSize()
	//Unknown Packet type somehow
	if err != nil {
		// Send size of whole object
		//return int(unsafe.Sizeof(u))
		panic(err)
	}
	return int(unsafe.Sizeof(u.Header)) + size
}

func (u *UdpMsg) PaylaodSize() (int, error) {
	var size int

	switch UDPMessageType(u.Header.HeaderType) {
	case SyncRequestMsg:
		return int(unsafe.Sizeof(u.SyncRequest)), nil
	case SyncReplyMsg:
		return int(unsafe.Sizeof(u.SyncReply)), nil
	case QualityReportMsg:
		return int(unsafe.Sizeof(u.QualityReport)), nil
	case QualityReplyMsg:
		return int(unsafe.Sizeof(u.QualityReply)), nil
	case InputAckMsg:
		return int(unsafe.Sizeof(u.InputAck)), nil
	case KeepAliveMsg:
		return 0, nil
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
		return size, nil
	}
	return 0, errors.New("ggthx UdpMsg PayloadSize: invalid packet type, could not find payload size")
}

// might just wanna make this a log function specifically, but this'll do for not
func (u UdpMsg) String() string {
	str := ""
	switch UDPMessageType(u.Header.HeaderType) {
	case SyncRequestMsg:
		str = fmt.Sprintf("sync-request (%d).\n", u.SyncRequest.RandomRequest)
	case SyncReplyMsg:
		str = fmt.Sprintf("sync-reply (%d).\n", u.SyncReply.RandomReply)
	case QualityReportMsg:
		str = "quality report.\n"
	case QualityReplyMsg:
		str = "quality reply.\n"
	case InputAckMsg:
		str = "input ack.\n"
	case KeepAliveMsg:
		str = "keep alive.\n"
	case InputMsg:
		str = fmt.Sprintf("game-compressed-input %d (+ %d bits).\n",
			u.Input.StartFrame, u.Input.NumBits)
	}
	return str
}

func EncodeMessage(packet UDPMessage) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	encErr := enc.Encode(&packet)
	if encErr != nil {
		return nil, encErr
	}
	return buf.Bytes(), nil
}

func DecodeMessage(buffer []byte) (UDPMessage, error) {
	buf := bytes.NewBuffer(buffer)
	dec := gob.NewDecoder(buf)
	var msg UDPMessage
	var err error

	if err = dec.Decode(&msg); err != nil {
		return nil, err
	}
	return msg, nil
}
