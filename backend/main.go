package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	var modePointer = flag.String("mode", "web", "")
	var wwwPointer = flag.String("www", programTemplate.webFilesPath, "")
	var allowOriginPointer = flag.String("allowOrigin", programTemplate.allowOrigin, "")
	var translatorApiPointer = flag.String("translatorApi", programTemplate.translatorApiUrl, "")
	flag.Parse()

	switch *modePointer {
	case "web":
		var theProgram = new(program).create()
		theProgram.webFilesPath = *wwwPointer
		theProgram.allowOrigin = *allowOriginPointer
		theProgram.runWeb()
	case "translate":
		var theProgram = new(program).create()
		theProgram.translatorApiUrl = *translatorApiPointer
		theProgram.runTranslate()
	default:
		log.Fatalf("Unknown mode: %v", *modePointer)
	}
}
