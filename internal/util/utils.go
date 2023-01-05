package util

import (
	"golang.org/x/exp/constraints"
	"io"
	"log"
)

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

var Log = DiscardLogger()

func DiscardLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}
