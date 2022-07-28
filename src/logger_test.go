package ggthx_test

import (
	"testing"

	ggthx "github.com/assemblaj/ggthx/src"
)

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
		t.Errorf("SetLoggerOutput did not return an error when trying to be used when the logger is disabled.")
	}
	ggthx.DisableLogger()
}
