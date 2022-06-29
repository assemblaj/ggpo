package ggthx

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)

const GameInputMaxBytes int = 9
const GameInputMaxPlayers int = 2
const NullFrame int = -1

type GameInput struct {
	Frame  int
	Size   int
	Bits   []byte
	Inputs [][]byte
}

// Will come back to this if the needing the offset becomes a thing
func NewGameInput(frame int, bits []byte, size int, offset ...int) GameInput {

	Assert(size > 0)
	/* Not useful for our purposes
	if len(offset) == 0 {
		Assert(size <= GAMEINPUT_MAX_BYTES*GAMEINPUT_MAX_PLAYERS)
	} else {
		Assert(size <= GAMEINPUT_MAX_BYTES)
	}*/
	return GameInput{
		Frame: frame,
		Size:  size,
		Bits:  bits,
	}

}

func (g GameInput) IsNull() bool {
	return g.Frame == NullFrame
}

func (g GameInput) Value(i int) bool {
	return (g.Bits[i/8] & (1 << (i % 8))) != 0
}

func (g *GameInput) Set(i int) {
	g.Bits[i/8] |= (1 << (i % 8))
}

func (g *GameInput) Clear(i int) {
	g.Bits[i/8] &= ^(1 << (i % 8))
}

func (g *GameInput) Erase() {
	for i := 0; i < len(g.Bits); i++ {
		g.Bits[i] = 0
	}
}

func (g GameInput) Log(prefix string, showFrame bool) {
	log.Printf("%s%s", prefix, g)
}

func (g GameInput) String() string {
	Assert(g.Size > 0)
	retval := fmt.Sprintf("(frame:%d size:%d", g.Frame, g.Size)
	builder := strings.Builder{}
	for i := 0; i < len(g.Bits); i++ {
		builder.WriteByte(g.Bits[i])
	}
	builder.WriteString(")")
	return retval + builder.String()
}

func (g GameInput) Equal(other *GameInput, bitsonly bool) bool {
	if !bitsonly && g.Frame != other.Frame {
		log.Printf("frames don't match: %d, %d\n", g.Frame, other.Frame)
	}
	if g.Size != other.Size {
		log.Printf("sizes don't match: %d, %d\n", g.Size, other.Size)
	}
	if !bytes.Equal(g.Bits, other.Bits) {
		log.Printf("bits don't match\n")
	}
	Assert(g.Size > 0 && other.Size > 0)
	return (bitsonly || g.Frame == other.Frame) &&
		g.Size == other.Size &&
		bytes.Equal(g.Bits, other.Bits)
}
