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
	articlePageCacheKey = "articlePage--"
)

type App struct {
	db           psql.Postgres
	cache        otter.Cache[string, templ.Component]
	articleCache otter.Cache[string, views.ArticlePageProps]
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

	articleCacheBuilder, err := otter.NewBuilder[string, views.ArticlePageProps](
		100,
	)
	if err != nil {
		panic(err)
	}

	articleCache, err := articleCacheBuilder.WithTTL(48 * time.Hour).Build()
	if err != nil {
		panic(err)
	}

	return App{db, pageCacher, articleCache}
}

func (a App) LandingPage(c echo.Context) error {
	// if value, ok := a.cache.Get(landingPageCacheKey); ok {
	// 	return views.HomePage(value).Render(renderArgs(c))
	// }

	articles, err := models.GetPublishedArticles(extractCtx(c), a.db.Pool)
	if err != nil {
		return err
	}

	var payload []views.HomeArticle
	for _, ar := range articles {
		tags := make([]string, len(ar.Tags))

		for i, t := range ar.Tags {
			tags[i] = t.Title
		}

		payload = append(payload, views.HomeArticle{
			Title:       ar.Title,
			Description: ar.Excerpt,
			Slug:        ar.Slug,
			PublishedAt: ar.FirstPublishedAt,
			Tags:        tags,
		})
	}

	return views.HomePage(views.Home(payload)).Render(renderArgs(c))
}

func (a App) AboutPage(c echo.Context) error {
	return views.AboutPage().Render(renderArgs(c))
}

func (a App) ArticlePage(c echo.Context) error {
	slug := c.Param("articleSlug")

	if value, ok := a.articleCache.Get(articlePageCacheKey + slug); ok {
		return views.ArticlePage(value).
			Render(renderArgs(c))
	}

	article, err := models.GetArticleBySlug(
		c.Request().Context(),
		a.db.Pool,
		slug,
	)
	if err != nil {
		return err
	}

	manager := NewManager()
	ar, e := manager.ParseContent(article.Content)
	if e != nil {
		return e
	}

	props := views.ArticlePageProps{
		Slug:      slug,
		MetaTitle: article.MetaTitle,
		MetaDesc:  article.MetaDescription,
		Content: views.Article(
			article.ImageLink,
			ar,
			article.Title,
			article.FirstPublishedAt,
			article.UpdatedAt,
		),
	}

	return views.ArticlePage(props).
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
