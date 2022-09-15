package mocks

import (
	"crypto/sha1"
	"fmt"
)

type FakeGame struct {
	Players []FakePlayer
}

func NewFakeGame() FakeGame {
	players := make([]FakePlayer, 2)
	players[0].PlayerNum = 1
	players[1].PlayerNum = 2
	return FakeGame{
		Players: players,
	}
}
func (f *FakeGame) clone() (result *FakeGame) {
	result = &FakeGame{}
	*result = *f
	result.Players = make([]FakePlayer, len(f.Players))
	for i := range f.Players {
		result.Players[i] = f.Players[i].clone()
	}
	return result
}

func (f *FakeGame) Checksum() int {
	h := sha1.New()
	h.Write([]byte(f.String()))
	toSum := h.Sum(nil)
	sum := 0
	for _, v := range toSum {
		sum += int(v)
	}
	return sum
}

func (g *FakeGame) String() string {
	return fmt.Sprintf("%s : %s ", g.Players[0].String(), g.Players[1].String())
}

func (g *FakeGame) UpdateByInputs(inputs [][]byte) {
	for i, input := range inputs {
		if len(input) > 0 {
			g.Players[i].X += float64(input[0])
			g.Players[i].Y += float64(input[1])
		}
	}
}

type FakePlayer struct {
	X         float64
	Y         float64
	PlayerNum int
}

func (p *FakePlayer) clone() FakePlayer {
	result := FakePlayer{}
	result.X = p.X
	result.Y = p.Y
	result.PlayerNum = p.PlayerNum
	return result
}

func (p *FakePlayer) String() string {
	return fmt.Sprintf("Player %d: X:%f Y:%f", p.PlayerNum, p.X, p.Y)
}
