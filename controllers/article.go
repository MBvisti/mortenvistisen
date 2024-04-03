package controllers

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/MBvisti/mortenvistisen/pkg/mail/templates"
	"github.com/MBvisti/mortenvistisen/pkg/queue"
	"github.com/MBvisti/mortenvistisen/pkg/telemetry"
	"github.com/MBvisti/mortenvistisen/pkg/tokens"
	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func (c *Controller) Article(ctx echo.Context) error {
	postSlug := ctx.Param("postSlug")

	post, err := c.db.GetPostBySlug(ctx.Request().Context(), postSlug)
	if err != nil {
		return err
	}

	postContent, err := c.postManager.Parse(post.Filename)
	if err != nil {
		return err
	}

	tags, err := c.db.GetTagsForPost(ctx.Request().Context(), post.ID)
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

	fiveRandomPosts, err := c.db.GetFiveRandomPosts(
		ctx.Request().Context(),
		post.ID,
	)
	if err != nil {
		return err
	}

	otherArticles := make(map[string]string, 5)

	for _, article := range fiveRandomPosts {
		otherArticles[article.Title] = c.buildURLFromSlug(
			"posts/" + article.Slug,
		)
	}

	return views.ArticlePage(views.ArticlePageData{
		Content:           postContent,
		Title:             post.Title,
		ReleaseDate:       post.ReleasedAt.Time,
		OtherArticleLinks: otherArticles,
		CsrfToken:         csrf.Token(ctx.Request()),
	}, views.Head{
		Title:       post.Title,
		Description: post.Excerpt,
		Slug:        c.buildURLFromSlug("posts/" + post.Slug),
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

type SubscriptionEventForm struct {
	Email string `form:"hero-input"`
	Title string `form:"article-title"`
}

func (c *Controller) SubscriptionEvent(ctx echo.Context) error {
	var form SubscriptionEventForm
	if err := ctx.Bind(&form); err != nil {
		if err := c.mail.Send(ctx.Request().Context(), "hi@mortenvistisen.com", "sub-blog@mortenvistisen.com",
			"Failed to subscribe", "sub_report", err.Error()); err != nil {
			telemetry.Logger.Error("Failed to send email", "error", err)
		}
		return ctx.String(200, "You're now subscribed!")
	}

	sub, err := c.db.InsertSubscriber(
		ctx.Request().Context(),
		database.InsertSubscriberParams{
			ID:        uuid.New(),
			CreatedAt: database.ConvertToPGTimestamptz(time.Now()),
			UpdatedAt: database.ConvertToPGTimestamptz(time.Now()),
			Email: sql.NullString{
				String: form.Email,
				Valid:  true,
			},
			SubscribedAt: database.ConvertToPGTimestamptz(time.Now()),
			Referer: sql.NullString{
				String: form.Title,
				Valid:  true,
			},
			IsVerified: pgtype.Bool{Bool: false, Valid: true},
		},
	)
	if err != nil {
		return err
	}

	generatedTkn, err := c.tknManager.GenerateToken()
	if err != nil {
		return err
	}

	activationToken := tokens.CreateActivationToken(
		generatedTkn.PlainTextToken,
		generatedTkn.HashedToken,
	)

	if err := c.db.InsertSubscriberToken(ctx.Request().Context(), database.InsertSubscriberTokenParams{
		ID:           uuid.New(),
		CreatedAt:    database.ConvertToPGTimestamptz(time.Now()),
		Hash:         activationToken.Hash,
		ExpiresAt:    database.ConvertToPGTimestamptz(activationToken.GetExpirationTime()),
		Scope:        activationToken.GetScope(),
		SubscriberID: sub.ID,
	}); err != nil {
		return err
	}

	newsletterMail := templates.NewsletterWelcomeMail{
		ConfirmationLink: fmt.Sprintf(
			"%s://%s/verify-subscriber?token=%s",
			c.cfg.App.AppScheme,
			c.cfg.App.AppHost,
			activationToken.GetPlainText(),
		),
		UnsubscribeLink: "",
	}
	textVersion, err := newsletterMail.GenerateTextVersion()
	if err != nil {
		return err
	}

	htmlVersion, err := newsletterMail.GenerateHtmlVersion()
	if err != nil {
		return err
	}
	_, err = c.queueClient.Insert(ctx.Request().Context(), queue.EmailJobArgs{
		To:          form.Email,
		From:        c.cfg.App.DefaultSenderSignature,
		Subject:     "Thanks for signing up!",
		TextVersion: textVersion,
		HtmlVersion: htmlVersion,
	}, nil)
	if err != nil {
		return err
	}

	isModal := ctx.QueryParam("is-modal")
	if isModal == "true" {
		return views.SubscribeModalResponse().
			Render(views.ExtractRenderDeps(ctx))
	} else {
		return ctx.String(200, "You're now subscribed!")
	}
}

func (c *Controller) RenderModal(ctx echo.Context) error {
	return views.SubscribeModal(
		csrf.Token(
			ctx.Request(),
		),
		ctx.QueryParam("article-name"),
	).Render(views.ExtractRenderDeps(ctx))
}
