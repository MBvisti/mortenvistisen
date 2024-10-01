package handlers

import (
	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/MBvisti/mortenvistisen/views"
	view "github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/labstack/echo/v4"
)

type Dashboard struct {
	base          Base
	articleSvc    models.ArticleService
	tagSvc        models.TagService
	postManager   posts.PostManager
	newsletterSvc models.NewsletterService
	subscriberSvc models.SubscriberService
	tokenService  services.Token
	emailService  services.Email
}

func NewDashboard(
	base Base,
	articleSvc models.ArticleService,
	tagSvc models.TagService,
	postManager posts.PostManager,
	newsletterSvc models.NewsletterService,
	subscriberSvc models.SubscriberService,
	tokenService services.Token,
	emailService services.Email,
) Dashboard {
	return Dashboard{
		base,
		articleSvc,
		tagSvc,
		postManager,
		newsletterSvc,
		subscriberSvc,
		tokenService,
		emailService,
	}
}

func (d Dashboard) Index(ctx echo.Context) error {
	data, err := d.articleSvc.List(ctx.Request().Context(), 0, 50)
	if err != nil {
		return err
	}

	posts := make([]views.Post, 0, len(data))
	for _, dd := range data {
		var tags []string
		for _, tag := range dd.Tags {
			tags = append(tags, tag.Name)
		}

		posts = append(posts, views.Post{
			ID:          dd.ID.String(),
			Title:       dd.Title,
			ReleaseDate: dd.ReleaseDate.String(),
			Excerpt:     dd.Excerpt,
			Tags:        tags,
			Slug:        d.base.FormatArticleSlug(dd.Slug),
		})
	}
	return view.Index(posts).Render(views.ExtractRenderDeps(ctx))
}
