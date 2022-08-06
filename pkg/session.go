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
	DisconnectPlayer(handle PlayerHandle) error
	GetNetworkStats(stats *protocol.NetworkStats, handle PlayerHandle) error
	SetFrameDelay(player PlayerHandle, delay int) error
	SetDisconnectTimeout(timeout int) error
	SetDisconnectNotifyStart(timeout int) error
	Close() error
	Start()
	InitializeConnection(c ...transport.Connection) error
}
