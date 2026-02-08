package server

import (
	"log"
	"net/http"
	"time"
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
	translatorApiUrl: "http://localhost:1235",
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
	assertError(http.ListenAndServe(me.netAddress, nil))
}

func (me *program) update() {
	me.database.init()
	me.updateTranslations()
	me.updateTitles()
	me.generateStatic(me.savedGoalsPath + "/static")
	me.uploadStatic()
}

func (me *program) updateTranslations() {
	var translator = translatorPreset
	if me.translatorApiUrl != "" {
		translator.apiUrl = me.translatorApiUrl + "/v1/chat/completions"
	}
	translator.db = me.database
	translator.run()
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

func (me *program) generatePrimeNumbers() {
	var primeNumbers = calculatePrimeNumbers(100_000)
	primeNumbers = primeNumbers[10_000:]
	var outputs []int
	for index, primeNumber := range primeNumbers {
		if (index % 10) == 0 {
			outputs = append(outputs, primeNumber)
		}
	}
	writeJsonFile(primeNumbersFileName, outputs)
}

func (me *program) generateStatic(folder string) {
	me.database.init()
	var webApp = &webApp{webPath: "/"}
	webApp.init(me.database)
	go func() {
		assertError(http.ListenAndServe(me.netAddress, nil))
	}()
	time.Sleep(1000 * time.Millisecond)

	var webStatic = new(webStaticGoals)
	webStatic.init("http://localhost:8080", me.database, folder)
	webStatic.run()
}
