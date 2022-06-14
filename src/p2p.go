package ggthx

type Peer1PeerBackend struct {
	callbacks     GGTHXSessionCallbacks
	poll          Poll
	sync          Sync
	udp           Udp
	endpoints     []UdpProtocol
	spectators    []UdpProtocol
	numSpectators int
	inputSize     int

	synchronizing        bool
	numPlayers           int
	nextRecommendedSleep int

	nextSpectatorFrame    int
	disconnectTimeout     int
	disconnectNotifyStart int

	localConnectStatus []UdpConnectStatus
}
