package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	var modePtr = flag.String("mode", "web", "")
	var wwwPtr = flag.String("www", programTemplate.webFilesPath, "")
	var allowOriginPtr = flag.String("allowOrigin", programTemplate.allowOrigin, "")
	var translatorApiPtr = flag.String("translatorApi", programTemplate.translatorApiUrl, "")
	flag.Parse()

	switch *modePtr {
	case "web":
		var theProgram = new(program).create()
		theProgram.webFilesPath = *wwwPtr
		theProgram.allowOrigin = *allowOriginPtr
		theProgram.runWeb()
	case "translate":
		var theProgram = new(program).create()
		theProgram.translatorApiUrl = *translatorApiPtr
		theProgram.translate()
	case "migrate":
		var theProgram = new(program).create()
		theProgram.migrate()
	default:
		log.Fatalf("Unknown mode: %v", *modePtr)
	}
}
