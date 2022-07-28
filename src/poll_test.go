package ggthx_test

import (
	"testing"

	ggthx "github.com/assemblaj/ggthx/src"
)

type funcTimeType func() int64

type FakeSink struct {
	used bool
}

func NewFakeSink() FakeSink {
	return FakeSink{}
}

func (f *FakeSink) OnLoopPoll(timeFunc ggthx.FuncTimeType) bool {
	f.used = true
	return true
}

type FakeFalseSink struct {
}

func NewFakeFalseSink() FakeFalseSink {
	return FakeFalseSink{}
}

func (f FakeFalseSink) OnLoopPoll(timeFunc ggthx.FuncTimeType) bool {
	return false
}

func TestRegisterollPanic(t *testing.T) {
	poll := ggthx.NewPoll()
	maxSinks := 16
	sink := NewFakeSink()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic when attempting to add more than the max static buffer of sinks.")
		}
	}()
	for i := 0; i < maxSinks+1; i++ {
		poll.RegisterLoop(&sink, nil)
	}
}

func TestPollPumpFalse(t *testing.T) {
	poll := ggthx.NewPoll()
	sink := NewFakeFalseSink()
	poll.RegisterLoop(sink, nil)
	want := true
	got := poll.Pump(ggthx.DefaultTime)
	if want != got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}
func TestPollPumpIteration(t *testing.T) {
	poll := ggthx.NewPoll()
	sink := NewFakeSink()
	poll.RegisterLoop(&sink, nil)
	poll.Pump(ggthx.DefaultTime)
	want := true
	got := sink.used
	if want != got {
		t.Errorf("expected '%#v' but got '%#v'", want, got)
	}
}

func TestPollPumpIterationMultiple(t *testing.T) {
	poll := ggthx.NewPoll()
	maxSinks := 15
	sinks := make([]FakeSink, maxSinks)
	for i := 0; i < maxSinks; i++ {
		newSink := NewFakeSink()
		sinks[i] = newSink
		poll.RegisterLoop(&sinks[i], nil)
	}
	poll.Pump(ggthx.DefaultTime)
	for i := 0; i < maxSinks; i++ {
		want := true
		got := sinks[i].used
		if want != got {
			t.Errorf("expected '%#v' but got '%#v'", want, got)
		}
	}
}
