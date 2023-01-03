package main

type LoopTimer struct {
	lastAdvantage      float32
	usPergameLoop      int
	usExtraToWait      int
	framesToSpreadWait int
	waitCount          int
	timeWait           int
}

func NewLoopTimer(fps uint32, framesToSpread uint32) LoopTimer {
	return LoopTimer{
		framesToSpreadWait: int(framesToSpread),
		usPergameLoop:      1000000 / int(fps),
	}
}

func (lt *LoopTimer) OnGGPOTimeSyncEvent(framesAhead float32) {
	lt.lastAdvantage = (1000.0 * framesAhead / 60.0)
	lt.lastAdvantage = lt.lastAdvantage / 2
	lt.usExtraToWait = int(lt.lastAdvantage * 1000)
	if lt.usExtraToWait > 0 {
		lt.usExtraToWait = lt.usExtraToWait / lt.framesToSpreadWait
		lt.waitCount = lt.framesToSpreadWait
	}
}

func (lt *LoopTimer) usToWaitThisLoop() int {
	lt.timeWait = lt.usPergameLoop
	if lt.waitCount > 0 {
		lt.timeWait += lt.usExtraToWait
		lt.waitCount--
		if lt.waitCount < 0 {
			lt.usExtraToWait = 0
		}
	}
	return lt.timeWait
}
