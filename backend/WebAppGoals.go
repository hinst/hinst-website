package main

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"golang.org/x/text/language"
)

type webAppGoals struct {
	savedGoalsPath        string
	goalDateStringMatcher *regexp.Regexp
}

const savedGoalHeaderFileName = "_header.json"
const publicPostsFileName = "public-posts.txt"
const cookieKeyAdminPassword = "adminPassword"

func (me *webAppGoals) init() []namedWebFunction {
	me.goalDateStringMatcher = regexp.MustCompile(`^\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d$`)
	return []namedWebFunction{
		{"/api/goals", me.getGoals},
		{"/api/goal", me.getGoal},
		{"/api/goalPosts", me.getGoalPosts},
		{"/api/goalPost", me.getGoalPost},
		{"/api/goalPost/images", me.getGoalPostImages},
	}
}

func (me *webAppGoals) getGoals(response http.ResponseWriter, request *http.Request) {
	var files = assertResultError(os.ReadDir(me.savedGoalsPath))
	var headers = make([]*goalHeader, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			var headerFilePath = filepath.Join(me.savedGoalsPath, file.Name(), savedGoalHeaderFileName)
			var header = readJsonFile(headerFilePath, &goalHeader{})
			headers = append(headers, header)
		}
	}
	response.Write(encodeJson(headers))
}

func (me *webAppGoals) getGoal(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalIdString(request.URL.Query().Get("id"))
	var headerFilePath = filepath.Join(me.savedGoalsPath, goalId, savedGoalHeaderFileName)
	var goal = readJsonFile(headerFilePath, &goalHeader{})
	setCacheAge(response, time.Minute)
	response.Write(encodeJson(goal))
}

func (me *webAppGoals) getGoalPosts(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalIdString(request.URL.Query().Get("id"))
	var allGoalPostsEnabled, _ = request.Cookie("allGoalPostsEnabled")
	var allEnabled = me.inputCheckAdminPassword(request) &&
		allGoalPostsEnabled != nil && allGoalPostsEnabled.Value == "1"
	var filePaths = me.getGoalPostFiles(goalId, allEnabled)
	var posts = readJsonFiles[smartPostHeader](filePaths, runtime.NumCPU())
	setCacheAge(response, time.Minute)
	response.Write(encodeJson(posts))
}

func (me *webAppGoals) getGoalPostFiles(goalId string, allEnabled bool) (filePaths []string) {
	var goalDirectory = filepath.Join(me.savedGoalsPath, goalId)
	var fileNames = getGoalFiles(goalDirectory)
	sort.Strings(fileNames)
	var allowedPosts = make(map[string]bool)
	if !allEnabled {
		allowedPosts = me.getAvailablePosts(goalId)
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

func (me *webAppGoals) getGoalPost(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalIdString(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var fileName = filepath.Join(goalId, postDateTime.Format(storedGoalFileTimeFormat)+".json")
	var filePath = filepath.Join(me.savedGoalsPath, fileName)
	var post = readJsonFile(filePath, &smartPostExtended{})

	var requestedLanguage = getWebLanguage(request)
	var translatedFilePath = translatorPresets.getTranslatedFilePath(filePath, requestedLanguage)
	if checkFileExists(translatedFilePath) {
		post.Msg = readTextFile(translatedFilePath)
		post.LanguageTag = requestedLanguage.String()
		post.LanguageName = getLanguageName(requestedLanguage)
		if requestedLanguage != language.Russian {
			post.IsAutoTranslated = true
		}
	} else {
		post.LanguageNamePending = getLanguageName(requestedLanguage)
	}

	post.Images = nil
	response.Write(encodeJson(post))
}

func (me *webAppGoals) getGoalPostImages(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalIdString(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var fileName = filepath.Join(goalId, postDateTime.Format(storedGoalFileTimeFormat)+".json")
	var filePath = filepath.Join(me.savedGoalsPath, fileName)
	var post = readJsonFile(filePath, &smartPostExtended{})
	setCacheAge(response, time.Minute)
	response.Write(encodeJson(post.Images))
}

func (me *webAppGoals) checkValidGoalIdString(goalId string) bool {
	return goalIdStringMatcher.MatchString(goalId)
}

func (me *webAppGoals) inputValidGoalIdString(goalId string) string {
	var createWebError = func() webError {
		return webError{"Need goal id. Received: " + goalId, http.StatusBadRequest}
	}
	assertCondition(me.checkValidGoalIdString(goalId), createWebError)
	return goalId
}

func (me *webAppGoals) inputValidPostDateTime(text string) time.Time {
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

func (me *webAppGoals) inputCheckAdminPassword(request *http.Request) bool {
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

func (me *webAppGoals) getAvailablePosts(goalId string) (dates map[string]bool) {
	dates = make(map[string]bool)
	var publicPostsFilePath = filepath.Join(me.savedGoalsPath, goalId, publicPostsFileName)
	if !checkFileExists(publicPostsFilePath) {
		return
	}
	var availablePostsText = readTextFile(publicPostsFilePath)
	var availablePosts = strings.Split(availablePostsText, "\n")
	for _, availablePost := range availablePosts {
		availablePost = strings.TrimSpace(availablePost)
		if len(availablePost) > 0 {
			dates[availablePost] = true
		}
	}
	return
}

func (me *webAppGoals) getAdminPassword() string {
	return os.Getenv("ADMIN_PASSWORD")
}
