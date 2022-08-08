package ggpo

type Session interface {
	SaveGameState(stateID int) int
	LoadGameState(stateID int)
	AdvanceFrame(flags int)
	OnEvent(info *Event)
	//SetBackend(backend Backend)
}

const DefaultChecksum = 0

type SessionCallbacks struct {
	BeginGame     beginGame
	SaveGameState saveGameState
	LoadGameState loadGameState
	LogGameState  logGameState
	FreeBuffer    freeBuffer
	AdvanceFrame  advanceFrame
	OnEvent       onEvent
}

type beginGame func(game string) bool

/*
type saveGameState func(len *int, checksum *int, frame int) ([]byte, bool)
type loadGameState func(buffer []byte, len int) bool
*/
type saveGameState func(stateID int) ([]byte, bool)
type loadGameState func(stateID int) bool

type logGameState func(fileName string, buffer []byte, len int) bool
type freeBuffer func(buffer []byte)
type advanceFrame func(flags int) bool
type onEvent func(info *Event) bool
