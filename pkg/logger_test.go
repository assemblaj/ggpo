package ggpo_test

import (
	"testing"

	ggpo "github.com/assemblaj/GGPO-Go/pkg"
)

type FakeWriter struct{}

func (f *FakeWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

func TestDisableLogger(t *testing.T) {
	ggpo.DisableLogger()
	err := ggpo.SetLoggerOutput(nil)
	if err == nil {
		t.Errorf("SetLoggerOutput did not return an error when trying to be used when the logger is disabled.")
	}
}
func TestEnabledLoggerNilInput(t *testing.T) {
	ggpo.EnableLogger()
	err := ggpo.SetLoggerOutput(nil)
	if err == nil {
		t.Errorf("SetLoggerOutput did not return an error when trying to be used with a nil io.Writer.")
	}
	ggpo.DisableLogger()
}

func TestSetLoggerOutputValid(t *testing.T) {
	ggpo.EnableLogger()
	err := ggpo.SetLoggerOutput(&FakeWriter{})
	if err != nil {
		t.Errorf("SetLoggerOutput returned an error when trying to be used with a valid writer.")
	}
	ggpo.DisableLogger()
}
