package html_util

import (
	"bytes"
	"strings"

	"github.com/hinst/go-gophers"
	"golang.org/x/net/html"
)

func Walk(root *html.Node, callback func(*html.Node)) {
	callback(root)
	for node := range root.Descendants() {
		callback(node)
	}
}

func FindElement(root *html.Node, predicate func(*html.Node) bool) (result *html.Node) {
	Walk(root, func(node *html.Node) {
		if result == nil && node.Type == html.ElementNode && predicate(node) {
			result = node
		}
	})
	return
}

func Attr(node *html.Node, key string) (result *html.Attribute) {
	for i := range node.Attr {
		if node.Attr[i].Key == key {
			result = &node.Attr[i]
			return
		}
	}
	return
}

func AttrValue(node *html.Node, key string) (result string) {
	if attr := Attr(node, key); attr != nil {
		result = attr.Val
	}
	return
}

func NodeHasClass(node *html.Node, className string) (result bool) {
	for _, class := range strings.Fields(AttrValue(node, "class")) {
		if class == className {
			return true
		}
	}
	return false
}

func NodeText(node *html.Node) (result string) {
	if node.Type == html.TextNode {
		return node.Data
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		result += NodeText(child)
	}
	return
}

func InnerHtml(node *html.Node) (result string) {
	var buffer bytes.Buffer
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		html.Render(&buffer, child)
	}
	return buffer.String()
}

func ParseHtmlFragment(htmlText string) (result *html.Node) {
	var buffer bytes.Buffer
	buffer.WriteString("<body>")
	buffer.WriteString(htmlText)
	buffer.WriteString("</body>")
	var document = gophers.AssertResultError(html.Parse(&buffer))
	var bodyNode = FindElement(document, func(node *html.Node) bool { return node.Data == "body" })
	if bodyNode == nil {
		panic("Cannot find <body>")
	}
	return bodyNode
}
