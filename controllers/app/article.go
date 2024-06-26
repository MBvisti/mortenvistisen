package app

import (
	"github.com/MBvisti/mortenvistisen/controllers"
	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

func Article(
	ctx echo.Context,
	articleModel models.ArticleService,
	postManager posts.PostManager,
) error {
	postSlug := ctx.Param("postSlug")

	post, err := articleModel.BySlug(ctx.Request().Context(), postSlug)
	if err != nil {
		return err
	}

	postContent, err := postManager.Parse(post.Filename)
	if err != nil {
		return err
	}

	var keywords string
	for i, tag := range post.Tags {
		if i == len(post.Tags)-1 {
			keywords = keywords + tag.Name
		} else {
			keywords = keywords + tag.Name + ", "
		}
	}

	// fiveRandomPosts, err := db.GetFiveRandomPosts(
	// 	ctx.Request().Context(),
	// 	post.ID,
	// )
	// if err != nil {
	// 	return err
	// }
	//
	// otherArticles := make(map[string]string, 5)
	//
	// for _, article := range fiveRandomPosts {
	// 	otherArticles[article.Title] = controllers.BuildURLFromSlug(
	// 		"posts/" + article.Slug,
	// 	)
	// }
	//
	return views.ArticlePage(views.ArticlePageData{
		Content:           postContent,
		HeaderTitle:       post.HeaderTitle,
		ReleaseDate:       post.ReleaseDate,
		OtherArticleLinks: nil,
		CsrfToken:         csrf.Token(ctx.Request()),
	}, views.Head{
		Title:       post.Title,
		Description: post.Excerpt,
		Slug:        controllers.BuildURLFromSlug("posts/" + post.Slug),
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
