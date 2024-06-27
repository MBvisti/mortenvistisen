package dashboard

import (
	"net/http"

	"github.com/MBvisti/mortenvistisen/controllers"
	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/repository/psql/database"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/components"
	"github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
)

func SubscribersIndex(
	ctx echo.Context,
	db database.Queries,
	subscriberSvc models.SubscriberService,
) error {
	page := ctx.QueryParam("page")
	pageLimit := 7

	offset, currentPage, err := controllers.GetOffsetAndCurrPage(page, pageLimit)
	if err != nil {
		return err
	}

	subscribers, err := subscriberSvc.List(ctx.Request().Context(), int32(offset), int32(pageLimit))
	if err != nil {
		return err
	}

	count, err := subscriberSvc.Count(ctx.Request().Context())
	if err != nil {
		return err
	}

	monthlySubscriberCount, err := subscriberSvc.NewForCurrentMonth(ctx.Request().Context())
	if err != nil {
		return err
	}

	unverifiedSubCount, err := subscriberSvc.UnverifiedCount(ctx.Request().Context())
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
		TotalPages:  controllers.CalculateNumberOfPages(int(count), 7),
	}

	return dashboard.Subscribers(int(count), int(len(monthlySubscriberCount)), int(unverifiedSubCount), viewData, pagination, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func ResendVerificationMail(
	ctx echo.Context,
	subscriberModel models.SubscriberService,
	tokenService services.TokenSvc,
	emailService services.Email,
) error {
	subscriberID := ctx.Param("id")

	subscriberUUID, err := uuid.Parse(subscriberID)
	if err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	subscriber, err := subscriberModel.ByID(ctx.Request().Context(), subscriberUUID)
	if err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	if subscriber.IsVerified {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	if err := tokenService.DeleteSubscriberToken(ctx.Request().Context(), subscriberUUID); err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	activationToken, err := tokenService.CreateSubscriptionToken(
		ctx.Request().Context(),
		subscriber.ID,
	)
	if err != nil {
		return err
	}

	unsubToken, err := tokenService.CreateUnsubscribeToken(
		ctx.Request().Context(),
		subscriber.ID,
	)
	if err != nil {
		return err
	}

	if err := emailService.SendNewSubscriberEmail(
		ctx.Request().Context(),
		subscriber.Email,
		activationToken,
		unsubToken,
	); err != nil {
		return err
	}

	return dashboard.SuccessMsg("Verification mail send").Render(views.ExtractRenderDeps(ctx))
}

func DeleteSubscriber(ctx echo.Context, subscriberModel models.SubscriberService) error {
	id := ctx.Param("ID")

	parsedID := uuid.MustParse(id)

	if err := subscriberModel.Delete(ctx.Request().Context(), parsedID); err != nil {
		return err
	}

	ctx.Response().Header().Add("HX-Refresh", "true")

	return ctx.Redirect(http.StatusSeeOther, "/dashboard/subscribers")
}
