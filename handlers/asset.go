package handlers

import (
	//nolint:gosec //only needed for browser caching
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/maypok86/otter"
	"github.com/mbvisti/mortenvistisen/assets"
	"github.com/mbvisti/mortenvistisen/config"
	"github.com/mbvisti/mortenvistisen/models"
	"github.com/mbvisti/mortenvistisen/psql"
	"github.com/mbvisti/mortenvistisen/router/routes"
	"gopkg.in/yaml.v2"
)

const (
	sitemapCacheKey  = "assets.sitemap"
	robotsCacheKey   = "assets.robots"
	weekInHours      = 168
	threeInHours     = 72
	threeMonthsCache = "63072000"
)

type Assets struct {
	sitemapCache otter.Cache[string, Sitemap]
	assetsCache  otter.Cache[string, string]
	db           psql.Postgres
}

func newAssets(
	db psql.Postgres,
) Assets {
	sitemapCacheBuilder, err := otter.NewBuilder[string, Sitemap](100)
	if err != nil {
		panic(err)
	}

	sitemapCache, err := sitemapCacheBuilder.WithTTL(weekInHours).Build()
	if err != nil {
		panic(err)
	}

	robotsCacheBuilder, err := otter.NewBuilder[string, string](1)
	if err != nil {
		panic(err)
	}

	robotsCache, err := robotsCacheBuilder.WithTTL(weekInHours).Build()
	if err != nil {
		panic(err)
	}

	return Assets{sitemapCache, robotsCache, db}
}

func (a Assets) enableCaching(c echo.Context, content []byte) echo.Context {
	if config.Cfg.Environment == config.PROD_ENVIRONMENT {
		//nolint:gosec //only needed for browser caching
		hash := md5.Sum(content)
		etag := fmt.Sprintf(`W/"%x-%x"`, hash, len(content))

		c.Response().Header().Del("Set-Cookie")

		c.Response().
			Header().
			Set(
				"Cache-Control",
				fmt.Sprintf("public, max-age=%s, immutable", threeMonthsCache),
			)
		c.Response().
			Header().
			Set("Vary", "Accept-Encoding")
		c.Response().
			Header().
			Set("ETag", etag)
	}

	return c
}

func (a Assets) Robots(c echo.Context) error {
	if value, ok := a.assetsCache.Get(robotsCacheKey); ok {
		return c.String(http.StatusOK, string(value))
	}

	type robotsTxt struct {
		UserAgent string `yaml:"User-agent"`
		Allow     string `yaml:"Allow"`
		Sitemap   string `yaml:"Sitemap"`
	}

	robots, err := yaml.Marshal(robotsTxt{
		UserAgent: "*",
		Allow:     "/",
		Sitemap: fmt.Sprintf(
			"%s%s",
			config.Cfg.GetFullDomain(),
			routes.Sitemap.Path,
		),
	})
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, string(robots))
}

func (a Assets) Sitemap(c echo.Context) error {
	if value, ok := a.sitemapCache.Get(sitemapCacheKey); ok {
		return c.XML(http.StatusOK, value)
	}

	articles, err := models.GetPublishedArticles(
		c.Request().Context(),
		a.db.Pool,
	)
	if err != nil {
		return err
	}

	newsletters, err := models.GetPublishedNewsletters(
		c.Request().Context(),
		a.db.Pool,
	)
	if err != nil {
		return err
	}

	sitemap, err := createSitemap(newsletters, articles)
	if err != nil {
		return err
	}

	if ok := a.sitemapCache.Set(sitemapCacheKey, sitemap); !ok {
		slog.ErrorContext(
			c.Request().Context(),
			"could not set sitemap cache",
			"error",
			err,
		)
	}

	return c.XML(http.StatusOK, sitemap)
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
	newsletters []models.Newsletter,
	articles []models.Article,
) (Sitemap, error) {
	baseUrl := config.Cfg.GetFullDomain()

	var urls []URL

	for _, r := range routes.App {
		if r.Path != "/redirect" && !strings.Contains(r.Path, ":") &&
			r.Method != http.MethodPost && r.Method != http.MethodPut {
			urls = append(urls, URL{
				Loc: fmt.Sprintf(
					"%s%s",
					baseUrl,
					r.Path,
				),
				ChangeFreq: "monthly",
			})
		}
	}

	for _, article := range articles {
		path := strings.Replace(
			routes.ArticlePage.Path,
			":articleSlug",
			article.Slug,
			1,
		)
		urls = append(urls, URL{
			Loc: fmt.Sprintf(
				"%s%s",
				baseUrl,
				path,
			),
			ChangeFreq: "monthly",
			LastMod:    article.UpdatedAt.Format("2006-01-02"),
		})
	}

	for _, newsletter := range newsletters {
		path := strings.Replace(
			routes.NewsletterPage.Path,
			":newsletterSlug",
			newsletter.Slug,
			1,
		)
		urls = append(urls, URL{
			Loc: fmt.Sprintf(
				"%s%s",
				baseUrl,
				path,
			),
			ChangeFreq: "monthly",
			LastMod:    newsletter.UpdatedAt.Format("2006-01-02"),
		})
	}

	sitemap := Sitemap{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URL:   urls,
	}

	return sitemap, nil
}

func (a Assets) Styles(c echo.Context) error {
	stylesheet, err := assets.Files.ReadFile(
		"css/styles.css",
	)
	if err != nil {
		return err
	}

	c = a.enableCaching(c, stylesheet)

	return c.Blob(http.StatusOK, "text/css", stylesheet)
}

func (a Assets) IndividualCSSFile(c echo.Context) error {
	filename := c.Param("file")
	stylesheet, err := assets.Files.ReadFile(
		fmt.Sprintf("css/%s", filename),
	)
	if err != nil {
		return err
	}

	c = a.enableCaching(c, stylesheet)

	return c.Blob(http.StatusOK, "text/css", stylesheet)
}

// func (a Assets) IndividualScript(c echo.Context) error {
// 	filename := strings.Split(c.Path(), routes.AssetsRoutePrefix)[1]
//
// 	script, err := assets.Files.ReadFile(
// 		strings.TrimPrefix(filename, "/"),
// 	)
// 	if err != nil {
// 		return err
// 	}
//
// 	c = a.enableCaching(c, script)
//
// 	return c.Blob(http.StatusOK, "application/javascript", script)
// }

func (a Assets) Scripts(c echo.Context) error {
	script, err := assets.Files.ReadFile(
		"js/script.js",
	)
	if err != nil {
		return err
	}

	c = a.enableCaching(c, script)

	return c.Blob(http.StatusOK, "application/javascript", script)
}

func (a Assets) ScriptsDashboard(c echo.Context) error {
	script, err := assets.Files.ReadFile(
		"js/dashboard_script.js",
	)
	if err != nil {
		return err
	}

	c = a.enableCaching(c, script)

	return c.Blob(http.StatusOK, "application/javascript", script)
}

func (a Assets) IndividualScript(c echo.Context) error {
	filename := c.Param("file")
	script, err := assets.Files.ReadFile(
		fmt.Sprintf("js/%s", filename),
	)
	if err != nil {
		return err
	}

	c = a.enableCaching(c, script)

	return c.Blob(http.StatusOK, "application/javascript", script)
}

func (a Assets) Favicon(c echo.Context) error {
	img, err := assets.Files.ReadFile(
		"images/favicon.ico",
	)
	if err != nil {
		return err
	}

	c = a.enableCaching(c, img)

	return c.Blob(http.StatusOK, "image/png", img)
}

func (a Assets) Favicon16(c echo.Context) error {
	img, err := assets.Files.ReadFile(
		"images/favicon-16.png",
	)
	if err != nil {
		return err
	}

	c = a.enableCaching(c, img)

	return c.Blob(http.StatusOK, "image/png", img)
}

func (a Assets) Favicon32(c echo.Context) error {
	img, err := assets.Files.ReadFile(
		"images/favicon-32.png",
	)
	if err != nil {
		return err
	}

	c = a.enableCaching(c, img)

	return c.Blob(http.StatusOK, "image/png", img)
}

func (a Assets) FaviconAppleTouch(c echo.Context) error {
	img, err := assets.Files.ReadFile(
		"images/apple-touch-icon.png",
	)
	if err != nil {
		return err
	}

	c = a.enableCaching(c, img)

	return c.Blob(http.StatusOK, "image/png", img)
}

func (a Assets) FaviconSiteManifest(c echo.Context) error {
	img, err := assets.Files.ReadFile(
		"images/site.webmanifest",
	)
	if err != nil {
		return err
	}

	c = a.enableCaching(c, img)

	return c.Blob(http.StatusOK, "image/png", img)
}

func (a Assets) LLM(c echo.Context) error {
	content, err := assets.Files.ReadFile(
		"files/llm.txt",
	)
	if err != nil {
		return err
	}

	c = a.enableCaching(c, content)

	return c.String(http.StatusOK, string(content))
}

func (a Assets) IndexNow(c echo.Context) error {
	content, err := assets.Files.ReadFile(
		"files/index_now.txt",
	)
	if err != nil {
		return err
	}

	c = a.enableCaching(c, content)

	return c.String(http.StatusOK, string(content))
}
