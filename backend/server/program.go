package server

import (
	"log"
	"net/http"
)

type program struct {
	netAddress       string
	webFilesPath     string
	savedGoalsPath   string
	allowOrigin      string
	translatorApiUrl string
	database         *database
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
	me.database.init(me.savedGoalsPath)

	var webApp = &webApp{
		allowOrigin: me.allowOrigin,
	}
	webApp.init(me.database)

	var fileServer = http.FileServer(http.Dir(me.webFilesPath))
	var filesPrefix = webApp.webPath + "/"
	http.Handle(filesPrefix, http.StripPrefix(filesPrefix, fileServer))

	log.Printf("Starting: netAddress=%v, webPath=%v, filesPath=%v",
		me.netAddress, webApp.webPath, me.webFilesPath)
	assertError(http.ListenAndServe(me.netAddress, nil))
}

func (me *program) translate() {
	me.database.init(me.savedGoalsPath)
	var translator = translatorPresets
	if me.translatorApiUrl != "" {
		translator.apiUrl = me.translatorApiUrl + "/v1/chat/completions"
	}
	translator.savedGoalsPath = me.savedGoalsPath
	translator.db = me.database
	translator.run()
}

func (me *program) migrate() {
	me.database.init(me.savedGoalsPath)
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
