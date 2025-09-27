package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/hinst/hinst-website/server/file_mode"
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
	me.database.init(me.savedGoalsPath)

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
	me.database.init(me.savedGoalsPath)
	me.updateTranslations()
	me.updateTitles()
	me.generateStatic()
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
	var gitBotName = requireEnvVar("GIT_BOT_NAME")
	var gitEmail = requireEnvVar("GIT_EMAIL")

	var staticGitPath = assertResultError(os.Getwd()) + "/static-git"
	assertError(os.MkdirAll(staticGitPath, file_mode.OS_USER_RW))
	var runner = &commandRunner{Dir: staticGitPath}
	runner.command("Git clone", true, "git", "clone", me.getStaticWebsiteGitUrl(), staticGitPath)

	runner.command("Git config", true, "git", "config", "core.fileMode", "false")
	runner.command("Git config", true, "git", "config", "core.autocrlf", "true")
	runner.command("Git config", true, "git", "config", "user.name", getQuotedString(gitBotName))
	runner.command("Git config", true, "git", "config", "user.email", getQuotedString(gitEmail))

	var preservedFiles = []string{".git", "posts"}
	for _, file := range assertResultError(os.ReadDir(staticGitPath)) {
		if !slices.Contains(preservedFiles, file.Name()) {
			assertError(os.RemoveAll(staticGitPath + "/" + file.Name()))
		}
	}

	var copyRunner = &commandRunner{Dir: assertResultError(os.Getwd())}
	for _, file := range assertResultError(os.ReadDir("./static")) {
		copyRunner.command("Copy", true, "cp", "-r", "./static/"+file.Name(), staticGitPath+"/")
	}

	runner.command("Git add", true, "git", "add", ".")
	runner.command("Git status", true, "git", "status")
	var commitOk = runner.command("Git commit", false, "git", "commit", "-m", "Automatic update")
	if !commitOk {
		log.Println("Nothing to commit")
		return
	}

	runner.command("Git push", true, "git", "push")
}

func (me *program) migrate() {
	me.database.init(me.savedGoalsPath)
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

func (me *program) generateStatic() {
	me.database.init(me.savedGoalsPath)
	var webApp = &webApp{webPath: "/"}
	webApp.init(me.database)
	go func() {
		assertError(http.ListenAndServe(me.netAddress, nil))
	}()
	time.Sleep(1000 * time.Millisecond)

	var webStatic = new(webStaticGoals)
	webStatic.init("http://localhost:8080", me.database)
	webStatic.run()
}

func (me *program) getStaticWebsiteGitUrl() string {
	return fmt.Sprintf("https://%v@github.com/hinst/hinst.github.io.git", requireEnvVar("GIT_TOKEN"))
}
