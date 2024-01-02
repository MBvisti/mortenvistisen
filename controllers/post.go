package controllers

import (
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/labstack/echo/v4"
)

func (c *Controller) Article(ctx echo.Context) error {
	postSlug := ctx.Param("postSlug")

	post, err := c.db.GetPostBySlug(ctx.Request().Context(), postSlug)
	if err != nil {
		return err
	}

	postContent, err := c.postManager.Parse(post.Filename)
	if err != nil {
		return err
	}

	return views.ArticlePage(views.ArticlePageData{
		Content:     postContent,
		Title:       post.Title,
		ReleaseDate: post.ReleasedAt.Time,
	}, views.Head{
		Title:       post.Title,
		Description: post.Excerpt,
		Slug:        c.buildURLFromSlug(post.Slug),
		MetaType:    "article",
	}).Render(views.ExtractRenderDeps(ctx))
}
