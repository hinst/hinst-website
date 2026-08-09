package server

import (
	"flag"
	"log"

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
		theProgram.runWeb()
	case "update":
		// All-in-one update: Update translations, generate titles, generate static files, upload static files.
		var theProgram = new(program).create()
		theProgram.translatorApiUrl = *translatorApiPtr
		theProgram.update()
	case "migrate":
		var theProgram = new(program).create()
		theProgram.migrate()
	case "generateStatic":
		// Generate static files, to be used for local testing
		var theProgram = new(program).create()
		theProgram.generateStatic("static")
	case "backup":
		var theProgram = new(program).create()
		theProgram.backup(*backupDirectoryPtr)
	default:
		log.Fatalf("Unknown mode: %v", *modePtr)
	}
}
