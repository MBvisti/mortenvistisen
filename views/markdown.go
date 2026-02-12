package views

import (
	"bytes"
	"html/template"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var markdownRenderer = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
)

func markdownToHTML(content string) string {
	var out bytes.Buffer
	if err := markdownRenderer.Convert([]byte(content), &out); err != nil {
		return template.HTMLEscapeString(content)
	}

	return out.String()
}
