package main

import (
	"image/color"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	argsWithoutProg := os.Args[1:]
	var err error
	currentPlayer, err = strconv.Atoi(argsWithoutProg[0])
	if err != nil {
		panic("Plase enter integer player")
	}

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

	start = int(time.Now().UnixMilli())
	now = start
	next = start

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
