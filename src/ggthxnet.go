package ggthx

import "fmt"

const (
	MaxPlayers             = 4
	MaxPredictionFrames    = 8
	MaxSpectators          = 32
	SpectatorInputInterval = 4
)

type PlayerHandle int

type PlayerType int

const (
	PlayerTypeLocal PlayerType = iota
	PlayerTypeRemote
	PlayerTypeSpectator
)

type Player struct {
	Size       int
	PlayerType PlayerType
	PlayerNum  int
	Remote     RemotePlayer
}

func NewLocalPlayer(size int, playerNum int) Player {
	return Player{
		Size:       size,
		PlayerNum:  playerNum,
		PlayerType: PlayerTypeLocal}
}

func NewRemotePlayer(size int, playerNum int, ipAdress string, port int) Player {
	return Player{
		Size:       size,
		PlayerNum:  playerNum,
		PlayerType: PlayerTypeRemote,
		Remote: RemotePlayer{
			IpAdress: ipAdress,
			Port:     port},
	}
}

type RemotePlayer struct {
	IpAdress string
	Port     int
}

type LocalEndpoint struct {
	playerNum int
}

type Error struct {
	Code ErrorCode
	Name string
}

func (e Error) Error() string {
	return fmt.Sprintf("ggthx: %s:%d", e.Name, e.Code)
}

type ErrorCode int

const (
	Ok                           ErrorCode = 0
	ErrorCodeSuccess             ErrorCode = 0
	ErrorCodeGeneralFailure      ErrorCode = -1
	ErrorCodeInvalidSession      ErrorCode = 1
	ErrorCodeInvalidPlayerHandle ErrorCode = 2
	ErrorCodePlayerOutOfRange    ErrorCode = 3
	ErrorCodePredictionThreshod  ErrorCode = 4
	ErrorCodeUnsupported         ErrorCode = 5
	ErrorCodeNotSynchronized     ErrorCode = 6
	ErrorCodeInRollback          ErrorCode = 7
	ErrorCodeInputDropped        ErrorCode = 8
	ErrorCodePlayerDisconnected  ErrorCode = 9
	ErrorCodeTooManySpectators   ErrorCode = 10
	ErrorCodeInvalidRequest      ErrorCode = 11
)

func Success(result ErrorCode) bool {
	return result == ErrorCodeSuccess
}

const InvalidHandle int = -1

type EventCode int

const (
	EventCodeConnectedToPeer       EventCode = 1000
	EventCodeSynchronizingWithPeer EventCode = 1001
	EventCodeSynchronizedWithPeer  EventCode = 1002
	EventCodeRunning               EventCode = 1003
	EventCodeDisconnectedFromPeer  EventCode = 1004
	EventCodeTimeSync              EventCode = 1005
	EventCodeConnectionInterrupted EventCode = 1006
	EventCodeConnectionResumed     EventCode = 1007
)

// the original had a union a named struct for each event type,
// but mostly each just had the PlayerHandle, so instead i'm
// including the extra specific fields in every event object
// Hopefully I remember I did this like this
type Event struct {
	Code              EventCode
	player            PlayerHandle
	count             int // synchronizing
	total             int // synchronizing
	framesAhead       int // timesync
	disconnectTimeout int // connection interrupted
}

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
type saveGameState func(len *int, checksum *int, frame int) ([]byte, bool)
type loadGameState func(buffer []byte, len int) bool
type logGameState func(fileName string, buffer []byte, len int) bool
type freeBuffer func(buffer []byte)
type advanceFrame func(flags int) bool
type onEvent func(info *Event) bool

type NetworkStats struct {
	network  NetworkNetworkStats
	timesync NetworkTimeSyncStats
}

type NetworkNetworkStats struct {
	sendQueueLen int
	recvQueueLen int
	ping         int
	kbpsSent     int
}
type NetworkTimeSyncStats struct {
	localFramesBehind  int
	remoteFramesBehind int
}
