package ggpo

type PlayerHandle int

type PlayerType int

const (
	PlayerTypeLocal PlayerType = iota
	PlayerTypeRemote
	PlayerTypeSpectator
)

const InvalidHandle int = -1

type Player struct {
	Size       int
	PlayerType PlayerType
	PlayerNum  int
	Remote     RemotePlayer
}

func NewLocalPlayer(size int, playerNum int) Player {
	return Player{
		Size:       size,
		PlayerNum:  playerNum,
		PlayerType: PlayerTypeLocal}
}

func NewRemotePlayer(size int, playerNum int, ipAdress string, port int) Player {
	return Player{
		Size:       size,
		PlayerNum:  playerNum,
		PlayerType: PlayerTypeRemote,
		Remote: RemotePlayer{
			IpAdress: ipAdress,
			Port:     port},
	}
}
func NewSpectatorPlayer(size int, ipAdress string, port int) Player {
	return Player{
		Size:       size,
		PlayerType: PlayerTypeSpectator,
		Remote: RemotePlayer{
			IpAdress: ipAdress,
			Port:     port},
	}
}

type RemotePlayer struct {
	IpAdress string
	Port     int
}

type LocalEndpoint struct {
	playerNum int
}
