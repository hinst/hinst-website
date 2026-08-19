# Refactoring

Move this code

```go
var _webApi huma.API

func webApi(path string) huma.API {
	var config = huma.DefaultConfig("Hinst-website API", "0.1.0")
	config.DocsPath = path + "/docs"
	if _webApi == nil {
		_webApi = humago.New(webRouter(), config)
	}
	return _webApi
}

var _webApiGroups map[string]*huma.Group = make(map[string]*huma.Group)

func webApiGrouped(path string) huma.API {
	if path == "" {
		return webApi(path)
	}
	var webApiGroup, contains = _webApiGroups[path]
	if !contains {
		webApiGroup = huma.NewGroup(webApi(path), path)
		_webApiGroups[path] = webApiGroup
	}
	return webApiGroup
}
```

from file `server\web.go` to file `server\webApp.go`

Both webApi and webApiGrouped should become fields of struct webApp.
There is no need for map of groups. Only one group shall exist within webApp.
