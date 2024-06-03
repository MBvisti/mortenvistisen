package dashboard

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MBvisti/mortenvistisen/controllers"
	"github.com/MBvisti/mortenvistisen/pkg/config"
	"github.com/MBvisti/mortenvistisen/pkg/mail/templates"
	"github.com/MBvisti/mortenvistisen/pkg/queue"
	"github.com/MBvisti/mortenvistisen/pkg/tokens"
	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/components"
	"github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/riverqueue/river"
)

func ResendVerificationMail(
	ctx echo.Context,
	db database.Queries,
	tknManager tokens.Manager,
	queueClient *river.Client[pgx.Tx],
	cfg config.Cfg,
) error {
	subscriberID := ctx.Param("id")

	subscriberUUID, err := uuid.Parse(subscriberID)
	if err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	subscriber, err := db.QuerySubscriber(ctx.Request().Context(), subscriberUUID)
	if err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	if subscriber.IsVerified.Bool {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	if err := db.DeleteSubscriberTokenBySubscriberID(ctx.Request().Context(), subscriberUUID); err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	generatedTkn, err := tknManager.GenerateToken()
	if err != nil {
		return dashboard.FailureMsg("Could not mail send").Render(views.ExtractRenderDeps(ctx))
	}

	activationToken := tokens.CreateActivationToken(
		generatedTkn.PlainTextToken,
		generatedTkn.HashedToken,
	)

	if err := db.InsertSubscriberToken(ctx.Request().Context(), database.InsertSubscriberTokenParams{
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
			cfg.App.AppScheme,
			cfg.App.AppHost,
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

	_, err = queueClient.Insert(ctx.Request().Context(), queue.EmailJobArgs{
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

func SubscribersIndex(ctx echo.Context, db database.Queries) error {
	page := ctx.QueryParam("page")

	var currentPage int
	if page == "" {
		currentPage = 1
	}
	if page != "" {
		cp, err := strconv.Atoi(page)
		if err != nil {
			return err
		}

		currentPage = cp
	}

	offset := 0
	if currentPage == 2 {
		offset = 7
	}

	if currentPage > 2 {
		offset = 7 * (currentPage - 1)
	}

	subs, err := db.QuerySubscribersInPages(
		ctx.Request().Context(),
		database.QuerySubscribersInPagesParams{
			Limit:  7,
			Offset: int32(offset),
		},
	)
	if err != nil {
		return err
	}

	subsCount, err := db.QuerySubscriberCount(ctx.Request().Context())
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

	pagination := components.PaginationProps{
		CurrentPage: currentPage,
		TotalPages:  controllers.CalculateNumberOfPages(int(subsCount), 7),
	}

	return dashboard.Subscribers(viewData, pagination, csrf.Token(ctx.Request())).
		Render(views.ExtractRenderDeps(ctx))
}

func DeleteSubscriber(ctx echo.Context, db database.Queries) error {
	id := ctx.Param("ID")

	parsedID := uuid.MustParse(id)

	if err := db.DeleteSubscriber(ctx.Request().Context(), parsedID); err != nil {
		return err
	}

	ctx.Response().Header().Add("HX-Refresh", "true")

	return ctx.Redirect(http.StatusSeeOther, "/dashboard/subscribers")
}
