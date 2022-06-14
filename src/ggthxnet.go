package ggthx

const GGTHX_MAX_PLAYERS int = 4
const GGTHX_MAX_PREDICTION_FRAMES int = 8
const GGTHX_MAX_SPECTATORS int = 32
const GGTHX_SPECTATOR_INPUT_INTERVAL int = 4

type GGTHXPlayerHandle int

type GGTHXPlayerType int

const (
	GGTHX_PLAYERTYPE_LOCAL GGTHXPlayerType = iota
	GGTHX_PLAYERTYPE_REMOTE
	GGTHX_PLAYERTYPE_SPECTATOR
)

type GGTHXPlayer struct {
	size       int
	playerType GGTHXPlayerType
	playerNum  int
	remote     GGTHXRemotePlayer
}
type GGTHXRemotePlayer struct {
	ipAdress []byte
	port     uint16
}

type GGTHXLocalEndpoint struct {
	playerNum int
}

type GGTHXErrorCode int

const (
	GGTHX_OK                              GGTHXErrorCode = 0
	GGTHX_ERRORCODE_SUCCESS               GGTHXErrorCode = 0
	GGTHX_ERRORCODE_GENERAL_FAILURE       GGTHXErrorCode = -1
	GGTHX_ERRORCODE_INVALID_SESSION       GGTHXErrorCode = 1
	GGTHX_ERRORCODE_INVALID_PLAYER_HANDLE GGTHXErrorCode = 2
	GGTHX_ERRORCODE_PLAYER_OUT_OF_RANGE   GGTHXErrorCode = 3
	GGTHX_ERRORCODE_PREDICTION_THRESHOLD  GGTHXErrorCode = 4
	GGTHX_ERRORCODE_UNSUPPORTED           GGTHXErrorCode = 5
	GGTHX_ERRORCODE_NOT_SYNCHRONIZED      GGTHXErrorCode = 6
	GGTHX_ERRORCODE_IN_ROLLBACK           GGTHXErrorCode = 7
	GGTHX_ERRORCODE_INPUT_DROPPED         GGTHXErrorCode = 8
	GGTHX_ERRORCODE_PLAYER_DISCONNECTED   GGTHXErrorCode = 9
	GGTHX_ERRORCODE_TOO_MANY_SPECTATORS   GGTHXErrorCode = 10
	GGTHX_ERRORCODE_INVALID_REQUEST       GGTHXErrorCode = 11
)

func GGTHX_SUCESS(result GGTHXErrorCode) bool {
	return result == GGTHX_ERRORCODE_SUCCESS
}

const GGTHX_INVALID_HANDLE int = -1

type GGTHXEventCode int

const (
	GGTHX_EVENTCODE_CONNECTED_TO_PEER       GGTHXEventCode = 1000
	GGTHX_EVENTCODE_SYNCHRONIZING_WITH_PEER GGTHXEventCode = 1001
	GGTHX_EVENTCODE_SYNCHRONIZED_WITH_PEER  GGTHXEventCode = 1002
	GGTHX_EVENTCODE_RUNNING                 GGTHXEventCode = 1003
	GGTHX_EVENTCODE_DISCONNECTED_FROM_PEER  GGTHXEventCode = 1004
	GGTHX_EVENTCODE_TIMESYNC                GGTHXEventCode = 1005
	GGTHX_EVENTCODE_CONNECTION_INTERRUPTED  GGTHXEventCode = 1006
	GGTHX_EVENTCODE_CONNECTION_RESUMED      GGTHXEventCode = 1007
)

// the original had a union a named struct for each event type,
// but mostly each just had the PlayerHandle, so instead i'm
// including the extra specific fields in every event object
// Hopefully I remember I did this like this
type GGTHXEvent struct {
	code              GGTHXEventCode
	player            GGTHXPlayerHandle
	count             int // synchronizing
	total             int // synchronizing
	framesAhead       int // timesync
	disconnectTimeout int // connection interrupted
}

type GGTHXSessionCallbacks struct {
	BeginGame     beginGame
	SaveGameState saveGameState
	LoadGameState loadGameState
	LogGameState  logGameState
	FreeBuffer    freeBuffer
	AdvanceFrame  advanceFrame
	OnEvent       onEvent
}

type beginGame func(game string) bool
type saveGameState func(buffer []byte, len *int, checksum *int, frame int) bool
type loadGameState func(buffer []byte, len int) bool
type logGameState func(fileName string, buffer []byte, len int) bool
type freeBuffer func(buffer []byte)
type advanceFrame func(flags int) bool
type onEvent func(info *GGTHXEvent) bool

type GGTHXNetworkStats struct {
	network  GGTHXNetworkNetworkStats
	timesync GGTHXNetworkTimeSyncStats
}

func (gns GGTHXNetworkStats) BuildGGTHXNetworkNetworkStats(
	sendQueueLen int,
	recvQueueLen int,
	ping int,
	kbpsSend int) {
	gns.network.sendQueueLen = sendQueueLen
	gns.network.recvQueueLen = recvQueueLen
	gns.network.ping = ping
	gns.network.kbpsSend = kbpsSend
}

func (gns GGTHXNetworkStats) BuildGGTHXNetworkTimeSyncStats(
	localFramesBehind int,
	remoteFramesBehind int) {
	gns.timesync.localFramesBehind = localFramesBehind
	gns.timesync.remoteFramesBehind = remoteFramesBehind
}

type GGTHXNetworkNetworkStats struct {
	sendQueueLen int
	recvQueueLen int
	ping         int
	kbpsSend     int
}
type GGTHXNetworkTimeSyncStats struct {
	localFramesBehind  int
	remoteFramesBehind int
}
