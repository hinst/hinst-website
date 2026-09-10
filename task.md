# Use array for localized titles

See these files in folder `C:\Dev\hinst-website\backend\server\db_objects`:
* `goalRow.go`
* `goalPostRow.go`

Currently the structures have fields named:
* `title`
* `titleEnglish`
* `titleGerman`

Also:
* `text`
* `textEnglish`
* `textGerman`

The goal is to refactor these fields: instead of three fields, use Postgres array type.

See list of supported languages in file: `C:\Dev\hinst-website\backend\server\base\language.go`:
```go
var SupportedLanguages = []language.Tag{language.Russian, language.English, language.German}
```


See also: `C:\Dev\hinst-website\backend\server\database.go`
```go
func (me *database) migrate() {
	// Put migration code here
}

