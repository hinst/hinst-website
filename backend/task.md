# Huma Rest migration

Migrating to new Golang framework.
Right now, only one API is under construction: get list of goals.
See `server\webAppGoals.go`
Nothing is finished yet, but right now we have to clarify ONE question:
why is the API available at path http://localhost:8080/api/goals, and not http://localhost:8080/hinst-website/api/goals ? I set prefix for it, see `webPath`
