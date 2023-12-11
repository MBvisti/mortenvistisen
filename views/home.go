package views

import (
	"github.com/MBvisti/grafto/views/internal/components"
	"github.com/MBvisti/grafto/views/internal/layouts"
	"github.com/MBvisti/grafto/views/internal/pages"
	"github.com/labstack/echo/v4"
)

type Post struct {
	Title       string
	ReleaseDate string
	Tags        []string
	Slug        string
}

type HomePageData struct {
	Posts []Post
}

func HomeIndex(ctx echo.Context, data HomePageData) error {
	posts := make([]pages.Post, 0, 5)

	for _, post := range data.Posts {
		posts = append(posts, pages.Post{
			Title:       post.Title,
			ReleaseDate: post.ReleaseDate,
			Tags:        post.Tags,
			Slug:        post.Slug,
		})
	}

	return layouts.Base(pages.HomeIndex(posts), components.Head{}.Default()).Render(extractRenderDeps(ctx))
}
