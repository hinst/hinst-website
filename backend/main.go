package main

import (
	"flag"
	"log"
)

func main() {
	var modePointer = flag.String("mode", "web", "")
	var wwwPointer = flag.String("www", "./www", "")
	var allowOriginPointer = flag.String("allowOrigin", "http://localhost:1234", "")
	flag.Parse()

	switch *modePointer {
	case "web":
		var theProgram = new(program).init()
		theProgram.webFilesPath = *wwwPointer
		theProgram.allowOrigin = *allowOriginPointer
		theProgram.runWeb()
	case "translate":
		new(program).init().runTranslate()
	default:
		log.Fatalf("Unknown mode: %v", *modePointer)
	}
}
