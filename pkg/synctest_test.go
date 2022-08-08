package ggpo_test

import (
	"bytes"
	"testing"

	"github.com/assemblaj/GGPO-Go/internal/mocks"

	ggpo "github.com/assemblaj/GGPO-Go/pkg"
)

func TestNewSyncTestBackend(t *testing.T) {
	session := mocks.NewFakeSession()
	player := ggpo.NewLocalPlayer(20, 1)
	stb := ggpo.NewSyncTestBackend(&session, "test", 1, 8, 4)
	var handle ggpo.PlayerHandle
	stb.AddPlayer(&player, &handle)

}

func TestSyncTestBackendAddPlayerOver(t *testing.T) {
	session := mocks.NewFakeSession()
	player := ggpo.NewLocalPlayer(20, 2)
	stb := ggpo.NewSyncTestBackend(&session, "test", 1, 8, 4)
	var handle ggpo.PlayerHandle
	err := stb.AddPlayer(&player, &handle)
	if err == nil {
		t.Errorf("There should be an error for adding a player greater than the total num Players.")
	}
}

func TestSyncTestBackendAddPlayerNegative(t *testing.T) {
	session := mocks.NewFakeSession()
	player := ggpo.NewLocalPlayer(20, -1)
	stb := ggpo.NewSyncTestBackend(&session, "test", 1, 8, 4)
	var handle ggpo.PlayerHandle
	err := stb.AddPlayer(&player, &handle)
	if err == nil {
		t.Errorf("There should be an error for adding a player with a negative number")
	}
}

func TestSyncTestBackendAddLocalInputError(t *testing.T) {
	session := mocks.NewFakeSession()
	player := ggpo.NewLocalPlayer(20, 1)
	stb := ggpo.NewSyncTestBackend(&session, "test", 1, 8, 4)
	var handle ggpo.PlayerHandle
	stb.AddPlayer(&player, &handle)
	err := stb.AddLocalInput(handle, []byte{1, 2, 3, 4}, 4)
	if err == nil {
		t.Errorf("There should be an error for adding local input when sync test isn't running yet.")
	}
}

func TestSyncTestBackendAddLocalInput(t *testing.T) {
	session := mocks.NewFakeSession()
	player := ggpo.NewLocalPlayer(20, 1)
	stb := ggpo.NewSyncTestBackend(&session, "test", 1, 8, 4)
	var handle ggpo.PlayerHandle
	stb.AddPlayer(&player, &handle)
	stb.Idle(0)
	err := stb.AddLocalInput(handle, []byte{1, 2, 3, 4}, 4)
	if err != nil {
		t.Errorf("There shouldn't be an error, adding local input should be successful.")
	}
}

func TestSyncTestBackendSyncInput(t *testing.T) {
	session := mocks.NewFakeSession()
	player := ggpo.NewLocalPlayer(20, 1)
	stb := ggpo.NewSyncTestBackend(&session, "test", 1, 8, 4)
	var handle ggpo.PlayerHandle
	stb.AddPlayer(&player, &handle)
	stb.Idle(0)
	inputBytes := []byte{1, 2, 3, 4}
	stb.AddLocalInput(handle, []byte{1, 2, 3, 4}, 4)
	var disconnectFlags int
	input, _ := stb.SyncInput(&disconnectFlags)
	got := input[0]
	want := inputBytes
	if !bytes.Equal(got, want) {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestSyncTestBackendIncrementFramePanic(t *testing.T) {
	session := mocks.NewFakeSession()
	player := ggpo.NewLocalPlayer(20, 1)
	checkDistance := 8
	stb := ggpo.NewSyncTestBackend(&session, "test", 1, checkDistance, 4)
	var handle ggpo.PlayerHandle
	stb.AddPlayer(&player, &handle)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic wen trying to load a state that hadn't been saved")
		}
	}()
	for i := 0; i < checkDistance; i++ {
		stb.AdvanceFrame()
	}
}

func TestSyncTestBackendIncrementFrameCharacterization(t *testing.T) {
	session := mocks.NewFakeSession()
	player := ggpo.NewLocalPlayer(20, 1)
	checkDistance := 8
	stb := ggpo.NewSyncTestBackend(&session, "test", 1, checkDistance, 4)
	var handle ggpo.PlayerHandle
	stb.AddPlayer(&player, &handle)
	stb.Idle(0)
	inputBytes := []byte{1, 2, 3, 4}
	var disconnectFlags int
	for i := 0; i < checkDistance-1; i++ {
		stb.AddLocalInput(handle, inputBytes, 4)
		stb.SyncInput(&disconnectFlags)
		stb.AdvanceFrame()
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic due to a SyncError")
		}
	}()
	stb.AdvanceFrame()

}

// Attempt to simulate STB workflow in a test harness. Always panics on check
// distance frame. Do not full understand why even though it works perfectly
// fine in real time.
func TestSyncTestBackendIncrementFrame(t *testing.T) {
	session := mocks.NewFakeSession()
	player := ggpo.NewLocalPlayer(20, 1)
	checkDistance := 8
	stb := ggpo.NewSyncTestBackend(&session, "test", 1, checkDistance, 4)
	var handle ggpo.PlayerHandle
	stb.AddPlayer(&player, &handle)
	inputBytes := []byte{1, 2, 3, 4}
	var disconnectFlags int
	var result error
	for i := 0; i < checkDistance-1; i++ {
		stb.Idle(0)
		result = stb.AddLocalInput(handle, inputBytes, 4)
		if result == nil {
			_, result = stb.SyncInput(&disconnectFlags)
			if result == nil {
				stb.AdvanceFrame()
			}
		}
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic due to a SyncError")
		}
	}()
	stb.AdvanceFrame()
}

/*Again, WIP, I don't know how to test that this is working, but it is. */
func TestSyncTestBackendChecksumCheck(t *testing.T) {
	session := mocks.NewFakeSessionWithBackend()
	var stb ggpo.SyncTestBackend
	session.SetBackend(&stb)
	player := ggpo.NewLocalPlayer(20, 1)
	checkDistance := 8
	stb = ggpo.NewSyncTestBackend(&session, "test", 1, checkDistance, 4)

	var handle ggpo.PlayerHandle
	stb.AddPlayer(&player, &handle)
	inputBytes := []byte{1, 2, 3, 4}
	var disconnectFlags int
	var result error

	for i := 0; i < checkDistance+1; i++ {
		stb.Idle(0)
		result = stb.AddLocalInput(handle, inputBytes, 4)
		if result == nil {
			_, result := stb.SyncInput(&disconnectFlags)
			if result == nil {
				stb.AdvanceFrame()
			}
		}
	}
}

// Unsupported functions
func TestSyncTestBackendDissconnectPlayerError(t *testing.T) {
	session := mocks.NewFakeSession()
	checkDistance := 8
	stb := ggpo.NewSyncTestBackend(&session, "test", 1, checkDistance, 4)
	err := stb.DisconnectPlayer(ggpo.PlayerHandle(1))
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSyncTestBackendGetNetworkStatsError(t *testing.T) {
	session := mocks.NewFakeSession()
	checkDistance := 8
	stb := ggpo.NewSyncTestBackend(&session, "test", 1, checkDistance, 4)
	_, err := stb.GetNetworkStats(ggpo.PlayerHandle(1))
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSyncTestBackendSetFrameDelayError(t *testing.T) {
	session := mocks.NewFakeSession()
	checkDistance := 8
	stb := ggpo.NewSyncTestBackend(&session, "test", 1, checkDistance, 4)
	err := stb.SetFrameDelay(ggpo.PlayerHandle(1), 20)
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSyncTestBackendSetDisconnectTimeoutError(t *testing.T) {
	session := mocks.NewFakeSession()
	checkDistance := 8
	stb := ggpo.NewSyncTestBackend(&session, "test", 1, checkDistance, 4)
	err := stb.SetDisconnectTimeout(20)
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSyncTestBackendSetDisconnectNotifyStartError(t *testing.T) {
	session := mocks.NewFakeSession()
	checkDistance := 8
	stb := ggpo.NewSyncTestBackend(&session, "test", 1, checkDistance, 4)
	err := stb.SetDisconnectNotifyStart(20)
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSyncTestBackendCloseError(t *testing.T) {
	session := mocks.NewFakeSession()
	checkDistance := 8
	stb := ggpo.NewSyncTestBackend(&session, "test", 1, checkDistance, 4)
	err := stb.Close()
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}
