package views

import (
	"github.com/MBvisti/grafto/views/internal/components"
	"github.com/MBvisti/grafto/views/internal/layouts"
	"github.com/MBvisti/grafto/views/internal/pages"
	"github.com/labstack/echo/v4"
)

type ArticleMetaData struct {
	Title       string
	Description string
	Image       string
	Slug        string
}

type ArticlePageData struct {
	Content string
	Meta    ArticleMetaData
}

func Article(ctx echo.Context, data ArticlePageData) error {
	head := components.Head{
		Title:       data.Meta.Title,
		Description: data.Meta.Description,
		Image:       data.Meta.Image,
		Slug:        data.Meta.Slug,
		MetaType:    "article",
	}

	return layouts.Base(pages.Article(data.Content), head).Render(extractRenderDeps(ctx))
}
