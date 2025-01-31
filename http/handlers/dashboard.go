package handlers

import (
	"log/slog"
	"strconv"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/psql"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/labstack/echo/v4"
)

type Dashboard struct {
	db psql.Postgres
}

func newDashboard(db psql.Postgres) Dashboard {
	return Dashboard{db}
}

func (d *Dashboard) Index(c echo.Context) error {
	newVerifiedSubsCount, err := models.GetNewVerifiedSubsCurrentMonth(
		c.Request().Context(),
		d.db.Pool,
	)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"could not get new verified sub count",
			"error",
			err,
		)
		return errorPage(c, views.ErrorPage())
	}

	verifiedSubsCount, err := models.GetVerifiedSubscribers(
		c.Request().Context(),
		d.db.Pool,
	)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"could not get verified sub count",
			"error",
			err,
		)
		return errorPage(c, views.ErrorPage())
	}

	unverifiedSubsCount, err := models.GetUnverifiedSubscribers(
		c.Request().Context(),
		d.db.Pool,
	)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"could not get unverified sub count",
			"error",
			err,
		)
		return errorPage(c, views.ErrorPage())
	}

	return dashboard.Home(dashboard.HomeProps{
		UnverifiedSubscribers: strconv.Itoa(len(unverifiedSubsCount)),
		VerifiedSubscribers:   strconv.Itoa(len(verifiedSubsCount)),
		NewSubscribers:        strconv.Itoa(int(newVerifiedSubsCount)),
	}).Render(renderArgs(c))
}
