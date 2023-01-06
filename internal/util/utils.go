package util

import (
	"io"
	"log"
	"math"

	"golang.org/x/exp/constraints"
)

// Log is the util.Log.Logger used when printing log messages.
var Log = DiscardLogger()

// DiscardLogger returns a util.Log.Logger which does not print anything.
func DiscardLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}

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

func MaxAbsFloat32(a, b float32) float32 {
	if math.Abs(float64(a)) > math.Abs(float64(b)) {
		return a
	} else {
		return b
	}
}
