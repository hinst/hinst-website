package server

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/hinst/go-gophers"
)

type staticFilesUpdate struct {
	db             *database
	savedGoalsPath string
}

func (me *staticFilesUpdate) run() {
	var staticGitPath = me.savedGoalsPath + "/static-git"
	var runner = new(commandRunner)
	if !gophers.CheckDirectoryExists(staticGitPath) {
		runner.command("Git clone", true, "git", "clone", me.getStaticWebsiteGitUrl(), staticGitPath)
	}
	me.flushFiles(staticGitPath)
	me.buildSiteMap()

	runner.Dir = staticGitPath
	runner.command("Git pull", true, "git", "pull")
	runner.command("Git config", true, "git", "config", "core.fileMode", "false")
	runner.command("Git config", true, "git", "config", "core.autocrlf", "false")
	runner.command("Git config", true, "git", "config", "user.name", gophers.GetQuotedString(me.getBotName()))
	runner.command("Git config", true, "git", "config", "user.email", gophers.GetQuotedString(me.getEmail()))

	runner.command("Git add", true, "git", "add", ".")
	runner.command("Git status", true, "git", "status")
	var commitOk = runner.command("Git commit", false, "git", "commit", "-m", "Automatic update")
	if commitOk {
		runner.command("Git push", true, "git", "push")
	} else {
		log.Println("Nothing to commit")
	}

	me.submitSiteMap()
}

// Copy old files from Git repository
// Copy new files into Git repository
func (me *staticFilesUpdate) flushFiles(staticGitPath string) {
	gophers.AssertError(os.RemoveAll(me.savedGoalsPath + "/static-old"))
	for _, file := range gophers.AssertResultError(os.ReadDir(staticGitPath)) {
		if !me.checkPreservedFile(file.Name()) {
			var filePath = staticGitPath + "/" + file.Name()
			var oldFilePath = me.savedGoalsPath + "/static-old/" + file.Name()
			if file.IsDir() {
				gophers.AssertError(os.CopyFS(oldFilePath, os.DirFS(filePath)))
			} else {
				gophers.CopyFile(oldFilePath, filePath)
			}
			gophers.AssertError(os.RemoveAll(filePath))
		}
	}
	gophers.AssertError(os.CopyFS(staticGitPath, os.DirFS(me.savedGoalsPath+"/static")))
}

func (staticFilesUpdate) checkPreservedFile(fileName string) bool {
	var preservedFiles = []string{".git", "posts", "robots.txt", "dynamic"}
	return slices.Contains(preservedFiles, fileName) || strings.HasPrefix(fileName, "googled")
}

func (me *staticFilesUpdate) buildSiteMap() {
	var builder = siteMapBuilder{
		webPath:      me.getPublicUrl(),
		newFilesPath: me.savedGoalsPath + "/static",
		oldFilesPath: me.savedGoalsPath + "/static-old",
	}
	builder.run()
	gophers.CopyFile(me.savedGoalsPath+"/static-git/sitemap.xml", me.savedGoalsPath+"/static/sitemap.xml")
}

func (me *staticFilesUpdate) submitSiteMap() {
	var submitter = siteMapSubmitter{
		db:          me.db,
		siteMapPath: me.savedGoalsPath + "/static/sitemap.xml",
	}
	submitter.run()
}

func (staticFilesUpdate) getStaticWebsiteGitUrl() string {
	return fmt.Sprintf("https://%v@github.com/hinst/hinst.github.io.git", gophers.RequireEnvVar("GIT_TOKEN"))
}

func (staticFilesUpdate) getPublicUrl() string {
	return "https://hinst.github.io"
}

func (me *staticFilesUpdate) getBotName() string {
	return gophers.RequireEnvVar("GIT_BOT_NAME")
}

func (me *staticFilesUpdate) getEmail() string {
	return gophers.RequireEnvVar("GIT_EMAIL")
}
