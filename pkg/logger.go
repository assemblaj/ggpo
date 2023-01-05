package ggpo

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var loggerEnabled = false

// logger disabled by default
func init() {
	DisableLogger()
}

func DisableLogger() {
	loggerEnabled = false
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
}

func EnableLogger() {
	loggerEnabled = true
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetOutput(os.Stdout)
}

func SetLoggerOutput(w io.Writer) error {
	if !loggerEnabled {
		return errors.New("logger must be enabled before using")
	}
	if w == nil {
		return errors.New("logger output cannot be nil")
	}
	log.SetOutput(w)
	return nil
}
