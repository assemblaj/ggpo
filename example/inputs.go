package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type Input struct {
	ButtonMap int
}

type InputBits int

func (i *InputBits) isButtonOn(button int) bool {
	return *i&(1<<button) > 0
}

func (i *InputBits) setButton(button int) {
	*i |= (1 << button)
}

func readI32(b []byte) int32 {
	if len(b) < 4 {
		return 0
	}
	return int32(b[0]) | int32(b[1])<<8 | int32(b[2])<<16 | int32(b[3])<<24
}

func writeI32(i32 int32) []byte {
	b := []byte{byte(i32), byte(i32 >> 8), byte(i32 >> 16), byte(i32 >> 24)}
	return b
}

func (i *Input) isButtonOn(button int) bool {
	return i.ButtonMap&(1<<button) > 0
}

func (i *Input) setButton(button int) {
	i.ButtonMap |= (1 << button)
}

func (i Input) String() string {
	return fmt.Sprintf("Input %d", i.ButtonMap)
}

func NewInput() Input {
	return Input{}
}

func encodeInputs(inputs InputBits) []byte {
	return writeI32(int32(inputs))
}

func decodeInputs(buffer [][]byte) []InputBits {
	var inputs = make([]InputBits, len(buffer))
	for i, b := range buffer {
		inputs[i] = InputBits(readI32(b))
	}
	return inputs
}

func decodeInputsGob(buffer [][]byte) []Input {
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
			log.Printf("inputs properly decoded: %s\n", inputs[i])
		}
	}
	return inputs
}

func encodeInputsGob(inputs Input) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(&inputs)
	if err != nil {
		log.Fatal("encode error ", err)
	}
	return buf.Bytes()
}
