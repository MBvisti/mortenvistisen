package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/MBvisti/mortenvistisen/config"
	"github.com/MBvisti/mortenvistisen/emails"
	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/psql"
	"github.com/MBvisti/mortenvistisen/queue/jobs"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/MBvisti/mortenvistisen/views/paths"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/riverqueue/river"
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
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
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
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
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

	recentSubs, err := models.GetRecentSubscribers(
		c.Request().Context(),
		d.db.Pool,
	)
	if err != nil {
		return errorPage(c, views.ErrorPage())
	}

	var recent []dashboard.RecentActivity
	for _, rs := range recentSubs {
		recent = append(recent, dashboard.RecentActivity{
			When:     rs.CreatedAt,
			Email:    rs.Email,
			Verified: rs.IsVerified,
		})
	}

	return dashboard.Home(dashboard.HomeProps{
		UnverifiedSubscribers: strconv.Itoa(len(unverifiedSubsCount)),
		VerifiedSubscribers:   strconv.Itoa(len(verifiedSubsCount)),
		NewSubscribers:        strconv.Itoa(int(newVerifiedSubsCount)),
		RecentActivities:      recent,
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
	return dashboard.NewsletterCreate(csrf.Token(c.Request()), false).
		Render(renderArgs(c))
}

func (d Dashboard) StoreNewsletter(c echo.Context) error {
	type newsletterPayload struct {
		Title   string `form:"title"`
		Content string `form:"content"`
	}
	var payload newsletterPayload
	if err := c.Bind(&payload); err != nil {
		return views.ErrorPage().Render(renderArgs(c))
	}

	tx, err := d.db.BeginTx(c.Request().Context())
	if err != nil {
		return views.ErrorPage().Render(renderArgs(c))
	}
	defer tx.Rollback(c.Request().Context())

	newsletter, err := models.NewNewsletter(
		c.Request().Context(),
		tx,
		models.NewNewsletterPayload{
			Title:      payload.Title,
			Content:    payload.Content,
			ReleasedAt: time.Now(),
			Released:   true,
		},
	)
	if err != nil {
		return views.ErrorPage().Render(renderArgs(c))
	}

	subscribers, err := models.GetVerifiedSubscribers(c.Request().Context(), tx)
	if err != nil {
		return views.ErrorPage().Render(renderArgs(c))
	}

	const (
		emailsPerDay = 50
	)

	totalDays := int(
		math.Ceil(float64(len(subscribers)) / float64(emailsPerDay)),
	)

	minutesBetweenEmails := 5 + rand.Intn(6) // Random number between 5-10

	startTime := time.Now()

	var insertMany []river.InsertManyParams
	for i, subscriber := range subscribers {
		dayOffset := i / emailsPerDay
		emailNumberForDay := i % emailsPerDay

		scheduleTime := startTime.
			Add(time.Duration(dayOffset) * 24 * time.Hour).
			Add(time.Duration(emailNumberForDay*minutesBetweenEmails) * time.Minute)

		unsubTkn, err := models.NewToken(
			c.Request().Context(),
			models.NewTokenPayload{
				Expiration: time.Now().Add(365 * (24 * time.Hour)),
				Meta: models.MetaInformation{
					Resource:   models.ResourceSubscriber,
					ResourceID: subscriber.ID,
					Scope:      models.ScopeUnsubscribe,
				},
			},
			tx,
		)
		if err != nil {
			slog.ErrorContext(
				c.Request().Context(),
				"failed to create unsubscribe token",
				"error", err,
				"subscriber_id", subscriber.ID,
			)
			return errorPage(c, views.ErrorPage())
		}

		html, txt, err := emails.NewsletterMail{
			Title:   newsletter.Title,
			Content: newsletter.Content,
			UnsubscribeLink: fmt.Sprintf(
				"%s%s?token=%s&email=%s",
				config.Cfg.GetFullDomain(),
				paths.Get(c.Request().Context(), paths.UnsubscribeEvent),
				url.QueryEscape(unsubTkn.Hash),
				url.QueryEscape(subscriber.Email),
			),
		}.Generate(c.Request().Context())
		if err != nil {
			slog.ErrorContext(
				c.Request().Context(),
				"failed to generate email content",
				"error", err,
				"subscriber_id", subscriber.ID,
			)
			return errorPage(c, views.ErrorPage())
		}

		insertMany = append(insertMany, river.InsertManyParams{
			Args: jobs.EmailJobArgs{
				To:          subscriber.Email,
				From:        "newsletter@mortenvistisen.com",
				Subject:     "newsletter - mortenvistisen.com",
				TextVersion: txt.String(),
				HtmlVersion: html.String(),
			},
			InsertOpts: &river.InsertOpts{
				ScheduledAt: scheduleTime,
			},
		})
		// if _, err := d.db.Queue.InsertManyTx(
		// 	c.Request().Context(),
		// 	tx,
		// 	jobs.EmailJobArgs{
		// 		To:          subscriber.Email,
		// 		From:        "newsletter@mortenvistisen.com",
		// 		Subject:     "Newsletter - mortenvistisen.com",
		// 		TextVersion: txt.String(),
		// 		HtmlVersion: html.String(),
		// 	},
		// 	&river.InsertOpts{
		// 		ScheduledAt: scheduleTime,
		// 	},
		// ); err != nil {
		// 	slog.ErrorContext(
		// 		c.Request().Context(),
		// 		"failed to schedule email",
		// 		"error", err,
		// 		"subscriber_id", subscriber.ID,
		// 		"scheduled_time", scheduleTime,
		// 	)
		// 	return errorPage(c, views.ErrorPage())
		// }
		//
		slog.InfoContext(
			c.Request().Context(),
			"scheduled newsletter email",
			"subscriber_id", subscriber.ID,
			"email", subscriber.Email,
			"scheduled_time", scheduleTime,
			"day", dayOffset+1,
			"total_days", totalDays,
		)
	}

	if _, err := d.db.Queue.InsertManyTx(c.Request().Context(), tx, insertMany); err != nil {
		return errorPage(c, views.ErrorPage())
	}

	if err := tx.Commit(c.Request().Context()); err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"failed to commit transaction",
			"error", err,
		)
		return errorPage(c, views.ErrorPage())
	}

	return dashboard.NewsletterCreate(csrf.Token(c.Request()), true).
		Render(renderArgs(c))
}
