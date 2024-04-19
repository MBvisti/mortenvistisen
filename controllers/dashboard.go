package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MBvisti/mortenvistisen/pkg/mail/templates"
	"github.com/MBvisti/mortenvistisen/pkg/queue"
	"github.com/MBvisti/mortenvistisen/pkg/tokens"
	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

func (c *Controller) DashboardIndex(ctx echo.Context) error {
	return dashboard.Index().Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardResendVerificationMail(ctx echo.Context) error {
	subscriberID := ctx.Param("id")

	subscriberUUID, err := uuid.Parse(subscriberID)
	if err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	subscriber, err := c.db.QuerySubscriber(ctx.Request().Context(), subscriberUUID)
	if err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	if subscriber.IsVerified.Bool {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	if err := c.db.DeleteSubscriberTokenBySubscriberID(ctx.Request().Context(), subscriberUUID); err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	generatedTkn, err := c.tknManager.GenerateToken()
	if err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
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
		SubscriberID: subscriberUUID,
	}); err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
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
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	htmlVersion, err := newsletterMail.GenerateHtmlVersion()
	if err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	_, err = c.queueClient.Insert(ctx.Request().Context(), queue.EmailJobArgs{
		To:          subscriber.Email.String,
		From:        "noreply@mortenvistisen.com",
		Subject:     "Thanks for signing up!",
		TextVersion: textVersion,
		HtmlVersion: htmlVersion,
	}, nil)
	if err != nil {
		return err
	}

	return dashboard.SuccessMsg("Verification mail send").Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardSubscribers(ctx echo.Context) error {
	subs, err := c.db.QueryAllSubscribers(ctx.Request().Context())
	if err != nil {
		return err
	}

	viewData := make([]dashboard.SubscriberViewData, 0, len(subs))
	for _, sub := range subs {
		viewData = append(viewData, dashboard.SubscriberViewData{
			Email:        sub.Email.String,
			ID:           sub.ID.String(),
			Verified:     sub.IsVerified.Bool,
			SubscribedAt: sub.SubscribedAt.Time.String(),
			Refererer:    sub.Referer.String,
		})
	}

	return dashboard.Subscribers(viewData, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardArticles(ctx echo.Context) error {
	articles, err := c.db.QueryAllPosts(ctx.Request().Context())
	if err != nil {
		return err
	}

	viewData := make([]dashboard.ArticleViewData, 0, len(articles))
	for _, article := range articles {
		viewData = append(viewData, dashboard.ArticleViewData{
			ID:         article.ID.String(),
			Title:      article.Title,
			Draft:      article.Draft,
			ReleasedAt: article.ReleasedAt.Time.String(),
			Slug:       article.Slug,
		})
	}

	return dashboard.Articles(viewData, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardArticleEdit(ctx echo.Context) error {
	slug := ctx.Param("slug")

	post, err := c.db.GetPostBySlug(ctx.Request().Context(), slug)
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

	return dashboard.ArticleEdit(dashboard.ArticleEditViewData{
		ID:          post.ID.String(),
		CreatedAt:   post.CreatedAt.Time,
		UpdatedAt:   post.UpdatedAt.Time,
		Title:       post.Title,
		HeaderTitle: post.HeaderTitle.String,
		Filename:    post.Filename,
		Slug:        post.Slug,
		Excerpt:     post.Excerpt,
		Draft:       post.Draft,
		ReleasedAt:  post.ReleasedAt.Time,
		ReadTime:    post.ReadTime.Int32,
		Content:     postContent,
		Keywords:    keywords,
	}, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DeleteSubscriber(ctx echo.Context) error {
	id := ctx.Param("ID")

	parsedID := uuid.MustParse(id)

	if err := c.db.DeleteSubscriber(ctx.Request().Context(), parsedID); err != nil {
		return err
	}

	ctx.Response().Header().Add("HX-Refresh", "true")

	return ctx.Redirect(http.StatusSeeOther, "/dashboard/subscribers")
}
