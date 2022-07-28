package ggthx_test

import (
	"bytes"
	"fmt"
	"math"
	"testing"
	"time"

	ggthx "github.com/assemblaj/ggthx/src"
)

type FakeP2PConnection struct {
	remoteHandler   ggthx.MessageHandler
	localPort       int
	remotePort      int
	localIP         string
	printOutput     bool
	lastSentMessage ggthx.UDPMessage
	messageHistory  []ggthx.UDPMessage
}

func (f *FakeP2PConnection) SendTo(msg ggthx.UDPMessage, remoteIp string, remotePort int) {
	if f.printOutput {
		fmt.Printf("f.localIP %s f.localPort %d msg %s size %d\n", f.localIP, f.localPort, msg, msg.PacketSize())
	}
	f.lastSentMessage = msg
	f.messageHistory = append(f.messageHistory, msg)
	f.remoteHandler.HandleMessage(f.localIP, f.localPort, msg, msg.PacketSize())
}

func (f *FakeP2PConnection) Read() {
}

func (f *FakeP2PConnection) Close() {

}

func NewFakeP2PConnection(remoteHandler ggthx.MessageHandler, localPort int, localIP string) FakeP2PConnection {
	f := FakeP2PConnection{}
	f.remoteHandler = remoteHandler
	f.localPort = localPort
	f.localIP = localIP
	f.messageHistory = make([]ggthx.UDPMessage, 10)
	return f
}

func TestP2PBackendAddPlayer(t *testing.T) {
	connection := NewFakeConnection()
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggthx.NewPeer2PeerBackend(&sessionCallbacks, "test", localPort, numPlayers, inputSize)
	p2p.InitalizeConnection(&connection)
	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
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
	connection := NewFakeConnection()
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggthx.NewPeer2PeerBackend(&sessionCallbacks, "test", localPort, numPlayers, inputSize)
	p2p.InitalizeConnection(&connection)
	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
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
	connection := NewFakeConnection()
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggthx.NewPeer2PeerBackend(&sessionCallbacks, "test", localPort, numPlayers, inputSize)
	p2p.InitalizeConnection(&connection)
	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
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
	connection := NewFakeConnection()
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggthx.NewPeer2PeerBackend(&sessionCallbacks, "test", localPort, numPlayers, inputSize)
	p2p.InitalizeConnection(&connection)
	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	err := p2p.AddPlayer(&player1, &p1Handle)
	if err != nil {
		t.Errorf("There was an error when adding player 1.")
	}
	err = p2p.AddPlayer(&player2, &p2Handle)
	if err != nil {
		t.Errorf("There was an error when adding player 2.")
	}
	err = p2p.IncrementFrame()
	if err != nil {
		t.Errorf("There was an error when incrementing the frame.")
	}
}

func TestP2PBackendSynchronizeInputs(t *testing.T) {
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
	connection := NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)
	inputBytes := []byte{1, 2, 3, 4}
	p2p.DoPoll(0)
	p2p2.DoPoll(0)

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
	connection := NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)
	inputBytes := []byte{1, 2, 3, 4}
	p2p.DoPoll(0)

	p2p2.DoPoll(0)
	p2p2.AddLocalInput(p2Handle, inputBytes, len(inputBytes))
	p2p2.DoPoll(0)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic due to an InputQueue error")
		}
	}()

	p2p2.AddLocalInput(p2Handle, inputBytes, len(inputBytes))

}

func TestP2PBackendPoll2PlayersDefault(t *testing.T) {
	connection := NewFakeConnection()
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 2
	inputSize := 4
	p2p := ggthx.NewPeer2PeerBackend(&sessionCallbacks, "test", localPort, numPlayers, inputSize)
	p2p.InitalizeConnection(&connection)
	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
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
	connection := NewFakeConnection()
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	localPort := 6000
	remotePort := 6001
	remoteIp := "127.2.1.1"
	numPlayers := 3
	inputSize := 4
	p2p := ggthx.NewPeer2PeerBackend(&sessionCallbacks, "test", localPort, numPlayers, inputSize)
	p2p.InitalizeConnection(&connection)
	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	player3 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort+1)
	var p2Handle ggthx.PlayerHandle
	var p3Handle ggthx.PlayerHandle

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
	connection := NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)
	cycles := 20
	inputBytes := []byte{1, 2, 3, 4}
	for i := 0; i < cycles; i++ {
		p2p2.DoPoll(0)
		p2p2.IncrementFrame()
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
	connection := NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)
	inputBytes := []byte{1, 2, 3, 4}
	p2p.DoPoll(0)
	p2p2.DoPoll(0)

	err := p2p2.AddLocalInput(p2Handle, inputBytes, len(inputBytes))
	if err != nil {
		t.Errorf("The backends didn't synchronize")
	}
}

func TestP2PBackendFullSession(t *testing.T) {
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
	connection := NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)
	inputBytes := []byte{1, 2, 3, 4}

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}

	for i := 0; i < ggthx.NumSyncPackets; i++ {
		p2p.DoPoll(0, advance)
		p2p2.DoPoll(0, advance)
	}

	for i := 0; i < 20; i++ {
		p2p2.DoPoll(0)
		err := p2p2.AddLocalInput(p2Handle, inputBytes, len(inputBytes))
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
	connection := NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	err := p2p2.DisconnectPlayer(p2handle2)
	if err != nil {
		t.Errorf("Had an error trying to disconnect the local player.")
	}
}

func TestP2PBackendDisconnectPlayerRemoteCharacterization(t *testing.T) {
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
	connection := NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
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
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	localPort := 6000
	numPlayers := 2
	inputSize := 4
	p2p := ggthx.NewPeer2PeerBackend(&sessionCallbacks, "test", localPort, numPlayers, inputSize)
	connection := NewFakeConnection()
	p2p.InitalizeConnection(&connection)

	err := p2p.DisconnectPlayer(ggthx.PlayerHandle(8))
	if err == nil {
		t.Errorf("The code should have created an error when passing an invalid player handle into DisconnectPlayer")
	}
}
func TestP2PBackendMockSynchronize(t *testing.T) {
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
	connection := NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	var err1, err2 error

	for i := 0; i < ggthx.NumSyncPackets; i++ {
		var disconnectFlags int
		p2p.DoPoll(0, advance)
		p2p2.DoPoll(0, advance)
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
	connection := NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		p2p.DoPoll(0, advance)
		p2p2.DoPoll(0, advance)
	}
	input1 := []byte{1, 2, 3, 4}
	input2 := []byte{5, 6, 7, 8}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic.")
		}
	}()

	for i := 0; i < 8; i++ {
		p2p.DoPoll(0, advance)
		p2p2.DoPoll(0, advance)
		p2p.AddLocalInput(p1Handle, input1, 4)
		p2p2.AddLocalInput(p2handle2, input2, 4)
		var disconnectFlags int
		p2p2.SyncInput(&disconnectFlags)
		p2p.SyncInput(&disconnectFlags)
		p2p.IncrementFrame()
		p2p2.IncrementFrame()
	}

}
func TestP2PBackendMoockInputExchange(t *testing.T) {
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
	connection := NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		p2p.DoPoll(0, advance)
		p2p2.DoPoll(0, advance)
	}
	input1 := []byte{1, 2, 3, 4}
	input2 := []byte{5, 6, 7, 8}

	for i := 0; i < 8; i++ {
		p2p.DoPoll(0, advance)
		p2p2.DoPoll(0, advance)
		p2p.AddLocalInput(p1Handle, input1, 4)
		p2p2.AddLocalInput(p2handle2, input2, 4)
		p2p.IncrementFrame()
		p2p2.IncrementFrame()
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
	connection := NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		p2p.DoPoll(0, advance)
		p2p2.DoPoll(0, advance)
	}
	input1 := []byte{1, 2, 3, 4}
	input2 := []byte{5, 6, 7, 8}
	doPollTimeOuts := 90
	for i := 0; i < 8; i++ {
		p2p.DoPoll(doPollTimeOuts, advance)
		p2p2.DoPoll(doPollTimeOuts, advance)
		p2p.AddLocalInput(p1Handle, input1, 4)
		p2p2.AddLocalInput(p2handle2, input2, 4)
		p2p.IncrementFrame()
		p2p2.IncrementFrame()
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
	connection := NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		p2p.DoPoll(0, advance)
		p2p2.DoPoll(0, advance)
	}
	input1 := []byte{1, 2, 3, 4}
	input2 := []byte{5, 6, 7, 8}
	doPollTimeOuts := 90
	for i := 0; i < 8; i++ {
		p2p.DoPoll(doPollTimeOuts, advance)
		p2p2.DoPoll(doPollTimeOuts, advance)
		p2p.AddLocalInput(p1Handle, input1, 4)
		p2p2.AddLocalInput(p2handle2, input2, 4)
		p2p.IncrementFrame()
		p2p2.IncrementFrame()
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
	connection := NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		p2p.DoPoll(0, advance)
		p2p2.DoPoll(0, advance)
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
		p2p.DoPoll(doPollTimeOuts, advance)
		p2p2.DoPoll(doPollTimeOuts, advance)
		p2p.AddLocalInput(p1Handle, input1, 4)
		p2p2.AddLocalInput(p2handle2, input2, 4)
		p2p.IncrementFrame()
		p2p2.IncrementFrame()
	}
	got := p2p.Poll2Players(iterations)
	want := iterations + 1
	if want != got {
		t.Errorf("wanted %d got %d ", want, got)
	}

}

func TestP2PBackendMoockDisconnectTimeoutCharacterization(t *testing.T) {
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
	connection := NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		p2p.DoPoll(0, advance)
		p2p2.DoPoll(0, advance)
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
		p2p.DoPoll(doPollTimeOuts, timeout)
		p2p2.DoPoll(doPollTimeOuts, timeout)
		//p2p.AddLocalInput(p1Handle, input1, 4)
		//p2p2.AddLocalInput(p2handle2, input2, 4)
		//p2p.IncrementFrame()
		//p2p2.IncrementFrame()
	}

}
func TestP2PBackendMoockDisconnectTimeoutCharacterization2(t *testing.T) {
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
	connection := NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		p2p.DoPoll(0, advance)
		p2p2.DoPoll(0, advance)
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

		p2p.DoPoll(doPollTimeOuts, timeout)
		p2p2.DoPoll(doPollTimeOuts, timeout)
		//p2p.AddLocalInput(p1Handle, input1, 4)
		//p2p2.AddLocalInput(p2handle2, input2, 4)
		p2p.IncrementFrame()
		p2p2.IncrementFrame()
	}
}
func TestP2PBackendMoockDisconnectTimeout(t *testing.T) {
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
	connection := NewFakeP2PConnection(&p2p2, localPort, remoteIp)
	connection2 := NewFakeP2PConnection(&p2p, remotePort, remoteIp)

	p2p.InitalizeConnection(&connection)
	p2p2.InitalizeConnection(&connection2)

	player1 := ggthx.NewLocalPlayer(20, 1)
	var p1Handle ggthx.PlayerHandle
	player2 := ggthx.NewRemotePlayer(20, 2, remoteIp, remotePort)
	var p2Handle ggthx.PlayerHandle
	p2p.AddPlayer(&player1, &p1Handle)
	p2p.AddPlayer(&player2, &p2Handle)

	player1 = ggthx.NewRemotePlayer(20, 1, remoteIp, localPort)
	player2 = ggthx.NewLocalPlayer(20, 2)
	var p2handle1 ggthx.PlayerHandle
	var p2handle2 ggthx.PlayerHandle
	p2p2.AddPlayer(&player1, &p2handle1)
	p2p2.AddPlayer(&player2, &p2handle2)

	advance := func() int64 {
		return time.Now().Add(time.Millisecond * 2000).UnixMilli()
	}
	for i := 0; i < ggthx.NumSyncPackets; i++ {
		p2p.DoPoll(0, advance)
		p2p2.DoPoll(0, advance)
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

	for i := 0; i < ggthx.MAX_PREDICTION_FRAMES; i++ {
		doPollTimeOuts = int(math.Max(0, float64(p1next-p1now-1)))
		p2p.DoPoll(doPollTimeOuts, currentTime)
		if p1now >= p1next {
			err := p2p.AddLocalInput(p1Handle, input1, 4)
			if err == nil {
				//_, err = p2p.SyncInput(&ignore)
				if err == nil {
					p2p.IncrementFrame()
				}
			}
			p1next = p1now + 1000/60
		}

		doPollTimeOuts = int(math.Max(0, float64(p2next-p2now-1)))
		p2p2.DoPoll(doPollTimeOuts, currentTime)
		if p2now >= p2next {
			err := p2p2.AddLocalInput(p2handle2, input2, 4)
			if err == nil {
				//_, err = p2p2.SyncInput(&ignore)
				if err == nil {
					p2p2.IncrementFrame()
				}
			}
			p2next = p2now + 1000/60
		}

		if i == ggthx.MAX_PREDICTION_FRAMES-2 {
			currentTime = timeout
		}
	}

	err := p2p.DisconnectPlayer(p2Handle)
	ggError := err.(ggthx.Error)
	if ggError.Code != ggthx.ErrorCodePlayerDisconnected {
		t.Errorf("The player should've been timed out and disconnected already.")
	}

	err = p2p2.DisconnectPlayer(p2handle1)
	ggError = err.(ggthx.Error)
	if ggError.Code != ggthx.ErrorCodePlayerDisconnected {
		t.Errorf("The player should've been timed out and disconnected already.")
	}
}
