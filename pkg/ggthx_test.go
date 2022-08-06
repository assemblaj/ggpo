package ggthx_test

import (
	"bytes"
	"encoding/gob"
	"log"

	ggthx "github.com/assemblaj/ggthx/pkg"
)

type FakeGame struct {
	hash string
}

func (f *FakeGame) clone() (result *FakeGame) {
	result = &FakeGame{}
	*result = *f
	return result
}

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

func (s *FakeSession) onEvent(info *ggthx.Event) bool {
	switch info.Code {
	case ggthx.EventCodeConnectedToPeer:
		log.Println("EventCodeConnectedToPeer")
	case ggthx.EventCodeSynchronizingWithPeer:
		log.Println("EventCodeSynchronizingWithPeer")
	case ggthx.EventCodeSynchronizedWithPeer:
		log.Println("EventCodeSynchronizedWithPeer")
	case ggthx.EventCodeRunning:
		log.Println("EventCodeRunning")
	case ggthx.EventCodeDisconnectedFromPeer:
		log.Println("EventCodeDisconnectedFromPeer")
	case ggthx.EventCodeTimeSync:
		log.Println("EventCodeTimeSync")
	case ggthx.EventCodeConnectionInterrupted:
		log.Println("EventCodeconnectionInterrupted")
	case ggthx.EventCodeConnectionResumed:
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
		//inputs, result := ggthx.SynchronizeInput(session, &discconectFlags)
		inputs, result := session.SyncInput(&discconectFlags)
		if result == nil {
			//log.Fatal("Error from GGTHXSynchronizeInput")
			input := decodeInputs(inputs)
			game.AdvanceFrame(input, discconectFlags)
		}*/

	return true
}

type FakeSessionWithBackend struct {
	backend    ggthx.Session
	game       FakeGame
	saveStates map[int]*FakeGame
}

func NewFakeSessionWithBackend() FakeSessionWithBackend {
	return FakeSessionWithBackend{
		game:       FakeGame{},
		saveStates: make(map[int]*FakeGame),
	}
}

func (f *FakeSessionWithBackend) SetBackend(backend ggthx.Session) {
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

func (f *FakeSessionWithBackend) onEvent(info *ggthx.Event) bool {
	switch info.Code {
	case ggthx.EventCodeConnectedToPeer:
		log.Println("EventCodeConnectedToPeer")
	case ggthx.EventCodeSynchronizingWithPeer:
		log.Println("EventCodeSynchronizingWithPeer")
	case ggthx.EventCodeSynchronizedWithPeer:
		log.Println("EventCodeSynchronizedWithPeer")
	case ggthx.EventCodeRunning:
		log.Println("EventCodeRunning")
	case ggthx.EventCodeDisconnectedFromPeer:
		log.Println("EventCodeDisconnectedFromPeer")
	case ggthx.EventCodeTimeSync:
		log.Println("EventCodeTimeSync")
	case ggthx.EventCodeConnectionInterrupted:
		log.Println("EventCodeconnectionInterrupted")
	case ggthx.EventCodeConnectionResumed:
		log.Println("EventCodeconnectionInterrupted")
	}
	return true
}
func (f *FakeSessionWithBackend) advanceFrame(flags int) bool {
	var discconectFlags int
	// Make sure we fetch the inputs from GGPO and use these to update
	// the game state instead of reading from the keyboard.
	//inputs, result := ggthx.SynchronizeInput(session, &discconectFlags)
	_, result := f.backend.SyncInput(&discconectFlags)
	if result == nil {
		f.backend.IncrementFrame()
	}

	return true
}

/*
func init() {
	var callbacks ggthx.SessionCallbacks

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
