package handlers

import (
	"log/slog"
	"math"
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

func (d Dashboard) Home(c echo.Context) error {
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

func (d Dashboard) Newsletters(c echo.Context) error {
	const pageSize = 10

	// Get page number from query params, default to 1
	page := 1
	if pageStr := c.QueryParam("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Get total count
	total, err := models.GetNewslettersCount(c.Request().Context(), d.db.Pool)
	if err != nil {
		return err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}

	// Get newsletters for current page
	newsletters, err := models.GetNewslettersPage(
		c.Request().Context(),
		d.db.Pool,
		models.QueryNewslettersParams{
			Limit:  pageSize,
			Offset: int32((page - 1) * pageSize),
		},
	)
	if err != nil {
		return err
	}

	data := dashboard.NewsletterPageData{
		Newsletters: newsletters,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  total,
	}

	return dashboard.Newsletters(data).
		Render(renderArgs(c))
}

func (d Dashboard) CreateNewsletters(c echo.Context) error {
	return dashboard.NewsletterCreate().Render(renderArgs(c))
}
