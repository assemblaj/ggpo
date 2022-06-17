package main

import (
	"image/color"
	"log"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var currentPlayer int = 1

type Game struct {
	players []Player
}

type Player struct {
	x         float64
	y         float64
	color     color.Color
	playerNum int
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.players[currentPlayer-1].y--
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.players[currentPlayer-1].y++
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.players[currentPlayer-1].x--
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.players[currentPlayer-1].x++
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.players {
		ebitenutil.DrawRect(screen, p.x, p.y, 50, 50, p.color)
	}
	ebitenutil.DebugPrint(screen, "Hello,, World!")
}

func (g *Game) Layout(outsideWidth, insideWidth int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	argsWithoutProg := os.Args[1:]
	var err error
	currentPlayer, err = strconv.Atoi(argsWithoutProg[0])
	if err != nil {
		panic("Plase enter integer player")
	}

	player1 := Player{
		x:         50,
		y:         50,
		color:     color.RGBA{255, 0, 0, 255},
		playerNum: 1}
	player2 := Player{
		x:         150,
		y:         50,
		color:     color.RGBA{0, 0, 255, 255},
		playerNum: 2}
	game := &Game{
		players: []Player{player1, player2}}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
