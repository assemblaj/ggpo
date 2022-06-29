package ggthx

/*
	Remember to
*/
type Session interface {
	DoPoll(timeout int) error
	AddPlayer(player *Player, handle *PlayerHandle) error
	AddLocalInput(player PlayerHandle, values []byte, size int) error
	SyncInput(disconnect_flags *int) ([][]byte, error)
	IncrementFrame() error
	Chat(text string) error
	DisconnectPlayer(handle PlayerHandle) error
	GetNetworkStats(stats *NetworkStats, handle PlayerHandle) error
	Logv(format string, args ...int) error
	SetFrameDelay(player PlayerHandle, delay int) error
	SetDisconnectTimeout(timeout int) error
	SetDisconnectNotifyStart(timeout int) error
	Close() error
}
