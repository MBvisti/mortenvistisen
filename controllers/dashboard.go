package controllers

import (
	"net/http"

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

	ctx.Response().Header().Add("HX-Refresh", "true")

	return ctx.Redirect(http.StatusSeeOther, "/dashboard/subscribers")
}

func (c *Controller) DashboardArticles(ctx echo.Context) error {
	return dashboard.Articles().Render(views.ExtractRenderDeps(ctx))
}
