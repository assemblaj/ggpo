package main

func GGTHXStartSession(
	session GGTHXSession,
	cb GGTHXSessionCallbacks,
	game string,
	numPlayers int,
	inputSize int,
	localPort uint16) GGTHXErrorCode {
	return GGTHX_OK
}

func GGTHXAddPlayer(ggthx GGTHXSession,
	player *GGTHXPlayer,
	handle *GGTHXPlayerHandle) GGTHXErrorCode {
	if ggthx == nil {
		return GGTHX_ERRORCODE_INVALID_SESSION
	}
	return ggthx.AddPlayer(player, handle)
}

func GGTHXStartSyncTest(ggthx GGTHXSession,
	cb *GGTHXSessionCallbacks,
	game string,
	numPlayers int,
	inputSize int,
	frames int) GGTHXErrorCode {
	//*ggthx = (GGTHXSession *)new SyncTestBackend(cb, game, frames, num_players);
	return GGTHX_OK
}

func GGTHXSetFrameDelay(ggthx GGTHXSession,
	player GGTHXPlayerHandle,
	frameDelay int) GGTHXErrorCode {
	if ggthx == nil {
		return GGTHX_ERRORCODE_INVALID_SESSION
	}
	return ggthx.SetFrameDelay(player, frameDelay)
}

func GGTHXIdle(ggthx GGTHXSession, timeout int) GGTHXErrorCode {
	if ggthx == nil {
		return GGTHX_ERRORCODE_INVALID_SESSION
	}
	return ggthx.DoPoll(timeout)
}

func main() {

}
