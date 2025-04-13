package main

import (
	"encoding/base64"
	"net/http"
	"regexp"
	"time"
)

type webAppGoals struct {
	webAppGoalsBase
}

func (me *webAppGoals) init(db *database) []namedWebFunction {
	me.db = db
	me.goalDateStringMatcher = regexp.MustCompile(`^\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d$`)
	return []namedWebFunction{
		{"/api/goals", me.getGoals},
		{"/api/goal", me.getGoal},
		{"/api/goalPosts", me.getGoalPosts},
		{"/api/goalPost", me.getGoalPost},
		{"/api/goalPost/images", me.getGoalPostImages},
		{"/api/goalPost/setPublic", me.setGoalPostPublic},
	}
}

func (me *webAppGoals) getGoals(response http.ResponseWriter, request *http.Request) {
	var records = me.db.getGoals()
	response.Write(encodeJson(records))
}

func (me *webAppGoals) getGoal(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.URL.Query().Get("id"))
	var goal = me.db.getGoal(goalId)
	setCacheAge(response, time.Minute)
	response.Write(encodeJson(goal))
}

func (me *webAppGoals) getGoalPosts(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalIdString(request.URL.Query().Get("id"))
	var goalManagerMode = me.inputCheckGoalManagerMode(request)
	var posts = me.db.getGoalPosts(getIntFromString(goalId), goalManagerMode)
	response.Write(encodeJson(posts))
}

func (me *webAppGoals) getGoalPost(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var goalManagerMode = me.inputCheckGoalManagerMode(request)

	// var requestedLanguage = getWebLanguage(request)
	var goalPostRow = me.db.getGoalPost(goalId, postDateTime)
	if !goalPostRow.IsPublic && !goalManagerMode {
		panic(webError{"Need goal manager access level", http.StatusUnauthorized})
	}
	var goalPostObject goalPostObject
	goalPostObject.GoalId = goalPostRow.GoalId
	goalPostObject.DateTime = goalPostRow.DateTime
	goalPostObject.Text = goalPostRow.Text
	response.Write(encodeJson(goalPostObject))
}

func (me *webAppGoals) getGoalPostImages(response http.ResponseWriter, request *http.Request) {
	var goalId = me.inputValidGoalId(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var images = me.db.getGoalPostImages(goalId, postDateTime)
	var imageStrings []string
	for _, image := range images {
		var imageString = "data:" + image.contentType + ";base64," +
			base64.StdEncoding.EncodeToString(image.file)
		imageStrings = append(imageStrings, imageString)
	}
	setCacheAge(response, time.Minute)
	response.Write(encodeJson(imageStrings))
}

func (me *webAppGoals) setGoalPostPublic(response http.ResponseWriter, request *http.Request) {
	var isAdmin = me.inputCheckAdminPassword(request)
	if !isAdmin {
		panic(webError{"Need administrator access", http.StatusUnauthorized})
	}
	var goalId = me.inputValidGoalIdString(request.URL.Query().Get("goalId"))
	var postDateTime = me.inputValidPostDateTime(request.URL.Query().Get("postDateTime"))
	var isPublic = request.URL.Query().Get("isPublic") == "true"
	var row = goalPostRow{GoalId: getInt64FromString(goalId), DateTime: postDateTime, IsPublic: isPublic}
	me.db.setGoalPostPublic(&row)
}
