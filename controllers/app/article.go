package app

import (
	"github.com/MBvisti/mortenvistisen/controllers/internal/utilities"
	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

func Article(ctx echo.Context, db database.Queries, postManager posts.PostManager) error {
	postSlug := ctx.Param("postSlug")

	post, err := db.GetPostBySlug(ctx.Request().Context(), postSlug)
	if err != nil {
		return err
	}

	postContent, err := postManager.Parse(post.Filename)
	if err != nil {
		return err
	}

	tags, err := db.GetTagsForPost(ctx.Request().Context(), post.ID)
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

	fiveRandomPosts, err := db.GetFiveRandomPosts(
		ctx.Request().Context(),
		post.ID,
	)
	if err != nil {
		return err
	}

	otherArticles := make(map[string]string, 5)

	for _, article := range fiveRandomPosts {
		otherArticles[article.Title] = utilities.BuildURLFromSlug(
			"posts/" + article.Slug,
		)
	}

	return views.ArticlePage(views.ArticlePageData{
		Content:           postContent,
		HeaderTitle:       post.HeaderTitle.String,
		ReleaseDate:       post.ReleasedAt.Time,
		OtherArticleLinks: otherArticles,
		CsrfToken:         csrf.Token(ctx.Request()),
	}, views.Head{
		Title:       post.Title,
		Description: post.Excerpt,
		Slug:        utilities.BuildURLFromSlug("posts/" + post.Slug),
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
