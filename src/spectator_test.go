package ggthx_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	ggthx "github.com/assemblaj/ggthx/src"
)

type FakeMultiplePeerConnection struct {
	remoteHandler []ggthx.MessageHandler
	localPort     int
	localIP       string
	printOutput   bool
}

func (f *FakeMultiplePeerConnection) SendTo(msg ggthx.UDPMessage, remoteIp string, remotePort int) {
	if f.printOutput {
		fmt.Printf("f.localIP %s f.localPort %d msg %s size %d\n", f.localIP, f.localPort, msg, msg.PacketSize())
	}
	for _, r := range f.remoteHandler {
		r.HandleMessage(f.localIP, f.localPort, msg, msg.PacketSize())
	}
}

func (f *FakeMultiplePeerConnection) Read() {
}

func (f *FakeMultiplePeerConnection) Close() {

}

func NewFakeMultiplePeerConnection(remoteHandler []ggthx.MessageHandler, localPort int, localIP string) FakeMultiplePeerConnection {
	f := FakeMultiplePeerConnection{}
	f.remoteHandler = remoteHandler
	f.localPort = localPort
	f.localIP = localIP
	return f
}

func TestNewSpectatorBackendSession(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggthx.NewPeer2PeerBackend(&sessionCallbacks, "test", localPort, numPlayers, inputSize)

	session2 := NewFakeSession()
	sessionCallbacks2 := makeSessionCallBacks(session2)
	p2p2 := ggthx.NewPeer2PeerBackend(&sessionCallbacks2, "test", remotePort, numPlayers, inputSize)

	hostIp := "127.2.1.1"
	specPort := 6005
	stb := ggthx.NewSpectatorBackend(&sessionCallbacks, "test", specPort, 2, 4, hostIp, localPort)

	connection := NewFakeMultiplePeerConnection([]ggthx.MessageHandler{&p2p2, &stb}, localPort, remoteIp)
	connection2 := NewFakeMultiplePeerConnection([]ggthx.MessageHandler{&p2p}, remotePort, remoteIp)
	connection3 := NewFakeMultiplePeerConnection([]ggthx.MessageHandler{&p2p}, specPort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)
	stb.InitalizeConnection(&connection3)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	spectator := ggthx.NewSpectator(20, remoteIp, specPort)
	var specHandle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)
	p2p.AddPlayer(&spectator, &specHandle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	stb.Start()
	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		p2p.DoPoll(0, advance)
		p2p2.DoPoll(0, advance)
		stb.DoPoll(0, advance)
	}
	inputBytes := []byte{1, 2, 3, 4}
	err := p2p.AddLocalInput(p1Handle, inputBytes, len(inputBytes))
	if err != nil {
		t.Errorf("Error when adding local input, %s", err)
	}

	err = p2p2.AddLocalInput(p2handle2, inputBytes, len(inputBytes))
	if err != nil {
		t.Errorf("Error when adding local input, %s", err)
	}
	var disconnectFlags int
	_, err = stb.SyncInput(&disconnectFlags)
	if err != nil {
		t.Errorf("Error when synchronizing input on spectator, %s", err)
	}

}
func TestNewSpectatorBackendInput(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggthx.NewPeer2PeerBackend(&sessionCallbacks, "test", localPort, numPlayers, inputSize)

	session2 := NewFakeSession()
	sessionCallbacks2 := makeSessionCallBacks(session2)
	p2p2 := ggthx.NewPeer2PeerBackend(&sessionCallbacks2, "test", remotePort, numPlayers, inputSize)

	hostIp := "127.2.1.1"
	specPort := 6005
	stb := ggthx.NewSpectatorBackend(&sessionCallbacks, "test", specPort, 2, 4, hostIp, localPort)

	connection := NewFakeMultiplePeerConnection([]ggthx.MessageHandler{&p2p2, &stb}, localPort, remoteIp)
	connection2 := NewFakeMultiplePeerConnection([]ggthx.MessageHandler{&p2p}, remotePort, remoteIp)
	connection3 := NewFakeMultiplePeerConnection([]ggthx.MessageHandler{&p2p}, specPort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)
	stb.InitalizeConnection(&connection3)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	spectator := ggthx.NewSpectator(20, remoteIp, specPort)
	var specHandle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)
	p2p.AddPlayer(&spectator, &specHandle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	stb.Start()
	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		p2p.DoPoll(0, advance)
		p2p2.DoPoll(0, advance)
		stb.DoPoll(0, advance)
	}
	inputBytes := []byte{1, 2, 3, 4}
	inputBytes2 := []byte{5, 6, 7, 8}

	for i := 0; i < 2; i++ {
		p2p2.DoPoll(0)
		err := p2p2.AddLocalInput(p2Handle, inputBytes2, len(inputBytes2))
		if err != nil {
			t.Errorf(" Error when adding local input to p2, %s", err)
		}
		p2p2.IncrementFrame()

		p2p.DoPoll(0)
		err = p2p.AddLocalInput(p1Handle, inputBytes, len(inputBytes))
		if err != nil {
			t.Errorf("Error when adding local input to p1, %s", err)
		}
		p2p.IncrementFrame()

		stb.DoPoll(0)
		stb.IncrementFrame()
	}
	var ignore int
	vals, err := stb.SyncInput(&ignore)
	if err != nil {
		t.Errorf("Error when spectator synchronize inputs. %s", err)
	}
	if !bytes.Equal(inputBytes, vals[0]) {
		t.Errorf("Returned p1 input %v doesn't equal given p1 input %v", vals[0], inputBytes)
	}
	if !bytes.Equal(inputBytes2, vals[1]) {
		t.Errorf("Returned p1 input %v doesn't equal given p1 input %v", vals[1], inputBytes2)
	}

}
func TestNewSpectatorBackendCharacterization(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggthx.NewPeer2PeerBackend(&sessionCallbacks, "test", localPort, numPlayers, inputSize)

	session2 := NewFakeSession()
	sessionCallbacks2 := makeSessionCallBacks(session2)
	p2p2 := ggthx.NewPeer2PeerBackend(&sessionCallbacks2, "test", remotePort, numPlayers, inputSize)

	hostIp := "127.2.1.1"
	specPort := 6005
	stb := ggthx.NewSpectatorBackend(&sessionCallbacks, "test", specPort, 2, 4, hostIp, localPort)

	connection := NewFakeMultiplePeerConnection([]ggthx.MessageHandler{&p2p2, &stb}, localPort, remoteIp)
	connection2 := NewFakeMultiplePeerConnection([]ggthx.MessageHandler{&p2p}, remotePort, remoteIp)
	connection3 := NewFakeMultiplePeerConnection([]ggthx.MessageHandler{&p2p}, specPort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)
	stb.InitalizeConnection(&connection3)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	spectator := ggthx.NewSpectator(20, remoteIp, specPort)
	var specHandle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)
	p2p.AddPlayer(&spectator, &specHandle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	stb.Start()
	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		p2p.DoPoll(0, advance)
		p2p2.DoPoll(0, advance)
		stb.DoPoll(0, advance)
	}
	inputBytes := []byte{1, 2, 3, 4}
	inputBytes2 := []byte{5, 6, 7, 8}
	var ignore int
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic when framecount hadn't been incremented..")
		}
	}()

	for i := 0; i < 2; i++ {
		p2p2.DoPoll(0, advance)
		err := p2p2.AddLocalInput(p2Handle, inputBytes2, len(inputBytes2))
		if err != nil {
			t.Errorf(" Error when adding local input to p2, %s", err)
		} else {
			_, err = p2p2.SyncInput(&ignore)
			if err == nil {
				p2p2.IncrementFrame()
			}
		}

		p2p.DoPoll(0, advance)
		err = p2p.AddLocalInput(p1Handle, inputBytes, len(inputBytes))
		if err != nil {
			t.Errorf("Error when adding local input to p1, %s", err)
		} else {
			_, err = p2p.SyncInput(&ignore)
			if err == nil {
				p2p.IncrementFrame()
			}
		}

		stb.DoPoll(0, advance)
		_, err = stb.SyncInput(&ignore)
		if err == nil {
			stb.IncrementFrame()
		}
	}
}

func TestNewSpectatorBackendNoInputYet(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggthx.NewPeer2PeerBackend(&sessionCallbacks, "test", localPort, numPlayers, inputSize)

	session2 := NewFakeSession()
	sessionCallbacks2 := makeSessionCallBacks(session2)
	p2p2 := ggthx.NewPeer2PeerBackend(&sessionCallbacks2, "test", remotePort, numPlayers, inputSize)

	hostIp := "127.2.1.1"
	specPort := 6005
	stb := ggthx.NewSpectatorBackend(&sessionCallbacks, "test", specPort, 2, 4, hostIp, localPort)

	connection := NewFakeMultiplePeerConnection([]ggthx.MessageHandler{&p2p2, &stb}, localPort, remoteIp)
	connection2 := NewFakeMultiplePeerConnection([]ggthx.MessageHandler{&p2p}, remotePort, remoteIp)
	connection3 := NewFakeMultiplePeerConnection([]ggthx.MessageHandler{&p2p}, specPort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)
	stb.InitalizeConnection(&connection3)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	spectator := ggthx.NewSpectator(20, remoteIp, specPort)
	var specHandle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)
	p2p.AddPlayer(&spectator, &specHandle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	stb.Start()
	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		p2p.DoPoll(0, advance)
		p2p2.DoPoll(0, advance)
		stb.DoPoll(0, advance)
	}

	var ignore int
	stb.SyncInput(&ignore)
	stb.DoPoll(0)
	_, err := stb.SyncInput(&ignore)
	ggErr := err.(ggthx.Error)
	if ggErr.Code != ggthx.ErrorCodePredictionThreshod {
		t.Errorf("Should have recieved ErrorCodePredictionThreshold from spectator because no inputs had been sent yet. ")
	}
}

func TestSpectatorBackendChatError(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	hostIp := "127.2.1.1"
	hostPort := 6001
	localPort := 6000
	stb := ggthx.NewSpectatorBackend(&sessionCallbacks, "test", localPort, 2, 4, hostIp, hostPort)
	err := stb.Chat("test")
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSpectatorBackendDissconnectPlayerError(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	hostIp := "127.2.1.1"
	hostPort := 6001
	localPort := 6000
	stb := ggthx.NewSpectatorBackend(&sessionCallbacks, "test", localPort, 2, 4, hostIp, hostPort)
	err := stb.DisconnectPlayer(ggthx.PlayerHandle(1))
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSpectatorBackendGetNetworkStatsError(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	hostIp := "127.2.1.1"
	hostPort := 6001
	localPort := 6000
	stb := ggthx.NewSpectatorBackend(&sessionCallbacks, "test", localPort, 2, 4, hostIp, hostPort)
	var status ggthx.NetworkStats
	err := stb.GetNetworkStats(&status, ggthx.PlayerHandle(1))
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSpectatorBackendLogvError(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	hostIp := "127.2.1.1"
	hostPort := 6001
	localPort := 6000
	stb := ggthx.NewSpectatorBackend(&sessionCallbacks, "test", localPort, 2, 4, hostIp, hostPort)
	err := stb.Logv("test")
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSpectatorBackendSetFrameDelayError(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	hostIp := "127.2.1.1"
	hostPort := 6001
	localPort := 6000
	stb := ggthx.NewSpectatorBackend(&sessionCallbacks, "test", localPort, 2, 4, hostIp, hostPort)
	err := stb.SetFrameDelay(ggthx.PlayerHandle(1), 20)
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSpectatorBackendSetDisconnectTimeoutError(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	hostIp := "127.2.1.1"
	hostPort := 6001
	localPort := 6000
	stb := ggthx.NewSpectatorBackend(&sessionCallbacks, "test", localPort, 2, 4, hostIp, hostPort)
	err := stb.SetDisconnectTimeout(20)
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSpectatorBackendCloseError(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	hostIp := "127.2.1.1"
	hostPort := 6001
	localPort := 6000
	stb := ggthx.NewSpectatorBackend(&sessionCallbacks, "test", localPort, 2, 4, hostIp, hostPort)
	err := stb.Close()
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}
func TestSpectatorBackendAddPlayerError(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	hostIp := "127.2.1.1"
	hostPort := 6001
	localPort := 6000
	stb := ggthx.NewSpectatorBackend(&sessionCallbacks, "test", localPort, 2, 4, hostIp, hostPort)
	var player ggthx.Player
	var playerHandle ggthx.PlayerHandle
	err := stb.AddPlayer(&player, &playerHandle)
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}
