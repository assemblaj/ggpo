package ggpo_test

import (
	"bytes"
	"math"
	"testing"
	"time"

	"github.com/assemblaj/GGPO-Go/internal/mocks"
	"github.com/assemblaj/GGPO-Go/internal/protocol"
	"github.com/assemblaj/GGPO-Go/pkg/transport"

	ggpo "github.com/assemblaj/GGPO-Go/pkg"
)

func TestNewSpectatorBackendSession(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer2PeerBackend(&session, "test", localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer2PeerBackend(&session2, "test", remotePort, numPlayers, inputSize)

	hostIp := "127.2.1.1"
	specPort := 6005
	stb := ggpo.NewSpectatorBackend(&session, "test", specPort, 2, 4, hostIp, localPort)

	connection := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p2, &stb}, localPort, remoteIp)
	connection2 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p}, remotePort, remoteIp)
	connection3 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p}, specPort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)
	stb.InitializeConnection(&connection3)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	spectator := ggpo.NewSpectator(20, remoteIp, specPort)
	var specHandle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)
	p2p.AddPlayer(&spectator, &specHandle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	stb.Start()
	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
		stb.Idle(0, advance)
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
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer2PeerBackend(&session, "test", localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer2PeerBackend(&session2, "test", remotePort, numPlayers, inputSize)

	hostIp := "127.2.1.1"
	specPort := 6005
	stb := ggpo.NewSpectatorBackend(&session, "test", specPort, 2, 4, hostIp, localPort)

	connection := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p2, &stb}, localPort, remoteIp)
	connection2 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p}, remotePort, remoteIp)
	connection3 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p}, specPort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)
	stb.InitializeConnection(&connection3)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	spectator := ggpo.NewSpectator(20, remoteIp, specPort)
	var specHandle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)
	p2p.AddPlayer(&spectator, &specHandle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	stb.Start()
	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
		stb.Idle(0, advance)
	}
	inputBytes := []byte{1, 2, 3, 4}
	inputBytes2 := []byte{5, 6, 7, 8}

	for i := 0; i < 2; i++ {
		p2p2.Idle(0)
		err := p2p2.AddLocalInput(p2Handle, inputBytes2, len(inputBytes2))
		if err != nil {
			t.Errorf(" Error when adding local input to p2, %s", err)
		}
		p2p2.AdvanceFrame()

		p2p.Idle(0)
		err = p2p.AddLocalInput(p1Handle, inputBytes, len(inputBytes))
		if err != nil {
			t.Errorf("Error when adding local input to p1, %s", err)
		}
		p2p.AdvanceFrame()

		stb.Idle(0)
		stb.AdvanceFrame()
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

/*WIP*/
func TestNewSpectatorBackendBehind(t *testing.T) {

	var p2p ggpo.Peer2PeerBackend
	session := mocks.NewFakeSessionWithBackend()
	session.SetBackend(&p2p)
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p = ggpo.NewPeer2PeerBackend(&session, "test", localPort, numPlayers, inputSize)

	var p2p2 ggpo.Peer2PeerBackend
	session2 := mocks.NewFakeSessionWithBackend()
	session2.SetBackend(&p2p2)

	p2p2 = ggpo.NewPeer2PeerBackend(&session2, "test", remotePort, numPlayers, inputSize)

	var stb ggpo.SpectatorBackend

	session3 := mocks.NewFakeSessionWithBackend()
	session3.SetBackend(&stb)

	hostIp := "127.2.1.1"
	specPort := 6005
	stb = ggpo.NewSpectatorBackend(&session3, "test", specPort, 2, 4, hostIp, localPort)

	connection := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p2, &stb}, localPort, remoteIp)
	connection2 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p}, remotePort, remoteIp)
	connection3 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p}, specPort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)
	stb.InitializeConnection(&connection3)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	spectator := ggpo.NewSpectator(20, remoteIp, specPort)
	var specHandle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)
	p2p.AddPlayer(&spectator, &specHandle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	stb.Start()
	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
		stb.Idle(0, advance)
	}
	inputBytes := []byte{1, 2, 3, 4}
	inputBytes2 := []byte{5, 6, 7, 8}
	var ignore int

	for i := 0; i < 10; i++ {
		p2p2.Idle(0)
		err := p2p2.AddLocalInput(p2Handle, inputBytes2, len(inputBytes2))
		if err != nil {
			t.Errorf(" Error when adding local input to p2, %s", err)
		}
		p2p2.SyncInput(&ignore)
		p2p2.AdvanceFrame()

		p2p.Idle(0)
		err = p2p.AddLocalInput(p1Handle, inputBytes, len(inputBytes))
		if err != nil {
			t.Errorf("Error when adding local input to p1, %s", err)
		}
		p2p.SyncInput(&ignore)
		p2p.AdvanceFrame()

		if i == 0 {
			stb.Idle(0)
			stb.SyncInput(&ignore)
			stb.AdvanceFrame()
		}
	}
	stb.Idle(0)
	stb.SyncInput(&ignore)
}

func TestNewSpectatorBackendCharacterization(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer2PeerBackend(&session, "test", localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer2PeerBackend(&session2, "test", remotePort, numPlayers, inputSize)

	hostIp := "127.2.1.1"
	specPort := 6005
	stb := ggpo.NewSpectatorBackend(&session, "test", specPort, 2, 4, hostIp, localPort)

	connection := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p2, &stb}, localPort, remoteIp)
	connection2 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p}, remotePort, remoteIp)
	connection3 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p}, specPort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)
	stb.InitializeConnection(&connection3)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	spectator := ggpo.NewSpectator(20, remoteIp, specPort)
	var specHandle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)
	p2p.AddPlayer(&spectator, &specHandle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	stb.Start()
	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
		stb.Idle(0, advance)
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
		p2p2.Idle(0, advance)
		err := p2p2.AddLocalInput(p2Handle, inputBytes2, len(inputBytes2))
		if err != nil {
			t.Errorf(" Error when adding local input to p2, %s", err)
		} else {
			_, err = p2p2.SyncInput(&ignore)
			if err == nil {
				p2p2.AdvanceFrame()
			}
		}

		p2p.Idle(0, advance)
		err = p2p.AddLocalInput(p1Handle, inputBytes, len(inputBytes))
		if err != nil {
			t.Errorf("Error when adding local input to p1, %s", err)
		} else {
			_, err = p2p.SyncInput(&ignore)
			if err == nil {
				p2p.AdvanceFrame()
			}
		}

		stb.Idle(0, advance)
		_, err = stb.SyncInput(&ignore)
		if err == nil {
			stb.AdvanceFrame()
		}
	}
}

func TestNewSpectatorBackendNoInputYet(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer2PeerBackend(&session, "test", localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer2PeerBackend(&session2, "test", remotePort, numPlayers, inputSize)

	hostIp := "127.2.1.1"
	specPort := 6005
	stb := ggpo.NewSpectatorBackend(&session, "test", specPort, 2, 4, hostIp, localPort)

	connection := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p2, &stb}, localPort, remoteIp)
	connection2 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p}, remotePort, remoteIp)
	connection3 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p}, specPort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)
	stb.InitializeConnection(&connection3)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	spectator := ggpo.NewSpectator(20, remoteIp, specPort)
	var specHandle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)
	p2p.AddPlayer(&spectator, &specHandle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	stb.Start()
	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
		stb.Idle(0, advance)
	}

	var ignore int
	stb.SyncInput(&ignore)
	stb.Idle(0)
	_, err := stb.SyncInput(&ignore)
	ggErr := err.(ggpo.Error)
	if ggErr.Code != ggpo.ErrorCodePredictionThreshod {
		t.Errorf("Should have recieved ErrorCodePredictionThreshold from spectator because no inputs had been sent yet. ")
	}
}

/* WIP, need to be able to test that the spectator is disconnected */
func TestNewSpectatorBackendDisconnect(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer2PeerBackend(&session, "test", localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer2PeerBackend(&session2, "test", remotePort, numPlayers, inputSize)

	hostIp := "127.2.1.1"
	specPort := 6005
	stb := ggpo.NewSpectatorBackend(&session, "test", specPort, 2, 4, hostIp, localPort)

	connection := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p2, &stb}, localPort, remoteIp)
	connection2 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p}, remotePort, remoteIp)
	connection3 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p}, specPort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)
	stb.InitializeConnection(&connection3)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	spectator := ggpo.NewSpectator(20, remoteIp, specPort)
	var specHandle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)
	p2p.AddPlayer(&spectator, &specHandle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	stb.Start()

	p2p.SetDisconnectTimeout(3000)
	p2p2.SetDisconnectTimeout(3000)
	p2p.SetDisconnectNotifyStart(1000)
	p2p2.SetDisconnectNotifyStart(1000)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
		stb.Idle(0, advance)
	}
	timeout := func() int64 {
		return time.Now().Add(time.Millisecond * 9500).UnixMilli()
	}
	var currentTime func() int64
	input1 := []byte{1, 2, 3, 4}
	input2 := []byte{5, 6, 7, 8}

	doPollTimeOuts := 0
	var p1now, p2now, p1next, p2next int
	p1now = int(time.Now().UnixMilli())
	p1next = p1now
	p2next = p1now
	p2now = p1now
	currentTime = advance
	for i := 0; i < ggpo.MaxPredictionFrames; i++ {
		doPollTimeOuts = int(math.Max(0, float64(p1next-p1now-1)))
		p2p.Idle(doPollTimeOuts, currentTime)
		if p1now >= p1next {
			err := p2p.AddLocalInput(p1Handle, input1, 4)
			if err == nil {
				//_, err = p2p.SyncInput(&ignore)
				if err == nil {
					p2p.AdvanceFrame()
				}
			}
			p1next = p1now + 1000/60
		}

		doPollTimeOuts = int(math.Max(0, float64(p2next-p2now-1)))
		p2p2.Idle(doPollTimeOuts, currentTime)
		if p2now >= p2next {
			err := p2p2.AddLocalInput(p2handle2, input2, 4)
			if err == nil {
				//_, err = p2p2.SyncInput(&ignore)
				if err == nil {
					p2p2.AdvanceFrame()
				}
			}
			p2next = p2now + 1000/60
		}

		if i == ggpo.MaxPredictionFrames-2 {
			currentTime = timeout
		}
	}

}

func TestNoAddingSpectatorAfterSynchronization(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer2PeerBackend(&session, "test", localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer2PeerBackend(&session2, "test", remotePort, numPlayers, inputSize)

	hostIp := "127.2.1.1"
	specPort := 6005
	stb := ggpo.NewSpectatorBackend(&session, "test", specPort, 2, 4, hostIp, localPort)

	connection := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p2, &stb}, localPort, remoteIp)
	connection2 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p}, remotePort, remoteIp)
	connection3 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p}, specPort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)
	stb.InitializeConnection(&connection3)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	spectator := ggpo.NewSpectator(20, remoteIp, specPort)
	var specHandle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	stb.Start()

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
	}

	err := p2p.AddPlayer(&spectator, &specHandle)
	if err == nil {
		t.Errorf("Spectators should not be able to be added after synchronization.")
	}
}

func TestSpectatorBackendDissconnectPlayerError(t *testing.T) {
	session := mocks.NewFakeSession()
	hostIp := "127.2.1.1"
	hostPort := 6001
	localPort := 6000
	stb := ggpo.NewSpectatorBackend(&session, "test", localPort, 2, 4, hostIp, hostPort)
	err := stb.DisconnectPlayer(ggpo.PlayerHandle(1))
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSpectatorBackendGetNetworkStatsError(t *testing.T) {
	session := mocks.NewFakeSession()
	hostIp := "127.2.1.1"
	hostPort := 6001
	localPort := 6000
	stb := ggpo.NewSpectatorBackend(&session, "test", localPort, 2, 4, hostIp, hostPort)
	_, err := stb.GetNetworkStats(ggpo.PlayerHandle(1))
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSpectatorBackendSetFrameDelayError(t *testing.T) {
	session := mocks.NewFakeSession()
	hostIp := "127.2.1.1"
	hostPort := 6001
	localPort := 6000
	stb := ggpo.NewSpectatorBackend(&session, "test", localPort, 2, 4, hostIp, hostPort)
	err := stb.SetFrameDelay(ggpo.PlayerHandle(1), 20)
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSpectatorBackendSetDisconnectTimeoutError(t *testing.T) {
	session := mocks.NewFakeSession()
	hostIp := "127.2.1.1"
	hostPort := 6001
	localPort := 6000
	stb := ggpo.NewSpectatorBackend(&session, "test", localPort, 2, 4, hostIp, hostPort)
	err := stb.SetDisconnectTimeout(20)
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}
func TestSpectatorBackendSetDisconnectNotifyStartError(t *testing.T) {
	session := mocks.NewFakeSession()
	hostIp := "127.2.1.1"
	hostPort := 6001
	localPort := 6000
	stb := ggpo.NewSpectatorBackend(&session, "test", localPort, 2, 4, hostIp, hostPort)
	err := stb.SetDisconnectNotifyStart(20)
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSpectatorBackendCloseError(t *testing.T) {
	session := mocks.NewFakeSession()
	hostIp := "127.2.1.1"
	hostPort := 6001
	localPort := 6000
	stb := ggpo.NewSpectatorBackend(&session, "test", localPort, 2, 4, hostIp, hostPort)
	err := stb.Close()
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}
func TestSpectatorBackendAddPlayerError(t *testing.T) {
	session := mocks.NewFakeSession()
	hostIp := "127.2.1.1"
	hostPort := 6001
	localPort := 6000
	stb := ggpo.NewSpectatorBackend(&session, "test", localPort, 2, 4, hostIp, hostPort)
	var player ggpo.Player
	var playerHandle ggpo.PlayerHandle
	err := stb.AddPlayer(&player, &playerHandle)
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}
