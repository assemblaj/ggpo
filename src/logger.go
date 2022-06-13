package main

import "log"

func init() {
	log.SetFlags(log.Ldate | log.Ltime)
}
