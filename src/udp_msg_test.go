package ggthx_test

import (
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
