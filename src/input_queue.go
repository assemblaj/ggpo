package ggthx

import (
	"errors"
	"log"
)

const (
	InputQueueLength = 128
	DefaultInputSize = 4
)

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
	var err error
	inputs := make([]GameInput, InputQueueLength)
	for i, _ := range inputs {
		inputs[i], err = NewGameInput(-1, nil, inputSize)
		if err != nil {
			panic(err)
		}
	}
	prediction, err := NewGameInput(NullFrame, nil, inputSize)
	if err != nil {
		panic(err)
	}
	return InputQueue{
		id:                  id,
		firstFrame:          true,
		lastUserAddedFrame:  NullFrame,
		firstIncorrectFrame: NullFrame,
		lastFrameRequested:  NullFrame,
		lastAddedFrame:      NullFrame,
		prediction:          prediction,
		inputs:              inputs,
	}

}

func (i *InputQueue) LastConfirmedFrame() int {
	log.Printf("returning last confirmed frame %d.\n", i.lastUserAddedFrame)
	return i.lastAddedFrame
}

func (i *InputQueue) FirstIncorrectFrame() int {
	return i.firstIncorrectFrame
}

func (i *InputQueue) DiscardConfirmedFrames(frame int) error {
	if frame < 0 {
		return errors.New("ggthx: InputQueue discardConfirmedFrames: frames <= 0")
	}

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
		if offset < 0 {
			return errors.New("ggthx: InputQueue discardConfirmedFrames: offet < 0")
		}

		i.tail = (i.tail + offset) % InputQueueLength
		i.length -= offset
	}

	log.Printf("after discarding, new tail is %d (frame:%d).\n", i.tail, i.inputs[i.tail].Frame)
	if i.length < 0 {
		return errors.New("ggthx: InputQueue discardConfirmedFrames: i.length < 0 ")
	}

	return nil
}

func (i *InputQueue) ResetPrediction(frame int) error {
	if !(i.firstIncorrectFrame == NullFrame || frame <= i.firstIncorrectFrame) {
		return errors.New("ggthx: InputQueue ResetPrediction: i.firstIncorrentFrame != NullFrame && frame > i.firstIncorrectFrame")
	}

	log.Printf("resetting all prediction errors back to frame %d.\n", frame)

	i.prediction.Frame = NullFrame
	i.firstIncorrectFrame = NullFrame
	i.lastFrameRequested = NullFrame
	return nil
}

func (i *InputQueue) GetConfirmedInput(requestedFrame int, input *GameInput) (bool, error) {
	if !(i.firstIncorrectFrame == NullFrame || requestedFrame < i.firstIncorrectFrame) {
		return false, errors.New("ggthx: InputQueue GetConfirmedInput : i.firstIncorrectFrame != NullFrame && requestedFrame >")
	}

	offset := requestedFrame % InputQueueLength
	if i.inputs[offset].Frame != requestedFrame {
		return false, nil
	}
	*input = i.inputs[offset]
	return true, nil
}

func (i *InputQueue) GetInput(requestedFrame int, input *GameInput) (bool, error) {
	log.Printf("requesting input frame %d.\n", requestedFrame)
	if i.firstIncorrectFrame != NullFrame {
		return false, errors.New("ggthx: InputQueue GetInput : i.firstIncorrectFrame != NullFrame")
	}

	i.lastFrameRequested = requestedFrame

	if requestedFrame < i.inputs[i.tail].Frame {
		return false, errors.New("ggthx: InputQueue GetInput : requestedFrame < i.inputs[i.tail].Frame")
	}

	if i.prediction.Frame == NullFrame {
		offset := requestedFrame - i.inputs[i.tail].Frame

		if offset < i.length {
			offset = (offset + i.tail) % InputQueueLength
			if i.inputs[offset].Frame != requestedFrame {
				return false, errors.New("ggthx: InputQueue GetInput : i.inputs[offset].Frame != requestedFrame")
			}
			*input = i.inputs[offset]
			log.Printf("returning confirmed frame number %d.\n", input.Frame)
			return true, nil
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

	if i.prediction.Frame < 0 {
		return false, errors.New("ggthx: InputQueue GetInput : i.prediction.Frame < 0")
	}
	*input = i.prediction
	input.Frame = requestedFrame
	log.Printf("returning prediction frame number %d (%d).\n", input.Frame, i.prediction.Frame)

	return false, nil
}

func (i *InputQueue) AddInput(input *GameInput) error {
	var newFrame int
	var err error
	log.Printf("adding input frame number %d to queue.\n", input.Frame)

	if !(i.lastUserAddedFrame == NullFrame || input.Frame == i.lastUserAddedFrame+1) {
		return errors.New("ggthx : InputQueue AddInput : !(i.lastUserAddedFrame == NullFrame || input.Frame == i.lastUserAddedFrame+1)")
	}
	i.lastUserAddedFrame = input.Frame

	newFrame, err = i.AdvanceQueueHead(input.Frame)
	if err != nil {
		panic(err)
	}

	if newFrame != NullFrame {
		i.AddDelayedInputToQueue(input, newFrame)
	}

	input.Frame = newFrame
	return nil
}

func (i *InputQueue) AddDelayedInputToQueue(input *GameInput, frameNumber int) error {
	log.Printf("adding delayed input frame number %d to queue.\n", frameNumber)

	// Assert(input.Size == i.prediction.Size) No
	if !(i.lastAddedFrame == NullFrame || frameNumber == i.lastAddedFrame+1) {
		return errors.New("ggthx: InputQueue AddDelayedInputToQueue : i.lastAddedFrame != NullFrame && frameNumber != i.lastAddedFrame+1")
	}
	if !(frameNumber == 0 || i.inputs[previousFrame(i.head)].Frame == frameNumber-1) {
		return errors.New("ggthx: InputQueue AddDelayedInputToQueue : frameNumber != 0 && i.inputs[previousFrame(i.head)].Frame == frameNumber-1")
	}
	/*
	 *	Add the frame to the back of the queue
	 */
	i.inputs[i.head] = *input
	i.inputs[i.head].Frame = frameNumber
	i.head = (i.head + 1) % InputQueueLength
	i.length++
	i.firstFrame = false

	i.lastAddedFrame = frameNumber
	if i.prediction.Frame != NullFrame {
		if frameNumber != i.prediction.Frame {
			return errors.New("ggthx: InputQueue AddDelayedInputToQueue : frameNumber != i.prediction.Frame")
		}
		equal, err := i.prediction.Equal(input, true)
		if err != nil {
			panic(err)
		}
		if i.firstIncorrectFrame == NullFrame && !equal {
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

	if i.length > InputQueueLength {
		return errors.New("ggthx: InputQueue AddDelayedInputToQueue : i.length > InputQueueLength")
	}
	return nil
}

func (i *InputQueue) AdvanceQueueHead(frame int) (int, error) {
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
		return NullFrame, nil
	}

	for expectedFrame < frame {
		log.Printf("Adding padding frame %d to account for change in frame delay.\n",
			expectedFrame)
		lastFrame := i.inputs[previousFrame(i.head)]
		err := i.AddDelayedInputToQueue(&lastFrame, expectedFrame)
		if err != nil {
			panic(err)
		}
		expectedFrame++
	}

	if !(frame == 0 || frame == i.inputs[previousFrame(i.head)].Frame+1) {
		return 0, errors.New("ggthx: InputQueue AdvanceQueueHead : frame != 0 && frame != i.inputs[previousFrame(i.head)].Frame+")
	}
	return frame, nil
}

func previousFrame(offset int) int {
	if offset == 0 {
		return InputQueueLength - 1
	} else {
		return offset - 1
	}
}

func (i *InputQueue) SetFrameDelay(delay int) {
	i.frameDelay = delay
}

func (i *InputQueue) Length() int {
	return i.length
}
