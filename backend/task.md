# Swagger OpenAPI

The of the task goal is to generate complete Swagger OpenAPI specification for this backend server
written in Golang.

Right now we already generate TypeScript structures using Tygo. See file: `tygo.yaml`
However, Tygo should be removed because the new full OpenAPI generation approach should replace Tygo completely in this project.

Right now the HTTP receivers are defined using standard Go library. See file `server\webAppGoals.go`
