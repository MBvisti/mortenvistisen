package handlers

import (
	"errors"
	"net/http"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/authentication"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

type App struct {
	base          Base
	articleSvc    models.ArticleService
	subscriberSvc models.SubscriberService
	postManager   posts.PostManager
	tokenService  services.Token
}

func NewApp(
	base Base,
	articleSvc models.ArticleService,
	subscriberSvc models.SubscriberService,
	postManager posts.PostManager,
	tokenService services.Token,
) App {
	return App{base, articleSvc, subscriberSvc, postManager, tokenService}
}

func (a App) Index(ctx echo.Context) error {
	data, err := a.articleSvc.List(ctx.Request().Context(), 0, 50)
	if err != nil {
		return err
	}

	posts := make([]views.Post, 0, len(data))
	for _, d := range data {
		var tags []string
		for _, tag := range d.Tags {
			tags = append(tags, tag.Name)
		}

		posts = append(posts, views.Post{
			Title:       d.Title,
			ReleaseDate: d.ReleaseDate.String(),
			Excerpt:     d.Excerpt,
			Tags:        tags,
			Slug:        a.base.FormatArticleSlug(d.Slug),
		})
	}

	return views.HomePage(posts).Render(views.ExtractRenderDeps(ctx))
}

func (a App) About(ctx echo.Context) error {
	return views.AboutPage(views.Head{
		Title:       "About",
		Description: "Contains information about the site owner, Morten Vistisen, the purpose of the site and the technologies used to build it.",
		Slug:        a.base.BuildURLFromSlug("about"),
		Image:       "https://mortenvistisen.com/static/images/mbv.png",
		MetaType:    "website",
	}).Render(views.ExtractRenderDeps(ctx))
}

func (a App) Newsletter(ctx echo.Context) error {
	return views.NewsletterPage(views.Head{
		Title:       "Newsletter",
		Description: "Signup page for joining Morten's newsletter where he shares his thoughts on software development, business and life.",
		Slug:        a.base.BuildURLFromSlug("newsletter"),
		MetaType:    "website",
		Image:       "https://mortenvistisen.com/static/images/mbv.png",
	}, csrf.Token(ctx.Request())).Render(views.ExtractRenderDeps(ctx))
}

func (a App) Projects(ctx echo.Context) error {
	return views.ProjectsPage(views.Head{
		Title:       "Projects",
		Description: "A collection of on-going and retired projects I've build over the years. This includes business projects, open source projects and personal projects.",
		Slug:        a.base.BuildURLFromSlug("projects"),
		MetaType:    "website",
		Image:       "https://mortenvistisen.com/static/images/mbv.png",
	}).Render(views.ExtractRenderDeps(ctx))
}

func (a App) SubscriptionEvent(ctx echo.Context) error {
	type subscriptionEventForm struct {
		Email string `form:"hero-input"`
		Title string `form:"article-title"`
	}

	var form subscriptionEventForm
	if err := ctx.Bind(&form); err != nil {
		return ctx.String(200, "You're now subscribed!")
	}

	param := ctx.QueryParam("book")
	var bookSub bool
	if param == "true" {
		bookSub = true
	}

	if err := a.subscriberSvc.New(ctx.Request().Context(), form.Email, form.Title, bookSub); err != nil {
		if errors.Is(err, models.ErrSubscriberExists) {
			return views.SubscribeModalResponse(bookSub, true).
				Render(views.ExtractRenderDeps(ctx))
		}
		if errors.Is(err, models.ErrValidation{}) {
			return nil
		}
		if errors.Is(err, models.ErrUnrecoverableEvent) {
			return ctx.JSON(http.StatusInternalServerError, "yo")
		}
	}

	return views.SubscribeModalResponse(bookSub, false).
		Render(views.ExtractRenderDeps(ctx))
}

func (a App) HowToStartFreelancing(ctx echo.Context) error {
	return views.HowToStartFreelancing(views.Head{
		Title:       "How To Start Freelancing Book",
		Description: "Want to start freelancing? Kickstart your journey to working on your own terms in the best possible way by signing up for the How To Start Freelancing Book.",
		Slug:        a.base.BuildURLFromSlug("books/" + "how-to-start-freelancing"),
		MetaType:    "website",
		Image:       "https://mortenvistisen.com/static/images/mbv.png",
		ExtraMeta: []views.MetaContent{
			{
				Content: "Morten Vistisen",
				Name:    "author",
			},
			{
				Content: "How to start freelancing book",
				Name:    "twitter:title",
			},
			{
				Content: "Want to start freelancing? Kickstart your journey to working on your own terms in the best possible way by signing up for the How To Start Freelancing Book.",
				Name:    "twitter:description",
			},
			{
				Content: "freelancing, contracting, developer, software engineer, book, tutorial",
				Name:    "keywords",
			},
		},
	},
		csrf.Token(ctx.Request()),
	).
		Render(views.ExtractRenderDeps(ctx))
}

func (a App) RenderModal(ctx echo.Context) error {
	return views.SubscribeModal(
		csrf.Token(
			ctx.Request(),
		),
		ctx.QueryParam("article-name"),
	).Render(views.ExtractRenderDeps(ctx))
}

type VerificationToken struct {
	Token string `query:"token"`
}

func (a App) SubscriberEmailVerification(ctx echo.Context) error {
	var payload VerificationToken
	if err := ctx.Bind(&payload); err != nil {
		ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
		ctx.Response().Writer.Header().Add("PreviousLocation", "/user/create")

		return a.base.InternalError(ctx)
	}

	if err := a.tokenService.Validate(ctx.Request().Context(), payload.Token, services.ScopeEmailVerification); err != nil {
		return err
	}

	subscriberID, err := a.tokenService.GetAssociatedSubscriberID(
		ctx.Request().Context(),
		payload.Token,
	)
	if err != nil {
		return err
	}

	if err := a.subscriberSvc.Verify(ctx.Request().Context(), subscriberID); err != nil {
		return err
	}

	if err := a.tokenService.Delete(ctx.Request().Context(), payload.Token); err != nil {
		return err
	}

	// hashedToken := tknManager.Hash(tkn.Token)
	//
	// token, err := db.QuerySubscriberTokenByHash(ctx.Request().Context(), hashedToken)
	// if err != nil {
	// 	if errors.Is(err, pgx.ErrNoRows) {
	// 		return authentication.VerifyEmailPage(true, views.Head{}).
	// 			Render(views.ExtractRenderDeps(ctx))
	// 	}
	//
	// 	ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
	// 	ctx.Response().Writer.Header().Add("PreviousLocation", "/login")
	//
	// 	slog.ErrorContext(ctx.Request().Context(), "could not query subscriber token", "error", err)
	// 	return misc.InternalError(ctx)
	// }
	//
	// if database.ConvertFromPGTimestamptzToTime(token.ExpiresAt).Before(time.Now()) &&
	// 	token.Scope != tokens.ScopeEmailVerification {
	// 	return authentication.VerifyEmailPage(true, views.Head{}).
	// 		Render(views.ExtractRenderDeps(ctx))
	// }
	//
	// if err := db.ConfirmSubscriberEmail(ctx.Request().Context(), database.ConfirmSubscriberEmailParams{
	// 	ID:        token.SubscriberID,
	// 	UpdatedAt: database.ConvertToPGTimestamptz(time.Now()),
	// }); err != nil {
	// 	ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
	// 	ctx.Response().Writer.Header().Add("PreviousLocation", "/login")
	//
	// 	slog.ErrorContext(ctx.Request().Context(), "could not confirm email", "error", err)
	// 	return misc.InternalError(ctx)
	// }

	return authentication.VerifySubscriberEmailPage(false, views.Head{}).
		Render(views.ExtractRenderDeps(ctx))
}

func (a App) SubscriberUnsub(ctx echo.Context) error {
	var payload VerificationToken
	if err := ctx.Bind(&payload); err != nil {
		ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
		ctx.Response().Writer.Header().Add("PreviousLocation", "/user/create")

		return a.base.InternalError(ctx)
	}

	subscriberID, err := a.tokenService.GetAssociatedSubscriberID(
		ctx.Request().Context(),
		payload.Token,
	)
	if err != nil {
		return err
	}

	if err := a.subscriberSvc.Delete(ctx.Request().Context(), subscriberID); err != nil {
		return err
	}

	if err := a.tokenService.Delete(ctx.Request().Context(), payload.Token); err != nil {
		return err
	}

	return ctx.String(http.StatusOK, "You've been unsubscribed")
}

func (a App) Article(ctx echo.Context) error {
	postSlug := ctx.Param("postSlug")

	post, err := a.articleSvc.BySlug(ctx.Request().Context(), postSlug)
	if err != nil {
		return err
	}

	postContent, err := a.postManager.Parse(post.Filename)
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
		Slug:        a.base.BuildURLFromSlug("posts/" + post.Slug),
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
