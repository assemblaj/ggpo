package ggthx_test

import (
	"bytes"
	"fmt"
	"testing"

	ggthx "github.com/assemblaj/ggthx/src"
)

func TestEncodeDecodeUDPMessage(t *testing.T) {
	want := &ggthx.SyncRequestPacket{}

	packetBuffer, err := ggthx.EncodeMessage(&ggthx.SyncRequestPacket{})
	if err != nil {
		t.Errorf("Error in EncodeMessage %s", err)
	}

	packet, err := ggthx.DecodeMessage(packetBuffer)
	got := packet.(*ggthx.SyncRequestPacket)
	if err != nil {
		t.Errorf("Error in DecodeMessage %s", err)
	}

	if *want != *got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func BenchmarkBinaryEncodeVsGob(b *testing.B) {
	msg := ggthx.NewUDPMessage(ggthx.SyncRequestMsg)

	buf := msg.ToBytes()
	fmt.Printf("Binary buffer size: %d\n", len(buf))

	buf, _ = ggthx.EncodeMessage(msg)
	fmt.Printf("Gob buffer size: %d\n", len(buf))
}

func BenchmarkBinaryEncode(b *testing.B) {
	msg := ggthx.NewUDPMessage(ggthx.SyncRequestMsg)
	for i := 0; i < b.N; i++ {
		msg.ToBytes()
	}
}

func BenchmarkGobEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ggthx.EncodeMessage(&ggthx.SyncRequestPacket{})
	}
}

func TestEncodeHeader(t *testing.T) {
	msg := ggthx.NewUDPMessage(ggthx.SyncRequestMsg)
	workingBuf := msg.ToBytes()
	testedBuf := ggthx.UDPHeader{HeaderType: 1}.ToBytes()
	want, err := ggthx.GetPacketTypeFromBuffer(workingBuf)
	if err != nil {
		t.Errorf("%s", err)
	}
	got, err := ggthx.GetPacketTypeFromBuffer(testedBuf)
	if err != nil {
		t.Errorf("%s", err)
	}
	if want != got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}
func TestEncodeDecodeHeader(t *testing.T) {
	want := ggthx.UDPHeader{HeaderType: uint8(ggthx.SyncRequestMsg), SequenceNumber: 8, Magic: 25}
	buf := want.ToBytes()
	got := ggthx.UDPHeader{}
	got.FromBytes(buf)

	if got != want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestEncodeDecodeUDPConnectionState(t *testing.T) {
	want := ggthx.UdpConnectStatus{Disconnected: true, LastFrame: -1}
	buf := want.ToBytes()
	got := ggthx.UdpConnectStatus{}
	got.FromBytes(buf)

	if got != want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestEncodeDecodeSyncRequestPacket(t *testing.T) {
	packet := ggthx.NewUDPMessage(ggthx.SyncRequestMsg)
	want := packet.(*ggthx.SyncRequestPacket)
	want.RandomRequest = 23
	want.RemoteEndpoint = 24
	want.RemoteMagic = 9000

	buf := want.ToBytes()

	got := ggthx.SyncRequestPacket{}
	got.FromBytes(buf)
	if got != *want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestEncodeDecodeSyncReplyPacket(t *testing.T) {
	packet := ggthx.NewUDPMessage(ggthx.SyncReplyMsg)
	want := packet.(*ggthx.SyncReplyPacket)
	want.RandomReply = 23

	buf := want.ToBytes()

	got := ggthx.SyncReplyPacket{}
	got.FromBytes(buf)
	if got != *want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}
func TestEncodeDecodeQualityReportPacket(t *testing.T) {
	packet := ggthx.NewUDPMessage(ggthx.QualityReportMsg)
	want := packet.(*ggthx.QualityReportPacket)
	want.FrameAdvantage = 90
	want.Ping = 202

	buf := want.ToBytes()

	got := ggthx.QualityReportPacket{}
	got.FromBytes(buf)
	if got != *want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestEncodeDecodeQualityReplytPacket(t *testing.T) {
	packet := ggthx.NewUDPMessage(ggthx.QualityReplyMsg)
	want := packet.(*ggthx.QualityReplyPacket)
	want.Pong = 23434

	buf := want.ToBytes()

	got := ggthx.QualityReplyPacket{}
	got.FromBytes(buf)
	if got != *want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestEncodeDecodeInputAckPacket(t *testing.T) {
	packet := ggthx.NewUDPMessage(ggthx.InputAckMsg)
	want := packet.(*ggthx.InputAckPacket)
	want.AckFrame = 112341

	buf := want.ToBytes()

	got := ggthx.InputAckPacket{}
	got.FromBytes(buf)
	if got != *want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}
func TestEncodeDecodeKeepAlivePacket(t *testing.T) {
	packet := ggthx.NewUDPMessage(ggthx.KeepAliveMsg)
	want := packet.(*ggthx.KeepAlivePacket)

	buf := want.ToBytes()

	got := ggthx.KeepAlivePacket{}
	got.FromBytes(buf)
	if got != *want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}
func TestEncodeInput(t *testing.T) {
	packet := ggthx.NewUDPMessage(ggthx.InputMsg)
	want := packet.(*ggthx.InputPacket)
	want.InputSize = 20
	want.StartFrame = 50
	want.NumBits = 643
	want.DisconectRequested = false
	want.PeerConnectStatus = make([]ggthx.UdpConnectStatus, 2)
	want.PeerConnectStatus[0] = ggthx.UdpConnectStatus{
		Disconnected: false,
		LastFrame:    -1,
	}
	want.PeerConnectStatus[1] = ggthx.UdpConnectStatus{
		Disconnected: true,
		LastFrame:    80,
	}
	want.Bits = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	buf := want.ToBytes()

	got := ggthx.InputPacket{}
	got.FromBytes(buf)
	//gobBuf, _ := ggthx.EncodeMessage(want)
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
This test made sense before ggthx.UDPMessage was a pointer
func TestPackets(t *testing.T) {
	packetTests := []struct {
		name       string
		packetType ggthx.UDPMessage
		want       ggthx.UDPMessage
	}{
		{name: "SyncRequest", packetType: &ggthx.SyncRequestPacket{}, want: &ggthx.SyncRequestPacket{}},
		{name: "SyncReply", packetType: &ggthx.SyncReplyPacket{}, want: &ggthx.SyncReplyPacket{}},
		{name: "QualityReport", packetType: &ggthx.QualityReportPacket{}, want: &ggthx.QualityReportPacket{}},
		{name: "QualityReply", packetType: &ggthx.QualityReplyPacket{}, want: &ggthx.QualityReplyPacket{}},
		//{name: "Input", packetType: Input{}, want: Input{}},
		{name: "InputAck", packetType: &ggthx.InputAckPacket{}, want: &ggthx.InputAckPacket{}},
		{name: "KeepAlive", packetType: &ggthx.KeepAlivePacket{}, want: &ggthx.KeepAlivePacket{}},
	}

	for _, tt := range packetTests {
		packetBuffer, err := ggthx.EncodeMessage(tt.packetType)
		if err != nil {
			t.Errorf("Error in EncodeMessage %s", err)
		}

		got, err := ggthx.DecodeMessage(packetBuffer)
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
	want := ggthx.InputPacket{
		StartFrame: 0,
		AckFrame:   12,
	}
	inputBytes, err := ggthx.EncodeMessage(&want)
	if err != nil {
		t.Errorf("Error in EncodeMessage %s", err)
	}

	packet, err := ggthx.DecodeMessage(inputBytes)
	if err != nil {
		t.Errorf("Error in DecodeMessage %s", err)
	}

	got := packet.(*ggthx.InputPacket)
	if want.StartFrame != got.StartFrame || want.AckFrame != got.AckFrame {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}

}

func TestNewUDPMessage(t *testing.T) {

	want := ggthx.InputPacket{
		AckFrame:   0,
		StartFrame: 5,
	}

	packet := ggthx.NewUDPMessage(ggthx.InputMsg)
	got := packet.(*ggthx.InputPacket)
	got.AckFrame = 0
	got.StartFrame = 5

	if want.StartFrame != got.StartFrame || want.AckFrame != got.AckFrame {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}

}
func TestDecodeMessageBinaryNil(t *testing.T) {
	_, err := ggthx.DecodeMessageBinary(nil)
	if err == nil {
		t.Errorf("A nil buffer should cause an error.")
	}
}

func TestDecodeMessageBinaryTooSmall(t *testing.T) {
	_, err := ggthx.DecodeMessageBinary([]byte{1})
	if err == nil {
		t.Errorf("A buffer with a size < 4 should cause an error.")
	}
}

func TestDecodeMessageBinaryInvalidBuffer(t *testing.T) {
	_, err := ggthx.DecodeMessageBinary([]byte{3, 3, 3, 3, 3, 3, 3, 3})
	if err == nil {
		t.Errorf("Invalid buffer should've created an error. ")
	}
}

func TestDecodeMessageBinaryAllInvalid(t *testing.T) {
	fakePacket := []byte{3, 3, 3, 3, 3, 3, 3, 3}
	for i := 1; i < 8; i++ {
		fakePacket[4] = byte(i)
		_, err := ggthx.DecodeMessageBinary(fakePacket)
		if ggthx.UDPMessageType(i) == ggthx.KeepAliveMsg {
			continue
		}
		if err == nil {
			t.Errorf("Invalid buffer should've created and error.")
		}
	}
}
