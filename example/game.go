package main

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	ggpo "github.com/assemblaj/GGPO-Go/pkg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GameSession struct {
	backend ggpo.Backend
	game    *Game
}

var backend ggpo.Backend
var player1 Player
var player2 Player
var game *Game
var start, next, now int
var playerInputs []Input
var saveStates map[int]*Game

var localPort, numPlayers int

const FRAME_DELAY int = 2

var currentPlayer int = 1

type Game struct {
	Players []Player
}

type Player struct {
	X         float64
	Y         float64
	Color     color.Color
	PlayerNum int
}

func (p *Player) clone() Player {
	result := Player{}
	result.X = p.X
	result.Y = p.Y
	result.Color = p.Color
	result.PlayerNum = p.PlayerNum
	return result
}

func (g *Game) clone() (result *Game) {
	result = &Game{}
	*result = *g

	result.Players = make([]Player, len(g.Players))
	for i := range g.Players {
		result.Players[i] = g.Players[i].clone()
	}
	return
}

func (g *Game) Update() error {
	now = int(time.Now().UnixMilli())
	fmt.Println("Idling ")
	err := backend.Idle(int(math.Max(0, float64(next-now-1))))
	if err != nil {
		panic(err)
	}
	fmt.Println("Idling Complete")
	if now >= next {
		g.RunFrame()
		next = now + 1000/60
	}
	return nil
}

func (g *Game) RunFrame() {
	input := g.ReadInputs()
	buffer := encodeInputs(input)

	fmt.Println("Attempting to add local inputs")
	result := backend.AddLocalInput(ggpo.PlayerHandle(currentPlayer), buffer, len(buffer))
	fmt.Println("Attempt to add local inputs complete")
	if result == nil {
		fmt.Println("Attempt to add local inputs was successful")
		var values [][]byte
		disconnectFlags := 0

		fmt.Println("Attempting to synchronize inputs")
		values, result = backend.SyncInput(&disconnectFlags)
		if result == nil {
			fmt.Println("Attempt synchronize inputs was sucessful")

			inputs := decodeInputs(values)
			fmt.Println("Advancing Frame from game loop")
			g.AdvanceFrame(inputs, disconnectFlags)
		} else {
			fmt.Printf("Attempt synchronize inputs was unsuccessful: %s\n", result)
		}
	} else {
		fmt.Printf("Attempt to add local inputs unsuccessful: %s\n", result)
	}
}

func (g *Game) AdvanceFrame(inputs []InputBits, disconnectFlags int) {
	g.UpdateByInputs(inputs)
	err := backend.AdvanceFrame()
	if err != nil {
		panic(err)
	}
}

func (g *Game) UpdateByInputs(inputs []InputBits) {
	for i, input := range inputs {
		if input.isButtonOn(int(ebiten.KeyArrowUp)) {
			g.Players[i].Y--
		}
		if input.isButtonOn(int(ebiten.KeyArrowDown)) {
			g.Players[i].Y++
		}
		if input.isButtonOn(int(ebiten.KeyArrowLeft)) {
			g.Players[i].X--
		}
		if input.isButtonOn(int(ebiten.KeyArrowRight)) {
			g.Players[i].X++
		}
	}
}

func (g *Game) ReadInputs() InputBits {
	var in InputBits

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		in.setButton(int(ebiten.KeyArrowUp))
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		in.setButton(int(ebiten.KeyArrowDown))
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		in.setButton(int(ebiten.KeyArrowLeft))
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		in.setButton(int(ebiten.KeyArrowRight))
	}
	return in
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.Players {
		ebitenutil.DrawRect(screen, p.X, p.Y, 50, 50, p.Color)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Player 1: X: %.2f Y:%.2f Player 2 X: %.2f Y: %.2f",
		g.Players[0].X, g.Players[0].Y, g.Players[1].X, g.Players[1].Y))
}

func (g *Game) Layout(outsideWidth, insideWidth int) (screenWidth, screenHeight int) {
	return 320, 240
}

func beginGame(game string) bool {
	log.Println("Starting Game!")
	return true
}

/*
func saveGameState(length *int, checksum *int, frame int) ([]byte, bool) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(&game)
	if err != nil {
		log.Fatal("encode error ", err)
	}
	buffer := make([]byte, len(buf.Bytes()))
	copy(buffer, buf.Bytes())

	*checksum = calculateChecksum(buffer)
	*length = len(buffer)
	return buffer, true
}
*/
func (g *GameSession) SaveGameState(stateID int) int {
	saveStates[stateID] = game.clone()
	checksum := calculateChecksum([]byte(saveStates[stateID].String()))
	return checksum
}

func calculateChecksum(buffer []byte) int {
	cSum := md5.Sum(buffer)
	checksum := 0
	for i := 0; i < len(cSum); i++ {
		checksum += int(cSum[i])
	}
	return checksum
}

/*
func loadGameState(buffer []byte, len int) bool {
	var buf bytes.Buffer = *bytes.NewBuffer(buffer)
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(&game)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	return true
}
*/
func (g *GameSession) LoadGameState(stateID int) {
	game = saveStates[stateID]
}

func (g *GameSession) LogGameState(fileName string, buffer []byte, len int) {
	var game2 Game
	var buf bytes.Buffer = *bytes.NewBuffer(buffer)
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(&game2)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	log.Printf("%s Game State: %s\n", fileName, game2.String())
}

func (g *GameSession) SetBackend(backend ggpo.Backend) {
}

func (g *Game) String() string {
	return fmt.Sprintf("%s : %s ", g.Players[0].String(), g.Players[1].String())
}

func (p *Player) String() string {
	return fmt.Sprintf("Player %d: X:%f Y:%f Color: %s", p.PlayerNum, p.X, p.Y, p.Color)
}

func freeBuffer(buffer []byte) {

}

func (g *GameSession) AdvanceFrame(flags int) {
	fmt.Println("Advancing frame from callback. ")
	var discconectFlags int

	// Make sure we fetch the inputs from GGPO and use these to update
	// the game state instead of reading from the keyboard.
	inputs, result := backend.SyncInput(&discconectFlags)
	if result == nil {
		input := decodeInputs(inputs)
		game.AdvanceFrame(input, discconectFlags)
	}
}

func (g *GameSession) OnEvent(info *ggpo.Event) {
	switch info.Code {
	case ggpo.EventCodeConnectedToPeer:
		log.Println("EventCodeConnectedToPeer")
	case ggpo.EventCodeSynchronizingWithPeer:
		log.Println("EventCodeSynchronizingWithPeer")
	case ggpo.EventCodeSynchronizedWithPeer:
		log.Println("EventCodeSynchronizedWithPeer")
	case ggpo.EventCodeRunning:
		log.Println("EventCodeRunning")
	case ggpo.EventCodeDisconnectedFromPeer:
		log.Println("EventCodeDisconnectedFromPeer")
	case ggpo.EventCodeTimeSync:
		log.Println("EventCodeTimeSync")
	case ggpo.EventCodeConnectionInterrupted:
		log.Println("EventCodeconnectionInterrupted")
	case ggpo.EventCodeConnectionResumed:
		log.Println("EventCodeconnectionInterrupted")
	}
}

func GameInitSpectator(localPort int, numPlayers int, hostIp string, hostPort int) {
	InitGameState()

	var inputBits InputBits = 0

	var inputSize int = len(encodeInputs(inputBits))
	session := NewGameSession()

	spectator := ggpo.NewSpectatorBackend(&session, "Test", localPort, numPlayers, inputSize, hostIp, hostPort)
	backend = &spectator
	spectator.InitializeConnection()
	spectator.Start()
}

func GameInit(localPort int, numPlayers int, players []ggpo.Player, numSpectators int) {
	var result error
	InitGameState()
	var inputBits InputBits = 0
	var inputSize int = len(encodeInputs(inputBits))

	session := NewGameSession()

	peer := ggpo.NewPeer2PeerBackend(&session, "Test", localPort, numPlayers, inputSize)
	//peer := ggpo.NewSyncTestBackend(&session, "Test", numPlayers, 8, inputSize)
	backend = &peer
	peer.InitializeConnection()
	peer.Start()

	//session.SetDisconnectTimeout(3000)
	//session.SetDisconnectNotifyStart(1000)

	for i := 0; i < numPlayers+numSpectators; i++ {
		var handle ggpo.PlayerHandle
		result = peer.AddPlayer(&players[i], &handle)
		if players[i].PlayerType == ggpo.PlayerTypeLocal {
			currentPlayer = int(handle)
		}
		if result != nil {
			log.Fatalf("There's an issue from AddPlayer")
		}
		if players[i].PlayerType == ggpo.PlayerTypeLocal {
			peer.SetFrameDelay(handle, FRAME_DELAY)
		}
	}
	peer.SetDisconnectTimeout(3000)
	peer.SetDisconnectNotifyStart(1000)
}

func NewGameSession() GameSession {
	g := GameSession{}
	return g
}

func InitGameState() {
	player1 = Player{
		X:         50,
		Y:         50,
		Color:     color.RGBA{255, 0, 0, 255},
		PlayerNum: 1}
	player2 = Player{
		X:         150,
		Y:         50,
		Color:     color.RGBA{0, 0, 255, 255},
		PlayerNum: 2}
	game = &Game{
		Players: []Player{player1, player2}}
	saveStates = make(map[int]*Game)

}

func init() {

	// have to register everything with gob, maybe automate this?
	// have to inititalize all arrays
	gob.Register(color.RGBA{})
	gob.Register(Input{})

}
