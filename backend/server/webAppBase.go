package server

import (
	"net/http"
	"os"
)

type webAppBase struct {
}

func (me *webAppBase) inputCheckAdminPassword(request *http.Request) bool {
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

func (me *webAppBase) inputAssertAdminPassword(request *http.Request) {
	if !me.inputCheckAdminPassword(request) {
		panic(webError{"Need admin password", http.StatusForbidden})
	}
}

func (me *webAppBase) inputCheckGoalManagerMode(request *http.Request) bool {
	var goalManagerModeCookie, _ = request.Cookie("goalManagerMode")
	return me.inputCheckAdminPassword(request) &&
		goalManagerModeCookie != nil && goalManagerModeCookie.Value == "1"
}

func (me *webAppBase) getAdminPassword() string {
	return os.Getenv("ADMIN_PASSWORD")
}

func (me *webAppBase) guardAdminFunction(function webFunction) webFunction {
	return func(response http.ResponseWriter, request *http.Request) {
		me.inputAssertAdminPassword(request)
		function(response, request)
	}
}
