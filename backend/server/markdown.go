package server

import (
	html_to_markdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/hinst/go-gophers"
)

func convertMarkdownToHtml(text string) string {
	var textBytes = []byte(text)
	text = string(convertMarkdownBytesToHtml(textBytes))
	return text
}

func convertMarkdownBytesToHtml(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func convertHtmlToMarkdown(text string) string {
	return gophers.AssertResultError(html_to_markdown.ConvertString(text))
}
