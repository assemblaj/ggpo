package ggthx

func StartSession(
	cb *SessionCallbacks,
	game string,
	numPlayers int,
	inputSize int,
	localPort int) *Peer2PeerBackend {
	ggthx := NewPeer2PeerBackend(cb, game, localPort, numPlayers, inputSize)
	return &ggthx
}

func AddPlayer(ggthx Session,
	player *Player,
	handle *PlayerHandle) error {
	if ggthx == nil {
		return Error{Code: ErrorCodeInvalidSession, Name: "ErrorCodeInvalidSession"}
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
	frameDelay int) error {
	if ggthx == nil {
		return Error{Code: ErrorCodeInvalidSession, Name: "ErrorCodeInvalidSession"}
	}
	return ggthx.SetFrameDelay(player, frameDelay)
}

func Idle(ggthx Session, timeout int) error {
	if ggthx == nil {
		return Error{Code: ErrorCodeInvalidSession, Name: "ErrorCodeInvalidSession"}
	}
	return ggthx.DoPoll(timeout)
}

func AddLocalInput(ggthx Session, player PlayerHandle, values []byte, size int) error {
	if ggthx == nil {
		return Error{Code: ErrorCodeInvalidSession, Name: "ErrorCodeInvalidSession"}
	}
	return ggthx.AddLocalInput(player, values, size)
}

func SynchronizeInput(ggthx Session, disconnectFlags *int) ([][]byte, error) {
	if ggthx == nil {
		return nil, Error{Code: ErrorCodeInvalidSession, Name: "ErrorCodeInvalidSession"}
	}
	return ggthx.SyncInput(disconnectFlags)
}

func DisconnectPlayer(ggthx Session, player PlayerHandle) error {
	if ggthx == nil {
		return Error{Code: ErrorCodeInvalidSession, Name: "ErrorCodeInvalidSession"}
	}
	return ggthx.DisconnectPlayer(player)
}

func AdvanceFrame(ggthx Session) error {
	if ggthx == nil {
		return Error{Code: ErrorCodeInvalidSession, Name: "ErrorCodeInvalidSession"}
	}
	return ggthx.IncrementFrame()
}

func ClientChat(ggthx Session, text string) error {
	if ggthx == nil {
		return Error{Code: ErrorCodeInvalidSession, Name: "ErrorCodeInvalidSession"}
	}
	return ggthx.Chat(text)
}

func GetNetworkStats(ggthx Session, player PlayerHandle, stats *NetworkStats) error {
	if ggthx == nil {
		return Error{Code: ErrorCodeInvalidSession, Name: "ErrorCodeInvalidSession"}
	}
	return ggthx.GetNetworkStats(stats, player)
}

func CloseSession(ggthx Session) error {
	if ggthx == nil {
		return Error{Code: ErrorCodeInvalidSession, Name: "ErrorCodeInvalidSession"}
	}
	ggthx.Close()
	return nil
}

func SetDisconnectTimeout(ggthx Session, timeout int) error {
	if ggthx == nil {
		return Error{Code: ErrorCodeInvalidSession, Name: "ErrorCodeInvalidSession"}
	}
	return ggthx.SetDisconnectTimeout(timeout)
}

func SetDisconnectNotifyStart(ggthx Session, timeout int) error {
	if ggthx == nil {
		return Error{Code: ErrorCodeInvalidSession, Name: "ErrorCodeInvalidSession"}
	}
	return ggthx.SetDisconnectNotifyStart(timeout)
}

func SartSpectating(cb *SessionCallbacks,
	gameName string, lcaolPort int, numPlayers int, inputSize int, hostIp string, hostPort int) (*SpectatorBackend, ErrorCode) {
	ggthx := NewSpectatorBackend(cb, gameName, lcaolPort, numPlayers, inputSize, hostIp, hostPort)
	return &ggthx, Ok
}
