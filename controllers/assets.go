package controllers

import (
	"context"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"

	"mortenvistisen/assets"
	"mortenvistisen/config"
	"mortenvistisen/internal/server"
	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
	"mortenvistisen/router/routes"

	"github.com/labstack/echo/v5"
)

const threeMonthsCache = "7776000"

type Assets struct {
	db    storage.Pool
	cache *Cache[string]
}

func NewAssets(db storage.Pool, cache *Cache[string]) Assets {
	return Assets{db, cache}
}

func (a Assets) enableCaching(etx *echo.Context, content []byte) *echo.Context {
	if config.Env == server.ProdEnvironment {
		//nolint:gosec //only needed for browser caching
		hash := md5.Sum(content)
		etag := fmt.Sprintf(`"%x-%x"`, hash, len(content))

		if match := etx.Request().Header.Get("If-None-Match"); match == etag {
			etx.Response().
				Header().
				Set("Cache-Control", fmt.Sprintf("public, max-age=%s, immutable", threeMonthsCache))
			etx.Response().
				Header().
				Set("ETag", etag)
			etx.NoContent(http.StatusNotModified)
			return etx
		}

		etx.Response().
			Header().
			Set("Cache-Control", fmt.Sprintf("public, max-age=%s, immutable", threeMonthsCache))
		etx.Response().
			Header().
			Set("Vary", "Accept-Encoding")
		etx.Response().
			Header().
			Set("ETag", etag)
	}

	return etx
}

func createRobotsTxt() (string, error) {
	robots := fmt.Sprintf(
		"User-agent: *\nDisallow: /admin/\nDisallow: /riverui\nDisallow: /login\nDisallow: /register\nDisallow: /confirm-email\nDisallow: /reset-password\nAllow: /\nSitemap: %s%s\n",
		config.BaseURL,
		routes.Sitemap.URL(),
	)
	return robots, nil
}

func (a Assets) Robots(etx *echo.Context) error {
	cacheKey := "assets:robots"

	robotsTxt, err := a.cache.Get(cacheKey, func() (string, error) {
		return createRobotsTxt()
	})
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to get robots.txt from cache",
			"error", err,
		)
		result, _ := createRobotsTxt()
		return etx.String(http.StatusOK, result)
	}

	return etx.String(http.StatusOK, robotsTxt)
}

func (a Assets) Sitemap(etx *echo.Context) error {
	cacheKey := "assets:sitemap"
	ctx := etx.Request().Context()

	sitemap, err := a.cache.Get(cacheKey, func() (string, error) {
		return a.createSitemap(ctx)
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to get sitemap from cache", "error", err)

		result, err := a.createSitemap(ctx)
		if err != nil {
			return err
		}

		return etx.Blob(http.StatusOK, "application/xml", []byte(result))
	}

	return etx.Blob(http.StatusOK, "application/xml", []byte(sitemap))
}

type URL struct {
	XMLName    xml.Name `xml:"url"`
	Loc        string   `xml:"loc"`
	ChangeFreq string   `xml:"changefreq"`
	LastMod    string   `xml:"lastmod,omitempty"`
	Priority   string   `xml:"priority,omitempty"`
}

type SitemapXML struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URL     []URL    `xml:"url"`
}

func (a Assets) createSitemap(ctx context.Context) (string, error) {
	baseURL := config.BaseURL
	conn := a.db.Conn()

	var urls []URL

	// Static pages
	urls = append(urls, URL{
		Loc:        baseURL,
		ChangeFreq: "monthly",
		Priority:   "1.0",
	})
	urls = append(urls, URL{
		Loc:        baseURL + routes.AboutPage.URL(),
		ChangeFreq: "monthly",
		Priority:   "0.8",
	})

	// Overview pages
	urls = append(urls, URL{
		Loc:        baseURL + routes.ArticleOverview.URL(),
		ChangeFreq: "weekly",
		Priority:   "0.9",
	})
	urls = append(urls, URL{
		Loc:        baseURL + routes.ProjectOverview.URL(),
		ChangeFreq: "weekly",
		Priority:   "0.9",
	})
	urls = append(urls, URL{
		Loc:        baseURL + routes.NewsletterOverview.URL(),
		ChangeFreq: "weekly",
		Priority:   "0.9",
	})

	// Published articles
	articles, err := models.AllPublishedArticles(ctx, conn)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch published articles for sitemap", "error", err)
	}
	for _, article := range articles {
		urls = append(urls, URL{
			Loc:        baseURL + routes.Article.URL(article.Slug),
			ChangeFreq: "monthly",
			LastMod:    article.UpdatedAt.Format("2006-01-02T15:04:05+00:00"),
			Priority:   "0.7",
		})
	}

	// Published projects
	projects, err := models.AllPublishedProjects(ctx, conn)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch published projects for sitemap", "error", err)
	}
	for _, project := range projects {
		urls = append(urls, URL{
			Loc:        baseURL + routes.Project.URL(project.Slug),
			ChangeFreq: "monthly",
			LastMod:    project.UpdatedAt.Format("2006-01-02T15:04:05+00:00"),
			Priority:   "0.7",
		})
	}

	// Published newsletters
	newsletters, err := models.AllPublishedNewsletters(ctx, conn)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch published newsletters for sitemap", "error", err)
	}
	for _, newsletter := range newsletters {
		urls = append(urls, URL{
			Loc:        baseURL + routes.Newsletter.URL(newsletter.Slug),
			ChangeFreq: "monthly",
			LastMod:    newsletter.UpdatedAt.Format("2006-01-02T15:04:05+00:00"),
			Priority:   "0.7",
		})
	}

	sitemap := SitemapXML{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URL:   urls,
	}

	xmlBytes, err := xml.MarshalIndent(sitemap, "", "  ")
	if err != nil {
		return "", err
	}

	return xml.Header + string(xmlBytes), nil
}

func (a Assets) Stylesheet(etx *echo.Context) error {
	stylesheet, err := assets.Files.ReadFile(
		"css/style.css",
	)
	if err != nil {
		return err
	}

	etx = a.enableCaching(etx, stylesheet)
	return etx.Blob(http.StatusOK, "text/css", stylesheet)
}

func (a Assets) Scripts(etx *echo.Context) error {
	stylesheet, err := assets.Files.ReadFile(
		"js/scripts.js",
	)
	if err != nil {
		return err
	}

	etx = a.enableCaching(etx, stylesheet)
	return etx.Blob(http.StatusOK, "application/javascript", stylesheet)
}

func (a Assets) Script(etx *echo.Context) error {
	param := etx.Param("file")
	stylesheet, err := assets.Files.ReadFile(
		fmt.Sprintf("js/%s", param),
	)
	if err != nil {
		return err
	}

	etx = a.enableCaching(etx, stylesheet)
	return etx.Blob(http.StatusOK, "application/javascript", stylesheet)
}

func (a Assets) Style(etx *echo.Context) error {
	param := etx.Param("file")
	stylesheet, err := assets.Files.ReadFile(
		fmt.Sprintf("css/%s", param),
	)
	if err != nil {
		return err
	}

	etx = a.enableCaching(etx, stylesheet)
	return etx.Blob(http.StatusOK, "text/css", stylesheet)
}
