package server

import (
	"net/http"

	"github.com/hinst/go-gophers"
)

type webAppBase struct{}

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

func (me *webAppBase) getAdminPassword() string {
	return gophers.ReadEnvVar("ADMIN_PASSWORD", "")
}

func (me *webAppBase) guardAdminFunction(function gophers.WebFunction) gophers.WebFunction {
	return func(response http.ResponseWriter, request *http.Request) {
		me.inputAssertAdminPassword(request)
		function(response, request)
	}
}
