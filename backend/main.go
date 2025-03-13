package main

import (
	"flag"
	"log"
)

func main() {
	var modePointer = flag.String("mode", "web", "")
	flag.Parse()
	switch *modePointer {
	case "web":
		new(program).init().runWeb()
	case "translate":
		new(program).init().runTranslate()
	default:
		log.Fatalf("Unknown mode: %v", *modePointer)
	}
}
