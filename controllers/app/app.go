package app

import (
	"github.com/MBvisti/mortenvistisen/controllers/internal/utilities"
	"github.com/MBvisti/mortenvistisen/pkg/telemetry"
	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

func Index(ctx echo.Context, db database.Queries) error {
	data, err := db.GetLatestPosts(ctx.Request().Context())
	if err != nil {
		telemetry.Logger.Error("failed to get posts", "error", err)
		return err
	}

	posts := make([]views.Post, 0, len(data))

	for _, d := range data {
		tagsData, err := db.GetTagsForPost(ctx.Request().Context(), d.ID)
		if err != nil {
			telemetry.Logger.Error("failed to get tags", "error", err)
			return err
		}

		tags := make([]string, 0, len(tagsData))
		for _, t := range tagsData {
			tags = append(tags, t.Name)
		}

		posts = append(posts, views.Post{
			Title:       d.HeaderTitle.String,
			ReleaseDate: d.ReleasedAt.Time.String(),
			Excerpt:     d.Excerpt,
			Tags:        tags,
			Slug:        utilities.FormatArticleSlug(d.Slug),
		})
	}

	return views.HomePage(posts).Render(views.ExtractRenderDeps(ctx))
}

func About(ctx echo.Context) error {
	return views.AboutPage(views.Head{
		Title:       "About",
		Description: "Contains information about the site owner, Morten Vistisen, the purpose of the site and the technologies used to build it.",
		Slug:        utilities.BuildURLFromSlug("about"),
		Image:       "https://mortenvistisen.com/static/images/mbv.png",
		MetaType:    "website",
	}).Render(views.ExtractRenderDeps(ctx))
}

func Newsletter(ctx echo.Context) error {
	return views.NewsletterPage(views.Head{
		Title:       "Newsletter",
		Description: "Signup page for joining Morten's newsletter where he shares his thoughts on software development, business and life.",
		Slug:        utilities.BuildURLFromSlug("newsletter"),
		MetaType:    "website",
		Image:       "https://mortenvistisen.com/static/images/mbv.png",
	}, csrf.Token(ctx.Request())).Render(views.ExtractRenderDeps(ctx))
}

func Projects(ctx echo.Context) error {
	return views.ProjectsPage(views.Head{
		Title:       "Projects",
		Description: "A collection of on-going and retired projects I've build over the years. This includes business projects, open source projects and personal projects.",
		Slug:        utilities.BuildURLFromSlug("projects"),
		MetaType:    "website",
		Image:       "https://mortenvistisen.com/static/images/mbv.png",
	}).Render(views.ExtractRenderDeps(ctx))
}
