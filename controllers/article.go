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

	tags, err := c.db.GetTagsForPost(ctx.Request().Context(), post.ID)
	if err != nil {
		return err
	}

	var keywords string
	for i, kw := range tags {
		if i == len(tags)-1 {
			keywords = keywords + kw.Name
		} else {
			keywords = keywords + kw.Name + ", "
		}
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
		ExtraMeta: []views.MetaContent{
			{
				Content: "Morten Vistisen",
				Name:    "author",
			},
			{
				Content: post.Title,
				Name:    "twitter:title",
			},
			{
				Content: post.Excerpt,
				Name:    "twitter:description",
			},
			{
				Content: keywords,
				Name:    "keywords",
			},
		},
	}).Render(views.ExtractRenderDeps(ctx))
}
