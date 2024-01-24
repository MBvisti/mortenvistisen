package controllers

import (
	"github.com/MBvisti/mortenvistisen/pkg/telemetry"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

func (c *Controller) HomeIndex(ctx echo.Context) error {
	data, err := c.db.GetLatestPosts(ctx.Request().Context())
	if err != nil {
		telemetry.Logger.Error("failed to get posts", "error", err)
		return err
	}

	posts := make([]views.Post, 0, len(data))

	for _, d := range data {
		tagsData, err := c.db.GetTagsForPost(ctx.Request().Context(), d.ID)
		if err != nil {
			telemetry.Logger.Error("failed to get tags", "error", err)
			return err
		}

		tags := make([]string, 0, len(tagsData))
		for _, t := range tagsData {
			tags = append(tags, t.Name)
		}

		posts = append(posts, views.Post{
			Title:       d.Title,
			ReleaseDate: d.ReleasedAt.Time.String(),
			Excerpt:     d.Excerpt,
			Tags:        tags,
			Slug:        c.formatArticleSlug(d.Slug),
		})
	}

	return views.HomePage(posts).Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) About(ctx echo.Context) error {
	return views.AboutPage(views.Head{
		Title:       "About",
		Description: "Contains information about the site owner, Morten Vistisen, the purpose of the site and the technologies used to build it.",
		Slug:        c.buildURLFromSlug("about"),
		Image:       "https://mortenvistisen.com/static/images/mbv.png",
		MetaType:    "website",
	}).Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) Newsletter(ctx echo.Context) error {
	return views.NewsletterPage(views.Head{
		Title:       "Newsletter",
		Description: "Signup page for joining Morten's newsletter where he shares his thoughts on software development, business and life.",
		Slug:        c.buildURLFromSlug("newsletter"),
		MetaType:    "website",
		Image:       "https://mortenvistisen.com/static/images/mbv.png",
	}, csrf.Token(ctx.Request())).Render(views.ExtractRenderDeps(ctx))
}
