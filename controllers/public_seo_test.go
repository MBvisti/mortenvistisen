package controllers

import (
	"context"
	"database/sql"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"mortenvistisen/config"
	"mortenvistisen/internal/hypermedia"
	"mortenvistisen/models"
	"mortenvistisen/router/routes"
	"mortenvistisen/views"

	"github.com/labstack/echo/v5"
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

func TestAssetsRSSReturnsAllPublishedPostsAndCachesResponse(t *testing.T) {
	publishedAt := time.Date(2026, time.July, 14, 12, 0, 0, 0, time.UTC)
	articles := []models.ArticleEntity{
		{
			Title:            "Newest & best",
			Slug:             "newest",
			Excerpt:          sql.NullString{String: "The newest post", Valid: true},
			FirstPublishedAt: sql.NullTime{Time: publishedAt, Valid: true},
			UpdatedAt:        publishedAt.Add(time.Hour),
		},
		{
			Title:            "Older post",
			Slug:             "older",
			Excerpt:          sql.NullString{String: "The older post", Valid: true},
			FirstPublishedAt: sql.NullTime{Time: publishedAt.Add(-time.Hour), Valid: true},
			UpdatedAt:        publishedAt.Add(-time.Hour),
		},
	}

	cache, err := NewCacheBuilder[string]().Build()
	if err != nil {
		t.Fatal(err)
	}
	loadCount := 0
	requestedLimit := -1
	assets := Assets{
		cache: cache,
		loadRSSArticles: func(_ context.Context, limit int) ([]models.ArticleEntity, error) {
			loadCount++
			requestedLimit = limit
			return articles, nil
		},
	}

	request := httptest.NewRequest(http.MethodGet, routes.RSS.URL(), nil)
	recorder := httptest.NewRecorder()
	etx := echo.New().NewContext(request, recorder)
	if err := assets.RSS(etx); err != nil {
		t.Fatal(err)
	}

	if recorder.Code != http.StatusOK {
		t.Fatalf("RSS status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Content-Type"); got != rssContentType {
		t.Fatalf("RSS Content-Type = %q, want %q", got, rssContentType)
	}
	if got := recorder.Header().Get("Cache-Control"); got != rssCacheControl {
		t.Fatalf("RSS Cache-Control = %q, want %q", got, rssCacheControl)
	}
	etag := recorder.Header().Get("ETag")
	if etag == "" {
		t.Fatal("RSS response is missing an ETag")
	}

	var feed RSSFeed
	if err := xml.Unmarshal(recorder.Body.Bytes(), &feed); err != nil {
		t.Fatalf("RSS response is invalid XML: %v", err)
	}
	if feed.Version != "2.0" {
		t.Fatalf("RSS version = %q, want 2.0", feed.Version)
	}
	if len(feed.Channel.Items) != len(articles) {
		t.Fatalf("RSS item count = %d, want %d", len(feed.Channel.Items), len(articles))
	}
	for i, article := range articles {
		item := feed.Channel.Items[i]
		wantURL := config.BaseURL + routes.PostShow.URL(article.Slug)
		if item.Title != article.Title || item.Link != wantURL || item.GUID.Value != wantURL {
			t.Fatalf("RSS item %d = %#v, want article %#v at %q", i, item, article, wantURL)
		}
		if item.Description != article.Excerpt.String {
			t.Fatalf("RSS item %d description = %q, want %q", i, item.Description, article.Excerpt.String)
		}
	}

	conditionalRequest := httptest.NewRequest(http.MethodGet, routes.RSS.URL(), nil)
	conditionalRequest.Header.Set("If-None-Match", etag)
	conditionalRecorder := httptest.NewRecorder()
	conditionalContext := echo.New().NewContext(conditionalRequest, conditionalRecorder)
	if err := assets.RSS(conditionalContext); err != nil {
		t.Fatal(err)
	}
	if conditionalRecorder.Code != http.StatusNotModified {
		t.Fatalf(
			"conditional RSS status = %d, want %d",
			conditionalRecorder.Code,
			http.StatusNotModified,
		)
	}
	if conditionalRecorder.Body.Len() != 0 {
		t.Fatalf("conditional RSS response body length = %d, want 0", conditionalRecorder.Body.Len())
	}
	if loadCount != 1 {
		t.Fatalf("RSS article loader called %d times, want 1", loadCount)
	}
	if requestedLimit != 0 {
		t.Fatalf("RSS article loader limit = %d, want 0 for all published posts", requestedLimit)
	}
}

func TestPublicPagesAdvertiseRSSFeed(t *testing.T) {
	html, err := hypermedia.RenderHTML(context.Background(), views.PageAbout())
	if err != nil {
		t.Fatal(err)
	}

	for _, attribute := range []string{
		`rel="alternate"`,
		`type="application/rss+xml"`,
		`href="` + routes.RSS.URL() + `"`,
	} {
		if !strings.Contains(html, attribute) {
			t.Fatalf("public page is missing RSS discovery attribute %s", attribute)
		}
	}
}
