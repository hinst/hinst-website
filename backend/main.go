package main

import (
	"runtime/debug"

	"github.com/hinst/hinst-website/server"
)

func main() {
	debug.SetGCPercent(25)
	server.Main()
}
