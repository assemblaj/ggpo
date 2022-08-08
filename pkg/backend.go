package ggpo

import (
	"github.com/assemblaj/GGPO-Go/internal/polling"
	"github.com/assemblaj/GGPO-Go/internal/protocol"
	"github.com/assemblaj/GGPO-Go/internal/transport"
)

/*
	Remember to
*/
type Backend interface {
	Idle(timeout int, timeFunc ...polling.FuncTimeType) error
	AddPlayer(player *Player, handle *PlayerHandle) error
	AddLocalInput(player PlayerHandle, values []byte, size int) error
	SyncInput(disconnectFlags *int) ([][]byte, error)
	AdvanceFrame() error
	DisconnectPlayer(handle PlayerHandle) error
	GetNetworkStats(handle PlayerHandle) (protocol.NetworkStats, error)
	SetFrameDelay(player PlayerHandle, delay int) error
	SetDisconnectTimeout(timeout int) error
	SetDisconnectNotifyStart(timeout int) error
	Close() error
	Start()
	InitializeConnection(c ...transport.Connection) error
}
