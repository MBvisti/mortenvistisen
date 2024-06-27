package app

import (
	"errors"
	"net/http"

	"github.com/MBvisti/mortenvistisen/controllers"
	"github.com/MBvisti/mortenvistisen/controllers/misc"
	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/authentication"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

func Index(ctx echo.Context, articleModel models.ArticleService) error {
	data, err := articleModel.List(ctx.Request().Context(), 0, 5)
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
			Slug:        controllers.FormatArticleSlug(d.Slug),
		})
	}

	return views.HomePage(posts).Render(views.ExtractRenderDeps(ctx))
}

func About(ctx echo.Context) error {
	return views.AboutPage(views.Head{
		Title:       "About",
		Description: "Contains information about the site owner, Morten Vistisen, the purpose of the site and the technologies used to build it.",
		Slug:        controllers.BuildURLFromSlug("about"),
		Image:       "https://mortenvistisen.com/static/images/mbv.png",
		MetaType:    "website",
	}).Render(views.ExtractRenderDeps(ctx))
}

func Newsletter(ctx echo.Context) error {
	return views.NewsletterPage(views.Head{
		Title:       "Newsletter",
		Description: "Signup page for joining Morten's newsletter where he shares his thoughts on software development, business and life.",
		Slug:        controllers.BuildURLFromSlug("newsletter"),
		MetaType:    "website",
		Image:       "https://mortenvistisen.com/static/images/mbv.png",
	}, csrf.Token(ctx.Request())).Render(views.ExtractRenderDeps(ctx))
}

func Projects(ctx echo.Context) error {
	return views.ProjectsPage(views.Head{
		Title:       "Projects",
		Description: "A collection of on-going and retired projects I've build over the years. This includes business projects, open source projects and personal projects.",
		Slug:        controllers.BuildURLFromSlug("projects"),
		MetaType:    "website",
		Image:       "https://mortenvistisen.com/static/images/mbv.png",
	}).Render(views.ExtractRenderDeps(ctx))
}

type SubscriptionEventForm struct {
	Email string `form:"hero-input"`
	Title string `form:"article-title"`
}

func SubscriptionEvent(
	ctx echo.Context,
	subscriberSvc models.SubscriberService,
) error {
	var form SubscriptionEventForm
	if err := ctx.Bind(&form); err != nil {
		// if err := mail.Send(ctx.Request().Context(), "hi@mortenvistisen.com", "sub-blog@mortenvistisen.com",
		// 	"Failed to subscribe", "sub_report", err.Error()); err != nil {
		// 	telemetry.Logger.Error("Failed to send email", "error", err)
		// }
		return ctx.String(200, "You're now subscribed!")
	}

	// sub, err := db.InsertSubscriber(
	// 	ctx.Request().Context(),
	// 	database.InsertSubscriberParams{
	// 		ID:        uuid.New(),
	// 		CreatedAt: database.ConvertToPGTimestamptz(time.Now()),
	// 		UpdatedAt: database.ConvertToPGTimestamptz(time.Now()),
	// 		Email: sql.NullString{
	// 			String: form.Email,
	// 			Valid:  true,
	// 		},
	// 		SubscribedAt: database.ConvertToPGTimestamptz(time.Now()),
	// 		Referer: sql.NullString{
	// 			String: form.Title,
	// 			Valid:  true,
	// 		},
	// 		IsVerified: pgtype.Bool{Bool: false, Valid: true},
	// 	},
	// )
	// if err != nil {
	// 	return err
	// }
	//
	// generatedTkn, err := tknManager.GenerateToken()
	// if err != nil {
	// 	return err
	// }
	//
	// activationToken := tokens.CreateActivationToken(
	// 	generatedTkn.PlainTextToken,
	// 	generatedTkn.HashedToken,
	// )
	//
	// if err := db.InsertSubscriberToken(ctx.Request().Context(), database.InsertSubscriberTokenParams{
	// 	ID:           uuid.New(),
	// 	CreatedAt:    database.ConvertToPGTimestamptz(time.Now()),
	// 	Hash:         activationToken.Hash,
	// 	ExpiresAt:    database.ConvertToPGTimestamptz(activationToken.GetExpirationTime()),
	// 	Scope:        activationToken.GetScope(),
	// 	SubscriberID: sub.ID,
	// }); err != nil {
	// 	return err
	// }
	//
	// newsletterMail := templates.NewsletterWelcomeMail{
	// 	ConfirmationLink: fmt.Sprintf(
	// 		"%s://%s/verify-subscriber?token=%s",
	// 		cfg.App.AppScheme,
	// 		cfg.App.AppHost,
	// 		activationToken.GetPlainText(),
	// 	),
	// 	UnsubscribeLink: "",
	// }
	// textVersion, err := newsletterMail.GenerateTextVersion()
	// if err != nil {
	// 	return err
	// }
	//
	// htmlVersion, err := newsletterMail.GenerateHtmlVersion()
	// if err != nil {
	// 	return err
	// }
	// _, err = queueClient.Insert(ctx.Request().Context(), queue.EmailJobArgs{
	// 	To:          form.Email,
	// 	From:        "noreply@mortenvistisen.com",
	// 	Subject:     "Thanks for signing up!",
	// 	TextVersion: textVersion,
	// 	HtmlVersion: htmlVersion,
	// }, nil)
	// if err != nil {
	// 	return err
	// }

	if err := subscriberSvc.New(ctx.Request().Context(), form.Email, form.Title); err != nil {
		if errors.Is(err, models.ErrValidation{}) {
			return nil
		}
		if errors.Is(err, models.ErrUnrecoverableEvent) {
			return ctx.JSON(http.StatusInternalServerError, "yo")
		}
	}

	return views.SubscribeModalResponse().
		Render(views.ExtractRenderDeps(ctx))
}

func RenderModal(ctx echo.Context) error {
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

func SubscriberEmailVerification(
	ctx echo.Context,
	subscriberModel models.SubscriberService,
	tokenService services.TokenSvc,
) error {
	var tkn VerificationToken
	if err := ctx.Bind(&tkn); err != nil {
		ctx.Response().Writer.Header().Add("HX-Redirect", "/500")
		ctx.Response().Writer.Header().Add("PreviousLocation", "/user/create")

		return misc.InternalError(ctx)
	}

	if err := tokenService.ValidateSubscriber(ctx.Request().Context(), tkn.Token); err != nil {
		return err
	}

	subscriberID, err := tokenService.GetAssociatedSubscriberID(ctx.Request().Context(), tkn.Token)
	if err != nil {
		return err
	}

	if err := subscriberModel.Verify(ctx.Request().Context(), subscriberID); err != nil {
		return err
	}

	if err := tokenService.DeleteSubscriberToken(ctx.Request().Context(), subscriberID); err != nil {
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
