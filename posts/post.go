package posts

import (
	"bytes"
	"embed"

	"github.com/MBvisti/grafto/pkg/telemetry"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

//go:embed *.md
var assets embed.FS

type PostManager struct {
	posts           embed.FS
	markdownHandler goldmark.Markdown
}

func NewPostManager() PostManager {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	return PostManager{
		posts:           assets,
		markdownHandler: md,
	}
}

func (pm *PostManager) GetPost(name string) (string, error) {
	source, err := pm.posts.ReadFile(name + ".md")
	if err != nil {
		telemetry.Logger.Info("failed to read markdown file", "error", err)
		return "", err
	}

	return string(source), nil
}

func (pm *PostManager) Parse(name string) (string, error) {
	source, err := pm.posts.ReadFile(name + ".md")
	if err != nil {
		telemetry.Logger.Info("failed to read markdown file", "error", err)
		return "", err
	}

	// Parse Markdown content
	var htmlOutput bytes.Buffer
	if err := pm.markdownHandler.Convert(source, &htmlOutput); err != nil {
		telemetry.Logger.Info("failed to parse markdown file", "error", err)
		return "", err
	}

	return string(htmlOutput.Bytes()), nil
}
