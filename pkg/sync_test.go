package ggthx_test

import (
	"bytes"
	"testing"

	"github.com/assemblaj/ggthx/internal/input"
	"github.com/assemblaj/ggthx/internal/mocks"
	"github.com/assemblaj/ggthx/internal/transport"
	ggthx "github.com/assemblaj/ggthx/pkg"
)

/*

	Characterization Tests Basically

*/
func TestNewSync(t *testing.T) {
	session := mocks.NewFakeSession()
	sessionCallbacks := mocks.MakeSessionCallBacks(session)

	peerConnection := []transport.UdpConnectStatus{
		{Disconnected: false, LastFrame: 12},
		{Disconnected: false, LastFrame: 13},
	}
	syncConfig := ggthx.NeweSyncConfig(
		sessionCallbacks, 8, 2, 4,
	)
	sync := ggthx.NewSync(peerConnection, &syncConfig)
	want := 0
	got := sync.FrameCount()
	if want != got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

// Using Trying to load a frame when it hasn't been saved
func TestSyncLoadFrameCharacterization(t *testing.T) {
	session := mocks.NewFakeSession()
	sessionCallbacks := mocks.MakeSessionCallBacks(session)

	peerConnection := []transport.UdpConnectStatus{
		{Disconnected: false, LastFrame: 12},
		{Disconnected: false, LastFrame: 13},
	}
	syncConfig := ggthx.NeweSyncConfig(
		sessionCallbacks, 8, 2, 4,
	)
	sync := ggthx.NewSync(peerConnection, &syncConfig)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic when load frame tried to load an unsaved frame")
		}
	}()
	sync.LoadFrame(6)

}
func TestSyncIncrementFrame(t *testing.T) {
	session := mocks.NewFakeSession()
	sessionCallbacks := mocks.MakeSessionCallBacks(session)

	peerConnection := []transport.UdpConnectStatus{
		{Disconnected: false, LastFrame: 12},
		{Disconnected: false, LastFrame: 13},
	}
	syncConfig := ggthx.NeweSyncConfig(
		sessionCallbacks, 8, 2, 4,
	)
	sync := ggthx.NewSync(peerConnection, &syncConfig)
	sync.IncrementFrame()
	sync.LoadFrame(0)
	want := 1
	got := sync.FrameCount()
	if got != want {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}
func TestSyncAdustSimulationPanicIfSeekToUnsavedFrame(t *testing.T) {
	session := mocks.NewFakeSession()
	sessionCallbacks := mocks.MakeSessionCallBacks(session)

	peerConnection := []transport.UdpConnectStatus{
		{Disconnected: false, LastFrame: 12},
		{Disconnected: false, LastFrame: 13},
	}
	syncConfig := ggthx.NeweSyncConfig(
		sessionCallbacks, 8, 2, 4,
	)
	sync := ggthx.NewSync(peerConnection, &syncConfig)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic when AdjustSimulation attempted to load an unsaved frame.")
		}
	}()
	sync.AdjustSimulation(18)

}
func TestSyncAjdustSimulationError(t *testing.T) {
	session := mocks.NewFakeSession()
	sessionCallbacks := mocks.MakeSessionCallBacks(session)

	peerConnection := []transport.UdpConnectStatus{
		{Disconnected: false, LastFrame: 12},
		{Disconnected: false, LastFrame: 13},
	}
	syncConfig := ggthx.NeweSyncConfig(
		sessionCallbacks, 8, 2, 4,
	)
	sync := ggthx.NewSync(peerConnection, &syncConfig)
	frameCount := 2
	for i := 0; i < frameCount; i++ {
		sync.IncrementFrame()
	}

	err := sync.AdjustSimulation(1)
	if err == nil {
		t.Errorf("The code did not error when AdustSimulation was seeking to the non current frame.")
	}
}

func TestSyncAdjustSimulationSucess(t *testing.T) {
	session := mocks.NewFakeSession()
	sessionCallbacks := mocks.MakeSessionCallBacks(session)

	peerConnection := []transport.UdpConnectStatus{
		{Disconnected: false, LastFrame: 12},
		{Disconnected: false, LastFrame: 13},
	}
	syncConfig := ggthx.NeweSyncConfig(
		sessionCallbacks, 8, 2, 4,
	)
	sync := ggthx.NewSync(peerConnection, &syncConfig)
	frameCount := 2
	for i := 0; i < frameCount; i++ {
		sync.IncrementFrame()
	}

	err := sync.AdjustSimulation(frameCount)
	if err != nil {
		t.Errorf("The code errored when AdustSimulation was seeking to the current frame.")
	}
}

func TestSyncAddLocalInput(t *testing.T) {
	session := mocks.NewFakeSession()
	sessionCallbacks := mocks.MakeSessionCallBacks(session)

	peerConnection := []transport.UdpConnectStatus{
		{Disconnected: false, LastFrame: 12},
		{Disconnected: false, LastFrame: 13},
	}
	syncConfig := ggthx.NeweSyncConfig(
		sessionCallbacks, 8, 2, 4,
	)
	sync := ggthx.NewSync(peerConnection, &syncConfig)
	queue := 0
	input := input.GameInput{}
	success := sync.AddLocalInput(queue, &input)
	want := 0
	got := input.Frame
	if success == false {
		t.Errorf("The AddLocalInput failed.")
	}
	if want != got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}
func TestSyncAddLocalInputAfterIncrementFrame(t *testing.T) {
	session := mocks.NewFakeSession()
	sessionCallbacks := mocks.MakeSessionCallBacks(session)

	peerConnection := []transport.UdpConnectStatus{
		{Disconnected: false, LastFrame: 12},
		{Disconnected: false, LastFrame: 13},
	}
	syncConfig := ggthx.NeweSyncConfig(
		sessionCallbacks, 8, 2, 4,
	)
	sync := ggthx.NewSync(peerConnection, &syncConfig)
	frameCount := 2
	for i := 0; i < frameCount; i++ {
		sync.IncrementFrame()
	}

	queue := 0
	input := input.GameInput{}
	success := sync.AddLocalInput(queue, &input)
	want := 2
	got := input.Frame
	if success == false {
		t.Errorf("The AddLocalInput failed.")
	}
	if want != got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}
func TestSyncSynchronizeInputsNoInput(t *testing.T) {
	session := mocks.NewFakeSession()
	sessionCallbacks := mocks.MakeSessionCallBacks(session)

	peerConnection := []transport.UdpConnectStatus{
		{Disconnected: false, LastFrame: 12},
		{Disconnected: false, LastFrame: 13},
	}
	syncConfig := ggthx.NeweSyncConfig(
		sessionCallbacks, 8, 2, 4,
	)
	sync := ggthx.NewSync(peerConnection, &syncConfig)
	inputs, disconnectFlags := sync.SynchronizeInputs()
	want := 2
	got := len(inputs)
	if want != got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}

	if len(inputs[0]) != 0 || len(inputs[1]) != 0 {
		t.Errorf("expected input lengths to be zero.")
	}

	want = 0
	got = disconnectFlags
	if want != got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestSyncSynchronizeInputsWithLocalInputs(t *testing.T) {
	session := mocks.NewFakeSession()
	sessionCallbacks := mocks.MakeSessionCallBacks(session)

	peerConnection := []transport.UdpConnectStatus{
		{Disconnected: false, LastFrame: 12},
		{Disconnected: false, LastFrame: 13},
	}
	syncConfig := ggthx.NeweSyncConfig(
		sessionCallbacks, 8, 2, 4,
	)
	sync := ggthx.NewSync(peerConnection, &syncConfig)
	queue := 0
	input := input.GameInput{Bits: []byte{1, 2, 3, 4}}
	sync.AddLocalInput(queue, &input)

	inputs, disconnectFlags := sync.SynchronizeInputs()
	want := 2
	got := len(inputs)
	if want != got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
	if len(inputs[0]) == 0 {
		t.Errorf("expected input for local not to be zero.")
	}

	want = 0
	got = disconnectFlags
	if want != got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestSyncSynchronizeInputsWithRemoteInputs(t *testing.T) {
	session := mocks.NewFakeSession()
	sessionCallbacks := mocks.MakeSessionCallBacks(session)

	peerConnection := []transport.UdpConnectStatus{
		{Disconnected: false, LastFrame: 12},
		{Disconnected: false, LastFrame: 13},
	}
	syncConfig := ggthx.NeweSyncConfig(
		sessionCallbacks, 8, 2, 4,
	)
	sync := ggthx.NewSync(peerConnection, &syncConfig)
	queue := 1
	input := input.GameInput{Bits: []byte{1, 2, 3, 4}}
	sync.AddRemoteInput(queue, &input)

	inputs, _ := sync.SynchronizeInputs()
	want := []byte{1, 2, 3, 4}
	got := inputs[1]
	if !bytes.Equal(want, got) {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
	if len(inputs[1]) == 0 {
		t.Errorf("expected input for remote not to be zero.")
	}

}

func TestSyncSynchronizeInputsWithBothInputs(t *testing.T) {
	session := mocks.NewFakeSession()
	sessionCallbacks := mocks.MakeSessionCallBacks(session)

	peerConnection := []transport.UdpConnectStatus{
		{Disconnected: false, LastFrame: 12},
		{Disconnected: false, LastFrame: 13},
	}
	syncConfig := ggthx.NeweSyncConfig(
		sessionCallbacks, 8, 2, 4,
	)
	sync := ggthx.NewSync(peerConnection, &syncConfig)
	queue := 1
	input1 := input.GameInput{Bits: []byte{1, 2, 3, 4}}
	input2 := input.GameInput{Bits: []byte{5, 6, 7, 8}}
	sync.AddRemoteInput(queue, &input1)
	sync.AddLocalInput(0, &input2)

	inputs, _ := sync.SynchronizeInputs()
	want := []byte{1, 2, 3, 4}
	got := inputs[1]
	if !bytes.Equal(want, got) {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}

	want = []byte{5, 6, 7, 8}
	got = inputs[0]

	if !bytes.Equal(want, got) {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}

}

func TestSyncGetConfirmedInputs(t *testing.T) {
	session := mocks.NewFakeSession()
	sessionCallbacks := mocks.MakeSessionCallBacks(session)

	peerConnection := []transport.UdpConnectStatus{
		{Disconnected: false, LastFrame: 12},
		{Disconnected: false, LastFrame: 13},
	}
	syncConfig := ggthx.NeweSyncConfig(
		sessionCallbacks, 8, 2, 4,
	)
	sync := ggthx.NewSync(peerConnection, &syncConfig)
	queue := 1
	input1 := input.GameInput{Bits: []byte{1, 2, 3, 4}}
	input2 := input.GameInput{Bits: []byte{5, 6, 7, 8}}
	sync.AddRemoteInput(queue, &input1)
	sync.AddLocalInput(0, &input2)

	inputs, _ := sync.GetConfirmedInputs(0)
	want := []byte{1, 2, 3, 4}
	got := inputs[1]
	if !bytes.Equal(want, got) {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}

	want = []byte{5, 6, 7, 8}
	got = inputs[0]

	if !bytes.Equal(want, got) {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}

}

// Characterization Test
func TestSyncAddLocalInputPanic(t *testing.T) {
	session := mocks.NewFakeSession()
	sessionCallbacks := mocks.MakeSessionCallBacks(session)

	peerConnection := []transport.UdpConnectStatus{
		{Disconnected: false, LastFrame: 12},
		{Disconnected: false, LastFrame: 13},
	}
	syncConfig := ggthx.NeweSyncConfig(
		sessionCallbacks, 8, 2, 4,
	)
	sync := ggthx.NewSync(peerConnection, &syncConfig)
	queue := 1
	input1 := input.GameInput{Bits: []byte{1, 2, 3, 4}}
	input2 := input.GameInput{Bits: []byte{5, 6, 7, 8}}
	sync.AddRemoteInput(queue, &input1)
	sync.AddLocalInput(0, &input2)
	//sync.SetLastConfirmedFrame(8)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic when AddLocalInput attempted to add an input even when the frame hadn't been incremented.")
		}
	}()
	sync.AddLocalInput(0, &input2)

}

func TestSyncAddLocalInputNoPanic(t *testing.T) {
	session := mocks.NewFakeSession()
	sessionCallbacks := mocks.MakeSessionCallBacks(session)

	peerConnection := []transport.UdpConnectStatus{
		{Disconnected: false, LastFrame: 12},
		{Disconnected: false, LastFrame: 13},
	}
	syncConfig := ggthx.NeweSyncConfig(
		sessionCallbacks, 8, 2, 4,
	)
	sync := ggthx.NewSync(peerConnection, &syncConfig)
	input1 := input.GameInput{Bits: []byte{1, 2, 3, 4}}
	input2 := input.GameInput{Bits: []byte{5, 6, 7, 8}}

	sync.AddLocalInput(0, &input1)
	sync.IncrementFrame()
	sync.AddLocalInput(0, &input2)
}

func TestSyncAddRemoteInputPanic(t *testing.T) {
	session := mocks.NewFakeSession()
	sessionCallbacks := mocks.MakeSessionCallBacks(session)

	peerConnection := []transport.UdpConnectStatus{
		{Disconnected: false, LastFrame: 12},
		{Disconnected: false, LastFrame: 13},
	}
	syncConfig := ggthx.NeweSyncConfig(
		sessionCallbacks, 8, 2, 4,
	)
	sync := ggthx.NewSync(peerConnection, &syncConfig)
	input1 := input.GameInput{Bits: []byte{1, 2, 3, 4}}
	input2 := input.GameInput{Bits: []byte{5, 6, 7, 8}}
	sync.AddRemoteInput(1, &input2)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic when AddRemoteInput attempted to add an input even when the frame hadn't been incremented.")
		}
	}()
	sync.AddRemoteInput(1, &input1)
}

// Characterization mocks. No idea why this works and the above doesn't.
func TestSyncAddRemoteInputNoPanic(t *testing.T) {
	session := mocks.NewFakeSession()
	sessionCallbacks := mocks.MakeSessionCallBacks(session)

	peerConnection := []transport.UdpConnectStatus{
		{Disconnected: false, LastFrame: 12},
		{Disconnected: false, LastFrame: 13},
	}
	syncConfig := ggthx.NeweSyncConfig(
		sessionCallbacks, 8, 2, 4,
	)
	sync := ggthx.NewSync(peerConnection, &syncConfig)
	input1 := input.GameInput{Bits: []byte{1, 2, 3, 4}}
	input2 := input.GameInput{Bits: []byte{5, 6, 7, 8}}
	sync.AddRemoteInput(1, &input2)
	sync.AddLocalInput(0, &input1)
	sync.IncrementFrame()
	sync.AddLocalInput(0, &input1)
	sync.AddRemoteInput(1, &input1)
}
func TestSyncAddFrameDelay(t *testing.T) {
	session := mocks.NewFakeSession()
	sessionCallbacks := mocks.MakeSessionCallBacks(session)

	peerConnection := []transport.UdpConnectStatus{
		{Disconnected: false, LastFrame: 12},
		{Disconnected: false, LastFrame: 13},
	}
	syncConfig := ggthx.NeweSyncConfig(
		sessionCallbacks, 8, 2, 4,
	)
	sync := ggthx.NewSync(peerConnection, &syncConfig)
	input1 := input.GameInput{Bits: []byte{1, 2, 3, 4}}
	input2 := input.GameInput{Bits: []byte{5, 6, 7, 8}}
	frameDelay := 5
	sync.SetFrameDelay(0, 5)

	want := sync.FrameCount() + frameDelay
	sync.AddRemoteInput(1, &input2)
	sync.AddLocalInput(0, &input1)
	got := input1.Frame

	if want != got {
		t.Errorf("The Input delay was not applied correctly, expected input to be at frame %d but got %d", want, got)
	}
	/*
		sync.IncrementFrame()
		sync.AddLocalInput(0, &input)
		sync.AddRemoteInput(1, &input)
	*/
}
func TestSyncUseAfterClose(t *testing.T) {
	session := mocks.NewFakeSession()
	sessionCallbacks := mocks.MakeSessionCallBacks(session)

	peerConnection := []transport.UdpConnectStatus{
		{Disconnected: false, LastFrame: 12},
		{Disconnected: false, LastFrame: 13},
	}
	syncConfig := ggthx.NeweSyncConfig(
		sessionCallbacks, 8, 2, 4,
	)
	sync := ggthx.NewSync(peerConnection, &syncConfig)
	queue := 1
	input1 := input.GameInput{Bits: []byte{1, 2, 3, 4}}
	input2 := input.GameInput{Bits: []byte{5, 6, 7, 8}}
	sync.AddRemoteInput(queue, &input1)
	sync.AddLocalInput(0, &input2)
	sync.Close()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic when attempting to use sync after Close")
		}
	}()

	sync.AddLocalInput(0, &input2)
}
