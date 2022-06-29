package ggthx

/*
	Remember to
*/
type Session interface {
	DoPoll(timeout int) ErrorCode
	AddPlayer(player *Player, handle *PlayerHandle) ErrorCode
	AddLocalInput(player PlayerHandle, values []byte, size int) ErrorCode
	SyncInput(disconnect_flags *int) ([][]byte, ErrorCode)
	IncrementFrame() ErrorCode
	Chat(text string) ErrorCode
	DisconnectPlayer(handle PlayerHandle) ErrorCode
	GetNetworkStats(stats *NetworkStats, handle PlayerHandle) ErrorCode
	Logv(format string, args ...int) ErrorCode
	SetFrameDelay(player PlayerHandle, delay int) ErrorCode
	SetDisconnectTimeout(timeout int) ErrorCode
	SetDisconnectNotifyStart(timeout int) ErrorCode
	Close() ErrorCode
}
