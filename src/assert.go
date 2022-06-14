package ggthx

import "log"

func Assert(condition bool) {
	if !condition {
		log.Panic()
	}
}
