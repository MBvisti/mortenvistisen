package views

import (
	"database/sql"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"mortenvistisen/models"
)

func TestPublicPageTitles(t *testing.T) {
	if got := (PostIndex{CurrentPage: 1}).pageTitle(); got != "Go and Software Engineering Articles | Morten Vistisen" {
		t.Fatalf("page one title = %q", got)
	}
	if got := (PostIndex{CurrentPage: 2}).pageTitle(); got != "Software Engineering Articles, Page 2 | Morten Vistisen" {
		t.Fatalf("page two title = %q", got)
	}
	for _, test := range []struct {
		slug string
		want string
	}{
		{slug: "andurel", want: "Andurel: Productive Go Web Framework | Morten Vistisen"},
		{slug: "deploy-crate", want: "Deploy Crate: Automated App Deployments | Morten Vistisen"},
	} {
		if got := projectPageTitle(models.ProjectEntity{Slug: test.slug}); got != test.want {
			t.Fatalf("project title for %q = %q, want %q", test.slug, got, test.want)
		}
	}
}

func TestPostStructuredDataUsesPersonAsAuthor(t *testing.T) {
	publishedAt := time.Date(2026, time.July, 14, 10, 0, 0, 0, time.UTC)
	data := HeadData{
		siteName:     "Morten Vistisen",
		Title:        "A practical Go article - Morten Vistisen",
		Description:  "An article about practical Go.",
		canonical:    "https://mortenvistisen.com",
		canonicalURL: "https://mortenvistisen.com/posts/practical-go",
		SchemaBuilders: []SchemaBuilder{PostShowSchema(models.ArticleEntity{
			Title:            "A practical Go article",
			Slug:             "practical-go",
			FirstPublishedAt: sql.NullTime{Time: publishedAt, Valid: true},
			UpdatedAt:        publishedAt,
		}, nil)},
	}

	encoded, err := json.Marshal(buildStructuredData(data))
	if err != nil {
		t.Fatal(err)
	}
	jsonLD := string(encoded)
	for _, expected := range []string{`"@type":"Person"`, `"@type":"BlogPosting"`, `"author":{"@id":"https://mortenvistisen.com/#person"}`} {
		if !strings.Contains(jsonLD, expected) {
			t.Fatalf("structured data does not contain %s: %s", expected, jsonLD)
		}
	}
	for _, unwanted := range []string{`"@type":"Organization"`, `"@type":"Service"`} {
		if strings.Contains(jsonLD, unwanted) {
			t.Fatalf("personal structured data contains %s: %s", unwanted, jsonLD)
		}
	}
}
