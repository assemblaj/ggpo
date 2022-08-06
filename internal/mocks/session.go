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

func (s *FakeSession) saveGameState(stateID int) ([]byte, bool) {
	s.saveStates[stateID] = s.game.clone()
	return []byte{}, true
}

func (s *FakeSession) loadGameState(stateID int) bool {
	s.game = *s.saveStates[stateID]
	return true
}

func (s *FakeSession) logGameState(fileName string, buffer []byte, len int) bool {
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

func (s *FakeSession) onEvent(info *ggpo.Event) bool {
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
func (s *FakeSession) advanceFrame(flags int) bool {
	/*
		fmt.Println("Advancing frame from callback. ")
		var discconectFlags int

		// Make sure we fetch the inputs from GGPO and use these to update
		// the game state instead of reading from the keyboard.
		//inputs, result := ggpo.SynchronizeInput(session, &discconectFlags)
		inputs, result := session.SyncInput(&discconectFlags)
		if result == nil {
			//log.Fatal("Error from GGTHXSynchronizeInput")
			input := decodeInputs(inputs)
			game.AdvanceFrame(input, discconectFlags)
		}*/

	return true
}

type FakeSessionWithBackend struct {
	backend    ggpo.Session
	game       FakeGame
	saveStates map[int]*FakeGame
}

func NewFakeSessionWithBackend() FakeSessionWithBackend {
	return FakeSessionWithBackend{
		game:       FakeGame{},
		saveStates: make(map[int]*FakeGame),
	}
}

func (f *FakeSessionWithBackend) SetBackend(backend ggpo.Session) {
	f.backend = backend
}

func (f *FakeSessionWithBackend) beginGame(game string) bool {
	log.Println("Starting Game!")
	return true
}

func (f *FakeSessionWithBackend) saveGameState(stateID int) ([]byte, bool) {
	f.saveStates[stateID] = f.game.clone()
	return []byte{}, true
}

func (f *FakeSessionWithBackend) loadGameState(stateID int) bool {
	f.game = *f.saveStates[stateID]
	return true
}

func (f *FakeSessionWithBackend) logGameState(fileName string, buffer []byte, len int) bool {
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

func (f *FakeSessionWithBackend) onEvent(info *ggpo.Event) bool {
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
func (f *FakeSessionWithBackend) advanceFrame(flags int) bool {
	var discconectFlags int
	// Make sure we fetch the inputs from GGPO and use these to update
	// the game state instead of reading from the keyboard.
	//inputs, result := ggpo.SynchronizeInput(session, &discconectFlags)
	_, result := f.backend.SyncInput(&discconectFlags)
	if result == nil {
		f.backend.IncrementFrame()
	}

	return true
}

/*
func init() {
	var callbacks ggpo.SessionCallbacks

	session := NewFakeSession()
	callbacks.AdvanceFrame = session.advanceFrame
	callbacks.BeginGame = session.beginGame
	callbacks.FreeBuffer = session.freeBuffer
	callbacks.LoadGameState = session.loadGameState
	callbacks.LogGameState = session.logGameState
	callbacks.OnEvent = session.onEvent
	callbacks.SaveGameState = session.saveGameState
}
*/
func MakeSessionCallBacks(session FakeSession) ggpo.SessionCallbacks {
	var sessionCallbacks ggpo.SessionCallbacks
	sessionCallbacks.AdvanceFrame = session.advanceFrame
	sessionCallbacks.BeginGame = session.beginGame
	sessionCallbacks.FreeBuffer = session.freeBuffer
	sessionCallbacks.LoadGameState = session.loadGameState
	sessionCallbacks.LogGameState = session.logGameState
	sessionCallbacks.OnEvent = session.onEvent
	sessionCallbacks.SaveGameState = session.saveGameState
	return sessionCallbacks
}

func MakeSessionCallBacksBackend(session FakeSessionWithBackend) ggpo.SessionCallbacks {
	var sessionCallbacks ggpo.SessionCallbacks
	sessionCallbacks.AdvanceFrame = session.advanceFrame
	sessionCallbacks.BeginGame = session.beginGame
	sessionCallbacks.FreeBuffer = session.freeBuffer
	sessionCallbacks.LoadGameState = session.loadGameState
	sessionCallbacks.LogGameState = session.logGameState
	sessionCallbacks.OnEvent = session.onEvent
	sessionCallbacks.SaveGameState = session.saveGameState
	return sessionCallbacks
}
