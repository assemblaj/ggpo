package ggthx

import "time"

type FuncTimeType func() int64

func DefaultTime() int64 {
	return time.Now().UnixMilli()
}

type Poll struct {
	startTime   int
	handleCount int
	loopSinks   StaticBuffer[PollSinkCb]
}

type Poller interface {
	RegisterLoop(sink PollSink, cookie []byte)
	Pump(timeFunc ...FuncTimeType) bool
}

type PollSinkCb struct {
	sink   PollSink
	cookie []byte
}

type PollSink interface {
	OnLoopPoll(timeFunc FuncTimeType) bool
}

func NewPoll() Poll {
	return Poll{
		startTime:   0,
		handleCount: 0,
		loopSinks:   NewStaticBuffer[PollSinkCb](16),
	}
}

func (p *Poll) RegisterLoop(sink PollSink, cookie []byte) {
	err := p.loopSinks.PushBack(
		PollSinkCb{
			sink:   sink,
			cookie: cookie})
	if err != nil {
		panic(err)
	}
}

func (p *Poll) Pump(timeFunc ...FuncTimeType) bool {
	finished := false
	if p.startTime == 0 {
		p.startTime = int(time.Now().UnixMilli())
	}
	for i := 0; i < p.loopSinks.Size(); i++ {
		cb, err := p.loopSinks.Get(i)
		if err != nil {
			panic(err)
		}
		if len(timeFunc) != 0 {
			finished = !(cb.sink.OnLoopPoll(timeFunc[0]) || finished)
		} else {
			finished = !(cb.sink.OnLoopPoll(DefaultTime) || finished)
		}
	}
	return finished

}
