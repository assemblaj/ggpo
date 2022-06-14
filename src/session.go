package ggthx

/*
	Remember to
*/
type GGTHXSession interface {
	DoPoll(timeout int) GGTHXErrorCode
	AddPlayer(player *GGTHXPlayer, handle *GGTHXPlayerHandle) GGTHXErrorCode
	AddLocalInput(player *GGTHXPlayerHandle, values []byte, size int) GGTHXErrorCode
	SyncInput(values []byte, size int, disconnect_flags *int) GGTHXErrorCode
	IncrementFrame() GGTHXErrorCode
	Chat(text string) GGTHXErrorCode
	DisconnectPlayer(handle *GGTHXPlayerHandle) GGTHXErrorCode
	GetNetworkStats(stats *GGTHXNetworkStats, handle GGTHXPlayerHandle) GGTHXErrorCode
	Logv(format string, args ...int) GGTHXErrorCode
	SetFrameDelay(player GGTHXPlayerHandle, delay int) GGTHXErrorCode
	SetDisconnectTimeout(timeout int) GGTHXErrorCode
	SetDisconnectNotifyStart(timeout int) GGTHXErrorCode
	Close() GGTHXErrorCode
}
