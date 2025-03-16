package main

import (
	"flag"
	"log"
)

func main() {
	var modePointer = flag.String("mode", "web", "")
	var wwwPointer = flag.String("www", "./www", "")
	flag.Parse()
	switch *modePointer {
	case "web":
		var theProgram = new(program).init()
		theProgram.webFilesPath = *wwwPointer
		theProgram.runWeb()
	case "translate":
		new(program).init().runTranslate()
	default:
		log.Fatalf("Unknown mode: %v", *modePointer)
	}
}
