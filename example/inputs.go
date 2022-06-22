package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Input struct {
	PlayerNum int
	Key       []ebiten.Key
}

func (i Input) String() string {
	return fmt.Sprintf("Player: %d:: Input %s", i.PlayerNum, i.Key)
}

func NewInput() Input {
	return Input{}
}

func decodeInputs(buffer [][]byte) []Input {
	var inputs = make([]Input, len(buffer))
	for i, b := range buffer {
		var buf bytes.Buffer = *bytes.NewBuffer(b)
		dec := gob.NewDecoder(&buf)
		err := dec.Decode(&inputs[i])
		if err != nil {
			log.Printf("decode error: %s. Returning empty input\n", err)
			// hack
			inputs[i] = NewInput()
			//panic("eof")
		} else {
			log.Printf("inputs properly encoded: %s\n", inputs[i])
		}
	}
	return inputs
}

func encodeInputs(inputs Input) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(&inputs)
	if err != nil {
		log.Fatal("encode error ", err)
	}
	return buf.Bytes()
}
