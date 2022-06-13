package main

import "time"

type Poll struct {
	startTime   int
	handleCount int
	loopSinks   StaticBuffer[PollSinkCb]
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
	p.loopSinks.PushBack(
		PollSinkCb{
			sink:   sink,
			cookie: cookie})
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
		cb := p.loopSinks.Get(i)
		finished = !(cb.sink.OnLoopPoll(cb.cookie) || finished)
	}
	return finished

}
