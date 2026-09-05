package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/hinst/go-gophers"
	"github.com/rs/cors"
)

type program struct {
	webFilesPath     string
	savedGoalsPath   string
	translatorApiUrl string

	database *database
}

var programTemplate = program{
	webFilesPath:     "./www",
	savedGoalsPath:   "./saved-goals",
	translatorApiUrl: "http://localhost:11434",
}

func (me *program) create() *program {
	*me = programTemplate
	me.database = new(database)
	return me
}

func (me *program) runWeb() {
	me.database.init()

	var webApp = new(webApp)
	webApp.init(me.database)

	var fileServer = http.FileServer(http.Dir(me.webFilesPath))
	var filesPrefix = webApp.webPath + "/"
	webRouter().Handle(filesPrefix, http.StripPrefix(filesPrefix, fileServer))

	var webHandler = cors.New(cors.Options{
		AllowedOrigins: []string{gophers.ReadEnvVar("ALLOW_ORIGIN", "")},
	}).Handler(webRouter())

	log.Printf("Starting: netAddress=%v, webPath=%v, webFilesPath=%v",
		me.netAddress(), webApp.webPath, me.webFilesPath)

	var terminatingContext, _ = signal.NotifyContext(context.Background(), os.Interrupt,
		syscall.SIGTERM, syscall.SIGINT)

	go func() {
		gophers.AssertError(http.ListenAndServe(me.netAddress(), webHandler))
	}()
	<-terminatingContext.Done()

	me.database.close()
}

func (me *program) importSmartProgress(goalIds []string) {
	me.database.init()
	var importer = smartProgressImporter{goalIds: goalIds}
	importer.database = me.database
	importer.run()
	me.database.close()
}

func (me *program) update() {
	me.database.init()
	me.updateTranslations()
	me.updateTitles()
	me.generateStatic(me.savedGoalsPath + "/static")
	me.uploadStatic()
	me.updateSearchIndexingStatus()
}

func (me *program) updateTranslations() {
	var theTranslator translator
	if me.translatorApiUrl != "" {
		theTranslator.apiUrl = me.translatorApiUrl + "/v1/chat/completions"
	}
	theTranslator.db = me.database
	theTranslator.run()
}

func (me *program) updateTitles() {
	var titleGenerator = titleGeneratorPreset
	if me.translatorApiUrl != "" {
		titleGenerator.apiUrl = me.translatorApiUrl + "/v1/chat/completions"
	}
	titleGenerator.db = me.database
	titleGenerator.run()
}

func (me *program) uploadStatic() {
	var staticFilesUpdate = &staticFilesUpdate{
		db:             me.database,
		savedGoalsPath: me.savedGoalsPath,
	}
	staticFilesUpdate.run()
}

func (me *program) updateSearchIndexingStatus() {
	var updater = &searchIndexingUpdater{db: me.database}
	updater.run()
}

func (me *program) migrate() {
	me.database.init()
	me.database.migrate()
}

func (me *program) generateStatic(folder string) {
	me.database.init()
	var webStatic = new(webStaticGoals)
	webStatic.init(me.database, folder)
	webStatic.run()
}

func (me *program) backup(directory string) {
	me.database.init()
	me.database.backup(directory)
}

// Register all API routes and write the OpenAPI schema to the specified file
func (me *program) generateSchema(filename string) {
	var webApp = new(webApp)
	webApp.init(me.database)
	var yaml, err = webApp.webApi.OpenAPI().YAML()
	if err != nil {
		log.Fatalf("Failed to marshal OpenAPI to YAML: %v", err)
	}
	if err = os.WriteFile(filename, yaml, 0644); err != nil {
		log.Fatalf("Failed to write %v: %v", filename, err)
	}
	log.Printf("Successfully saved OpenAPI schema to %v", filename)
}

func (program) netAddress() string {
	return ":8080"
}
