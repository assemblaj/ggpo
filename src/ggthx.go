package ggthx

func StartSession(
	cb *SessionCallbacks,
	game string,
	numPlayers int,
	inputSize int,
	localPort int) (*Peer2PeerBackend, ErrorCode) {
	ggthx := NewPeer2PeerBackend(cb, game, localPort, numPlayers, inputSize)
	return &ggthx, Ok
}

func AddPlayer(ggthx Session,
	player *Player,
	handle *PlayerHandle) ErrorCode {
	if ggthx == nil {
		return ErrorCodeInvalidSession
	}
	return ggthx.AddPlayer(player, handle)
}

func StartSyncTest(
	cb *SessionCallbacks,
	game string,
	numPlayers int,
	inputSize int,
	frames int) (*SyncTestBackend, ErrorCode) {
	ggthx := NewSyncTestBackend(cb, game, numPlayers, frames, inputSize)
	return &ggthx, Ok
}

func SetFrameDelay(ggthx Session,
	player PlayerHandle,
	frameDelay int) ErrorCode {
	if ggthx == nil {
		return ErrorCodeInvalidSession
	}
	return ggthx.SetFrameDelay(player, frameDelay)
}

func Idle(ggthx Session, timeout int) ErrorCode {
	if ggthx == nil {
		return ErrorCodeInvalidSession
	}
	return ggthx.DoPoll(timeout)
}

func AddLocalInput(ggthx Session, player PlayerHandle, values []byte, size int) ErrorCode {
	if ggthx == nil {
		return ErrorCodeInvalidSession
	}
	return ggthx.AddLocalInput(player, values, size)
}

func SynchronizeInput(ggthx Session, disconnectFlags *int) ([][]byte, ErrorCode) {
	if ggthx == nil {
		return nil, ErrorCodeInvalidSession
	}
	return ggthx.SyncInput(disconnectFlags)
}

func DisconnectPlayer(ggthx Session, player PlayerHandle) ErrorCode {
	if ggthx == nil {
		return ErrorCodeInvalidSession
	}
	return ggthx.DisconnectPlayer(player)
}

func AdvanceFrame(ggthx Session) ErrorCode {
	if ggthx == nil {
		return ErrorCodeInvalidSession
	}
	return ggthx.IncrementFrame()
}

func ClientChat(ggthx Session, text string) ErrorCode {
	if ggthx == nil {
		return ErrorCodeInvalidSession
	}
	return ggthx.Chat(text)
}

func GetNetworkStats(ggthx Session, player PlayerHandle, stats *NetworkStats) ErrorCode {
	if ggthx == nil {
		return ErrorCodeInvalidSession
	}
	return ggthx.GetNetworkStats(stats, player)
}

func CloseSession(ggthx Session) ErrorCode {
	if ggthx == nil {
		return ErrorCodeInvalidSession
	}
	ggthx.Close()
	return Ok
}

func SetDisconnectTimeout(ggthx Session, timeout int) ErrorCode {
	if ggthx == nil {
		return ErrorCodeInvalidSession
	}
	return ggthx.SetDisconnectTimeout(timeout)
}

func SetDisconnectNotifyStart(ggthx Session, timeout int) ErrorCode {
	if ggthx == nil {
		return ErrorCodeInvalidSession
	}
	return ggthx.SetDisconnectNotifyStart(timeout)
}

func SartSpectating(cb *SessionCallbacks,
	gameName string, lcaolPort int, numPlayers int, inputSize int, hostIp string, hostPort int) (*SpectatorBackend, ErrorCode) {
	ggthx := NewSpectatorBackend(cb, gameName, lcaolPort, numPlayers, inputSize, hostIp, hostPort)
	return &ggthx, Ok
}
