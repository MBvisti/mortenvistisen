package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"math"
	"strconv"
	"time"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/psql"
	"github.com/MBvisti/mortenvistisen/queue/jobs"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/MBvisti/mortenvistisen/views/dashboard"
	"github.com/dromara/carbon/v2"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type Dashboard struct {
	db psql.Postgres
}

func newDashboard(db psql.Postgres) Dashboard {
	return Dashboard{db}
}

func (d Dashboard) Home(c echo.Context) error {
	end := carbon.Now().EndOfHour()
	start := end.SubHours(24)

	dailyViews, err := models.GetSiteViewsByDate(
		c.Request().Context(),
		d.db.Pool,
		start.StdTime(),
		end.StdTime(),
	)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "ERROR", "error_value", err)
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

	recentSubs, err := models.GetRecentSubscribers(
		c.Request().Context(),
		d.db.Pool,
	)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "ERROR", "error_value", err)
		return errorPage(c, views.ErrorPage())
	}

	var recent []dashboard.RecentActivity
	for _, rs := range recentSubs {
		recent = append(recent, dashboard.RecentActivity{
			ID:       rs.ID,
			When:     rs.CreatedAt,
			Email:    rs.Email,
			Verified: rs.IsVerified,
		})
	}

	var stats []dashboard.HourlyStat

	oneDayAgo := carbon.Now(carbon.Berlin).SubDay()
	for i := range 24 {
		h := oneDayAgo.StartOfHour().
			AddHours(i + 1).
			ToKitchenString(carbon.Berlin)
		stat := dashboard.HourlyStat{
			Hour: h,
		}
		var visi int64
		var vies int64
		for _, dv := range dailyViews {
			kitchenTime := carbon.CreateFromStdTime(dv.CreatedAt, carbon.Berlin).
				StartOfHour().
				ToKitchenString()
			if kitchenTime == h {
				visi++
				vies++
			}
		}

		stat.Visits = visi
		stat.Views = vies

		stats = append(stats, stat)
	}

	mStats, err := json.Marshal(stats)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "ERROR", "error_value", err)
		return errorPage(c, views.ErrorPage())
	}

	uniqueVisitors := make(map[uuid.UUID]struct{})
	for _, dv := range dailyViews {
		if dv.VisitorID != uuid.Nil {
			uniqueVisitors[dv.VisitorID] = struct{}{}
		}
	}

	return dashboard.Home(dashboard.HomeProps{
		HourlyStats:         string(mStats),
		DailyVisits:         strconv.Itoa(len(uniqueVisitors)),
		VerifiedSubscribers: strconv.Itoa(len(verifiedSubsCount)),
		DailyViews:          strconv.Itoa(len(dailyViews)),
		RecentActivities:    recent,
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
		slog.ErrorContext(
			c.Request().Context(),
			"StoreNewsletter",
			"error",
			err,
		)
		return views.ErrorPage().Render(renderArgs(c))
	}

	tx, err := d.db.BeginTx(c.Request().Context())
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"StoreNewsletter",
			"error",
			err,
		)
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
		slog.ErrorContext(
			c.Request().Context(),
			"StoreNewsletter",
			"error",
			err,
		)
		return views.ErrorPage().Render(renderArgs(c))
	}

	if _, err := d.db.Queue.InsertTx(c.Request().Context(), tx, jobs.ScheduleNewsletterRelease{
		NewsletterID: newsletter.ID,
	}, nil); err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"StoreNewsletter",
			"error",
			err,
		)
		return views.ErrorPage().Render(renderArgs(c))
	}

	if err := tx.Commit(c.Request().Context()); err != nil {
		return views.ErrorPage().Render(renderArgs(c))
	}

	return dashboard.NewsletterCreate(csrf.Token(c.Request()), true).
		Render(renderArgs(c))
}

func (d Dashboard) ShowSubscriber(c echo.Context) error {
	subID := c.Param("id")
	id, err := uuid.Parse(subID)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "ShowSubscriber", "error", err)
		return errorPage(c, views.ErrorPage())
	}

	subcriber, err := models.GetSubscriberByID(
		c.Request().Context(),
		d.db.Pool,
		id,
	)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "ShowSubscriber", "error", err)
		return errorPage(c, views.ErrorPage())
	}

	return dashboard.ShowSubscriber(dashboard.ShowSubscriberProps{
		Email:        subcriber.Email,
		SubscribedAt: subcriber.SubscribedAt,
		Referere:     subcriber.Referer,
		Verified:     subcriber.IsVerified,
	}).
		Render(renderArgs(c))
}
