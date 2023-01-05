package messages_test

import (
	"bytes"
	"fmt"
	"testing"

	messages "github.com/assemblaj/ggpo/internal/messages"
)

func TestEncodeDecodeUDPMessage(t *testing.T) {
	want := &messages.SyncRequestPacket{}

	packetBuffer, err := messages.EncodeMessage(&messages.SyncRequestPacket{})
	if err != nil {
		t.Errorf("Error in EncodeMessage %s", err)
	}

	packet, err := messages.DecodeMessage(packetBuffer)
	got := packet.(*messages.SyncRequestPacket)
	if err != nil {
		t.Errorf("Error in DecodeMessage %s", err)
	}

	if *want != *got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func BenchmarkBinaryEncodeVsGob(b *testing.B) {
	msg := messages.NewUDPMessage(messages.SyncRequestMsg)

	buf := msg.ToBytes()
	fmt.Printf("Binary buffer size: %d\n", len(buf))

	buf, _ = messages.EncodeMessage(msg)
	fmt.Printf("Gob buffer size: %d\n", len(buf))
}

func BenchmarkBinaryEncode(b *testing.B) {
	msg := messages.NewUDPMessage(messages.SyncRequestMsg)
	for i := 0; i < b.N; i++ {
		msg.ToBytes()
	}
}

func BenchmarkGobEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		messages.EncodeMessage(&messages.SyncRequestPacket{})
	}
}

func TestEncodeHeader(t *testing.T) {
	msg := messages.NewUDPMessage(messages.SyncRequestMsg)
	workingBuf := msg.ToBytes()
	testedBuf := messages.UDPHeader{HeaderType: 1}.ToBytes()
	want, err := messages.GetPacketTypeFromBuffer(workingBuf)
	if err != nil {
		t.Errorf("%s", err)
	}
	got, err := messages.GetPacketTypeFromBuffer(testedBuf)
	if err != nil {
		t.Errorf("%s", err)
	}
	if want != got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}
func TestEncodeDecodeHeader(t *testing.T) {
	want := messages.UDPHeader{HeaderType: uint8(messages.SyncRequestMsg), SequenceNumber: 8, Magic: 25}
	buf := want.ToBytes()
	got := messages.UDPHeader{}
	got.FromBytes(buf)

	if got != want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestEncodeDecodeUDPConnectionState(t *testing.T) {
	want := messages.UdpConnectStatus{Disconnected: true, LastFrame: -1}
	buf := want.ToBytes()
	got := messages.UdpConnectStatus{}
	got.FromBytes(buf)

	if got != want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestEncodeDecodeSyncRequestPacket(t *testing.T) {
	packet := messages.NewUDPMessage(messages.SyncRequestMsg)
	want := packet.(*messages.SyncRequestPacket)
	want.RandomRequest = 23
	want.RemoteEndpoint = 24
	want.RemoteMagic = 9000

	buf := want.ToBytes()

	got := messages.SyncRequestPacket{}
	got.FromBytes(buf)
	if got != *want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestEncodeDecodeSyncReplyPacket(t *testing.T) {
	packet := messages.NewUDPMessage(messages.SyncReplyMsg)
	want := packet.(*messages.SyncReplyPacket)
	want.RandomReply = 23

	buf := want.ToBytes()

	got := messages.SyncReplyPacket{}
	got.FromBytes(buf)
	if got != *want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}
func TestEncodeDecodeQualityReportPacket(t *testing.T) {
	packet := messages.NewUDPMessage(messages.QualityReportMsg)
	want := packet.(*messages.QualityReportPacket)
	want.FrameAdvantage = 90
	want.Ping = 202

	buf := want.ToBytes()

	got := messages.QualityReportPacket{}
	got.FromBytes(buf)
	if got != *want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestEncodeDecodeQualityReplytPacket(t *testing.T) {
	packet := messages.NewUDPMessage(messages.QualityReplyMsg)
	want := packet.(*messages.QualityReplyPacket)
	want.Pong = 23434

	buf := want.ToBytes()

	got := messages.QualityReplyPacket{}
	got.FromBytes(buf)
	if got != *want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestEncodeDecodeInputAckPacket(t *testing.T) {
	packet := messages.NewUDPMessage(messages.InputAckMsg)
	want := packet.(*messages.InputAckPacket)
	want.AckFrame = 112341

	buf := want.ToBytes()

	got := messages.InputAckPacket{}
	got.FromBytes(buf)
	if got != *want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}
func TestEncodeDecodeKeepAlivePacket(t *testing.T) {
	packet := messages.NewUDPMessage(messages.KeepAliveMsg)
	want := packet.(*messages.KeepAlivePacket)

	buf := want.ToBytes()

	got := messages.KeepAlivePacket{}
	got.FromBytes(buf)
	if got != *want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}
func TestEncodeInput(t *testing.T) {
	packet := messages.NewUDPMessage(messages.InputMsg)
	want := packet.(*messages.InputPacket)
	want.InputSize = 20
	want.StartFrame = 50
	want.NumBits = 643
	want.Checksum = 98790
	want.DisconectRequested = false
	want.PeerConnectStatus = make([]messages.UdpConnectStatus, 2)
	want.PeerConnectStatus[0] = messages.UdpConnectStatus{
		Disconnected: false,
		LastFrame:    -1,
	}
	want.PeerConnectStatus[1] = messages.UdpConnectStatus{
		Disconnected: true,
		LastFrame:    80,
	}
	want.Bits = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	buf := want.ToBytes()

	got := messages.InputPacket{}
	got.FromBytes(buf)
	//gobBuf, _ := messages.EncodeMessage(want)
	//fmt.Printf("Gob input size: %d Binary input size %d\n", len(gobBuf), len(buf))

	if got.InputSize != want.InputSize {
		t.Errorf("expected Input Size '%d' but got '%d'", want.InputSize, got.InputSize)
	}
	if got.StartFrame != want.StartFrame {
		t.Errorf("expected Start Frame '%d' but got '%d'", want.StartFrame, got.StartFrame)
	}
	if got.Checksum != want.Checksum {
		t.Errorf("expected Checksum '%d' but got '%d'", want.Checksum, got.Checksum)
	}
	if got.NumBits != want.NumBits {
		t.Errorf("expected Num Bits '%d' but got '%d'", want.NumBits, got.NumBits)
	}
	if got.DisconectRequested != want.DisconectRequested {
		t.Errorf("expected Disconnect Requested '%t' but got '%t'", want.DisconectRequested, got.DisconectRequested)
	}
	if len(got.PeerConnectStatus) != len(want.PeerConnectStatus) {
		t.Errorf("expected Peer Connect Status length'%#v' but got '%#v'", len(want.PeerConnectStatus), len(got.PeerConnectStatus))
	}

	for i := 0; i < len(want.PeerConnectStatus); i++ {
		if got.PeerConnectStatus[i] != want.PeerConnectStatus[i] {
			t.Errorf("expected Peer Connect Status '%#v' but got '%#v'", want.PeerConnectStatus[i], got.PeerConnectStatus[i])
		}
	}

	if len(got.Bits) != len(want.Bits) {
		t.Errorf("expected Bits length'%#v' but got '%#v'", len(want.Bits), len(got.Bits))
	}

	if !bytes.Equal(got.Bits, want.Bits) {
		t.Errorf("expected Bits Slice '%#v' but got '%#v'", want.Bits, got.Bits)
	}

}

/*
/*
This test made sense before messages.UDPMessage was a pointer
func TestPackets(t *testing.T) {
	packetTests := []struct {
		name       string
		packetType messages.UDPMessage
		want       messages.UDPMessage
	}{
		{name: "SyncRequest", packetType: &messages.SyncRequestPacket{}, want: &messages.SyncRequestPacket{}},
		{name: "SyncReply", packetType: &messages.SyncReplyPacket{}, want: &messages.SyncReplyPacket{}},
		{name: "QualityReport", packetType: &messages.QualityReportPacket{}, want: &messages.QualityReportPacket{}},
		{name: "QualityReply", packetType: &messages.QualityReplyPacket{}, want: &messages.QualityReplyPacket{}},
		//{name: "Input", packetType: Input{}, want: Input{}},
		{name: "InputAck", packetType: &messages.InputAckPacket{}, want: &messages.InputAckPacket{}},
		{name: "KeepAlive", packetType: &messages.KeepAlivePacket{}, want: &messages.KeepAlivePacket{}},
	}

	for _, tt := range packetTests {
		packetBuffer, err := messages.EncodeMessage(tt.packetType)
		if err != nil {
			t.Errorf("Error in EncodeMessage %s", err)
		}

		got, err := messages.DecodeMessage(packetBuffer)
		if err != nil {
			t.Errorf("Error in DecodeMessage %s", err)
		}

		if got != tt.want {
			t.Errorf("%s got %v want %v", tt.name, got, tt.want)
		}
	}

}
*/

func TestExtractInputFromBytes(t *testing.T) {
	want := messages.InputPacket{
		StartFrame: 0,
		AckFrame:   12,
		Checksum:   5002,
	}
	inputBytes, err := messages.EncodeMessage(&want)
	if err != nil {
		t.Errorf("Error in EncodeMessage %s", err)
	}

	packet, err := messages.DecodeMessage(inputBytes)
	if err != nil {
		t.Errorf("Error in DecodeMessage %s", err)
	}

	got := packet.(*messages.InputPacket)
	if want.StartFrame != got.StartFrame || want.AckFrame != got.AckFrame || want.Checksum != got.Checksum {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}

}

func TestNewUDPMessage(t *testing.T) {

	want := messages.InputPacket{
		AckFrame:   0,
		StartFrame: 5,
		Checksum:   16,
	}

	packet := messages.NewUDPMessage(messages.InputMsg)
	got := packet.(*messages.InputPacket)
	got.AckFrame = 0
	got.StartFrame = 5
	got.Checksum = 16

	if want.StartFrame != got.StartFrame || want.AckFrame != got.AckFrame || want.Checksum != got.Checksum {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}

}
func TestDecodeMessageBinaryNil(t *testing.T) {
	_, err := messages.DecodeMessageBinary(nil)
	if err == nil {
		t.Errorf("A nil buffer should cause an error.")
	}
}

func TestDecodeMessageBinaryTooSmall(t *testing.T) {
	_, err := messages.DecodeMessageBinary([]byte{1})
	if err == nil {
		t.Errorf("A buffer with a size < 4 should cause an error.")
	}
}

func TestDecodeMessageBinaryInvalidBuffer(t *testing.T) {
	_, err := messages.DecodeMessageBinary([]byte{3, 3, 3, 3, 3, 3, 3, 3})
	if err == nil {
		t.Errorf("Invalid buffer should've created an error. ")
	}
}

func TestDecodeMessageBinaryAllInvalid(t *testing.T) {
	fakePacket := []byte{3, 3, 3, 3, 3, 3, 3, 3}
	for i := 1; i < 8; i++ {
		fakePacket[4] = byte(i)
		_, err := messages.DecodeMessageBinary(fakePacket)
		if messages.UDPMessageType(i) == messages.KeepAliveMsg {
			continue
		}
		if err == nil {
			t.Errorf("Invalid buffer should've created and error.")
		}
	}
}
