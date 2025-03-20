package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	var modePointer = flag.String("mode", "web", "")
	var wwwPointer = flag.String("www", "./www", "")
	var allowOriginPointer = flag.String("allowOrigin", "http://localhost:1234", "")
	var translatorApiPointer = flag.String("translatorApi", "http://localhost:1235", "")
	flag.Parse()

	switch *modePointer {
	case "web":
		var theProgram = new(program).init()
		theProgram.webFilesPath = *wwwPointer
		theProgram.allowOrigin = *allowOriginPointer
		theProgram.runWeb()
	case "translate":
		var theProgram = new(program).init()
		theProgram.translatorApiUrl = *translatorApiPointer
		theProgram.runTranslate()
	default:
		log.Fatalf("Unknown mode: %v", *modePointer)
	}
}
