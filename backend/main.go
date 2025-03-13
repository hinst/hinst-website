package main

import (
	"flag"
	"log"
)

func main() {
	var modePointer = flag.String("mode", "web", "mode: web, translate")
	flag.Parse()
	switch *modePointer {
	case "web":
		new(program).init().runWeb()
	case "translate":
		new(program).init().runTranslate()
	case "migrate":
		new(program).init().runMigrate()
	default:
		log.Fatalf("Unknown mode: %v", *modePointer)
	}
}
