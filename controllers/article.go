package controllers

import (
	"log"

	"github.com/MBvisti/mortenvistisen/pkg/telemetry"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/gorilla/csrf"
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

	latestArticles, err := c.db.GetLatestPosts(ctx.Request().Context())
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
		Title:       post.Title,
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

type SubscriptionEventForm struct {
	Email string `form:"hero-input"`
	Title string `form:"article-title"`
}

func (c *Controller) SubscriptionEvent(ctx echo.Context) error {
	var form SubscriptionEventForm
	if err := ctx.Bind(&form); err != nil {
		if err := c.mail.Send(ctx.Request().Context(), "hi@mortenvistisen.com", "sub-blog@mortenvistisen.com", "Failed to subscribe", "sub_report", err.Error()); err != nil {
			telemetry.Logger.Error("Failed to send email", "error", err)
		}
		return ctx.String(200, "You're now subscribed!")
	}
	log.Println(form.Email)
	if err := c.mail.Send(ctx.Request().Context(), "hi@mortenvistisen.com", "sub-blog@mortenvistisen.com", "New subscriber", "sub_report", form); err != nil {
		telemetry.Logger.Error("Failed to send email", "error", err)
	}

	return ctx.String(200, "You're now subscribed!")
}
