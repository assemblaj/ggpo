package transport_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/assemblaj/ggthx/internal/transport"
)

func TestEncodeDecodeUDPMessage(t *testing.T) {
	want := &transport.SyncRequestPacket{}

	packetBuffer, err := transport.EncodeMessage(&transport.SyncRequestPacket{})
	if err != nil {
		t.Errorf("Error in EncodeMessage %s", err)
	}

	packet, err := transport.DecodeMessage(packetBuffer)
	got := packet.(*transport.SyncRequestPacket)
	if err != nil {
		t.Errorf("Error in DecodeMessage %s", err)
	}

	if *want != *got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func BenchmarkBinaryEncodeVsGob(b *testing.B) {
	msg := transport.NewUDPMessage(transport.SyncRequestMsg)

	buf := msg.ToBytes()
	fmt.Printf("Binary buffer size: %d\n", len(buf))

	buf, _ = transport.EncodeMessage(msg)
	fmt.Printf("Gob buffer size: %d\n", len(buf))
}

func BenchmarkBinaryEncode(b *testing.B) {
	msg := transport.NewUDPMessage(transport.SyncRequestMsg)
	for i := 0; i < b.N; i++ {
		msg.ToBytes()
	}
}

func BenchmarkGobEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		transport.EncodeMessage(&transport.SyncRequestPacket{})
	}
}

func TestEncodeHeader(t *testing.T) {
	msg := transport.NewUDPMessage(transport.SyncRequestMsg)
	workingBuf := msg.ToBytes()
	testedBuf := transport.UDPHeader{HeaderType: 1}.ToBytes()
	want, err := transport.GetPacketTypeFromBuffer(workingBuf)
	if err != nil {
		t.Errorf("%s", err)
	}
	got, err := transport.GetPacketTypeFromBuffer(testedBuf)
	if err != nil {
		t.Errorf("%s", err)
	}
	if want != got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}
func TestEncodeDecodeHeader(t *testing.T) {
	want := transport.UDPHeader{HeaderType: uint8(transport.SyncRequestMsg), SequenceNumber: 8, Magic: 25}
	buf := want.ToBytes()
	got := transport.UDPHeader{}
	got.FromBytes(buf)

	if got != want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestEncodeDecodeUDPConnectionState(t *testing.T) {
	want := transport.UdpConnectStatus{Disconnected: true, LastFrame: -1}
	buf := want.ToBytes()
	got := transport.UdpConnectStatus{}
	got.FromBytes(buf)

	if got != want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestEncodeDecodeSyncRequestPacket(t *testing.T) {
	packet := transport.NewUDPMessage(transport.SyncRequestMsg)
	want := packet.(*transport.SyncRequestPacket)
	want.RandomRequest = 23
	want.RemoteEndpoint = 24
	want.RemoteMagic = 9000

	buf := want.ToBytes()

	got := transport.SyncRequestPacket{}
	got.FromBytes(buf)
	if got != *want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestEncodeDecodeSyncReplyPacket(t *testing.T) {
	packet := transport.NewUDPMessage(transport.SyncReplyMsg)
	want := packet.(*transport.SyncReplyPacket)
	want.RandomReply = 23

	buf := want.ToBytes()

	got := transport.SyncReplyPacket{}
	got.FromBytes(buf)
	if got != *want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}
func TestEncodeDecodeQualityReportPacket(t *testing.T) {
	packet := transport.NewUDPMessage(transport.QualityReportMsg)
	want := packet.(*transport.QualityReportPacket)
	want.FrameAdvantage = 90
	want.Ping = 202

	buf := want.ToBytes()

	got := transport.QualityReportPacket{}
	got.FromBytes(buf)
	if got != *want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestEncodeDecodeQualityReplytPacket(t *testing.T) {
	packet := transport.NewUDPMessage(transport.QualityReplyMsg)
	want := packet.(*transport.QualityReplyPacket)
	want.Pong = 23434

	buf := want.ToBytes()

	got := transport.QualityReplyPacket{}
	got.FromBytes(buf)
	if got != *want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestEncodeDecodeInputAckPacket(t *testing.T) {
	packet := transport.NewUDPMessage(transport.InputAckMsg)
	want := packet.(*transport.InputAckPacket)
	want.AckFrame = 112341

	buf := want.ToBytes()

	got := transport.InputAckPacket{}
	got.FromBytes(buf)
	if got != *want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}
func TestEncodeDecodeKeepAlivePacket(t *testing.T) {
	packet := transport.NewUDPMessage(transport.KeepAliveMsg)
	want := packet.(*transport.KeepAlivePacket)

	buf := want.ToBytes()

	got := transport.KeepAlivePacket{}
	got.FromBytes(buf)
	if got != *want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}
func TestEncodeInput(t *testing.T) {
	packet := transport.NewUDPMessage(transport.InputMsg)
	want := packet.(*transport.InputPacket)
	want.InputSize = 20
	want.StartFrame = 50
	want.NumBits = 643
	want.DisconectRequested = false
	want.PeerConnectStatus = make([]transport.UdpConnectStatus, 2)
	want.PeerConnectStatus[0] = transport.UdpConnectStatus{
		Disconnected: false,
		LastFrame:    -1,
	}
	want.PeerConnectStatus[1] = transport.UdpConnectStatus{
		Disconnected: true,
		LastFrame:    80,
	}
	want.Bits = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	buf := want.ToBytes()

	got := transport.InputPacket{}
	got.FromBytes(buf)
	//gobBuf, _ := transport.EncodeMessage(want)
	//fmt.Printf("Gob input size: %d Binary input size %d\n", len(gobBuf), len(buf))

	if got.InputSize != want.InputSize {
		t.Errorf("expected Input Size '%d' but got '%d'", want.InputSize, got.InputSize)
	}
	if got.StartFrame != want.StartFrame {
		t.Errorf("expected Start Frame '%d' but got '%d'", want.StartFrame, got.StartFrame)
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
This test made sense before transport.UDPMessage was a pointer
func TestPackets(t *testing.T) {
	packetTests := []struct {
		name       string
		packetType transport.UDPMessage
		want       transport.UDPMessage
	}{
		{name: "SyncRequest", packetType: &transport.SyncRequestPacket{}, want: &transport.SyncRequestPacket{}},
		{name: "SyncReply", packetType: &transport.SyncReplyPacket{}, want: &transport.SyncReplyPacket{}},
		{name: "QualityReport", packetType: &transport.QualityReportPacket{}, want: &transport.QualityReportPacket{}},
		{name: "QualityReply", packetType: &transport.QualityReplyPacket{}, want: &transport.QualityReplyPacket{}},
		//{name: "Input", packetType: Input{}, want: Input{}},
		{name: "InputAck", packetType: &transport.InputAckPacket{}, want: &transport.InputAckPacket{}},
		{name: "KeepAlive", packetType: &transport.KeepAlivePacket{}, want: &transport.KeepAlivePacket{}},
	}

	for _, tt := range packetTests {
		packetBuffer, err := transport.EncodeMessage(tt.packetType)
		if err != nil {
			t.Errorf("Error in EncodeMessage %s", err)
		}

		got, err := transport.DecodeMessage(packetBuffer)
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
	want := transport.InputPacket{
		StartFrame: 0,
		AckFrame:   12,
	}
	inputBytes, err := transport.EncodeMessage(&want)
	if err != nil {
		t.Errorf("Error in EncodeMessage %s", err)
	}

	packet, err := transport.DecodeMessage(inputBytes)
	if err != nil {
		t.Errorf("Error in DecodeMessage %s", err)
	}

	got := packet.(*transport.InputPacket)
	if want.StartFrame != got.StartFrame || want.AckFrame != got.AckFrame {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}

}

func TestNewUDPMessage(t *testing.T) {

	want := transport.InputPacket{
		AckFrame:   0,
		StartFrame: 5,
	}

	packet := transport.NewUDPMessage(transport.InputMsg)
	got := packet.(*transport.InputPacket)
	got.AckFrame = 0
	got.StartFrame = 5

	if want.StartFrame != got.StartFrame || want.AckFrame != got.AckFrame {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}

}
func TestDecodeMessageBinaryNil(t *testing.T) {
	_, err := transport.DecodeMessageBinary(nil)
	if err == nil {
		t.Errorf("A nil buffer should cause an error.")
	}
}

func TestDecodeMessageBinaryTooSmall(t *testing.T) {
	_, err := transport.DecodeMessageBinary([]byte{1})
	if err == nil {
		t.Errorf("A buffer with a size < 4 should cause an error.")
	}
}

func TestDecodeMessageBinaryInvalidBuffer(t *testing.T) {
	_, err := transport.DecodeMessageBinary([]byte{3, 3, 3, 3, 3, 3, 3, 3})
	if err == nil {
		t.Errorf("Invalid buffer should've created an error. ")
	}
}

func TestDecodeMessageBinaryAllInvalid(t *testing.T) {
	fakePacket := []byte{3, 3, 3, 3, 3, 3, 3, 3}
	for i := 1; i < 8; i++ {
		fakePacket[4] = byte(i)
		_, err := transport.DecodeMessageBinary(fakePacket)
		if transport.UDPMessageType(i) == transport.KeepAliveMsg {
			continue
		}
		if err == nil {
			t.Errorf("Invalid buffer should've created and error.")
		}
	}
}
