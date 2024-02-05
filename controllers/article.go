package controllers

import (
	"fmt"

	"github.com/MBvisti/mortenvistisen/views"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

func (c *Controller) Article(ctx echo.Context) error {
	postSlug := ctx.Param("postSlug")

	post, err := c.db.QueryPostBySlug(ctx.Request().Context(), postSlug)
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

	latestArticles, err := c.db.QueryLatestPosts(ctx.Request().Context())
	if err != nil {
		return err
	}

	otherArticles := make(map[string]string, len(latestArticles))

	for _, article := range latestArticles {
		if article.Slug != post.Slug {
			otherArticles[article.Title] = c.buildURLFromSlug("posts/" + article.Slug)
		}
	}

	return views.ArticlePage(views.ArticlePageData{
		Content:           postContent,
		Title:             post.Title,
		ReleaseDate:       post.ReleasedAt.Time,
		OtherArticleLinks: otherArticles,
		CsrfToken:         csrf.Token(ctx.Request()),
	}, views.Head{
		Title:       fmt.Sprintf("mortenvistisen: %v", post.Title),
		Description: post.Excerpt,
		Slug:        c.buildURLFromSlug("posts/" + post.Slug),
		MetaType:    "article",
		Image:       "https://mortenvistisen.com/static/images/mbv.png",
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
