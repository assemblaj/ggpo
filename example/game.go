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

	ggthx "github.com/assemblaj/ggthx/src"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var session *ggthx.SyncTestBackend
var callbacks ggthx.GGTHXSessionCallbacks
var player1 Player
var player2 Player
var game *Game
var start, next, now int
var playerInputs []Input

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

func (g *Game) Update() error {
	now = int(time.Now().UnixMilli())
	ggthx.Idle(session, int(math.Max(0, float64(next-now-1))))
	if now >= next {
		g.RunFrame()
		next = now + 1000/60
	}
	return nil
}

func (g *Game) RunFrame() {
	input := g.ReadInputs()
	buffer := encodeInputs(input)

	result := ggthx.AddLocalInput(session, ggthx.GGTHXPlayerHandle(currentPlayer), buffer, len(buffer))
	if ggthx.GGTHX_SUCESS(result) {
		var values []byte
		disconnectFlags := 0

		values, result = ggthx.SynchronizeInput(session, &disconnectFlags)
		if ggthx.GGTHX_SUCESS(result) {
			var buf bytes.Buffer = *bytes.NewBuffer(values)
			dec := gob.NewDecoder(&buf)
			err := dec.Decode(&input)
			if err != nil {
				log.Fatal("decode error:", err)
			}
			g.AdvanceFrame(input, disconnectFlags)
		}
	}
}

func (g *Game) AdvanceFrame(inputs Input, disconnectFlags int) {
	g.UpdateByInputs(inputs)
	ggthx.AdvanceFrame(session)
}

func (g *Game) UpdateByInputs(inputs Input) {
	for _, v := range inputs.Key {
		if ebiten.Key(v) == ebiten.KeyArrowUp {
			g.Players[inputs.PlayerNum-1].Y--
		}
		if ebiten.Key(v) == ebiten.KeyArrowDown {
			g.Players[inputs.PlayerNum-1].Y++
		}
		if ebiten.Key(v) == ebiten.KeyArrowLeft {
			g.Players[inputs.PlayerNum-1].X--
		}
		if ebiten.Key(v) == ebiten.KeyArrowRight {
			g.Players[inputs.PlayerNum-1].X++
		}

	}

}

func (g *Game) ReadInputs() Input {
	in := Input{
		PlayerNum: currentPlayer,
		Key:       make([]ebiten.Key, 0)}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		in.Key = append(in.Key, ebiten.KeyArrowUp)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		in.Key = append(in.Key, ebiten.KeyArrowDown)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		in.Key = append(in.Key, ebiten.KeyArrowLeft)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		in.Key = append(in.Key, ebiten.KeyArrowRight)
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

func calculateChecksum(buffer []byte) int {
	cSum := md5.Sum(buffer)
	checksum := 0
	for i := 0; i < len(cSum); i++ {
		checksum += int(cSum[i])
	}
	return checksum
}

func loadGameState(buffer []byte, len int) bool {
	var buf bytes.Buffer = *bytes.NewBuffer(buffer)
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(&game)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	return true
}

func logGameState(fileName string, buffer []byte, len int) bool {
	var game2 Game
	var buf bytes.Buffer = *bytes.NewBuffer(buffer)
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(&game2)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	log.Printf("%s Game State: %s\n", fileName, game2)
	return true
}

func (g Game) String() string {
	return fmt.Sprintf("%s : %s ", g.Players[0], g.Players[1])
}

func (p Player) String() string {
	return fmt.Sprintf("Player %d: X:%f Y:%f Color: %s", p.PlayerNum, p.X, p.Y, p.Color)
}

func freeBuffer(buffer []byte) {

}

func advanceFrame(flags int) bool {
	var discconectFlags int

	// Make sure we fetch the inputs from GGPO and use these to update
	// the game state instead of reading from the keyboard.
	inputs, result := ggthx.SynchronizeInput(session, &discconectFlags)
	if result != ggthx.GGTHX_OK {
		log.Fatal("Error from GGTHXSynchronizeInput")
	}

	input := decodeInputs(inputs)
	game.AdvanceFrame(input, discconectFlags)

	return true
}

func onEvent(info *ggthx.GGTHXEvent) bool {
	switch info.Code {
	case ggthx.GGTHX_EVENTCODE_CONNECTED_TO_PEER:
		log.Println("GGTHX_EVENTCODE_CONNECTED_TO_PEER")
		break
	case ggthx.GGTHX_EVENTCODE_SYNCHRONIZING_WITH_PEER:
		log.Println("GGTHX_EVENTCODE_SYNCHRONIZING_WITH_PEER")
		break
	case ggthx.GGTHX_EVENTCODE_SYNCHRONIZED_WITH_PEER:
		log.Println("GGTHX_EVENTCODE_SYNCHRONIZED_WITH_PEER")
		break
	case ggthx.GGTHX_EVENTCODE_RUNNING:
		log.Println("GGTHX_EVENTCODE_RUNNING")
		break
	case ggthx.GGTHX_EVENTCODE_DISCONNECTED_FROM_PEER:
		log.Println("GGTHX_EVENTCODE_DISCONNECTED_FROM_PEER")
		break
	case ggthx.GGTHX_EVENTCODE_TIMESYNC:
		log.Println("GGTHX_EVENTCODE_TIMESYNC")
		break
	case ggthx.GGTHX_EVENTCODE_CONNECTION_INTERRUPTED:
		log.Println("GGTHX_EVENTCODE_CONNECTION_INTERRUPTED")
		break
	case ggthx.GGTHX_EVENTCODE_CONNECTION_RESUMED:
		log.Println("GGTHX_EVENTCODE_CONNECTION_RESUMED")
		break
	}
	return true
}

func init() {

	// have to register everything with gob, maybe automate this?
	// have to inititalize all arrays
	gob.Register(color.RGBA{})
	gob.Register(Input{})
	callbacks.AdvanceFrame = advanceFrame
	callbacks.BeginGame = beginGame
	callbacks.FreeBuffer = freeBuffer
	callbacks.LoadGameState = loadGameState
	callbacks.LogGameState = logGameState
	callbacks.OnEvent = onEvent
	callbacks.SaveGameState = saveGameState
	var result ggthx.GGTHXErrorCode
	session, result = ggthx.StartSyncTest(&callbacks, "Test", 2, 4, 1)
	if result != ggthx.GGTHX_OK {
		log.Fatalf("There's an issue / \n ")
	}

	ggthx.SetDisconnectTimeout(session, 3000)
	ggthx.SetDisconnectNotifyStart(session, 1000)
	var handle ggthx.GGTHXPlayerHandle
	var player ggthx.GGTHXPlayer
	ggthx.AddPlayer(session, &player, &handle)
	ggthx.SetFrameDelay(session, handle, 2)
}
