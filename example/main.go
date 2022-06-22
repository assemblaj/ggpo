package main

import (
	"flag"
	"image/color"
	"log"
	"os"
	"strconv"
	"time"

	ggthx "github.com/assemblaj/ggthx/src"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 5 {
		panic("Must enter <port> <num players> ('local' |IP adress) ('local' |IP adress) currentPlayer")
	}
	var port, numPlayers int
	var err error
	port, err = strconv.Atoi(argsWithoutProg[0])
	if err != nil {
		panic("Plase enter integer port")
	}

	numPlayers, err = strconv.Atoi(argsWithoutProg[1])
	if err != nil {
		panic("Please enter integer numPlayers")
	}

	ipAddress := []string{argsWithoutProg[2], argsWithoutProg[3]}

	currentPlayer, err = strconv.Atoi(argsWithoutProg[4])
	if err != nil {
		panic("Please enter integer currentPlayer")
	}

	players := make([]ggthx.GGTHXPlayer, numPlayers)
	for i := 0; i < numPlayers; i++ {
		if ipAddress[i] == "local" {
			players[i] = ggthx.NewLocalPlayer(20, i+1)
		} else {
			players[i] = ggthx.NewRemotePlayer(20, i+1, ipAddress[i], port)
		}
	}

	flag.Parse()

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

	Game_Init(port, numPlayers, players, 0)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
