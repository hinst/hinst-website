package server

import (
	"bytes"
	"html/template"
)

func executeTemplateFile(filePath string, data any) string {
	var text = readTextFile(filePath)
	var compiledTemplate = assertResultError(template.New("page").Parse(text))
	var buffer = &bytes.Buffer{}
	assertError(compiledTemplate.Execute(buffer, data))
	return buffer.String()
}
