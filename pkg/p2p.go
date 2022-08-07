package ggpo

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/assemblaj/GGPO-Go/internal/input"
	"github.com/assemblaj/GGPO-Go/internal/polling"
	"github.com/assemblaj/GGPO-Go/internal/protocol"
	"github.com/assemblaj/GGPO-Go/internal/transport"
	"github.com/assemblaj/GGPO-Go/internal/util"
)

const (
	RecommendationInterval       = 240
	DefaultDisconnectTimeout     = 5000
	DefaultDisconnectNotifyStart = 750
)

type Peer2PeerBackend struct {
	session       Session
	poll          polling.Poller
	sync          Sync
	connection    transport.Connection
	endpoints     []protocol.UdpProtocol
	spectators    []protocol.UdpProtocol
	numSpectators int
	inputSize     int

	synchronizing        bool
	numPlayers           int
	nextRecommendedSleep int

	nextSpectatorFrame    int
	disconnectTimeout     int
	disconnectNotifyStart int

	localConnectStatus []transport.UdpConnectStatus

	localPort int
}

func NewPeer2PeerBackend(cb Session, gameName string,
	localPort int, numPlayers int, inputSize int) Peer2PeerBackend {
	p := Peer2PeerBackend{}
	p.numPlayers = numPlayers
	p.inputSize = inputSize
	p.session = cb
	p.synchronizing = true
	p.disconnectTimeout = DefaultDisconnectTimeout
	p.disconnectNotifyStart = DefaultDisconnectNotifyStart
	var poll polling.Poll = polling.NewPoll()
	p.poll = &poll

	//p.udp = NewUdp(&p, localPort)
	p.localPort = localPort

	p.localConnectStatus = make([]transport.UdpConnectStatus, transport.UDPMsgMaxPlayers)
	for i := 0; i < len(p.localConnectStatus); i++ {
		p.localConnectStatus[i].LastFrame = -1
	}
	var config SyncConfig = NeweSyncConfig(
		p.session, MaxPredictionFrames, p.numPlayers, p.inputSize)
	config.numPlayers = numPlayers
	config.inputSize = inputSize
	config.session = p.session
	config.numPredictionFrames = MaxPredictionFrames
	p.sync = NewSync(p.localConnectStatus, &config)
	p.endpoints = make([]protocol.UdpProtocol, numPlayers)
	p.spectators = make([]protocol.UdpProtocol, MaxSpectators)
	//messages := make(chan UdpPacket)
	//p.poll.RegisterLoop(&p.udp, nil )
	//go p.udp.Read()
	//go p.udp.ReadMsg(messages)
	//go p.OnMsg(messages)
	return p
}

func (p *Peer2PeerBackend) Close() error {
	for _, e := range p.endpoints {
		e.Close()
	}
	for _, s := range p.spectators {
		s.Close()
	}
	return nil
}
func (p *Peer2PeerBackend) Idle(timeout int, timeFunc ...polling.FuncTimeType) error {
	if !p.sync.InRollback() {
		if len(timeFunc) == 0 {
			p.poll.Pump()
		} else {
			p.poll.Pump(timeFunc[0])
		}
		p.PollUdpProtocolEvents()

		if !p.synchronizing {
			p.sync.CheckSimulation(timeout)

			// notify all of our endpoints of their local frame number for their
			// next connection quality report
			currentFrame := p.sync.FrameCount()
			for i := 0; i < p.numPlayers; i++ {
				p.endpoints[i].SetLocalFrameNumber(currentFrame)
			}

			var totalMinConfirmed int
			if p.numPlayers <= 2 {
				totalMinConfirmed = p.Poll2Players(currentFrame)
			} else {
				totalMinConfirmed = p.PollNPlayers(currentFrame)
			}

			log.Printf("last confirmed frame in p2p backend is %d.\n", totalMinConfirmed)
			if totalMinConfirmed >= 0 {
				if totalMinConfirmed == math.MaxInt {
					return Error{Code: ErrorCodeGeneralFailure, Name: "ErrorCodeGeneralFailure"}
				}
				if p.numSpectators > 0 {
					for p.nextSpectatorFrame <= totalMinConfirmed {
						log.Printf("pushing frame %d to spectators.\n", p.nextSpectatorFrame)

						var input input.GameInput
						var inputs [][]byte
						input.Frame = p.nextSpectatorFrame
						input.Size = p.inputSize * p.numPlayers
						inputs, _ = p.sync.GetConfirmedInputs(p.nextSpectatorFrame)
						for i, _ := range inputs {
							input.Bits = append(input.Bits, inputs[i]...)
						}
						for i := 0; i < p.numSpectators; i++ {
							p.spectators[i].SendInput(&input)
						}
						p.nextSpectatorFrame++
					}
				}
				log.Printf("setting confirmed frame in sync to %d.\n", totalMinConfirmed)
				p.sync.SetLastConfirmedFrame(totalMinConfirmed)
			}

			// send timesync notifications if now is the proper time
			if currentFrame > p.nextRecommendedSleep {
				interval := 0
				for i := 0; i < p.numPlayers; i++ {
					interval = util.Max(interval, p.endpoints[i].RecommendFrameDelay())
				}

				if interval > 0 {
					var info Event
					info.Code = EventCodeTimeSync
					info.framesAhead = interval
					p.session.OnEvent(&info)
					p.nextRecommendedSleep = currentFrame + RecommendationInterval
				}
			}
			// because GGPO had this
			if timeout > 0 {
				time.Sleep(time.Millisecond)
			}
		}
	}
	return nil
}

// Checks each endpoint to see if it's running (i.e we haven't called Synchronize or Disconnect
// on it lately )
// Again check if it's disconnected (GetPeerConnectionStatus which would check if the user has
// sent a disconnect message.
// Try to get the minimum last confirmed frame across all the inputs
// If by chance the user sent a disconnect request and we haven't disconnected them yet,
// We disconnect them.
// The value returned by this function (the total minimum confirmed frame across all the inputs)
// is used in Idle to set the last confirmed frame in the sync backend, which tells
// the input queue to discard the frame before that.
func (p *Peer2PeerBackend) Poll2Players(currentFrame int) int {
	totalMinConfirmed := int32(math.MaxInt32)
	for i := 0; i < p.numPlayers; i++ {
		queueConnected := true
		if p.endpoints[i].IsRunning() {
			var ignore int32
			queueConnected = p.endpoints[i].GetPeerConnectStatus(i, &ignore)
		}
		if !p.localConnectStatus[i].Disconnected {
			totalMinConfirmed = util.Min(p.localConnectStatus[i].LastFrame, totalMinConfirmed)
		}
		log.Printf("  local endp: connected = %t, last_received = %d, total_min_confirmed = %d.\n",
			!p.localConnectStatus[i].Disconnected, p.localConnectStatus[i].LastFrame, totalMinConfirmed)
		if !queueConnected && !p.localConnectStatus[i].Disconnected {
			log.Printf("disconnecting i %d by remote request.\n", i)
			p.DisconnectPlayerQueue(i, int(totalMinConfirmed))
		}
		log.Printf("  total_min_confirmed = %d.\n", totalMinConfirmed)
	}
	return int(totalMinConfirmed)
}

// Just for parity with GGPO. Don't care to actually use this.
func (p *Peer2PeerBackend) PollNPlayers(currentFrame int) int {
	var i, queue int
	var lastRecieved int32

	totalMinConfirmed := int32(math.MaxInt32)
	for queue = 0; queue < p.numPlayers; queue++ {
		queueConnected := true
		queueMinConfirmed := int32(math.MaxInt32)
		log.Printf("considering queue %d.\n", queue)
		for i = 0; i < p.numPlayers; i++ {
			// we're going to do a lot of logic here in consideration of endpoint i.
			// keep accumulating the minimum confirmed point for all n*n packets and
			// throw away the rest. -pond3r
			if p.endpoints[i].IsRunning() {
				connected := p.endpoints[i].GetPeerConnectStatus(queue, &lastRecieved)

				queueConnected = queueConnected && connected
				queueMinConfirmed = util.Min(lastRecieved, queueMinConfirmed)
				log.Printf("  endpoint %d: connected = %t, last_received = %d, queue_min_confirmed = %d.\n",
					i, connected, lastRecieved, queueMinConfirmed)
			} else {
				log.Printf("  endpoint %d: ignoring... not running.\n", i)
			}
		}
		// merge in our local status only if we're still connected!
		if !p.localConnectStatus[queue].Disconnected {
			queueMinConfirmed = util.Min(p.localConnectStatus[queue].LastFrame, queueMinConfirmed)
		}
		log.Printf("  local endp: connected = %t, last_received = %d, queue_min_confirmed = %d.\n",
			!p.localConnectStatus[queue].Disconnected, p.localConnectStatus[queue].LastFrame, queueMinConfirmed)

		if queueConnected {
			totalMinConfirmed = util.Min(queueMinConfirmed, totalMinConfirmed)
		} else {
			// check to see if this disconnect notification is further back than we've been before.  If
			// so, we need to re-adjust.  This can happen when we detect our own disconnect at frame n
			// and later receive a disconnect notification for frame n-1.
			if !p.localConnectStatus[queue].Disconnected || p.localConnectStatus[queue].LastFrame > queueMinConfirmed {
				log.Printf("disconnecting queue %d by remote request.\n", queue)
				p.DisconnectPlayerQueue(queue, int(queueMinConfirmed))
			}
		}
		log.Printf("  total_min_confirmed = %d.\n", totalMinConfirmed)
	}
	return int(totalMinConfirmed)
}

/*
	Giving each spectator and remote player their own UDP object (which GGPO didn't do)
	a copy of the poll, a copy of our localConnectStatus (might want to send a pointer?)
	Setting the default disconnect timeout and disconnect notify
	And calling the synchronize method, which sends a sync request to that endpoint.
*/
func (p *Peer2PeerBackend) AddRemotePlayer(ip string, port int, queue int) {
	p.synchronizing = true
	p.endpoints[queue] = protocol.NewUdpProtocol(p.connection, queue, ip, port, &p.localConnectStatus)
	// have to reqgister the loop from here or else the Poll won't see changed state
	// that we've initiated.
	p.poll.RegisterLoop(&(p.endpoints[queue]), nil)

	// actually this Idle wouldn't run at all if it wasn't called from here.
	//p.poll.RegisterLoop(&udp, nil)
	p.endpoints[queue].SetDisconnectTimeout(p.disconnectTimeout)
	p.endpoints[queue].SetDisconnectNotifyStart(p.disconnectNotifyStart)
	p.endpoints[queue].Synchronize()
}

func (p *Peer2PeerBackend) AddSpectator(ip string, port int) error {
	if p.numSpectators == MaxSpectators {
		return Error{Code: ErrorCodeTooManySpectators, Name: "ErrorCodeTooManySpectators"}
	}
	// Currently, we can only add spectators before the game starts.
	if !p.synchronizing {
		return Error{Code: ErrorCodeInvalidRequest, Name: "ErrorCodeInvalidRequest"}
	}
	queue := p.numSpectators
	p.numSpectators++
	p.spectators[queue] = protocol.NewUdpProtocol(p.connection, queue+1000, ip, port, &p.localConnectStatus)
	p.poll.RegisterLoop(&(p.spectators[queue]), nil)
	p.spectators[queue].SetDisconnectTimeout(p.disconnectTimeout)
	p.spectators[queue].SetDisconnectNotifyStart(p.disconnectNotifyStart)
	p.spectators[queue].Synchronize()

	return nil
}

// Adds player or spectator
// Maps to top level API function
func (p *Peer2PeerBackend) AddPlayer(player *Player, handle *PlayerHandle) error {
	if player.PlayerType == PlayerTypeSpectator {
		return p.AddSpectator(player.Remote.IpAdress, player.Remote.Port)
	}

	queue := player.PlayerNum - 1
	if player.PlayerNum < 1 || player.PlayerNum > p.numPlayers {
		return Error{Code: ErrorCodePlayerOutOfRange, Name: "ErrorCodePlayerOutOfRange"}
	}
	*handle = p.QueueToPlayerHandle(queue)

	if player.PlayerType == PlayerTypeRemote {
		p.AddRemotePlayer(player.Remote.IpAdress, player.Remote.Port, queue)
	}

	return nil
}

// Sends input to the synchronization layer and to all the endpoints.
// Which adds it to those respective queues (inputQueue for synchronization layer, pendingOutput
// for endpoint)
// Maps to top level API function.
func (p *Peer2PeerBackend) AddLocalInput(player PlayerHandle, values []byte, size int) error {
	var queue int
	var localInput input.GameInput
	var err error

	if p.sync.InRollback() {
		return Error{Code: ErrorCodeInRollback, Name: "ErrorCodeInRollback"}
	}
	if p.synchronizing {
		return Error{Code: ErrorCodeNotSynchronized, Name: "ErrorCodeNotSynchronized"}
	}

	result := p.PlayerHandleToQueue(player, &queue)
	if result != nil {
		return result
	}

	localInput, err = input.NewGameInput(-1, values, size)
	if err != nil {
		panic(err)
	}

	// Feed the input for the current frame into the synchronization layer.
	if !p.sync.AddLocalInput(queue, &localInput) {
		return Error{Code: ErrorCodePredictionThreshod, Name: "ErrorCodePredictionThreshod"}
	}

	if localInput.Frame != input.NullFrame {
		// Update the local connect status state to indicate that we've got a
		// confirmed local frame for this player.  this must come first so it
		// gets incorporated into the next packet we send.
		// - pond3r

		log.Printf("setting local connect status for local queue %d to %d", queue, localInput.Frame)
		p.localConnectStatus[queue].LastFrame = int32(localInput.Frame)

		// Send the input to all the remote players.
		for i := 0; i < p.numPlayers; i++ {
			if p.endpoints[i].IsInitialized() {
				p.endpoints[i].SendInput(&localInput)
			}
		}
	}

	return nil
}

// Maps to top level API function
// Returns input from all players (both local and remote) from the input queue
// Which can be a prediction (i.e the last frame) or the latest input recieved
// Also used to fetch inputs from GGPO  to update the game states during the advance
// frame callback
func (p *Peer2PeerBackend) SyncInput(disconnectFlags *int) ([][]byte, error) {
	// Wait until we've started to return inputs.
	if p.synchronizing {
		return nil, Error{Code: ErrorCodeNotSynchronized, Name: "ErrorCodeNotSynchronized"}
	}
	values, flags := p.sync.SynchronizeInputs()

	if disconnectFlags != nil {
		*disconnectFlags = flags
	}
	return values, nil
}

// Maps to top level API AdvanceFrame function
// sync.IncremenetFrame increments the frame count and saves the
// current state via user provided callback
// Do Poll Not only runs everything in the system that's registered to poll
// it... well does everything. I'll get ti it when I get to it.
func (p *Peer2PeerBackend) AdvanceFrame() error {
	log.Printf("End of frame (%d)...\n", p.sync.FrameCount())
	p.sync.AdvanceFrame()
	err := p.Idle(0)
	if err != nil {
		panic(err)
	}

	return nil
}

// Handles all the events  for all spactors and players. Done OnPoll
func (p *Peer2PeerBackend) PollUdpProtocolEvents() {
	for i := 0; i < p.numPlayers; i++ {
		for {
			evt, err := p.endpoints[i].GetEvent()
			if err != nil {
				break
			} else {
				err := p.OnUdpProtocolPeerEvent(evt, i)
				if err != nil {
					panic(err)
				}
			}
		}
	}
	for i := 0; i < p.numSpectators; i++ {
		for {
			evt, err := p.spectators[i].GetEvent()
			if err != nil {
				break
			} else {
				p.OnUdpProtocolSpectatorEvent(evt, i)
			}
		}
	}
}

// Takes events that come from UdpProtocol and if they are input,
// Sends remote input to Sync which adds it to its InputQueue
// Also updates that lastFrame for this endpoint as the most
// recently recived frame.
// Disconnects if necesary
func (p *Peer2PeerBackend) OnUdpProtocolPeerEvent(evt *protocol.UdpProtocolEvent, queue int) error {
	handle := p.QueueToPlayerHandle(queue)
	p.OnUdpProtocolEvent(evt, handle)
	switch evt.Type() {
	case protocol.InputEvent:
		if !p.localConnectStatus[queue].Disconnected {
			currentRemoteFrame := p.localConnectStatus[queue].LastFrame
			newRemoteFrame := evt.Input.Frame
			if !(currentRemoteFrame == -1 || int32(newRemoteFrame) == (currentRemoteFrame+1)) {
				return errors.New("ggpo Peer2PeerBackend OnUdpProtocolPeerEvent : !(currentRemoteFrame == -1 || newRemoteFrame == (currentRemoteFrame+1)) ")
			}

			p.sync.AddRemoteInput(queue, &evt.Input)
			// Notify the other endpoints which frame we received from a peer
			log.Printf("setting remote connect status for queue %d to %d\n", queue,
				evt.Input.Frame)
			p.localConnectStatus[queue].LastFrame = int32(evt.Input.Frame)
		}
	case protocol.DisconnectedEvent:
		err := p.DisconnectPlayer(handle)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

// Every Idle, every endpoint and spectator goes through its event queue
// handles each event and pops it from the queue.  Though most of the logic
// for handling these events is the same (see: OnUdpProtocolEvent ), spectators
// and peers handle certain events differently
func (p *Peer2PeerBackend) OnUdpProtocolSpectatorEvent(evt *protocol.UdpProtocolEvent, queue int) {
	handle := p.QueueToSpectatorHandle(queue)
	p.OnUdpProtocolEvent(evt, handle)

	var info Event
	switch evt.Type() {
	case protocol.DisconnectedEvent:
		p.spectators[queue].Disconnect()

		info.Code = EventCodeDisconnectedFromPeer
		info.player = handle
		p.session.OnEvent(&info)
	}
}

// Logic for parsing UdpProtocol events and sending them up to the user via callbacks.
// In P2P Backend, called by OnUdpProtocolSpectatorEvent and OnUdpProtocolPeerEvent,
// which themselves are called by PollUdpProtocolEvents, which happens every Idle
func (p *Peer2PeerBackend) OnUdpProtocolEvent(evt *protocol.UdpProtocolEvent, handle PlayerHandle) {
	var info Event

	switch evt.Type() {
	case protocol.ConnectedEvent:
		info.Code = EventCodeConnectedToPeer
		info.player = handle
		p.session.OnEvent(&info)

	case protocol.SynchronizingEvent:
		info.Code = EventCodeSynchronizingWithPeer
		info.player = handle
		info.count = evt.Count
		info.total = evt.Total
		p.session.OnEvent(&info)

	case protocol.SynchronziedEvent:
		info.Code = EventCodeSynchronizedWithPeer
		info.player = handle
		p.session.OnEvent(&info)

		p.CheckInitialSync()

	case protocol.NetworkInterruptedEvent:
		info.Code = EventCodeConnectionInterrupted
		info.player = handle
		info.disconnectTimeout = evt.DisconnectTimeout
		p.session.OnEvent(&info)

	case protocol.NetworkResumedEvent:
		info.Code = EventCodeConnectionResumed
		info.player = handle
		p.session.OnEvent(&info)
	}
}

/*
 * Called only as the result of a local decision to disconnect.  The remote
 * decisions to disconnect are a result of us parsing the peer_connect_settings
 * blob in every endpoint periodically.
 * - pond3r
	This is the function that's called when the UdpProtocol endpoint recogniizes
	a disconnect (that lastRecvTime + disconnectTimeout < now) and sends that event
	up to the backend.
    Also maps to API function
*/
func (p *Peer2PeerBackend) DisconnectPlayer(player PlayerHandle) error {
	var queue int

	result := p.PlayerHandleToQueue(player, &queue)
	if result != nil {
		return result
	}

	if p.localConnectStatus[queue].Disconnected {
		return Error{Code: ErrorCodePlayerDisconnected, Name: "ErrorCodePlayerDisconnected"}
	}

	if !p.endpoints[queue].IsInitialized() {
		currentFrame := p.sync.FrameCount()
		// xxx: we should be tracking who the local player is, but for now assume
		// that if the endpoint is not initalized, this must be the local player. - pond3r

		// 100% sure this assumption won't be applicable for me, but we'll see.
		// Not that it matters lol
		log.Printf("Disconnecting local player %d at frame %d by user request.\n",
			queue, p.localConnectStatus[queue].LastFrame)
		// Disconnecting all the other players too
		for i := 0; i < p.numPlayers; i++ {
			if p.endpoints[i].IsInitialized() {
				p.DisconnectPlayerQueue(i, currentFrame)
			}
		}
	} else {
		log.Printf("Disconnecting queue %d at frame %d by user request.\n",
			queue, p.localConnectStatus[queue].LastFrame)
		p.DisconnectPlayerQueue(queue, int(p.localConnectStatus[queue].LastFrame))
	}
	return nil
}

/*
	Sets the enpoints' state to disconnected.
	Also sets localConnectStatus to disconnected which is used all over the p2p backend to
	verify if an endpoint is connected
	And adjusts the simulation to get to syncto frames
	Then sends the Disconnect Event up to te user.
*/
func (p *Peer2PeerBackend) DisconnectPlayerQueue(queue int, syncto int) {
	var info Event
	frameCount := p.sync.FrameCount()

	p.endpoints[queue].Disconnect()

	log.Printf("Changing queue %d local connect status for last frame from %d to %d on disconnect request (current: %d).\n",
		queue, p.localConnectStatus[queue].LastFrame, syncto, frameCount)

	p.localConnectStatus[queue].Disconnected = true
	p.localConnectStatus[queue].LastFrame = int32(syncto)

	if syncto < frameCount {
		log.Printf("adjusting simulation to account for the fact that %d disconnected @ %d.\n", queue, syncto)
		p.sync.AdjustSimulation(syncto)
		log.Printf("Finished adjusting simulation.\n")
	}

	info.Code = EventCodeDisconnectedFromPeer
	info.player = p.QueueToPlayerHandle(queue)
	p.session.OnEvent(&info)

	p.CheckInitialSync()
}

/*
	Gets network stats for that specific play from their UdpProtocol Endpoint
	Includes ping, sendQueLen, kbpsSent, remoteFramesBehind and remoteFrameAdvantage
	All coming from the UdpProtocol Endpoint
	Maps to top level API function.
*/
func (p *Peer2PeerBackend) GetNetworkStats(stats *protocol.NetworkStats, player PlayerHandle) error {
	var queue int

	result := p.PlayerHandleToQueue(player, &queue)
	if result != nil {
		return result
	}

	p.endpoints[queue].GetNetworkStats(stats)

	return nil
}

/*
	Sets frame delay for that specific player's input queue in Sync.
	Frame delay is used in the input queue, when remote inputs are recieved from
	the UdpProtocol and sent to Sync, sync then adds those inputs to the input queue
	for that specific player, and sort of artificially corrects the frame that player
	should be on by increasing it frameDelay amount
	Maps to top level API function
*/
func (p *Peer2PeerBackend) SetFrameDelay(player PlayerHandle, delay int) error {
	var queue int

	result := p.PlayerHandleToQueue(player, &queue)
	if result != nil {
		return result
	}
	p.sync.SetFrameDelay(queue, delay)
	return nil
}

/*
	Propagates the disconnect timeout to all of the endpoints.
    lastRecvTime + disconnectTimeout < now means the endpoint has stopped
	recieving packets and we are now disconnecting, effectively timing out.
	The Udp endpoint propogates the Disconnect Event up to the backend.
	Which, in the P2P Backend, Disconnects the Player from the backend,
	then sends the event upward.
	Mapped to top level API function
*/
func (p *Peer2PeerBackend) SetDisconnectTimeout(timeout int) error {
	p.disconnectTimeout = timeout
	for i := 0; i < p.numPlayers; i++ {
		if p.endpoints[i].IsInitialized() {
			p.endpoints[i].SetDisconnectTimeout(p.disconnectTimeout)
		}
	}
	return nil
}

/*
	Propagates the disconnect notify start to all of the endpoints
	lastRecTime + disconnectNotifyStart < now  means the endpoint has
	stopped recieving packets. The udp endpoint starts sending a NetworkInterrupted event
	up to the backend, check sends it up to the user via the API's callbacks.
	Mapped to top level Api function
*/
func (p *Peer2PeerBackend) SetDisconnectNotifyStart(timeout int) error {
	p.disconnectNotifyStart = timeout
	for i := 0; i < p.numPlayers; i++ {
		if p.endpoints[i].IsInitialized() {
			p.endpoints[i].SetDisconnectNotifyStart(p.disconnectNotifyStart)
		}
	}
	return nil
}

/*
	Used for getting the index for PlayerHandle mapped arrays such as
	endpoints, localConnectStatus, andinputQueues in Sync
*/
func (p *Peer2PeerBackend) PlayerHandleToQueue(player PlayerHandle, queue *int) error {
	offset := int(player) - 1
	if offset < 0 || offset >= p.numPlayers {
		return Error{Code: ErrorCodeInvalidPlayerHandle, Name: "ErrorCodeInvalidPlayerHandle"}
	}
	*queue = offset
	return nil
}

/*
	Does the inverse of the above. Turns index index into human readible number
*/
func (p *Peer2PeerBackend) QueueToPlayerHandle(queue int) PlayerHandle {
	return PlayerHandle(queue + 1)
}

func (p *Peer2PeerBackend) QueueToSpectatorHandle(queue int) PlayerHandle {
	return PlayerHandle(queue + 1000) /* out of range of the player array, basically  - pond3r*/
}

/*
	Propogates messages to all endpoints and spectators (?)
	As of right now it hands the message off to the first endpoint that
	handles it then returns?
*/
func (p *Peer2PeerBackend) HandleMessage(ipAddress string, port int, msg transport.UDPMessage, length int) {
	for i := 0; i < p.numPlayers; i++ {
		if p.endpoints[i].HandlesMsg(ipAddress, port) {
			p.endpoints[i].OnMsg(msg, length)
			return
		}
	}
	for i := 0; i < p.numSpectators; i++ {
		if p.spectators[i].HandlesMsg(ipAddress, port) {
			p.spectators[i].OnMsg(msg, length)
			return
		}
	}
}

/*
	Checks if all endpoints and spectators are initialized and synchronized
	and sends an event when they are.
*/
func (p *Peer2PeerBackend) CheckInitialSync() {
	var i int

	if p.synchronizing {
		// Check to see if everyone is now synchronized. If so,
		// go and tell the client that we're ok to accept in.
		for i = 0; i < p.numPlayers; i++ {
			// xxx: IsInitialized() must go... we're actually using it as a proxy for "represents the local player" -pond3r
			if p.endpoints[i].IsInitialized() && !p.endpoints[i].IsSynchronized() && !p.localConnectStatus[i].Disconnected {
				return
			}
		}

		for i = 0; i < p.numSpectators; i++ {
			if p.spectators[i].IsInitialized() && !p.spectators[i].IsSynchronized() {
				return
			}
		}

		var info Event
		info.Code = EventCodeRunning
		p.session.OnEvent(&info)
		p.synchronizing = false
	}
}

func (p *Peer2PeerBackend) InitializeConnection(t ...transport.Connection) error {
	if len(t) == 0 {
		p.connection = transport.NewUdp(p, p.localPort)
		return nil
	}
	p.connection = t[0]
	return nil
}

func (p *Peer2PeerBackend) Start() {
	//messages := make(chan UdpPacket)
	//go p.udp.ReadMsg(messages)
	//go p.OnMsg(messages)
	//p.udp.messageHandler = p
	fmt.Printf("In peer2peerbackend start udp: %v\n", p.connection)
	go p.connection.Read()
}
