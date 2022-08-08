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

func slice2dEqual(a [][]byte, b [][]byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !bytes.Equal(a[i], b[i]) {
			return false
		}
	}
	return true
}

func TestP2PBackendAddPlayer(t *testing.T) {
	connection := mocks.NewFakeConnection()
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)
	p2p.InitializeConnection(&connection)
	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	err := p2p.AddPlayer(&player1, &p1Handle)
	if err != nil {
		t.Errorf("There was an error when adding player 1.")
	}
	err = p2p.AddPlayer(&player2, &p2Handle)
	if err != nil {
		t.Errorf("There was an error when adding player 2.")
	}

}

func TestP2PBackendAddLocalInputError(t *testing.T) {
	connection := mocks.NewFakeConnection()
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)
	p2p.InitializeConnection(&connection)
	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	err := p2p.AddPlayer(&player1, &p1Handle)
	if err != nil {
		t.Errorf("There was an error when adding player 1.")
	}
	err = p2p.AddPlayer(&player2, &p2Handle)
	if err != nil {
		t.Errorf("There was an error when adding player 2.")
	}
	inputBytes := []byte{1, 2, 3, 4}
	err = p2p.AddLocalInput(p1Handle, inputBytes, len(inputBytes))
	if err == nil {
		t.Errorf("There should be an error when attempting to add local input while still synchronizing")
	}
}

func TestP2PBackendSyncInputError(t *testing.T) {
	connection := mocks.NewFakeConnection()
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)
	p2p.InitializeConnection(&connection)
	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	err := p2p.AddPlayer(&player1, &p1Handle)
	if err != nil {
		t.Errorf("There was an error when adding player 1.")
	}
	err = p2p.AddPlayer(&player2, &p2Handle)
	if err != nil {
		t.Errorf("There was an error when adding player 2.")
	}
	var disconnectFlags int
	_, err = p2p.SyncInput(&disconnectFlags)
	if err == nil {
		t.Errorf("There should be an error when attempting to synchrinoze input while still synchronizing")
	}
}

func TestP2PBackendIncrementFrame(t *testing.T) {
	connection := mocks.NewFakeConnection()
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)
	p2p.InitializeConnection(&connection)
	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	err := p2p.AddPlayer(&player1, &p1Handle)
	if err != nil {
		t.Errorf("There was an error when adding player 1.")
	}
	err = p2p.AddPlayer(&player2, &p2Handle)
	if err != nil {
		t.Errorf("There was an error when adding player 2.")
	}
	err = p2p.AdvanceFrame()
	if err != nil {
		t.Errorf("There was an error when incrementing the frame.")
	}
}

func TestP2PBackendSynchronizeInputs(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)
	connection := mocks.NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := mocks.NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)
	inputBytes := []byte{1, 2, 3, 4}
	p2p.Idle(0)
	p2p2.Idle(0)

	err := p2p2.AddLocalInput(p2Handle, inputBytes, len(inputBytes))
	if err != nil {
		t.Errorf("The backends didn't synchronize")
	}
	var discconectFlags int
	vals, err := p2p2.SyncInput(&discconectFlags)
	if err != nil {
		t.Errorf("Synchronize Input returned an error ")
	}
	want := inputBytes
	got := vals[1]
	if !bytes.Equal(want, got) {
		t.Errorf("expected '%v' but got '%v'", want, got)
	}
}

func TestP2PBackendCharacterizationAddLocalInput(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)
	connection := mocks.NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := mocks.NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)
	inputBytes := []byte{1, 2, 3, 4}
	p2p.Idle(0)

	p2p2.Idle(0)
	p2p2.AddLocalInput(p2Handle, inputBytes, len(inputBytes))
	p2p2.Idle(0)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic due to an InputQueue error")
		}
	}()

	p2p2.AddLocalInput(p2Handle, inputBytes, len(inputBytes))

}

func TestP2PBackendPoll2PlayersDefault(t *testing.T) {
	connection := mocks.NewFakeConnection()
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)
	p2p.InitializeConnection(&connection)
	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)
	result := p2p.Poll2Players(0)
	want := -1
	got := result
	if want != got {
		t.Errorf("expected '%d' but got '%d'", want, got)
	}
}

func TestP2PBackendPollNPlayersDefault(t *testing.T) {
	connection := mocks.NewFakeConnection()
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 3
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)
	p2p.InitializeConnection(&connection)
	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	player3 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort+1)
	var p2Handle ggpo.PlayerHandle
	var p3Handle ggpo.PlayerHandle

	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)
	p2p.AddPlayer(&player3, &p3Handle)
	result := p2p.PollNPlayers(0)
	want := -1
	got := result
	if want != got {
		t.Errorf("expected '%d' but got '%d'", want, got)
	}
}

func TestP2PBackendAddLocalInputMultiple(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)
	connection := mocks.NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := mocks.NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)
	cycles := 20
	inputBytes := []byte{1, 2, 3, 4}
	for i := 0; i < cycles; i++ {
		p2p2.Idle(0)
		p2p2.AdvanceFrame()
		p2p2.AddLocalInput(p2handle2, inputBytes, len(inputBytes))
	}
	var discconectFlags int
	values, err := p2p2.SyncInput(&discconectFlags)
	if err != nil {
		t.Errorf("Got an error when synchronizing input.")
	}
	want := inputBytes
	got := values[1]
	if !bytes.Equal(want, got) {
		t.Errorf("expected '%v' but got '%v'", want, got)
	}
}

func TestP2PBackendSynchronize(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)
	connection := mocks.NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := mocks.NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)
	inputBytes := []byte{1, 2, 3, 4}
	p2p.Idle(0)
	p2p2.Idle(0)

	err := p2p2.AddLocalInput(p2Handle, inputBytes, len(inputBytes))
	if err != nil {
		t.Errorf("The backends didn't synchronize")
	}
}

func TestP2PBackendFullSession(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)
	connection := mocks.NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := mocks.NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)
	inputBytes := []byte{1, 2, 3, 4}

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}

	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
	}

	for i := 0; i < 20; i++ {
		p2p2.Idle(0)
		err := p2p2.AddLocalInput(p2Handle, inputBytes, len(inputBytes))
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
	}
	var disconnectFlags int
	vals, err := p2p.SyncInput(&disconnectFlags)
	if err != nil {
		t.Errorf("Error when synchronizing input on p1, %s", err)
	}
	vals2, err2 := p2p2.SyncInput(&disconnectFlags)
	if err != nil {
		t.Errorf("Error when synchronizing input on p2, %s", err2)
	}
	if len(vals) != len(vals2) {
		t.Errorf("Error, lengths of synchronized input not equal.")
	}
	for i := 0; i < len(vals); i++ {
		if !bytes.Equal(vals[i], vals2[i]) {
			t.Errorf("Error, Expected synchronized input to be the same, input %d for p1 is %v, p2 %v",
				i+1, vals[i], vals[2])
		}
	}
	err = p2p2.DisconnectPlayer(p2handle1)
	if err != nil {
		t.Errorf("Disconnecting player caused error %s ", err)
	}
	err = p2p2.DisconnectPlayer(p2handle2)
	if err != nil {
		t.Errorf("Disconnecting player caused error %s ", err)
	}
	err = p2p.DisconnectPlayer(p1Handle)
	if err != nil {
		t.Errorf("Disconnecting player caused error %s ", err)
	}
}

func TestP2PBackendDisconnectPlayerLocal(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)
	connection := mocks.NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := mocks.NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	err := p2p2.DisconnectPlayer(p2handle2)
	if err != nil {
		t.Errorf("Had an error trying to disconnect the local player.")
	}
}

func TestP2PBackendDisconnectPlayerRemoteCharacterization(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)
	connection := mocks.NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := mocks.NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic due attempting to load frames when none had been saved.")
		}
	}()
	p2p2.DisconnectPlayer(p2handle1)
}

func TestP2PBackendDisconnectPlayerError(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)
	connection := mocks.NewFakeConnection()
	p2p.InitializeConnection(&connection)

	err := p2p.DisconnectPlayer(ggpo.PlayerHandle(8))
	if err == nil {
		t.Errorf("The code should have created an error when passing an invalid player handle into DisconnectPlayer")
	}
}
func TestP2PBackendMockSynchronize(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)
	connection := mocks.NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := mocks.NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	var err1, err2 error

	for i := 0; i < protocol.NumSyncPackets; i++ {
		var disconnectFlags int
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
		_, err1 = p2p.SyncInput(&disconnectFlags)
		if err1 != nil {
			continue
		}
		_, err2 = p2p2.SyncInput(&disconnectFlags)
		if err2 != nil {
			continue
		}
		if err2 == nil && err1 == nil {
			break
		}
	}
	if err1 != nil {
		t.Errorf("The players (specifically p1) did not synchronize during the sync period.")
	}
	if err2 != nil {
		t.Errorf("The players (specifically p2) did not synchronize during the sync period")
	}
}
func TestP2PBackendMoockInputExchangeCharacterization(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)
	connection := mocks.NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := mocks.NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
	}
	input1 := []byte{1, 2, 3, 4}
	input2 := []byte{5, 6, 7, 8}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic.")
		}
	}()

	for i := 0; i < 8; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
		p2p.AddLocalInput(p1Handle, input1, 4)
		p2p2.AddLocalInput(p2handle2, input2, 4)
		var disconnectFlags int
		p2p2.SyncInput(&disconnectFlags)
		p2p.SyncInput(&disconnectFlags)
		p2p.AdvanceFrame()
		p2p2.AdvanceFrame()
	}

}
func TestP2PBackendMoockInputExchange(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)
	connection := mocks.NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := mocks.NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
	}
	input1 := []byte{1, 2, 3, 4}
	input2 := []byte{5, 6, 7, 8}

	for i := 0; i < 8; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
		p2p.AddLocalInput(p1Handle, input1, 4)
		p2p2.AddLocalInput(p2handle2, input2, 4)
		p2p.AdvanceFrame()
		p2p2.AdvanceFrame()
	}
	var disconnectFlags int
	vals, _ := p2p2.SyncInput(&disconnectFlags)
	vals2, _ := p2p.SyncInput(&disconnectFlags)
	if len(vals) != len(vals2) {
		t.Errorf("Inputs should be synchronized between the 2 inputs")
	}

	for i := 0; i < len(vals2); i++ {
		if !bytes.Equal(vals2[i], vals[i]) {
			t.Errorf("Expected %v and %v to be equal", vals2[i], vals[i])
		}
	}
}
func TestP2PBackendMoockInputExchangeWithTimeout(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)
	connection := mocks.NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := mocks.NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
	}
	input1 := []byte{1, 2, 3, 4}
	input2 := []byte{5, 6, 7, 8}
	doPollTimeOuts := 90
	for i := 0; i < 8; i++ {
		p2p.Idle(doPollTimeOuts, advance)
		p2p2.Idle(doPollTimeOuts, advance)
		p2p.AddLocalInput(p1Handle, input1, 4)
		p2p2.AddLocalInput(p2handle2, input2, 4)
		p2p.AdvanceFrame()
		p2p2.AdvanceFrame()
	}
	var disconnectFlags int
	vals, _ := p2p2.SyncInput(&disconnectFlags)
	vals2, _ := p2p.SyncInput(&disconnectFlags)
	if len(vals) != len(vals2) {
		t.Errorf("Inputs should be synchronized between the 2 inputs")
	}

	for i := 0; i < len(vals2); i++ {
		if !bytes.Equal(vals2[i], vals[i]) {
			t.Errorf("Expected %v and %v to be equal", vals2[i], vals[i])
		}
	}
}
func TestP2PBackendMoockInputExchangePol2Players(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)
	connection := mocks.NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := mocks.NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
	}
	input1 := []byte{1, 2, 3, 4}
	input2 := []byte{5, 6, 7, 8}
	doPollTimeOuts := 90
	for i := 0; i < 8; i++ {
		p2p.Idle(doPollTimeOuts, advance)
		p2p2.Idle(doPollTimeOuts, advance)
		p2p.AddLocalInput(p1Handle, input1, 4)
		p2p2.AddLocalInput(p2handle2, input2, 4)
		p2p.AdvanceFrame()
		p2p2.AdvanceFrame()
	}
	want := 7
	got := p2p2.Poll2Players(8)
	if want != got {
		t.Errorf("wanted %d got %d t", want, got)

	}

	got = p2p.Poll2Players(8)
	if want != got {
		t.Errorf("wanted %d got %d ", want, got)
	}

}
func TestP2PBackendMoockInputDelay(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)
	connection := mocks.NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := mocks.NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
	}

	p2p.SetFrameDelay(p1Handle, 2)
	p2p2.SetFrameDelay(p2handle2, 2)

	/*
		p2p.SetDisconnectTimeout(3000)
		p2p.SetDisconnectNotifyStart(1000)
		p2p2.SetDisconnectTimeout(3000)
		p2p2.SetDisconnectNotifyStart(1000) */

	input1 := []byte{1, 2, 3, 4}
	input2 := []byte{5, 6, 7, 8}
	doPollTimeOuts := 90
	iterations := 6
	for i := 0; i < iterations; i++ {
		p2p.Idle(doPollTimeOuts, advance)
		p2p2.Idle(doPollTimeOuts, advance)
		p2p.AddLocalInput(p1Handle, input1, 4)
		p2p2.AddLocalInput(p2handle2, input2, 4)
		p2p.AdvanceFrame()
		p2p2.AdvanceFrame()
	}
	got := p2p.Poll2Players(iterations)
	want := iterations + 1
	if want != got {
		t.Errorf("wanted %d got %d ", want, got)
	}

}

func TestP2PBackendMoockDisconnectTimeoutCharacterization(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)
	connection := mocks.NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := mocks.NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
	}

	timeout := func() int64 {
		return time.Now().Add(time.Millisecond * 3000).UnixMilli()
	}

	p2p.SetDisconnectTimeout(3000)
	//	p2p.SetDisconnectNotifyStart(1000)
	p2p2.SetDisconnectTimeout(3000)
	//	p2p2.SetDisconnectNotifyStart(1000) */

	doPollTimeOuts := 90
	iterations := 6
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic.")
		}
	}()
	for i := 0; i < iterations; i++ {
		p2p.Idle(doPollTimeOuts, timeout)
		p2p2.Idle(doPollTimeOuts, timeout)
		//p2p.AddLocalInput(p1Handle, input1, 4)
		//p2p2.AddLocalInput(p2handle2, input2, 4)
		//p2p.AdvanceFrame()
		//p2p2.AdvanceFrame()
	}

}
func TestP2PBackendMoockDisconnectTimeoutCharacterization2(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)
	connection := mocks.NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := mocks.NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
	}

	timeout := func() int64 {
		return time.Now().Add(time.Millisecond * 4000).UnixMilli()
	}

	p2p.SetDisconnectTimeout(3000)
	p2p2.SetDisconnectTimeout(3000)
	doPollTimeOuts := 0
	iterations := 2
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic due to trying to load frame that hadn't been saved.")
		}
	}()
	for i := 0; i < iterations; i++ {

		p2p.Idle(doPollTimeOuts, timeout)
		p2p2.Idle(doPollTimeOuts, timeout)
		//p2p.AddLocalInput(p1Handle, input1, 4)
		//p2p2.AddLocalInput(p2handle2, input2, 4)
		p2p.AdvanceFrame()
		p2p2.AdvanceFrame()
	}
}
func TestP2PBackendMoockDisconnectTimeout(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)
	connection := mocks.NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := mocks.NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
	}

	timeout := func() int64 {
		return time.Now().Add(time.Millisecond * 9500).UnixMilli()
	}
	var currentTime func() int64

	currentTime = advance
	p2p.SetDisconnectTimeout(3000)
	p2p2.SetDisconnectTimeout(3000)
	doPollTimeOuts := 0
	//var ignore int
	input1 := []byte{1, 2, 3, 4}
	input2 := []byte{5, 6, 7, 8}
	var p1now, p2now, p1next, p2next int
	p1now = int(time.Now().UnixMilli())
	p1next = p1now
	p2next = p1now
	p2now = p1now

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

	err := p2p.DisconnectPlayer(p2Handle)
	ggError := err.(ggpo.Error)
	if ggError.Code != ggpo.ErrorCodePlayerDisconnected {
		t.Errorf("The player should've been timed out and disconnected already.")
	}

	err = p2p2.DisconnectPlayer(p2handle1)
	ggError = err.(ggpo.Error)
	if ggError.Code != ggpo.ErrorCodePlayerDisconnected {
		t.Errorf("The player should've been timed out and disconnected already.")
	}
}
func TestP2PBackendNPlayersSynchronize(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 3
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)

	session3 := mocks.NewFakeSession()

	p3port := 6005
	p2p3 := ggpo.NewPeer(&session3, p3port, numPlayers, inputSize)

	connection := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p2, &p2p3}, localPort, remoteIp)
	connection2 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p, &p2p3}, remotePort, remoteIp)
	connection3 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p, &p2p2}, p3port, remoteIp)
	//ggpo.EnableLogger()
	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)
	p2p3.InitializeConnection(&connection3)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	player3 := ggpo.NewRemotePlayer(20, 3, remoteIp, p3port)
	var p3Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)
	p2p.AddPlayer(&player3, &p3Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	player3 = ggpo.NewRemotePlayer(20, 3, remoteIp, p3port)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	var p2handle3 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)
	p2p2.AddPlayer(&player3, &p2handle3)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	player3 = ggpo.NewLocalPlayer(20, 3)
	var p3handle1 ggpo.PlayerHandle
	var p3handle2 ggpo.PlayerHandle
	var p3handle3 ggpo.PlayerHandle
	p2p3.AddPlayer(&player1, &p3handle1)
	p2p3.AddPlayer(&player2, &p3handle2)
	p2p3.AddPlayer(&player3, &p3handle3)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
		p2p3.Idle(0, advance)
	}
	input1 := []byte{1, 2, 3, 4}
	input2 := []byte{5, 6, 7, 8}
	input3 := []byte{9, 10, 11, 12}
	err := p2p.AddLocalInput(p1Handle, input1, 4)
	if err != nil {
		t.Errorf("Peer 1 of 3 did not synchronize.")
	}
	err = p2p2.AddLocalInput(p2handle2, input2, 4)
	if err != nil {
		t.Errorf("Peer 2 of 3 did not synchronize.")
	}
	err = p2p3.AddLocalInput(p3handle3, input3, 4)
	if err != nil {
		t.Errorf("Peer 3 of 3 did not synchronize.")
	}
}

func TestP2PBackendNPlayersShareInput(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 3
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)

	session3 := mocks.NewFakeSession()

	p3port := 6005
	p2p3 := ggpo.NewPeer(&session3, p3port, numPlayers, inputSize)

	connection := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p2, &p2p3}, localPort, remoteIp)
	connection2 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p, &p2p3}, remotePort, remoteIp)
	connection3 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p, &p2p2}, p3port, remoteIp)
	//ggpo.EnableLogger()
	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)
	p2p3.InitializeConnection(&connection3)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	player3 := ggpo.NewRemotePlayer(20, 3, remoteIp, p3port)
	var p3Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)
	p2p.AddPlayer(&player3, &p3Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	player3 = ggpo.NewRemotePlayer(20, 3, remoteIp, p3port)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	var p2handle3 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)
	p2p2.AddPlayer(&player3, &p2handle3)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	player3 = ggpo.NewLocalPlayer(20, 3)
	var p3handle1 ggpo.PlayerHandle
	var p3handle2 ggpo.PlayerHandle
	var p3handle3 ggpo.PlayerHandle
	p2p3.AddPlayer(&player1, &p3handle1)
	p2p3.AddPlayer(&player2, &p3handle2)
	p2p3.AddPlayer(&player3, &p3handle3)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
		p2p3.Idle(0, advance)
	}
	input1 := []byte{1, 2, 3, 4}
	input2 := []byte{5, 6, 7, 8}
	input3 := []byte{9, 10, 11, 12}
	var err error
	var ignore int

	for i := 0; i < 2; i++ {
		p2p.Idle(0, advance)
		err = p2p.AddLocalInput(p1Handle, input1, 4)
		if err == nil {
			//	_, err = p2p.SyncInput(&ignore)
			if err == nil {
				p2p.AdvanceFrame()
			}
		}
		p2p2.Idle(0, advance)
		err = p2p2.AddLocalInput(p2handle2, input2, 4)
		if err == nil {
			//	_, err = p2p2.SyncInput(&ignore)
			if err == nil {
				p2p2.AdvanceFrame()
			}
		}
		p2p3.Idle(0, advance)
		err = p2p3.AddLocalInput(p3handle3, input3, 4)
		if err == nil {
			//	_, err = p2p3.SyncInput(&ignore)
			if err == nil {
				p2p3.AdvanceFrame()
			}
		}
	}
	val1, _ := p2p.SyncInput(&ignore)
	val2, _ := p2p2.SyncInput(&ignore)
	val3, _ := p2p3.SyncInput(&ignore)
	if !slice2dEqual(val1, val2) || !slice2dEqual(val2, val3) || !slice2dEqual(val1, val3) {
		t.Errorf("All peers did not recieve the inputs. P1 inputs %v P2 inputs %v P3 inputs %v", val1, val2, val3)
	}
}

func TestP2PBackend4PlayerSynchronize(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 4
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)

	session3 := mocks.NewFakeSession()

	p3port := 6005
	p2p3 := ggpo.NewPeer(&session3, p3port, numPlayers, inputSize)

	session4 := mocks.NewFakeSession()
	p4port := 6006
	p2p4 := ggpo.NewPeer(&session4, p4port, numPlayers, inputSize)

	connection := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p2, &p2p3, &p2p4}, localPort, remoteIp)
	connection2 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p, &p2p3, &p2p4}, remotePort, remoteIp)
	connection3 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p, &p2p2, &p2p4}, p3port, remoteIp)
	connection4 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p, &p2p2, &p2p3}, p4port, remoteIp)
	//ggpo.EnableLogger()
	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)
	p2p3.InitializeConnection(&connection3)
	p2p4.InitializeConnection(&connection4)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	player3 := ggpo.NewRemotePlayer(20, 3, remoteIp, p3port)
	var p3Handle ggpo.PlayerHandle
	player4 := ggpo.NewRemotePlayer(20, 4, remoteIp, p4port)
	var p4Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)
	p2p.AddPlayer(&player3, &p3Handle)
	p2p.AddPlayer(&player4, &p4Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	player3 = ggpo.NewRemotePlayer(20, 3, remoteIp, p3port)
	player4 = ggpo.NewRemotePlayer(20, 4, remoteIp, p4port)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	var p2handle3 ggpo.PlayerHandle
	var p2handle4 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)
	p2p2.AddPlayer(&player3, &p2handle3)
	p2p2.AddPlayer(&player4, &p2handle4)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	player3 = ggpo.NewLocalPlayer(20, 3)
	player4 = ggpo.NewRemotePlayer(20, 4, remoteIp, p4port)
	var p3handle1 ggpo.PlayerHandle
	var p3handle2 ggpo.PlayerHandle
	var p3handle3 ggpo.PlayerHandle
	var p3handle4 ggpo.PlayerHandle
	p2p3.AddPlayer(&player1, &p3handle1)
	p2p3.AddPlayer(&player2, &p3handle2)
	p2p3.AddPlayer(&player3, &p3handle3)
	p2p3.AddPlayer(&player4, &p3handle4)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	player3 = ggpo.NewRemotePlayer(20, 3, remoteIp, p3port)
	player4 = ggpo.NewLocalPlayer(20, 4)
	var p4handle1 ggpo.PlayerHandle
	var p4handle2 ggpo.PlayerHandle
	var p4handle3 ggpo.PlayerHandle
	var p4handle4 ggpo.PlayerHandle
	p2p4.AddPlayer(&player1, &p4handle1)
	p2p4.AddPlayer(&player2, &p4handle2)
	p2p4.AddPlayer(&player3, &p4handle3)
	p2p4.AddPlayer(&player4, &p4handle4)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
		p2p3.Idle(0, advance)
		p2p4.Idle(0, advance)
	}
	var err error
	var ignore int
	_, err = p2p.SyncInput(&ignore)
	if err != nil {
		t.Errorf("Player 1 didn't synchronize.")
	}
	_, err = p2p2.SyncInput(&ignore)
	if err != nil {
		t.Errorf("Player 2 didn't synchronize.")
	}

	_, err = p2p3.SyncInput(&ignore)
	if err != nil {
		t.Errorf("Player 3 didn't synchronize.")
	}
	_, err = p2p4.SyncInput(&ignore)
	if err != nil {
		t.Errorf("Player 4 didn't synchronize.")
	}

}

func TestP2PBackend4PlayerShareInput(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 4
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)

	session3 := mocks.NewFakeSession()

	p3port := 6005
	p2p3 := ggpo.NewPeer(&session3, p3port, numPlayers, inputSize)

	session4 := mocks.NewFakeSession()
	p4port := 6006
	p2p4 := ggpo.NewPeer(&session4, p4port, numPlayers, inputSize)

	connection := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p2, &p2p3, &p2p4}, localPort, remoteIp)
	connection2 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p, &p2p3, &p2p4}, remotePort, remoteIp)
	connection3 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p, &p2p2, &p2p4}, p3port, remoteIp)
	connection4 := mocks.NewFakeMultiplePeerConnection([]transport.MessageHandler{&p2p, &p2p2, &p2p3}, p4port, remoteIp)
	//ggpo.EnableLogger()
	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)
	p2p3.InitializeConnection(&connection3)
	p2p4.InitializeConnection(&connection4)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	player3 := ggpo.NewRemotePlayer(20, 3, remoteIp, p3port)
	var p3Handle ggpo.PlayerHandle
	player4 := ggpo.NewRemotePlayer(20, 4, remoteIp, p4port)
	var p4Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)
	p2p.AddPlayer(&player3, &p3Handle)
	p2p.AddPlayer(&player4, &p4Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	player3 = ggpo.NewRemotePlayer(20, 3, remoteIp, p3port)
	player4 = ggpo.NewRemotePlayer(20, 4, remoteIp, p4port)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	var p2handle3 ggpo.PlayerHandle
	var p2handle4 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)
	p2p2.AddPlayer(&player3, &p2handle3)
	p2p2.AddPlayer(&player4, &p2handle4)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	player3 = ggpo.NewLocalPlayer(20, 3)
	player4 = ggpo.NewRemotePlayer(20, 4, remoteIp, p4port)
	var p3handle1 ggpo.PlayerHandle
	var p3handle2 ggpo.PlayerHandle
	var p3handle3 ggpo.PlayerHandle
	var p3handle4 ggpo.PlayerHandle
	p2p3.AddPlayer(&player1, &p3handle1)
	p2p3.AddPlayer(&player2, &p3handle2)
	p2p3.AddPlayer(&player3, &p3handle3)
	p2p3.AddPlayer(&player4, &p3handle4)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	player3 = ggpo.NewRemotePlayer(20, 3, remoteIp, p3port)
	player4 = ggpo.NewLocalPlayer(20, 4)
	var p4handle1 ggpo.PlayerHandle
	var p4handle2 ggpo.PlayerHandle
	var p4handle3 ggpo.PlayerHandle
	var p4handle4 ggpo.PlayerHandle
	p2p4.AddPlayer(&player1, &p4handle1)
	p2p4.AddPlayer(&player2, &p4handle2)
	p2p4.AddPlayer(&player3, &p4handle3)
	p2p4.AddPlayer(&player4, &p4handle4)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
		p2p3.Idle(0, advance)
		p2p4.Idle(0, advance)
	}
	input1 := []byte{1, 2, 3, 4}
	input2 := []byte{5, 6, 7, 8}
	input3 := []byte{9, 10, 11, 12}
	input4 := []byte{13, 14, 15, 16}
	var err error
	var ignore int
	var p1now, p2now, p1next, p2next, p3next, p3now, p4next, p4now int
	p1now = int(time.Now().UnixMilli())
	p1next = p1now
	p2next = p1now
	p2now = p1now
	p3next = p1now
	p3now = p1now
	p4next = p1now
	p4now = p1now

	for i := 0; i < 4; i++ {
		doPollTimeOuts := int(math.Max(0, float64(p1next-p1now-1)))

		p2p.Idle(doPollTimeOuts, advance)
		if p1next >= p1now {
			err = p2p.AddLocalInput(p1Handle, input1, 4)
			if err == nil {
				//_, err = p2p.SyncInput(&ignore)
				if err == nil {
					p2p.AdvanceFrame()
				}
			}
			p1next = p1now + 1000/60
		}

		doPollTimeOuts = int(math.Max(0, float64(p2next-p2now-1)))
		p2p2.Idle(doPollTimeOuts, advance)
		if p2next >= p2now {
			err = p2p2.AddLocalInput(p2handle2, input2, 4)
			if err == nil {
				//_, err = p2p2.SyncInput(&ignore)
				if err == nil {
					p2p2.AdvanceFrame()
				}
			}
			p2next = p2now + 1000/60
		}

		doPollTimeOuts = int(math.Max(0, float64(p3next-p3now-1)))
		p2p3.Idle(doPollTimeOuts, advance)
		if p3next >= p3now {
			err = p2p3.AddLocalInput(p3handle3, input3, 4)
			if err == nil {
				//_, err = p2p3.SyncInput(&ignore)
				if err == nil {
					p2p3.AdvanceFrame()
				}
			}
			p3next = p3now + 1000/60
		}

		doPollTimeOuts = int(math.Max(0, float64(p4next-p4now-1)))
		p2p4.Idle(doPollTimeOuts, advance)
		if p4next >= p4now {
			err = p2p4.AddLocalInput(p4handle4, input4, 4)
			if err == nil {
				//_, err = p2p4.SyncInput(&ignore)
				if err == nil {
					p2p4.AdvanceFrame()
				}
			}
			p4next = p4now + 1000/60
		}
	}

	val1, _ := p2p.SyncInput(&ignore)
	val2, _ := p2p2.SyncInput(&ignore)
	val3, _ := p2p3.SyncInput(&ignore)
	val4, _ := p2p4.SyncInput(&ignore)
	if !slice2dEqual(val1, val2) || !slice2dEqual(val2, val3) || !slice2dEqual(val1, val3) ||
		!slice2dEqual(val2, val4) || !slice2dEqual(val1, val4) || !slice2dEqual(val4, val3) {
		t.Errorf("All peers did not recieve the inputs. P1 inputs %v P2 inputs %v P3 inputs %v P4 inputs %v", val1, val2, val3, val4)
	}
}

func TestP2PBackendGetNetworkStats(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)
	connection := mocks.NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := mocks.NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggpo.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggpo.NewLocalPlayer(20, 2)
	var p2handle1 ggpo.PlayerHandle
	var p2handle2 ggpo.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < protocol.NumSyncPackets; i++ {
		p2p.Idle(0, advance)
		p2p2.Idle(0, advance)
	}
	p1stats := make([]protocol.NetworkStats, numPlayers)
	for i := 0; i < numPlayers; i++ {
		p1stats[i], _ = p2p.GetNetworkStats(ggpo.PlayerHandle(i + 1))
	}
	p2stats := make([]protocol.NetworkStats, numPlayers)
	for i := 0; i < numPlayers; i++ {
		p2stats[i], _ = p2p2.GetNetworkStats(ggpo.PlayerHandle(i + 1))
	}
	if p2stats[0].Timesync.LocalFramesBehind != p1stats[1].Timesync.LocalFramesBehind {
		t.Errorf("Remote local frames behind for both endpoints should be -1 at this state.")
	}
}
func TestP2PBackendGetNetworkStatsInvalid(t *testing.T) {
	session := mocks.NewFakeSession()
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggpo.NewPeer(&session, localPort, numPlayers, inputSize)

	session2 := mocks.NewFakeSession()
	p2p2 := ggpo.NewPeer(&session2, remotePort, numPlayers, inputSize)
	connection := mocks.NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := mocks.NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitializeConnection(&connection)
	p2p2.InitializeConnection(&connection2)

	player1 := ggpo.NewLocalPlayer(20, 1)
	var p1Handle ggpo.PlayerHandle
	player2 := ggpo.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggpo.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	_, err := p2p.GetNetworkStats(ggpo.PlayerHandle(29))
	if err == nil {
		t.Errorf("Trying to create stats for an invalid player handle should return an error.")
	}
}
