package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/hinst/go-gophers"
	"github.com/rs/cors"
)

type program struct {
	webFilesPath     string
	savedGoalsPath   string
	translatorApiUrl string

	_db *database
}

var programTemplate = program{
	webFilesPath:     "./www",
	savedGoalsPath:   "./saved-goals",
	translatorApiUrl: "http://localhost:11434",
}

func (me *program) create() *program {
	*me = programTemplate
	return me
}

func (me *program) runWeb() {
	var webApp = new(webApp)
	webApp.init(me.db())

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
}

func (me *program) update() {
	me.importSmartProgress()
	me.updateTranslations()
	me.updateTitles()
	me.generateStatic(me.savedGoalsPath + "/static")
	me.uploadStatic()
	me.updateSearchIndexingStatus()
}

func (me *program) importSmartProgress() {
	var goalIds = strings.Split(gophers.RequireEnvVar("GOAL_IDS"), ",")
	var importer = smartProgressImporter{goalIds: goalIds}
	importer.database = me.db()
	importer.run()
}

func (me *program) updateTranslations() {
	var theTranslator translator
	if me.translatorApiUrl != "" {
		theTranslator.apiUrl = me.translatorApiUrl + "/v1/chat/completions"
	}
	theTranslator.db = me.db()
	theTranslator.run()
}

func (me *program) updateTitles() {
	var titleGenerator = titleGeneratorPreset
	if me.translatorApiUrl != "" {
		titleGenerator.apiUrl = me.translatorApiUrl + "/v1/chat/completions"
	}
	titleGenerator.db = me.db()
	titleGenerator.run()
}

func (me *program) uploadStatic() {
	var staticFilesUpdate = &staticFilesUpdate{
		db:             me.db(),
		savedGoalsPath: me.savedGoalsPath,
	}
	staticFilesUpdate.run()
}

func (me *program) updateSearchIndexingStatus() {
	var updater = &searchIndexingUpdater{db: me.db()}
	updater.run()
}

func (me *program) migrate() {
	me.db().migrate()
}

func (me *program) generateStatic(folder string) {
	var webStatic = new(webStaticGoals)
	webStatic.init(me.db(), folder)
	webStatic.run()
}

func (me *program) backup(directory string) {
	me.db().saveTables(directory)
}

// Register all API routes and write the OpenAPI schema to the specified file
func (me *program) generateSchema(filename string) {
	var webApp = new(webApp)
	webApp.init(nil)
	var yaml, err = webApp.webApi.OpenAPI().YAML()
	if err != nil {
		log.Fatalf("Failed to marshal OpenAPI to YAML: %v", err)
	}
	if err = os.WriteFile(filename, yaml, 0644); err != nil {
		log.Fatalf("Failed to write %v: %v", filename, err)
	}
	log.Printf("Successfully saved OpenAPI schema to %v", filename)
}

func (me *program) db() *database {
	if me._db == nil {
		me._db = new(database)
		me._db.init()
	}
	return me._db
}

func (me *program) close() {
	if me._db != nil {
		me._db.close()
	}
	me._db = nil
}

func (program) netAddress() string {
	return ":8080"
}
