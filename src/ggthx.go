package ggthx

func StartSession(
	session GGTHXSession,
	cb GGTHXSessionCallbacks,
	game string,
	numPlayers int,
	inputSize int,
	localPort uint16) GGTHXErrorCode {
	return GGTHX_OK
}

func AddPlayer(ggthx GGTHXSession,
	player *GGTHXPlayer,
	handle *GGTHXPlayerHandle) GGTHXErrorCode {
	if ggthx == nil {
		return GGTHX_ERRORCODE_INVALID_SESSION
	}
	return ggthx.AddPlayer(player, handle)
}

func StartSyncTest(
	cb *GGTHXSessionCallbacks,
	game string,
	numPlayers int,
	inputSize int,
	frames int) (*SyncTestBackend, GGTHXErrorCode) {
	ggthx := NewSyncTestBackend(cb, game, numPlayers, frames, inputSize)
	return &ggthx, GGTHX_OK
}

func SetFrameDelay(ggthx GGTHXSession,
	player GGTHXPlayerHandle,
	frameDelay int) GGTHXErrorCode {
	if ggthx == nil {
		return GGTHX_ERRORCODE_INVALID_SESSION
	}
	return ggthx.SetFrameDelay(player, frameDelay)
}

func Idle(ggthx GGTHXSession, timeout int) GGTHXErrorCode {
	if ggthx == nil {
		return GGTHX_ERRORCODE_INVALID_SESSION
	}
	return ggthx.DoPoll(timeout)
}

func AddLocalInput(ggthx GGTHXSession, player GGTHXPlayerHandle, values []byte, size int) GGTHXErrorCode {
	if ggthx == nil {
		return GGTHX_ERRORCODE_INVALID_SESSION
	}
	return ggthx.AddLocalInput(player, values, size)
}

func SynchronizeInput(ggthx GGTHXSession, disconnectFlags *int) ([]byte, GGTHXErrorCode) {
	if ggthx == nil {
		return nil, GGTHX_ERRORCODE_INVALID_SESSION
	}
	return ggthx.SyncInput(disconnectFlags)
}

func DisconnectPlayer(ggthx GGTHXSession, player GGTHXPlayerHandle) GGTHXErrorCode {
	if ggthx == nil {
		return GGTHX_ERRORCODE_INVALID_SESSION
	}
	return ggthx.DisconnectPlayer(&player)
}

func AdvanceFrame(ggthx GGTHXSession) GGTHXErrorCode {
	if ggthx == nil {
		return GGTHX_ERRORCODE_INVALID_SESSION
	}
	return ggthx.IncrementFrame()
}

func ClientChat(ggthx GGTHXSession, text string) GGTHXErrorCode {
	if ggthx == nil {
		return GGTHX_ERRORCODE_INVALID_SESSION
	}
	return ggthx.Chat(text)
}

func GetNetworkStats(ggthx GGTHXSession, player GGTHXPlayerHandle, stats *GGTHXNetworkStats) GGTHXErrorCode {
	if ggthx == nil {
		return GGTHX_ERRORCODE_INVALID_SESSION
	}
	return ggthx.GetNetworkStats(stats, player)
}

func CloseSession(ggthx GGTHXSession) GGTHXErrorCode {
	if ggthx == nil {
		return GGTHX_ERRORCODE_INVALID_SESSION
	}
	ggthx.Close()
	return GGTHX_OK
}

func SetDisconnectTimeout(ggthx GGTHXSession, timeout int) GGTHXErrorCode {
	if ggthx == nil {
		return GGTHX_ERRORCODE_INVALID_SESSION
	}
	return ggthx.SetDisconnectTimeout(timeout)
}

func SetDisconnectNotifyStart(ggthx GGTHXSession, timeout int) GGTHXErrorCode {
	if ggthx == nil {
		return GGTHX_ERRORCODE_INVALID_SESSION
	}
	return ggthx.SetDisconnectNotifyStart(timeout)
}

func SartSpectating(cb *GGTHXSessionCallbacks,
	gameName string, lcaolPort int, numPlayers int, inputSize int, hostIp string, hostPort int) (*SpectatorBackend, GGTHXErrorCode) {
	ggthx := NewSpectatorBackend(cb, gameName, lcaolPort, numPlayers, inputSize, hostIp, hostPort)
	return &ggthx, GGTHX_OK
}
