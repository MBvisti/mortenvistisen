package controllers

import (
	"crypto/md5"
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"strings"
	"time"

	"mortenvistisen/assets"
	"mortenvistisen/config"
	"mortenvistisen/internal/server"
	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
	"mortenvistisen/router"
	"mortenvistisen/router/routes"

	"github.com/labstack/echo/v5"
)

const threeMonthsCache = "7776000"

type Assets struct {
	cache *Cache[string]
	db    storage.Pool
}

func NewAssets(cache *Cache[string], db storage.Pool) Assets {
	return Assets{cache: cache, db: db}
}

func (a Assets) RegisterRoutes(r *router.Router) error {
	errs := []error{}

	_, err := r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.Robots.Path(),
		Name:    routes.Robots.Name(),
		Handler: a.Robots,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.Sitemap.Path(),
		Name:    routes.Sitemap.Name(),
		Handler: a.Sitemap,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.Stylesheet.Path(),
		Name:    routes.Stylesheet.Name(),
		Handler: a.Stylesheet,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.Scripts.Path(),
		Name:    routes.Scripts.Name(),
		Handler: a.Scripts,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.Script.Path(),
		Name:    routes.Script.Name(),
		Handler: a.Script,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.ViteBuild.Path(),
		Name:    routes.ViteBuild.Name(),
		Handler: a.ViteBuild,
	})
	if err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
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

func createRobotsTxt() string {
	return fmt.Sprintf(
		"User-agent: *\nAllow: /\nSitemap: %s%s\n",
		config.BaseURL,
		routes.Sitemap.URL(),
	)
}

func (a Assets) Robots(etx *echo.Context) error {
	cacheKey := "assets:robots"

	robotsTxt, err := a.cache.Get(cacheKey, func() (string, error) {
		return createRobotsTxt(), nil
	})
	if err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to get robots.txt from cache",
			"error", err,
		)
		return etx.String(http.StatusOK, createRobotsTxt())
	}

	return etx.String(http.StatusOK, robotsTxt)
}

func (a Assets) Sitemap(etx *echo.Context) error {
	cacheKey := "assets:sitemap"
	ctx := etx.Request().Context()

	sitemap, err := a.cache.Get(cacheKey, func() (string, error) {
		db := a.db.Executor()
		articles, err := models.Article.LatestPublished(ctx, db, 0)
		if err != nil {
			return "", fmt.Errorf("load sitemap articles: %w", err)
		}
		newsletters, err := models.Newsletter.LatestPublished(ctx, db, 0)
		if err != nil {
			return "", fmt.Errorf("load sitemap newsletters: %w", err)
		}
		projects, err := models.Project.LatestPublished(ctx, db, 0)
		if err != nil {
			return "", fmt.Errorf("load sitemap projects: %w", err)
		}
		return createSitemap(articles, newsletters, projects)
	})
	if err != nil {
		return err
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

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URL     []URL    `xml:"url"`
}

func createSitemap(
	articles []models.ArticleEntity,
	newsletters []models.NewsletterEntity,
	projects []models.ProjectEntity,
) (string, error) {
	absoluteURL := func(value string) string {
		return strings.TrimRight(config.BaseURL, "/") + "/" + strings.TrimLeft(value, "/")
	}
	lastModified := func(value time.Time) string {
		if value.IsZero() {
			return ""
		}
		return value.Format(time.RFC3339)
	}

	urls := []URL{
		{Loc: absoluteURL(routes.HomePage.URL()), ChangeFreq: "monthly", Priority: "1"},
		{Loc: absoluteURL(routes.AboutPage.URL()), ChangeFreq: "monthly", Priority: "0.8"},
		{Loc: absoluteURL(routes.PostIndex.URL()), ChangeFreq: "weekly", Priority: "0.9"},
		{Loc: absoluteURL(routes.NewsletterIndex.URL()), ChangeFreq: "weekly", Priority: "0.9"},
		{Loc: absoluteURL(routes.ProjectIndex.URL()), ChangeFreq: "weekly", Priority: "0.9"},
	}
	for _, article := range articles {
		urls = append(urls, URL{
			Loc:        absoluteURL(routes.PostShow.URL(article.Slug)),
			ChangeFreq: "monthly",
			LastMod:    lastModified(article.UpdatedAt),
			Priority:   "0.7",
		})
	}
	for _, newsletter := range newsletters {
		urls = append(urls, URL{
			Loc:        absoluteURL(routes.NewsletterShow.URL(newsletter.Slug.String)),
			ChangeFreq: "monthly",
			LastMod:    lastModified(newsletter.UpdatedAt),
			Priority:   "0.7",
		})
	}
	for _, project := range projects {
		urls = append(urls, URL{
			Loc:        absoluteURL(routes.ProjectShow.URL(project.Slug)),
			ChangeFreq: "monthly",
			LastMod:    lastModified(project.UpdatedAt),
			Priority:   "0.7",
		})
	}

	sitemap := Sitemap{
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
func contentTypeByExt(name string) string {
	switch {
	case strings.HasSuffix(name, ".js"):
		return "application/javascript"
	case strings.HasSuffix(name, ".css"):
		return "text/css"
	case strings.HasSuffix(name, ".json"):
		return "application/json"
	case strings.HasSuffix(name, ".svg"):
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}

func (a Assets) ViteBuild(etx *echo.Context) error {
	param := strings.TrimPrefix(etx.Param("*"), "/")
	param = path.Clean(param)
	if param == "." || param == "" || strings.HasPrefix(param, "../") {
		return echo.NewHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	data, err := assets.Files.ReadFile(fmt.Sprintf("dist/%s", param))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	etx = a.enableCaching(etx, data)
	return etx.Blob(http.StatusOK, contentTypeByExt(param), data)
}
