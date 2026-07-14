package views

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type HeadingInfo struct {
	Text string
	ID   string
}

var markdownRenderer = goldmark.New(
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
		parser.WithAttribute(),
	),
)

func RenderMarkdown(source string) (templ.Component, []HeadingInfo, error) {
	input := []byte(source)
	document := markdownRenderer.Parser().Parse(text.NewReader(input))
	headings := make([]HeadingInfo, 0)

	if err := ast.Walk(document, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if image, ok := node.(*ast.Image); ok {
			image.SetAttributeString("loading", "lazy")
			image.SetAttributeString("decoding", "async")
			return ast.WalkContinue, nil
		}

		heading, ok := node.(*ast.Heading)
		if !ok || heading.Level != 2 {
			return ast.WalkContinue, nil
		}
		value, ok := heading.AttributeString("id")
		if !ok {
			return ast.WalkContinue, nil
		}
		id, ok := value.([]byte)
		if !ok || len(id) == 0 {
			return ast.WalkContinue, nil
		}

		heading.SetAttributeString(
			"data-on-intersect__full",
			fmt.Appendf(nil, "$activeHeading=%q", id),
		)
		headings = append(headings, HeadingInfo{Text: headingText(heading, input), ID: string(id)})
		return ast.WalkContinue, nil
	}); err != nil {
		return nil, nil, fmt.Errorf("extract markdown headings: %w", err)
	}

	var output bytes.Buffer
	if err := markdownRenderer.Renderer().Render(&output, input, document); err != nil {
		return nil, nil, fmt.Errorf("render markdown: %w", err)
	}

	return templ.Raw(output.String()), headings, nil
}

func headingText(heading *ast.Heading, source []byte) string {
	var value strings.Builder
	_ = ast.Walk(heading, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch node := node.(type) {
		case *ast.Text:
			value.Write(node.Value(source))
			if node.SoftLineBreak() {
				value.WriteByte(' ')
			}
		case *ast.String:
			value.Write(node.Value)
		case *ast.AutoLink:
			value.Write(node.Label(source))
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})
	return strings.TrimSpace(value.String())
}
