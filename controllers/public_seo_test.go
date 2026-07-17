package controllers

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	"mortenvistisen/models"
	"mortenvistisen/router/routes"
)

func TestParsePublicPage(t *testing.T) {
	for _, test := range []struct {
		value string
		page  int64
		ok    bool
	}{
		{value: "", page: 1, ok: true},
		{value: "1", page: 1, ok: true},
		{value: "2", page: 2, ok: true},
		{value: "0", page: 0, ok: false},
		{value: "03", page: 3, ok: false},
		{value: "invalid", page: 0, ok: false},
	} {
		page, ok := parsePublicPage(test.value)
		if page != test.page || ok != test.ok {
			t.Fatalf(
				"parsePublicPage(%q) = (%d, %t), want (%d, %t)",
				test.value,
				page,
				ok,
				test.page,
				test.ok,
			)
		}
	}
	if publicPageExists(2, 1) || !publicPageExists(1, 0) {
		t.Fatal("public page range validation failed")
	}
}

func TestCreateSitemapIncludesPublishedContent(t *testing.T) {
	updatedAt := time.Date(2026, time.July, 14, 12, 0, 0, 0, time.UTC)
	sitemap, err := createSitemap(
		[]models.ArticleEntity{{Slug: "article", UpdatedAt: updatedAt}},
		[]models.NewsletterEntity{
			{Slug: sql.NullString{String: "newsletter", Valid: true}, UpdatedAt: updatedAt},
		},
		[]models.ProjectEntity{{Slug: "project", UpdatedAt: updatedAt}},
	)
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{
		routes.HomePage.URL(),
		routes.AboutPage.URL(),
		routes.PostIndex.URL(),
		routes.NewsletterIndex.URL(),
		routes.ProjectIndex.URL(),
		routes.PostShow.URL("article"),
		routes.NewsletterShow.URL("newsletter"),
		routes.ProjectShow.URL("project"),
	} {
		if !strings.Contains(sitemap, path) {
			t.Fatalf("sitemap is missing %q: %s", path, sitemap)
		}
	}
	if !strings.Contains(sitemap, updatedAt.Format(time.RFC3339)) {
		t.Fatalf("sitemap is missing lastmod: %s", sitemap)
	}
}
