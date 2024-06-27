package posts

import (
	"bytes"
	"embed"

	"github.com/MBvisti/mortenvistisen/pkg/telemetry"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
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
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("gruvbox"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
					chromahtml.TabWidth(4),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)

	return PostManager{
		posts:           assets,
		markdownHandler: md,
	}
}

func CheckIfFileExist(name string) error {
	_, err := assets.ReadFile(name)
	if err != nil {
		return err
	}

	return nil
}

func GetAllFiles() ([]string, error) {
	entries, err := assets.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() != "post.go" {
			filenames = append(filenames, entry.Name())
		}
	}

	return filenames, nil
}

func (pm *PostManager) GetPost(name string) (string, error) {
	source, err := pm.posts.ReadFile(name)
	if err != nil {
		telemetry.Logger.Info("failed to read markdown file", "error", err)
		return "", err
	}

	return string(source), nil
}

func (pm *PostManager) Parse(name string) (string, error) {
	source, err := pm.posts.ReadFile(name)
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

	return htmlOutput.String(), nil
}
