package controllers

import (
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/pkg/tokens"
	"github.com/MBvisti/mortenvistisen/posts"
	"github.com/MBvisti/mortenvistisen/repository/database"
	"github.com/MBvisti/mortenvistisen/repository/psql"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/MBvisti/mortenvistisen/views"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/riverqueue/river"
)

// Actions:
// - index | GET
// - create | GET
// - store | POST
// - show | GET
// - edit | GET
// - update | PUT/PATCH
// - destroy | DELETE

type Dependencies struct {
	DB              database.Queries
	TknManager      tokens.Manager
	QueueClient     *river.Client[pgx.Tx]
	PostManager     posts.PostManager
	EmailSvc        services.EmailSvc
	AuthStore       *sessions.CookieStore
	NewsletterModel models.NewsletterService
	SubscriberModel models.SubscriberService
	Database        psql.Postgres
}

func NewDependencies(
	db database.Queries,
	tknManager tokens.Manager,
	queueClient *river.Client[pgx.Tx],
	postManager posts.PostManager,
	emailSvc services.EmailSvc,
	authStore *sessions.CookieStore,
	newsletterModel models.NewsletterService,
	subscriberModel models.SubscriberService,
	psql psql.Postgres,
) Dependencies {
	return Dependencies{
		db,
		tknManager,
		queueClient,
		postManager,
		emailSvc,
		authStore,
		newsletterModel,
		subscriberModel,
		psql,
	}
}

func RedirectHx(w http.ResponseWriter, url string) error {
	slog.Info("redirecting", "url", url)
	w.Header().Set("HX-Redirect", url)
	w.WriteHeader(http.StatusSeeOther)

	return nil
}

func Redirect(w http.ResponseWriter, r *http.Request, url string) error {
	http.Redirect(w, r, url, http.StatusSeeOther)

	return nil
}

func InternalError(ctx echo.Context) error {
	from := "/"

	return views.InternalServerErr(ctx, views.InternalServerErrData{
		FromLocation: from,
	})
}

func CalculateNumberOfPages(totalItems, pageSize int) int {
	return int(math.Ceil(float64(totalItems) / float64(pageSize)))
}

func GetOffsetAndCurrPage(page string, limit int) (int, int, error) {
	var currentPage int
	if page == "" {
		currentPage = 1
	}
	if page != "" {
		cp, err := strconv.Atoi(page)
		if err != nil {
			return 0, 0, err
		}

		currentPage = cp
	}

	offset := 0
	if currentPage == 2 {
		offset = limit
	}

	if currentPage > 2 {
		offset = limit * (currentPage - 1)
	}

	return offset, currentPage, nil
}

func FormatArticleSlug(slug string) string {
	return fmt.Sprintf("posts/%s", slug)
}

func BuildURLFromSlug(slug string) string {
	return fmt.Sprintf("%s://%s/%s", os.Getenv("APP_SCHEME"), os.Getenv("APP_HOST"), slug)
}
