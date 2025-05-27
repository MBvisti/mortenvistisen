package handlers

import (
	"bytes"
	"embed"
	"log/slog"
	"time"

	"github.com/a-h/templ"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/labstack/echo/v4"
	"github.com/maypok86/otter"
	"github.com/mbvisti/mortenvistisen/models"
	"github.com/mbvisti/mortenvistisen/psql"
	"github.com/mbvisti/mortenvistisen/router/routes"

	"github.com/mbvisti/mortenvistisen/views"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	landingPageCacheKey = "landingPage"
	articlePageCacheKey = "articlePage"
)

type App struct {
	db    psql.Postgres
	cache otter.Cache[string, templ.Component]
}

func newApp(
	db psql.Postgres,
) App {
	cacheBuilder, err := otter.NewBuilder[string, templ.Component](20)
	if err != nil {
		panic(err)
	}

	pageCacher, err := cacheBuilder.WithTTL(48 * time.Hour).Build()
	if err != nil {
		panic(err)
	}

	return App{db, pageCacher}
}

func (a App) LandingPage(c echo.Context) error {
	if value, ok := a.cache.Get(landingPageCacheKey); ok {
		return views.HomePage(value).Render(renderArgs(c))
	}

	return views.HomePage(views.Home()).Render(renderArgs(c))
}

func (a App) AboutPage(c echo.Context) error {
	return views.AboutPage().Render(renderArgs(c))
}

func (a App) ArticlePage(c echo.Context) error {
	slug := c.Param("articleSlug")

	article, err := models.GetArticleBySlug(
		c.Request().Context(),
		a.db.Pool,
		slug,
	)
	if err != nil {
		return err
	}

	slog.Info(
		"$$$$$$$$$$$$$$$$$$$$$$$$",
		"article",
		article.Content,
		"id",
		article.ID,
	)

	manager := NewManager()
	ar, e := manager.ParseContent(article.Content)
	if e != nil {
		return err
	}

	// if value, ok := a.cache.Get(articlePageCacheKey); ok {
	// 	return views.ArticlePage("A love letter to Go", slug, "desc", value).
	// 		Render(renderArgs(c))
	// }

	return views.ArticlePage("A love letter to Go", slug, "desc", views.Article(ar)).
		Render(renderArgs(c))
}

func (a App) Redirect(c echo.Context) error {
	to := c.QueryParam("to")
	for _, r := range routes.AllRoutes {
		if to == r.Path {
			return redirectHx(c.Response(), to)
		}
	}

	slog.InfoContext(c.Request().Context(),
		"security warning: someone tried to missue open redirect",
		"to", to,
		"ip", c.RealIP(),
	)

	return redirect(c.Response(), c.Request(), "/")
}

type Manager struct {
	posts           embed.FS
	markdownHandler goldmark.Markdown
}

func NewManager() Manager {
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

	return Manager{
		markdownHandler: md,
	}
}

func (pm *Manager) GetPost(name string) (string, error) {
	source, err := pm.posts.ReadFile(name)
	if err != nil {
		slog.Error("failed to read markdown file", "error", err)
		return "", err
	}

	return string(source), nil
}

func (pm *Manager) Parse(name string) (string, error) {
	source, err := pm.posts.ReadFile(name)
	if err != nil {
		slog.Error("failed to read markdown file", "error", err)
		return "", err
	}

	// Parse Markdown content
	var htmlOutput bytes.Buffer
	if err := pm.markdownHandler.Convert(source, &htmlOutput); err != nil {
		slog.Error("failed to parse markdown file", "error", err)
		return "", err
	}

	return htmlOutput.String(), nil
}

func (pm *Manager) ParseContent(content string) (string, error) {
	// Parse Markdown content
	var htmlOutput bytes.Buffer
	if err := pm.markdownHandler.Convert([]byte(content), &htmlOutput); err != nil {
		slog.Error("failed to parse markdown file", "error", err)
		return "", err
	}

	return htmlOutput.String(), nil
}
