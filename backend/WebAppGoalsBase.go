package main

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

type webAppGoalsBase struct {
	db                    *database
	savedGoalsPath        string
	goalDateStringMatcher *regexp.Regexp
}

const publicPostsFileName = "public-posts.txt"

func (me *webAppGoalsBase) getGoalPostFiles(goalId string, allEnabled bool) (filePaths []string) {
	var goalDirectory = filepath.Join(me.savedGoalsPath, goalId)
	var fileNames = getGoalFiles(goalDirectory)
	sort.Strings(fileNames)
	var allowedPosts = make(map[string]bool)
	if !allEnabled {
		allowedPosts = me.db.getAvailablePosts(getIntFromString(goalId))
	}
	for _, fileName := range fileNames {
		var isAllowed = false
		if allEnabled {
			isAllowed = true
		} else {
			var fileNameBase = getFileNameWithoutExtension(fileName)
			var date = assertResultError(parseStoredGoalFileDate(fileNameBase))
			var dateText = date.Format(smartProgressTimeFormat)
			isAllowed = allowedPosts[dateText]
		}
		if !isAllowed {
			continue
		}
		var filePath = filepath.Join(me.savedGoalsPath, goalId, fileName)
		filePaths = append(filePaths, filePath)
	}
	return
}

func (me *webAppGoalsBase) checkValidGoalIdString(goalId string) bool {
	return goalIdStringMatcher.MatchString(goalId)
}

func (me *webAppGoalsBase) inputValidGoalIdString(goalId string) string {
	var createWebError = func() webError {
		return webError{"Need goal id. Received: " + goalId, http.StatusBadRequest}
	}
	assertCondition(me.checkValidGoalIdString(goalId), createWebError)
	return goalId
}

func (me *webAppGoalsBase) inputValidPostDateTime(text string) time.Time {
	var postDateTime, postDateTimeError = parseSmartProgressDate(text)
	var createWebError = func() webError {
		return webError{
			"Need valid postDateTime. Format: " + smartProgressTimeFormat,
			http.StatusBadRequest,
		}
	}
	assertCondition(nil == postDateTimeError, createWebError)
	return postDateTime
}

func (me *webAppGoalsBase) inputCheckAdminPassword(request *http.Request) bool {
	var actualAdminPassword = me.getAdminPassword()
	if actualAdminPassword == "" {
		return false
	}
	var adminPassword, _ = request.Cookie(cookieKeyAdminPassword)
	if adminPassword != nil {
		return adminPassword.Value == actualAdminPassword
	}
	return false
}

func (me *webAppGoalsBase) inputCheckGoalManagerMode(request *http.Request) bool {
	var goalManagerModeCookie, _ = request.Cookie("goalManagerMode")
	return me.inputCheckAdminPassword(request) &&
		goalManagerModeCookie != nil && goalManagerModeCookie.Value == "1"
}

func (me *webAppGoalsBase) getAdminPassword() string {
	return os.Getenv("ADMIN_PASSWORD")
}
