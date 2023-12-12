package controllers

import (
	"fmt"

	"github.com/MBvisti/grafto/views"
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

	return views.Article(ctx, views.ArticlePageData{
		Content:     postContent,
		ReleaseDate: post.ReleasedAt.Time,
		Meta: views.ArticleMetaData{
			Title:       post.Title,
			Description: post.Excerpt,
			Image:       "",
			Slug:        fmt.Sprintf("http://localhost:8000/%s/%s", "posts", post.Slug),
		},
	})
}
