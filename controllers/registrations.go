package controllers

import (
	"log/slog"
	"net/http"

	"mortenvistisen/config"
	"mortenvistisen/internal/storage"
	"mortenvistisen/queue"
	"mortenvistisen/router/cookies"
	"mortenvistisen/router/routes"
	"mortenvistisen/services"
	"mortenvistisen/views"

	"github.com/labstack/echo/v5"

	"mortenvistisen/internal/hypermedia"
)

type Registrations struct {
	db         storage.Pool
	insertOnly queue.InsertOnly
	cfg        config.Config
}

func NewRegistrations(
	db storage.Pool,
	insertOnly queue.InsertOnly,
	cfg config.Config,
) Registrations {
	return Registrations{db, insertOnly, cfg}
}

func (r Registrations) New(etx *echo.Context) error {
	return render(etx, views.RegistrationForm())
}

func (r Registrations) Create(etx *echo.Context) error {
	var payload struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse signup form payload",
			"error",
			err,
		)
		return render(etx, views.BadRequest())
	}

	if err := services.RegisterUser(
		etx.Request().Context(),
		r.db,
		r.insertOnly,
		r.cfg.Auth.Pepper,
		services.RegisterUserData{
			Email:           payload.Email,
			Password:        payload.Password,
			ConfirmPassword: payload.ConfirmPassword,
		},
	); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to register user",
			"error",
			err,
		)

		if flashErr := cookies.AddFlash(etx, cookies.FlashError, "Failed to register user"); flashErr != nil {
			return render(etx, views.InternalError())
		}

		return etx.Redirect(http.StatusSeeOther, routes.RegistrationNew.URL())
	}

	return hypermedia.Redirect(etx, routes.ConfirmationNew.URL())
}
