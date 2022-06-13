package main

import "log"

func Assert(condition bool) {
	if !condition {
		log.Panic()
	}
}
