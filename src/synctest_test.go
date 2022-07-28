package ggthx_test

import (
	"bytes"
	"testing"

	ggthx "github.com/assemblaj/ggthx/src"
)

func TestNewSyncTestBackend(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	player := ggthx.NewLocalPlayer(20, 1)
	stb := ggthx.NewSyncTestBackend(&sessionCallbacks, "test", 1, 8, 4)
	var handle ggthx.PlayerHandle
	stb.AddPlayer(&player, &handle)

}

func TestSyncTestBackendAddPlayerOver(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	player := ggthx.NewLocalPlayer(20, 2)
	stb := ggthx.NewSyncTestBackend(&sessionCallbacks, "test", 1, 8, 4)
	var handle ggthx.PlayerHandle
	err := stb.AddPlayer(&player, &handle)
	if err == nil {
		t.Errorf("There should be an error for adding a player greater than the total num Players.")
	}
}

func TestSyncTestBackendAddPlayerNegative(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	player := ggthx.NewLocalPlayer(20, -1)
	stb := ggthx.NewSyncTestBackend(&sessionCallbacks, "test", 1, 8, 4)
	var handle ggthx.PlayerHandle
	err := stb.AddPlayer(&player, &handle)
	if err == nil {
		t.Errorf("There should be an error for adding a player with a negative number")
	}
}

func TestSyncTestBackendAddLocalInputError(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	player := ggthx.NewLocalPlayer(20, 1)
	stb := ggthx.NewSyncTestBackend(&sessionCallbacks, "test", 1, 8, 4)
	var handle ggthx.PlayerHandle
	stb.AddPlayer(&player, &handle)
	err := stb.AddLocalInput(handle, []byte{1, 2, 3, 4}, 4)
	if err == nil {
		t.Errorf("There should be an error for adding local input when sync test isn't running yet.")
	}
}

func TestSyncTestBackendAddLocalInput(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	player := ggthx.NewLocalPlayer(20, 1)
	stb := ggthx.NewSyncTestBackend(&sessionCallbacks, "test", 1, 8, 4)
	var handle ggthx.PlayerHandle
	stb.AddPlayer(&player, &handle)
	stb.DoPoll(0)
	err := stb.AddLocalInput(handle, []byte{1, 2, 3, 4}, 4)
	if err != nil {
		t.Errorf("There shouldn't be an error, adding local input should be successful.")
	}
}

func TestSyncTestBackendSyncInput(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	player := ggthx.NewLocalPlayer(20, 1)
	stb := ggthx.NewSyncTestBackend(&sessionCallbacks, "test", 1, 8, 4)
	var handle ggthx.PlayerHandle
	stb.AddPlayer(&player, &handle)
	stb.DoPoll(0)
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
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	player := ggthx.NewLocalPlayer(20, 1)
	checkDistance := 8
	stb := ggthx.NewSyncTestBackend(&sessionCallbacks, "test", 1, checkDistance, 4)
	var handle ggthx.PlayerHandle
	stb.AddPlayer(&player, &handle)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic wen trying to load a state that hadn't been saved")
		}
	}()
	for i := 0; i < checkDistance; i++ {
		stb.IncrementFrame()
	}
}

func TestSyncTestBackendIncrementFrameCharacterization(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	player := ggthx.NewLocalPlayer(20, 1)
	checkDistance := 8
	stb := ggthx.NewSyncTestBackend(&sessionCallbacks, "test", 1, checkDistance, 4)
	var handle ggthx.PlayerHandle
	stb.AddPlayer(&player, &handle)
	stb.DoPoll(0)
	inputBytes := []byte{1, 2, 3, 4}
	var disconnectFlags int
	for i := 0; i < checkDistance-1; i++ {
		stb.AddLocalInput(handle, inputBytes, 4)
		stb.SyncInput(&disconnectFlags)
		stb.IncrementFrame()
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic due to a SyncError")
		}
	}()
	stb.IncrementFrame()

}

// Attempt to simulate STB workflow in a test harness. Always panics on check
// distance frame. Do not full understand why even though it works perfectly
// fine in real time.
func TestSyncTestBackendIncrementFrame(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	player := ggthx.NewLocalPlayer(20, 1)
	checkDistance := 8
	stb := ggthx.NewSyncTestBackend(&sessionCallbacks, "test", 1, checkDistance, 4)
	var handle ggthx.PlayerHandle
	stb.AddPlayer(&player, &handle)
	inputBytes := []byte{1, 2, 3, 4}
	var disconnectFlags int
	var result error
	for i := 0; i < checkDistance-1; i++ {
		stb.DoPoll(0)
		result = stb.AddLocalInput(handle, inputBytes, 4)
		if result == nil {
			_, result = stb.SyncInput(&disconnectFlags)
			if result == nil {
				stb.IncrementFrame()
			}
		}
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic due to a SyncError")
		}
	}()
	stb.IncrementFrame()
}

// Unsupported functions
func TestSyncTestBackendChatError(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	checkDistance := 8
	stb := ggthx.NewSyncTestBackend(&sessionCallbacks, "test", 1, checkDistance, 4)
	err := stb.Chat("test")
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSyncTestBackendDissconnectPlayerError(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	checkDistance := 8
	stb := ggthx.NewSyncTestBackend(&sessionCallbacks, "test", 1, checkDistance, 4)
	err := stb.DisconnectPlayer(ggthx.PlayerHandle(1))
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSyncTestBackendGetNetworkStatsError(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	checkDistance := 8
	stb := ggthx.NewSyncTestBackend(&sessionCallbacks, "test", 1, checkDistance, 4)
	var status ggthx.NetworkStats
	err := stb.GetNetworkStats(&status, ggthx.PlayerHandle(1))
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSyncTestBackendLogvError(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	checkDistance := 8
	stb := ggthx.NewSyncTestBackend(&sessionCallbacks, "test", 1, checkDistance, 4)
	err := stb.Logv("test")
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSyncTestBackendSetFrameDelayError(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	checkDistance := 8
	stb := ggthx.NewSyncTestBackend(&sessionCallbacks, "test", 1, checkDistance, 4)
	err := stb.SetFrameDelay(ggthx.PlayerHandle(1), 20)
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSyncTestBackendSetDisconnectTimeoutError(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	checkDistance := 8
	stb := ggthx.NewSyncTestBackend(&sessionCallbacks, "test", 1, checkDistance, 4)
	err := stb.SetDisconnectTimeout(20)
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}

func TestSyncTestBackendCloseError(t *testing.T) {
	session := NewFakeSession()
	sessionCallbacks := makeSessionCallBacks(session)
	checkDistance := 8
	stb := ggthx.NewSyncTestBackend(&sessionCallbacks, "test", 1, checkDistance, 4)
	err := stb.Close()
	if err == nil {
		t.Errorf("The code did not error when using an unsupported Feature.")
	}
}
