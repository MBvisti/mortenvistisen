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
	tokenService  services.TokenSvc
	emailService  services.Email
}

func NewDashboard(
	base Base,
	articleSvc models.ArticleService,
	tagSvc models.TagService,
	postManager posts.PostManager,
	newsletterSvc models.NewsletterService,
	subscriberSvc models.SubscriberService,
	tokenService services.TokenSvc,
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
	return view.Index().Render(views.ExtractRenderDeps(ctx))
}
