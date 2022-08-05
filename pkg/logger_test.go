package temp_test

import (
	"testing"

	ggthx "github.com/assemblaj/ggthx/src"
)

type FakeWriter struct{}

func (f *FakeWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

func TestDisableLogger(t *testing.T) {
	ggthx.DisableLogger()
	err := ggthx.SetLoggerOutput(nil)
	if err == nil {
		t.Errorf("SetLoggerOutput did not return an error when trying to be used when the logger is disabled.")
	}
}
func TestEnabledLoggerNilInput(t *testing.T) {
	ggthx.EnableLogger()
	err := ggthx.SetLoggerOutput(nil)
	if err == nil {
		t.Errorf("SetLoggerOutput did not return an error when trying to be used with a nil io.Writer.")
	}
	ggthx.DisableLogger()
}

func TestSetLoggerOutputValid(t *testing.T) {
	ggthx.EnableLogger()
	err := ggthx.SetLoggerOutput(&FakeWriter{})
	if err != nil {
		t.Errorf("SetLoggerOutput returned an error when trying to be used with a valid writer.")
	}
	ggthx.DisableLogger()
}
