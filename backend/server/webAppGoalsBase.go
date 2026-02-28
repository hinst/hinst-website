package server

import (
	"net/http"
	"strconv"
	"time"
)

type webAppGoalsBase struct {
	webAppBase
	db *database
}

func (me *webAppGoalsBase) inputValidGoalId(goalId string) int64 {
	var result, parseError = strconv.ParseInt(goalId, 10, 64)
	var createWebError = func() webError {
		return webError{"Need goal id. Received: " + goalId, http.StatusBadRequest}
	}
	AssertCondition(parseError == nil, createWebError)
	return result
}

func (me *webAppGoalsBase) inputValidPostDateTime(text string) time.Time {
	var unixEpochSeconds, parseIntError = strconv.ParseInt(text, 10, 64)
	var createWebError = func() webError {
		return webError{
			"Need valid post date time. Format: unix epoch seconds, number",
			http.StatusBadRequest,
		}
	}
	AssertCondition(nil == parseIntError, createWebError)
	return time.Unix(unixEpochSeconds, 0)
}
