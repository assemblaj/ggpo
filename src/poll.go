package ggthx

import "time"

type Poll struct {
	startTime   int
	handleCount int
	loopSinks   StaticBuffer[PollSinkCb]
}

type Poller interface {
	RegisterLoop(sink PollSink, cookie []byte)
	Run()
	Pump(timeout int) bool
}

type PollSinkCb struct {
	sink   PollSink
	cookie []byte
}

type PollSink interface {
	OnLoopPoll([]byte) bool
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

func (p *Poll) Run() {
	for p.Pump(100) {
		continue
	}
}

func (p *Poll) Pump(timeout int) bool {
	finished := false
	if p.startTime == 0 {
		p.startTime = int(time.Now().UnixMilli())
	}
	for i := 0; i < p.loopSinks.Size(); i++ {
		cb, err := p.loopSinks.Get(i)
		if err != nil {
			panic(err)
		}
		finished = !(cb.sink.OnLoopPoll(cb.cookie) || finished)
	}
	return finished

}
