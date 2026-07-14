package views

import (
	"context"
	"mortenvistisen/models"
	"strings"
	"testing"
)

func TestPublicPostTagsKeepsHiddenTagsInDOM(t *testing.T) {
	tags := []models.TagEntity{
		{Title: "one"},
		{Title: "two"},
		{Title: "three"},
		{Title: "four"},
		{Title: "five"},
	}

	var output strings.Builder
	if err := publicPostTags(tags).Render(context.Background(), &output); err != nil {
		t.Fatal(err)
	}
	for _, tag := range tags {
		if !strings.Contains(output.String(), tag.Title) {
			t.Fatalf("tag %q is missing from the DOM", tag.Title)
		}
	}
	if !strings.Contains(output.String(), "<details") || !strings.Contains(output.String(), "Show 2 more") {
		t.Fatalf("expected native show-more disclosure: %s", output.String())
	}
}

func TestRenderMarkdownExtractsHeadings(t *testing.T) {
	content, headings, err := RenderMarkdown("## First *heading*\n\nBody.\n\n![Diagram](https://example.com/diagram.png)\n\n## First heading\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(headings) != 2 {
		t.Fatalf("expected 2 headings, got %d", len(headings))
	}
	if headings[0].Text != "First heading" || headings[0].ID != "first-heading" {
		t.Fatalf("unexpected first heading: %#v", headings[0])
	}
	if headings[1].ID != "first-heading-1" {
		t.Fatalf("expected unique second heading ID, got %q", headings[1].ID)
	}

	var output strings.Builder
	if err := content.Render(context.Background(), &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), `id="first-heading" data-on-intersect__full=`) {
		t.Fatalf("rendered heading is missing its ID or intersection signal: %s", output.String())
	}
	if !strings.Contains(output.String(), `loading="lazy" decoding="async"`) {
		t.Fatalf("rendered image is missing lazy-loading attributes: %s", output.String())
	}
}
