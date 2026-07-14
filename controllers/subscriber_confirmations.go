package controllers

import (
	"errors"
	"log/slog"
	"net/http"

	"mortenvistisen/internal/inertia"
	"mortenvistisen/internal/validation"
	"mortenvistisen/router"
	"mortenvistisen/router/cookies"
	"mortenvistisen/router/middleware"
	"mortenvistisen/router/routes"
	"mortenvistisen/services"

	"github.com/labstack/echo/v5"
)

type SubscriberConfirmations struct {
	confirmation services.SubscriberConfirmation
}

func NewSubscriberConfirmations(
	confirmation services.SubscriberConfirmation,
) SubscriberConfirmations {
	return SubscriberConfirmations{confirmation: confirmation}
}

func (s SubscriberConfirmations) RegisterRoutes(r *router.Router) error {
	var errs []error

	_, err := r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.SubscriberConfirmationNew.Path(),
		Name:    routes.SubscriberConfirmationNew.Name(),
		Handler: s.New,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodPost,
		Path:    routes.SubscriberConfirmationCreate.Path(),
		Name:    routes.SubscriberConfirmationCreate.Name(),
		Handler: s.Create,
		Middlewares: []echo.MiddlewareFunc{
			middleware.IPRateLimiter(10, routes.SubscriberConfirmationNew),
		},
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (s SubscriberConfirmations) New(etx *echo.Context) error {
	etx.Response().Header().Set("Referrer-Policy", "no-referrer")
	return inertia.Page(etx, "Subscription/ConfirmEmail", inertia.Props{})
}

func (s SubscriberConfirmations) Create(etx *echo.Context) error {
	var payload struct {
		Code string `json:"code"`
	}
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse subscriber confirmation payload",
			"error", err,
		)
		return inertia.Page(etx, "Errors/BadRequest", inertia.Props{})
	}

	_, err := s.confirmation.Confirm(
		etx.Request().Context(),
		services.ConfirmSubscriberEmailData{Code: payload.Code},
	)
	if err != nil {
		if validationErrors, ok := validation.As(err); ok {
			return inertia.Page(
				etx,
				"Subscription/ConfirmEmail",
				inertia.Props{},
				inertia.WithValidationErrors(validationErrors.ToMap()),
			)
		}

		slog.ErrorContext(
			etx.Request().Context(),
			"failed to confirm subscriber email",
			"error", err,
		)

		message := "That confirmation code is invalid"
		if errors.Is(err, services.ErrExpiredSubscriberConfirmation) {
			message = "That confirmation code has expired"
		}
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, message); flashErr != nil {
			return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
		}
		return inertia.Redirect(
			etx,
			routes.SubscriberConfirmationNew.URL(),
			http.StatusSeeOther,
		)
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Your email address has been confirmed",
	); flashErr != nil {
		return inertia.Page(etx, "Errors/InternalError", inertia.Props{})
	}
	return inertia.Location(etx, routes.HomePage.URL())
}
