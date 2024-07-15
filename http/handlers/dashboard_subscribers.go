package handlers

import (
	"net/http"

	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/components"
	"github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

func (d Dashboard) SubscribersIndex(ctx echo.Context) error {
	page := ctx.QueryParam("page")
	pageLimit := 7

	offset, currentPage, err := d.base.GetOffsetAndCurrPage(page, pageLimit)
	if err != nil {
		return err
	}

	subscribers, err := d.subscriberSvc.List(
		ctx.Request().Context(),
		int32(offset),
		int32(pageLimit),
	)
	if err != nil {
		return err
	}

	count, err := d.subscriberSvc.Count(ctx.Request().Context())
	if err != nil {
		return err
	}

	monthlySubscriberCount, err := d.subscriberSvc.NewForCurrentMonth(ctx.Request().Context())
	if err != nil {
		return err
	}

	unverifiedSubCount, err := d.subscriberSvc.UnverifiedCount(ctx.Request().Context())
	if err != nil {
		return err
	}

	viewData := make([]dashboard.SubscriberViewData, 0, len(subscribers))
	for _, sub := range subscribers {
		viewData = append(viewData, dashboard.SubscriberViewData{
			Email:        sub.Email,
			ID:           sub.ID.String(),
			Verified:     sub.IsVerified,
			SubscribedAt: sub.SubscribedAt.String(),
			Refererer:    sub.Referer,
		})
	}

	pagination := components.PaginationProps{
		CurrentPage: currentPage,
		TotalPages:  d.base.CalculateNumberOfPages(int(count), 7),
	}

	return dashboard.Subscribers(int(count), int(len(monthlySubscriberCount)), int(unverifiedSubCount), viewData, pagination, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func (d Dashboard) ResendVerificationMail(ctx echo.Context) error {
	subscriberID := ctx.Param("id")

	subscriberUUID, err := uuid.Parse(subscriberID)
	if err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	subscriber, err := d.subscriberSvc.ByID(ctx.Request().Context(), subscriberUUID)
	if err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	if subscriber.IsVerified {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	// if err := d.tokenService.DeleteSubscriberToken(ctx.Request().Context(), subscriberUUID); err != nil {
	// 	return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	// }

	activationToken, err := d.tokenService.CreateSubscriberEmailValidation(
		ctx.Request().Context(),
		subscriber.ID,
	)
	if err != nil {
		return err
	}

	unsubToken, err := d.tokenService.CreateUnsubscribeToken(
		ctx.Request().Context(),
		subscriber.ID,
	)
	if err != nil {
		return err
	}

	if err := d.emailService.SendNewSubscriberEmail(
		ctx.Request().Context(),
		subscriber.Email,
		activationToken,
		unsubToken,
	); err != nil {
		return err
	}

	return dashboard.SuccessMsg("Verification mail send").Render(views.ExtractRenderDeps(ctx))
}

func (d Dashboard) DeleteSubscriber(ctx echo.Context) error {
	id := ctx.Param("ID")

	parsedID := uuid.MustParse(id)

	if err := d.subscriberSvc.Delete(ctx.Request().Context(), parsedID); err != nil {
		return err
	}

	ctx.Response().Header().Add("HX-Refresh", "true")

	return ctx.Redirect(http.StatusSeeOther, "/dashboard/subscribers")
}
