package mocks

import (
	"bytes"
	"encoding/gob"

	"github.com/assemblaj/GGPO-Go/internal/util"
	ggpo "github.com/assemblaj/GGPO-Go/pkg"
)

type FakeSession struct {
	Game       FakeGame
	SaveStates map[int]*FakeGame
}

func NewFakeSession() FakeSession {
	return FakeSession{
		Game:       NewFakeGame(),
		SaveStates: make(map[int]*FakeGame),
	}
}
func (s *FakeSession) beginGame(game string) bool {
	util.Log.Println("Starting Game!")
	return true
}

func (s *FakeSession) SaveGameState(stateID int) int {
	s.SaveStates[stateID] = s.Game.clone()
	return ggpo.DefaultChecksum
}

func (s *FakeSession) LoadGameState(stateID int) {
	s.Game = *s.SaveStates[stateID]
}

func (s *FakeSession) LogGameState(fileName string, buffer []byte, len int) bool {
	var game2 FakeGame
	var buf bytes.Buffer = *bytes.NewBuffer(buffer)
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(&game2)
	if err != nil {
		util.Log.Fatal("decode error:", err)
	}
	util.Log.Printf("%s Game State: %s\n", fileName, &game2)
	return true
}

func (s *FakeSession) freeBuffer(buffer []byte) {

}

func (s *FakeSession) OnEvent(info *ggpo.Event) {
	switch info.Code {
	case ggpo.EventCodeConnectedToPeer:
		util.Log.Println("EventCodeConnectedToPeer")
	case ggpo.EventCodeSynchronizingWithPeer:
		util.Log.Println("EventCodeSynchronizingWithPeer")
	case ggpo.EventCodeSynchronizedWithPeer:
		util.Log.Println("EventCodeSynchronizedWithPeer")
	case ggpo.EventCodeRunning:
		util.Log.Println("EventCodeRunning")
	case ggpo.EventCodeDisconnectedFromPeer:
		util.Log.Println("EventCodeDisconnectedFromPeer")
	case ggpo.EventCodeTimeSync:
		util.Log.Println("EventCodeTimeSync")
	case ggpo.EventCodeConnectionInterrupted:
		util.Log.Println("EventCodeconnectionInterrupted")
	case ggpo.EventCodeConnectionResumed:
		util.Log.Println("EventCodeconnectionInterrupted")
	}
}
func (s *FakeSession) AdvanceFrame(flags int) {
}

func (s *FakeSession) SetBackend(backend ggpo.Backend) {}

type FakeSessionWithBackend struct {
	backend    ggpo.Backend
	Game       FakeGame
	SaveStates map[int]*FakeGame
}

func NewFakeSessionWithBackend() FakeSessionWithBackend {
	return FakeSessionWithBackend{
		Game:       NewFakeGame(),
		SaveStates: make(map[int]*FakeGame),
	}
}

func (f *FakeSessionWithBackend) SetBackend(backend ggpo.Backend) {
	f.backend = backend
}

func (f *FakeSessionWithBackend) beginGame(game string) bool {
	util.Log.Println("Starting Game!")
	return true
}

func (f *FakeSessionWithBackend) SaveGameState(stateID int) int {
	f.SaveStates[stateID] = f.Game.clone()
	return f.SaveStates[stateID].Checksum()
}

func (f *FakeSessionWithBackend) LoadGameState(stateID int) {
	f.Game = *f.SaveStates[stateID]
}

func (f *FakeSessionWithBackend) LogGameState(fileName string, buffer []byte, len int) {
	var game2 FakeGame
	var buf bytes.Buffer = *bytes.NewBuffer(buffer)
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(&game2)
	if err != nil {
		util.Log.Fatal("decode error:", err)
	}
	util.Log.Printf("%s Game State: %s\n", fileName, &game2)
}

func (f *FakeSessionWithBackend) freeBuffer(buffer []byte) {

}

func (f *FakeSessionWithBackend) OnEvent(info *ggpo.Event) {
	switch info.Code {
	case ggpo.EventCodeConnectedToPeer:
		util.Log.Println("EventCodeConnectedToPeer")
	case ggpo.EventCodeSynchronizingWithPeer:
		util.Log.Println("EventCodeSynchronizingWithPeer")
	case ggpo.EventCodeSynchronizedWithPeer:
		util.Log.Println("EventCodeSynchronizedWithPeer")
	case ggpo.EventCodeRunning:
		util.Log.Println("EventCodeRunning")
	case ggpo.EventCodeDisconnectedFromPeer:
		util.Log.Println("EventCodeDisconnectedFromPeer")
	case ggpo.EventCodeTimeSync:
		util.Log.Println("EventCodeTimeSync")
	case ggpo.EventCodeConnectionInterrupted:
		util.Log.Println("EventCodeconnectionInterrupted")
	case ggpo.EventCodeConnectionResumed:
		util.Log.Println("EventCodeconnectionInterrupted")
	}
}
func (f *FakeSessionWithBackend) AdvanceFrame(flags int) {
	var discconectFlags int
	// Make sure we fetch the inputs from GGPO and use these to update
	// the game state instead of reading from the keyboard.
	vals, result := f.backend.SyncInput(&discconectFlags)
	if result == nil {
		f.Game.UpdateByInputs(vals)
		f.backend.AdvanceFrame(ggpo.DefaultChecksum)
	}
}

func MakeSessionCallBacks(session FakeSession) ggpo.SessionCallbacks {
	var sessionCallbacks ggpo.SessionCallbacks
	//sessionCallbacks.AdvanceFrame = session.AdvanceFrame
	sessionCallbacks.BeginGame = session.beginGame
	sessionCallbacks.FreeBuffer = session.freeBuffer
	//sessionCallbacks.LoadGameState = session.LoadGameState
	sessionCallbacks.LogGameState = session.LogGameState
	//sessionCallbacks.OnEvent = session.OnEvent
	//sessionCallbacks.SaveGameState = session.SaveGameState
	return sessionCallbacks
}

func MakeSessionCallBacksBackend(session FakeSessionWithBackend) ggpo.SessionCallbacks {
	var sessionCallbacks ggpo.SessionCallbacks
	//sessionCallbacks.AdvanceFrame = session.AdvanceFrame
	sessionCallbacks.BeginGame = session.beginGame
	sessionCallbacks.FreeBuffer = session.freeBuffer
	//sessionCallbacks.LoadGameState = session.LoadGameState
	//sessionCallbacks.LogGameState = session.LogGameState
	//sessionCallbacks.OnEvent = session.OnEvent
	//sessionCallbacks.SaveGameState = session.SaveGameState
	return sessionCallbacks
}
