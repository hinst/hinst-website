package server

import (
	"flag"
	"log"
	"os"

	"github.com/hinst/go-gophers"
	"github.com/joho/godotenv"
)

func Main() {
	if gophers.CheckFileExists(".env") {
		gophers.AssertError(godotenv.Load())
	}
	var modePtr = flag.String("mode", "web", "")
	var wwwPtr = flag.String("www", programTemplate.webFilesPath, "")
	var translatorApiPtr = flag.String("translatorApi", programTemplate.translatorApiUrl, "")
	var backupDirectoryPtr = flag.String("backup-directory", programTemplate.savedGoalsPath+"/backup", "")
	flag.Parse()

	switch *modePtr {
	case "web":
		var theProgram = new(program).create()
		theProgram.webFilesPath = *wwwPtr
		if envVal := os.Getenv("ALLOW_ORIGIN"); envVal != "" {
			theProgram.allowOrigin = envVal
		}
		theProgram.runWeb()
	case "update":
		var theProgram = new(program).create()
		theProgram.translatorApiUrl = *translatorApiPtr
		theProgram.update()
	case "migrate":
		var theProgram = new(program).create()
		theProgram.migrate()
	case "generateStatic":
		var theProgram = new(program).create()
		theProgram.generateStatic("static")
	case "backup":
		var theProgram = new(program).create()
		theProgram.backup(*backupDirectoryPtr)
	default:
		log.Fatalf("Unknown mode: %v", *modePtr)
	}
}
