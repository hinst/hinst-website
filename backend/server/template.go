package server

import (
	"bytes"
	"html/template"

	"github.com/hinst/go-common"
)

var templateFunctions = template.FuncMap{
	"dict": func(pairs ...any) map[string]any {
		result := make(map[string]any)
		for i := 0; i+1 < len(pairs); i += 2 {
			result[pairs[i].(string)] = pairs[i+1]
		}
		return result
	},
}

func executeTemplateFile(filePath string, data any) string {
	var text = common.ReadTextFile(filePath)
	var compiledTemplate = common.AssertResultError(template.New("page").Funcs(templateFunctions).Parse(text))
	var buffer = &bytes.Buffer{}
	common.AssertError(compiledTemplate.Execute(buffer, data))
	return buffer.String()
}
