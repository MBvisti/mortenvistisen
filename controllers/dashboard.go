package controllers

import (
	"fmt"

	"github.com/MBvisti/mortenvistisen/pkg/mail"
	"github.com/MBvisti/mortenvistisen/pkg/telemetry"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

func (c *Controller) DashboardIndex(ctx echo.Context) error {
	return dashboard.Index().Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardSubscribers(ctx echo.Context) error {
	subs, err := c.db.QueryAllSubscribers(ctx.Request().Context())
	if err != nil {
		return err
	}

	viewData := make([]dashboard.SubscriberViewData, 0, len(subs))
	for _, sub := range subs {
		viewData = append(viewData, dashboard.SubscriberViewData{
			Email:    sub.Email.String,
			ID:       sub.ID.String(),
			Verified: sub.IsVerified.Bool,
		})
	}

	return dashboard.Subscribers(viewData, csrf.Token(ctx.Request())).Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DeleteSubscriber(ctx echo.Context) error {
	id := ctx.Param("ID")
	parsedID := uuid.MustParse(id)

	if err := c.db.DeleteSubscriber(ctx.Request().Context(), parsedID); err != nil {
		return err
	}

	subs, err := c.db.QueryAllSubscribers(ctx.Request().Context())
	if err != nil {
		return err
	}

	viewData := make([]dashboard.SubscriberViewData, 0, len(subs))
	for _, sub := range subs {
		viewData = append(viewData, dashboard.SubscriberViewData{
			Email:    sub.Email.String,
			ID:       sub.ID.String(),
			Verified: sub.IsVerified.Bool,
		})
	}

	return dashboard.SubscriberTable(viewData).Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardArticles(ctx echo.Context) error {
	articles, err := c.db.QueryAllPost(ctx.Request().Context())
	if err != nil {
		return err
	}

	viewData := make([]dashboard.ArticleViewData, 0, len(articles))
	for _, article := range articles {
		viewData = append(viewData, dashboard.ArticleViewData{
			Slug:    article.Slug,
			Title:   article.Title,
			Tags:    article.Tags,
			InDraft: article.Draft,
		})
	}

	sess, err := flashStore.Get(ctx.Request(), "flash")
	if err != nil {
		return err
	}

	var flashesToDisplay []string
	if flashes := sess.Flashes("NotifyStatusMessage"); len(flashes) > 0 {
		for _, flash := range flashes {
			f, ok := (flash).(string)
			if !ok {
				return err
			}
			flashesToDisplay = append(flashesToDisplay, f)
		}
	}

	return dashboard.Articles(viewData, "", flashesToDisplay).Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardArticleDetails(ctx echo.Context) error {
	article, err := c.db.QueryPostBySlug(ctx.Request().Context(), ctx.Param("slug"))
	if err != nil {
		return err
	}

	viewData := dashboard.ArticleDetailsViewData{
		ID:      article.ID.String(),
		Title:   article.Title,
		Excerpt: article.Excerpt,
		InDraft: article.Draft,
		Tags:    article.Tags,
		Slug:    article.Slug,
	}

	csrftoken := csrf.Token(ctx.Request())

	return dashboard.ArticleDetails(viewData, csrftoken).Render(views.ExtractRenderDeps(ctx))
}

func (c *Controller) DashboardNotifySubscribers(ctx echo.Context) error {
	ctx.Response().Header().Set("HX-Redirect", "/dashboard/articles")

	article, err := c.db.QueryPostBySlug(ctx.Request().Context(), ctx.Param("slug"))
	if err != nil {
		return err
	}

	subs, err := c.db.QueryAllSubscribers(ctx.Request().Context())
	if err != nil {
		return err
	}

	var emailList []string
	for _, sub := range subs {
		if sub.IsVerified.Bool {
			emailList = append(emailList, sub.Email.String)
		}
	}

	sess, err := flashStore.Get(ctx.Request(), "flash")
	if err != nil {
		return err
	}
	sess.Options.MaxAge = 5

	for _, email := range emailList {
		if err := c.mail.Send(ctx.Request().Context(), email, "newsletter@mortenvistisen.com", "MBV Newsletter: I just released a new article", "notify_subscribers", mail.ArticleNotification{
			Title: article.Title,
			Slug:  c.buildURLFromSlug("posts/" + article.Slug),
			Email: email,
		}); err != nil {
			telemetry.Logger.Error("failed to send article notification", "error", err, "email", email)

			sess.AddFlash("Could not notify subscribers", "NotifyStatusMessage")
			sess.AddFlash(fmt.Sprintf("error: %v", err), "NotifyStatusMessage")

			if err := sess.Save(ctx.Request(), ctx.Response()); err != nil {
				return err
			}

			return nil
		}
	}

	sess.AddFlash("All subscribers have been notified", "NotifyStatusMessage")
	if err := sess.Save(ctx.Request(), ctx.Response()); err != nil {
		return err
	}

	return nil
}
