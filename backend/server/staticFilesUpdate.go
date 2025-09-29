package server

import (
	"fmt"
	"log"
	"os"
	"slices"
)

type staticFilesUpdate struct {
	savedGoalsPath string
}

func (me *staticFilesUpdate) run() {
	var staticGitPath = me.savedGoalsPath + "/static-git"
	var runner = new(commandRunner)
	if !checkDirectoryExists(staticGitPath) {
		runner.command("Git clone", true, "git", "clone", me.getStaticWebsiteGitUrl(), staticGitPath)
	}
	me.flushFiles(staticGitPath)
	me.buildSitemap()

	runner.Dir = staticGitPath
	runner.command("Git fetch", true, "git", "fetch")
	runner.command("Git config", true, "git", "config", "core.fileMode", "false")
	runner.command("Git config", true, "git", "config", "core.autocrlf", "true")
	runner.command("Git config", true, "git", "config", "user.name", getQuotedString(me.getBotName()))
	runner.command("Git config", true, "git", "config", "user.email", getQuotedString(me.getEmail()))

	runner.command("Git add", true, "git", "add", ".")
	runner.command("Git status", true, "git", "status")
	var commitOk = runner.command("Git commit", false, "git", "commit", "-m", "Automatic update")
	if !commitOk {
		log.Println("Nothing to commit")
		return
	}

	if false {
		runner.command("Git push", true, "git", "push")
	}
}

// Copy old files from Git repository
// Copy new files into Git repository
func (me *staticFilesUpdate) flushFiles(staticGitPath string) {
	assertError(os.RemoveAll(me.savedGoalsPath + "/static-old"))
	var preservedFiles = []string{".git", "posts"}
	for _, file := range assertResultError(os.ReadDir(staticGitPath)) {
		if !slices.Contains(preservedFiles, file.Name()) {
			var filePath = staticGitPath + "/" + file.Name()
			var oldFilePath = me.savedGoalsPath + "/static-old/" + file.Name()
			if file.IsDir() {
				os.CopyFS(oldFilePath, os.DirFS(filePath))
			} else {
				copyFile(oldFilePath, filePath)
			}
			assertError(os.RemoveAll(filePath))
		}
	}
	assertError(os.CopyFS(staticGitPath, os.DirFS(me.savedGoalsPath+"/static")))
}

func (me *staticFilesUpdate) buildSitemap() {
	var builder = siteMapBuilder{
		newFilesPath: me.savedGoalsPath + "/static",
		oldFilesPath: me.savedGoalsPath + "/static-old",
	}
	builder.run()
}

func (me *staticFilesUpdate) getStaticWebsiteGitUrl() string {
	return fmt.Sprintf("https://%v@github.com/hinst/hinst.github.io.git", requireEnvVar("GIT_TOKEN"))
}

func (me *staticFilesUpdate) getBotName() string {
	return requireEnvVar("GIT_BOT_NAME")
}

func (me *staticFilesUpdate) getEmail() string {
	return requireEnvVar("GIT_EMAIL")
}
