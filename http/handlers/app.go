package handlers

import (
	"errors"
	"log/slog"
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

func (a App) SubscriptionEvent(c echo.Context) error {
	ctx, subscriptionEventSpan := a.base.Tracer.CreateSpan(c.Request().Context(), "article/store")
	subscriptionEventSpan.AddEvent("AppHandler/SubscriptionEvent")

	type subscriptionEventForm struct {
		Email string `form:"hero-input"`
		Title string `form:"article-title"`
	}

	var form subscriptionEventForm
	if err := c.Bind(&form); err != nil {
		slog.ErrorContext(ctx, "could not bind subscriptionEventForm", "error", err)
		return c.String(200, "You're now subscribed!")
	}
	slog.InfoContext(ctx, "starting to create new subscriber from handler", "email", form.Email)

	param := c.QueryParam("book")
	var bookSub bool
	if param == "true" {
		bookSub = true
	}

	newSubscriberCtx, newSubscriberSpan := a.base.Tracer.CreateChildSpan(
		ctx,
		subscriptionEventSpan,
		"subscriberService/New",
	)
	newSubscriberSpan.AddEvent("New/Start")
	if err := a.subscriberSvc.New(newSubscriberCtx, form.Email, form.Title, bookSub); err != nil {
		if errors.Is(err, models.ErrSubscriberExists) {
			return views.SubscribeModalResponse(bookSub, true).
				Render(views.ExtractRenderDeps(c))
		}
		if errors.Is(err, models.ErrValidation{}) {
			return nil
		}
		if errors.Is(err, models.ErrUnrecoverableEvent) {
			return c.JSON(http.StatusInternalServerError, "yo")
		}
	}
	newSubscriberSpan.End()

	subscriptionEventSpan.End()
	return views.SubscribeModalResponse(bookSub, false).
		Render(views.ExtractRenderDeps(c))
}

func (a App) HowToStartFreelancing(ctx echo.Context) error {
	return views.HowToStartFreelancing(views.Head{
		Title:       "How To Start Freelancing Book",
		Description: "Want to start freelancing? Kickstart your journey to working on your own terms in the best possible way by signing up for the How To Start Freelancing Book.",
		Slug:        a.base.BuildURLFromSlug("books/" + "how-to-start-freelancing"),
		MetaType:    "website",
		Image:       "https://mortenvistisen.com/static/images/book-cover-ebook-wip.png",
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

func (a App) SubscriberEmailVerification(c echo.Context) error {
	ctx, span := a.base.Tracer.CreateSpan(
		c.Request().Context(),
		"AppHandler/SubscriberEmailVerification",
	)
	span.AddEvent("SubscriberEmailVerification/start")

	var payload VerificationToken
	if err := c.Bind(&payload); err != nil {
		slog.ErrorContext(ctx, "could not bind verification token", "error", err)
		c.Response().Writer.Header().Add("HX-Redirect", "/500")
		c.Response().Writer.Header().Add("PreviousLocation", "/user/create")

		return a.base.InternalError(c)
	}

	if err := a.tokenService.Validate(ctx, payload.Token, services.ScopeEmailVerification); err != nil {
		return err
	}

	subscriberID, err := a.tokenService.GetAssociatedSubscriberID(
		ctx,
		payload.Token,
	)
	if err != nil {
		return err
	}

	subVerifyCtx, subVerifySpan := a.base.Tracer.CreateChildSpan(
		ctx,
		span,
		"SubscriberService/Verify",
	)
	subVerifySpan.AddEvent("Verify/Start")
	if err := a.subscriberSvc.Verify(subVerifyCtx, subscriberID); err != nil {
		return err
	}
	subVerifySpan.End()

	tokenDeleteCtx, tokenDeleteSpan := a.base.Tracer.CreateChildSpan(
		ctx,
		span,
		"TokenService/Delete",
	)
	tokenDeleteSpan.AddEvent("Delete/Start")
	if err := a.tokenService.Delete(tokenDeleteCtx, payload.Token); err != nil {
		return err
	}
	tokenDeleteSpan.End()

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

	span.End()
	return authentication.VerifySubscriberEmailPage(false, views.Head{}).
		Render(views.ExtractRenderDeps(c))
}

func (a App) SubscriberUnsub(c echo.Context) error {
	ctx, span := a.base.Tracer.CreateSpan(
		c.Request().Context(),
		"AppHandler/SubscriberUnsub",
	)
	span.AddEvent("SubscriberUnsub/start")
	var payload VerificationToken
	if err := c.Bind(&payload); err != nil {
		slog.ErrorContext(ctx, "could not bind verification token", "error", err)
		c.Response().Writer.Header().Add("HX-Redirect", "/500")
		c.Response().Writer.Header().Add("PreviousLocation", "/user/create")

		return a.base.InternalError(c)
	}

	subscriberID, err := a.tokenService.GetAssociatedSubscriberID(
		ctx,
		payload.Token,
	)
	if err != nil {
		return err
	}

	subscriberDeleteCtx, subscriberDeleteSpan := a.base.Tracer.CreateChildSpan(
		ctx,
		span,
		"SubscriberService/Delete",
	)
	subscriberDeleteSpan.AddEvent("Delete/Start")
	if err := a.subscriberSvc.Delete(subscriberDeleteCtx, subscriberID); err != nil {
		return err
	}

	tokenDeleteCtx, tokenDeleteSpan := a.base.Tracer.CreateChildSpan(
		ctx,
		span,
		"TokenService/Delete",
	)
	tokenDeleteSpan.AddEvent("Delete/Start")
	if err := a.tokenService.Delete(tokenDeleteCtx, payload.Token); err != nil {
		return err
	}
	tokenDeleteSpan.End()

	span.End()
	return c.String(http.StatusOK, "You've been unsubscribed")
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

	fiveRandomPosts, err := a.base.DB.GetFiveRandomPosts(
		ctx.Request().Context(),
		post.ID,
	)
	if err != nil {
		return err
	}

	otherArticles := make(map[string]string, 5)

	for _, article := range fiveRandomPosts {
		otherArticles[article.Title] = a.base.BuildURLFromSlug(
			"posts/" + article.Slug,
		)
	}

	return views.ArticlePage(views.ArticlePageData{
		Content:           postContent,
		HeaderTitle:       post.HeaderTitle,
		ReleaseDate:       post.ReleaseDate,
		OtherArticleLinks: otherArticles,
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
