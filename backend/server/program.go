package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/hinst/go-gophers"
)

type program struct {
	netAddress       string
	webFilesPath     string
	savedGoalsPath   string
	allowOrigin      string
	translatorApiUrl string

	database *database
}

var programTemplate = program{
	netAddress:       ":8080",
	webFilesPath:     "./www",
	savedGoalsPath:   "./saved-goals",
	allowOrigin:      "http://localhost:1234",
	translatorApiUrl: "http://localhost:11434",
}

func (me *program) create() *program {
	*me = programTemplate
	me.database = new(database)
	return me
}

func (me *program) runWeb() {
	me.database.init()

	var webApp = &webApp{
		allowOrigin: me.allowOrigin,
	}
	webApp.init(me.database)

	var fileServer = http.FileServer(http.Dir(me.webFilesPath))
	var filesPrefix = webApp.webPath + "/"
	http.Handle(filesPrefix, http.StripPrefix(filesPrefix, fileServer))

	log.Printf("Starting: netAddress=%v, webPath=%v, webFilesPath=%v",
		me.netAddress, webApp.webPath, me.webFilesPath)

	var terminatingContext, _ = signal.NotifyContext(context.Background(), os.Interrupt,
		syscall.SIGTERM, syscall.SIGINT)
	go func() {
		gophers.AssertError(http.ListenAndServe(me.netAddress, nil))
	}()
	<-terminatingContext.Done()

	me.database.close()
}

func (me *program) update() {
	me.database.init()
	me.updateTranslations()
	me.updateTitles()
	me.generateStatic(me.savedGoalsPath + "/static")
	me.uploadStatic()
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

func (me *program) migrate() {
	me.database.init()
	me.database.migrate()
}

func (me *program) generateStatic(folder string) {
	me.database.init()
	var webStatic = new(webStaticGoals)
	webStatic.init("http://localhost:8080", me.database, folder)
	webStatic.run()
}
