package ggthx

import (
	"github.com/assemblaj/ggthx/internal/polling"
	"github.com/assemblaj/ggthx/internal/protocol"
	"github.com/assemblaj/ggthx/internal/transport"
)

/*
	Remember to
*/
type Session interface {
	DoPoll(timeout int, timeFunc ...polling.FuncTimeType) error
	AddPlayer(player *Player, handle *PlayerHandle) error
	AddLocalInput(player PlayerHandle, values []byte, size int) error
	SyncInput(disconnect_flags *int) ([][]byte, error)
	IncrementFrame() error
	Chat(text string) error
	DisconnectPlayer(handle PlayerHandle) error
	GetNetworkStats(stats *protocol.NetworkStats, handle PlayerHandle) error
	Logv(format string, args ...int) error
	SetFrameDelay(player PlayerHandle, delay int) error
	SetDisconnectTimeout(timeout int) error
	SetDisconnectNotifyStart(timeout int) error
	Close() error
	Start()
	InitalizeConnection(c ...transport.Connection) error
}
