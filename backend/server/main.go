package server

import (
	"flag"
	"log"

	"github.com/hinst/go-common"
	"github.com/joho/godotenv"
)

func Main() {
	if checkFileExists(".env") {
		common.AssertError(godotenv.Load())
	}
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
	case "update":
		var theProgram = new(program).create()
		theProgram.translatorApiUrl = *translatorApiPtr
		theProgram.update()
	case "migrate":
		var theProgram = new(program).create()
		theProgram.migrate()
	case "generatePrimeNumbers":
		var theProgram = new(program).create()
		theProgram.generatePrimeNumbers()
	case "generateStatic":
		var theProgram = new(program).create()
		theProgram.generateStatic("static")
	default:
		log.Fatalf("Unknown mode: %v", *modePtr)
	}
}
