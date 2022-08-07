package mocks

import (
	"bytes"
	"encoding/gob"
	"log"

	ggpo "github.com/assemblaj/GGPO-Go/pkg"
)

type FakeSession struct {
	game       FakeGame
	saveStates map[int]*FakeGame
}

func NewFakeSession() FakeSession {
	return FakeSession{
		game:       FakeGame{},
		saveStates: make(map[int]*FakeGame),
	}
}
func (s *FakeSession) beginGame(game string) bool {
	log.Println("Starting Game!")
	return true
}

func (s *FakeSession) SaveGameState(stateID int) ([]byte, bool) {
	s.saveStates[stateID] = s.game.clone()
	return []byte{}, true
}

func (s *FakeSession) LoadGameState(stateID int) bool {
	s.game = *s.saveStates[stateID]
	return true
}

func (s *FakeSession) LogGameState(fileName string, buffer []byte, len int) bool {
	var game2 FakeGame
	var buf bytes.Buffer = *bytes.NewBuffer(buffer)
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(&game2)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	log.Printf("%s Game State: %s\n", fileName, game2)
	return true
}

func (s *FakeSession) freeBuffer(buffer []byte) {

}

func (s *FakeSession) OnEvent(info *ggpo.Event) bool {
	switch info.Code {
	case ggpo.EventCodeConnectedToPeer:
		log.Println("EventCodeConnectedToPeer")
	case ggpo.EventCodeSynchronizingWithPeer:
		log.Println("EventCodeSynchronizingWithPeer")
	case ggpo.EventCodeSynchronizedWithPeer:
		log.Println("EventCodeSynchronizedWithPeer")
	case ggpo.EventCodeRunning:
		log.Println("EventCodeRunning")
	case ggpo.EventCodeDisconnectedFromPeer:
		log.Println("EventCodeDisconnectedFromPeer")
	case ggpo.EventCodeTimeSync:
		log.Println("EventCodeTimeSync")
	case ggpo.EventCodeConnectionInterrupted:
		log.Println("EventCodeconnectionInterrupted")
	case ggpo.EventCodeConnectionResumed:
		log.Println("EventCodeconnectionInterrupted")
	}
	return true
}
func (s *FakeSession) AdvanceFrame(flags int) bool {
	return true
}

func (s *FakeSession) SetBackend(backend ggpo.Backend) {}

type FakeSessionWithBackend struct {
	backend    ggpo.Backend
	game       FakeGame
	saveStates map[int]*FakeGame
}

func NewFakeSessionWithBackend() FakeSessionWithBackend {
	return FakeSessionWithBackend{
		game:       FakeGame{},
		saveStates: make(map[int]*FakeGame),
	}
}

func (f *FakeSessionWithBackend) SetBackend(backend ggpo.Backend) {
	f.backend = backend
}

func (f *FakeSessionWithBackend) beginGame(game string) bool {
	log.Println("Starting Game!")
	return true
}

func (f *FakeSessionWithBackend) SaveGameState(stateID int) ([]byte, bool) {
	f.saveStates[stateID] = f.game.clone()
	return []byte{}, true
}

func (f *FakeSessionWithBackend) LoadGameState(stateID int) bool {
	f.game = *f.saveStates[stateID]
	return true
}

func (f *FakeSessionWithBackend) LogGameState(fileName string, buffer []byte, len int) bool {
	var game2 FakeGame
	var buf bytes.Buffer = *bytes.NewBuffer(buffer)
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(&game2)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	log.Printf("%s Game State: %s\n", fileName, game2)
	return true
}

func (f *FakeSessionWithBackend) freeBuffer(buffer []byte) {

}

func (f *FakeSessionWithBackend) OnEvent(info *ggpo.Event) bool {
	switch info.Code {
	case ggpo.EventCodeConnectedToPeer:
		log.Println("EventCodeConnectedToPeer")
	case ggpo.EventCodeSynchronizingWithPeer:
		log.Println("EventCodeSynchronizingWithPeer")
	case ggpo.EventCodeSynchronizedWithPeer:
		log.Println("EventCodeSynchronizedWithPeer")
	case ggpo.EventCodeRunning:
		log.Println("EventCodeRunning")
	case ggpo.EventCodeDisconnectedFromPeer:
		log.Println("EventCodeDisconnectedFromPeer")
	case ggpo.EventCodeTimeSync:
		log.Println("EventCodeTimeSync")
	case ggpo.EventCodeConnectionInterrupted:
		log.Println("EventCodeconnectionInterrupted")
	case ggpo.EventCodeConnectionResumed:
		log.Println("EventCodeconnectionInterrupted")
	}
	return true
}
func (f *FakeSessionWithBackend) AdvanceFrame(flags int) bool {
	var discconectFlags int
	// Make sure we fetch the inputs from GGPO and use these to update
	// the game state instead of reading from the keyboard.
	_, result := f.backend.SyncInput(&discconectFlags)
	if result == nil {
		f.backend.AdvanceFrame()
	}
	return true
}

func MakeSessionCallBacks(session FakeSession) ggpo.SessionCallbacks {
	var sessionCallbacks ggpo.SessionCallbacks
	sessionCallbacks.AdvanceFrame = session.AdvanceFrame
	sessionCallbacks.BeginGame = session.beginGame
	sessionCallbacks.FreeBuffer = session.freeBuffer
	sessionCallbacks.LoadGameState = session.LoadGameState
	sessionCallbacks.LogGameState = session.LogGameState
	sessionCallbacks.OnEvent = session.OnEvent
	sessionCallbacks.SaveGameState = session.SaveGameState
	return sessionCallbacks
}

func MakeSessionCallBacksBackend(session FakeSessionWithBackend) ggpo.SessionCallbacks {
	var sessionCallbacks ggpo.SessionCallbacks
	sessionCallbacks.AdvanceFrame = session.AdvanceFrame
	sessionCallbacks.BeginGame = session.beginGame
	sessionCallbacks.FreeBuffer = session.freeBuffer
	sessionCallbacks.LoadGameState = session.LoadGameState
	sessionCallbacks.LogGameState = session.LogGameState
	sessionCallbacks.OnEvent = session.OnEvent
	sessionCallbacks.SaveGameState = session.SaveGameState
	return sessionCallbacks
}
