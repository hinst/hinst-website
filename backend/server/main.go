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
	var backupDirectoryPtr = flag.String("backup-directory", programTemplate.savedGoalsPath+"/backup", "")
	flag.Parse()

	var theProgram = new(program).create()
	defer theProgram.close()
	switch *modePtr {
	case "web":
		theProgram.webFilesPath = *wwwPtr
		theProgram.runWeb()
	case "importSmartProgress":
		// Import blog posts from SmartProgress
		theProgram.importSmartProgress()
	case "update":
		// All-in-one update: Update translations, generate titles, generate static files, upload static files.
		theProgram.translatorApiUrl = gophers.ReadEnvVar("AI_URL", programTemplate.translatorApiUrl)
		theProgram.update()
	case "migrate":
		theProgram.migrate()
	case "generateStatic":
		// Generate static files, to be used for local testing
		theProgram.generateStatic("static")
	case "backup":
		theProgram.backup(*backupDirectoryPtr)
	case "generateSchema":
		// Save the OpenAPI schema to disk
		theProgram.generateSchema("schema.yaml")
	default:
		log.Fatalf("Unknown mode: %v", *modePtr)
	}
}
