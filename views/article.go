package views

import (
	"fmt"
	"time"

	"github.com/MBvisti/grafto/views/internal/components"
	"github.com/MBvisti/grafto/views/internal/layouts"
	"github.com/MBvisti/grafto/views/internal/pages"
	"github.com/golang-module/carbon/v2"
	"github.com/labstack/echo/v4"
)

type ArticleMetaData struct {
	Title       string
	Description string
	Image       string
	Slug        string
}

type ArticlePageData struct {
	Content     string
	ReleaseDate time.Time
	Meta        ArticleMetaData
}

func Article(ctx echo.Context, data ArticlePageData) error {
	head := components.Head{
		Title:       data.Meta.Title,
		Description: data.Meta.Description,
		Image:       data.Meta.Image,
		Slug:        data.Meta.Slug,
		MetaType:    "article",
	}
	releasedAt := fmt.Sprintf("%s %v, %v", carbon.CreateFromStdTime(data.ReleaseDate).ToShortMonthString(),
		carbon.CreateFromStdTime(data.ReleaseDate).DayOfMonth(), carbon.CreateFromStdTime(data.ReleaseDate).Year())

	return layouts.Base(pages.Article(data.Meta.Title, releasedAt, data.Content), head).Render(extractRenderDeps(ctx))
}
