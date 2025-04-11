package main

import (
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

type webAppGoalsBase struct {
	db                    *database
	savedGoalsPath        string
	goalDateStringMatcher *regexp.Regexp
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
	var unixEpochSeconds, parseIntError = strconv.ParseInt(text, 10, 64)
	var createWebError = func() webError {
		return webError{
			"Need valid postDateTime. Format: unix epoch seconds, number",
			http.StatusBadRequest,
		}
	}
	assertCondition(nil == parseIntError, createWebError)
	return time.Unix(unixEpochSeconds, 0)
}

func (me *webAppGoalsBase) inputCheckAdminPassword(request *http.Request) bool {
	var actualAdminPassword = me.getAdminPassword()
	if actualAdminPassword == "" {
		return false
	}
	var adminPassword, _ = request.Cookie("adminPassword")
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
