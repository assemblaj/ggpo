package main

import (
	"log"
)

const INPUT_QUEUE_LENGTH int = 128
const DEFAULT_INPUT_SIZE int = 4

type InputQueue struct {
	id         int
	head       int
	tail       int
	length     int
	firstFrame bool

	lastUserAddedFrame  int
	lastAddedFrame      int
	firstIncorrectFrame int
	lastFrameRequested  int

	frameDelay int

	inputs     []GameInput
	prediction GameInput
}

func NewInputQueue(id int, inputSize int) InputQueue {
	inputs := make([]GameInput, INPUT_QUEUE_LENGTH)
	for _, v := range inputs {
		v.Size = inputSize
	}

	return InputQueue{
		id:                  id,
		firstFrame:          true,
		lastUserAddedFrame:  NullFrame,
		firstIncorrectFrame: NullFrame,
		lastFrameRequested:  NullFrame,
		lastAddedFrame:      NullFrame,
		prediction:          NewGameInput(NullFrame, nil, inputSize),
		inputs:              inputs,
	}

}

func (i InputQueue) GetLastConfirmedFrame() int {
	log.Printf("returning last confirmed frame %d.\n", i.lastUserAddedFrame)
	return i.lastAddedFrame
}

func (i InputQueue) GetFirstIncorrectFrame() int {
	return i.firstIncorrectFrame
}

func (i *InputQueue) DiscardConfirmedFrames(frame int) {
	Assert(frame > 0)
	if i.lastFrameRequested != NullFrame {
		frame = Min(frame, i.lastFrameRequested)
	}
	log.Printf("discarding confirmed frames up to %d (last_added:%d length:%d [head:%d tail:%d]).\n",
		frame, i.lastAddedFrame, i.length, i.head, i.tail)
	if frame >= i.lastAddedFrame {
		i.tail = i.head
	} else {
		offset := frame - i.inputs[i.tail].Frame + 1

		log.Printf("difference of %d frames.\n", offset)
		Assert(offset >= 0)

		i.tail = (i.tail + offset) % INPUT_QUEUE_LENGTH
		i.length -= offset
	}

	log.Printf("after discarding, new tail is %d (frame:%d).\n", i.tail, i.inputs[i.tail].Frame)
	Assert(i.length >= 0)
}

func (i *InputQueue) ResetPrediction(frame int) {
	Assert(i.firstIncorrectFrame == NullFrame || frame <= i.firstIncorrectFrame)

	log.Printf("resetting all prediction errors back to frame %d.\n", frame)

	i.prediction.Frame = NullFrame
	i.firstIncorrectFrame = NullFrame
	i.lastFrameRequested = NullFrame

}

func (i InputQueue) GetConfirmedInput(requestedFrame int, input *GameInput) bool {
	Assert(i.firstIncorrectFrame == NullFrame || requestedFrame < i.firstIncorrectFrame)
	offset := requestedFrame % INPUT_QUEUE_LENGTH
	if i.inputs[offset].Frame != requestedFrame {
		return false
	}
	*input = i.inputs[offset]
	return true
}

func (i *InputQueue) GetInput(requestedFrame int, input *GameInput) bool {
	log.Printf("requesting input frame %d.\n", requestedFrame)
	Assert(i.firstIncorrectFrame == NullFrame)

	i.lastFrameRequested = requestedFrame

	Assert(requestedFrame >= i.inputs[i.tail].Frame)

	if i.prediction.Frame == NullFrame {
		offset := requestedFrame - i.inputs[i.tail].Frame

		if offset < i.length {
			offset = (offset + i.tail) % INPUT_QUEUE_LENGTH
			Assert(i.inputs[offset].Frame == requestedFrame)
			*input = i.inputs[offset]
			log.Printf("returning confirmed frame number %d.\n", input.Frame)
			return true
		}

		if requestedFrame == 0 {
			log.Printf("basing new prediction frame from nothing, you're client wants frame 0.\n")
			i.prediction.Erase()
		} else if i.lastAddedFrame == NullFrame {
			log.Printf("basing new prediction frame from nothing, since we have no frames yet.\n")
			i.prediction.Erase()
		} else {
			log.Printf("basing new prediction frame from previously added frame (queue entry:%d, frame:%d).\n",
				previousFrame(i.head), i.inputs[previousFrame(i.head)].Frame)
			i.prediction = i.inputs[previousFrame(i.head)]
		}
		i.prediction.Frame++
	}

	Assert(i.prediction.Frame >= 0)

	*input = i.prediction
	input.Frame = requestedFrame
	log.Printf("returning prediction frame number %d (%d).\n", input.Frame, i.prediction.Frame)

	return false
}

func (i *InputQueue) AddInput(input *GameInput) {
	var newFrame int

	log.Printf("adding input frame number %d to queue.\n", input.Frame)

	Assert(i.lastUserAddedFrame == NullFrame || input.Frame == i.lastUserAddedFrame+1)
	i.lastUserAddedFrame = input.Frame

	newFrame = i.AdvanceQueueHead(input.Frame)
	if newFrame != NullFrame {
		i.AddDelayedInputToQueue(input, newFrame)
	}

	input.Frame = newFrame
}

func (i *InputQueue) AddDelayedInputToQueue(input *GameInput, frameNumber int) {
	log.Printf("adding delayed input frame number %d to queue.\n", frameNumber)

	Assert(input.Size == i.prediction.Size)
	Assert(i.lastAddedFrame == NullFrame || frameNumber == i.lastAddedFrame+1)
	Assert(frameNumber == 0 || i.inputs[previousFrame(i.head)].Frame == frameNumber-1)

	i.inputs[i.head] = *input
	i.inputs[i.head].Frame = frameNumber
	i.length++
	i.firstFrame = false

	i.lastAddedFrame = frameNumber

	if i.prediction.Frame != NullFrame {
		Assert(frameNumber == i.prediction.Frame)

		if i.firstIncorrectFrame == NullFrame && !i.prediction.Equal(input, true) {
			log.Printf("frame %d does not match prediction.  marking error.\n", frameNumber)
			i.firstIncorrectFrame = frameNumber
		}

		if i.prediction.Frame == i.lastFrameRequested && i.firstIncorrectFrame == NullFrame {
			log.Printf("prediction is correct!  dumping out of prediction mode.\n")
			i.prediction.Frame = NullFrame
		} else {
			i.prediction.Frame++
		}
	}

	Assert(i.length <= INPUT_QUEUE_LENGTH)
}

func (i *InputQueue) AdvanceQueueHead(frame int) int {
	log.Printf("advancing queue head to frame %d.\n", frame)

	var expectedFrame int
	if i.firstFrame {
		expectedFrame = 0
	} else {
		expectedFrame = i.inputs[previousFrame(i.head)].Frame + 1
	}

	frame += i.frameDelay
	if expectedFrame > frame {
		log.Printf("Dropping input frame %d (expected next frame to be %d).\n",
			frame, expectedFrame)
		return NullFrame
	}

	for expectedFrame < frame {
		log.Printf("Adding padding frame %d to account for change in frame delay.\n",
			expectedFrame)
		lastFrame := i.inputs[previousFrame(i.head)]
		i.AddDelayedInputToQueue(&lastFrame, expectedFrame)
		expectedFrame++
	}

	Assert(frame == 0 || frame == i.inputs[previousFrame(i.head)].Frame+1)
	return frame
}

func previousFrame(offset int) int {
	if offset == 0 {
		return INPUT_QUEUE_LENGTH - 1
	} else {
		return offset - 1
	}
}

func (i *InputQueue) SetFrameDelay(delay int) {
	i.frameDelay = delay
}
