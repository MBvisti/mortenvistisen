package handlers

import (
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/MBvisti/mortenvistisen/pkg/telemetry"
	"github.com/MBvisti/mortenvistisen/repository/psql/database"
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

type Base struct {
	DB database.Queries
	// TknManager      tokens.Manager
	QueueClient *river.Client[pgx.Tx]
	// PostManager     posts.PostManager
	// EmailService    services.Email
	// AuthService     services.Auth
	// TokenService    services.TokenSvc
	// NewsletterModel models.NewsletterService
	// SubscriberModel models.SubscriberService
	// ArticleModel    models.ArticleService
	// TagModel        models.TagService
	// UserModel       models.UserService
	// Database        psql.Postgres
	CookieStore CookieStore
	Tracer      telemetry.Tracer
}

func NewDependencies(
	db database.Queries,
	// tknManager tokens.Manager,
	queueClient *river.Client[pgx.Tx],
	// postManager posts.PostManager,
	// emailService services.Email,
	// tokenService services.TokenSvc,
	// authService services.Auth,
	// newsletterModel models.NewsletterService,
	// subscriberModel models.SubscriberService,
	// articleModel models.ArticleService,
	// tagModel models.TagService,
	// userModel models.UserService,
	// psql psql.Postgres,
	cookieStore CookieStore,
	tracer telemetry.Tracer,
) Base {
	return Base{
		db,
		// tknManager,
		queueClient,
		// postManager,
		// emailService,
		// authService,
		// tokenService,
		// newsletterModel,
		// subscriberModel,
		// articleModel,
		// tagModel,
		// userModel,
		// psql,
		cookieStore,
		tracer,
	}
}

func (bd Base) RedirectHx(w http.ResponseWriter, url string) error {
	slog.Info("redirecting", "url", url)
	w.Header().Set("HX-Redirect", url)
	w.WriteHeader(http.StatusSeeOther)

	return nil
}

func (bd Base) Redirect(w http.ResponseWriter, r *http.Request, url string) error {
	http.Redirect(w, r, url, http.StatusSeeOther)

	return nil
}

func (bd Base) InternalError(ctx echo.Context) error {
	from := "/"

	return views.InternalServerErr(ctx, views.InternalServerErrData{
		FromLocation: from,
	})
}

func (bd Base) CalculateNumberOfPages(totalItems, pageSize int) int {
	return int(math.Ceil(float64(totalItems) / float64(pageSize)))
}

func (bd Base) GetOffsetAndCurrPage(page string, limit int) (int, int, error) {
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

func (bd Base) FormatArticleSlug(slug string) string {
	return fmt.Sprintf("posts/%s", slug)
}

func (bd Base) BuildURLFromSlug(slug string) string {
	return fmt.Sprintf("%s://%s/%s", os.Getenv("APP_SCHEME"), os.Getenv("APP_HOST"), slug)
}

type CookieStore struct {
	store *sessions.CookieStore
}

func NewCookieStore(sessionKey string) CookieStore {
	store := sessions.NewCookieStore([]byte(sessionKey))
	return CookieStore{store}
}

func (cs CookieStore) CreateFlashMsg(
	r *http.Request,
	rw http.ResponseWriter,
	key string,
	args ...string,
) error {
	s, err := cs.store.Get(r, "flashMsg")
	if err != nil {
		return err
	}

	s.AddFlash(key, args...)
	if err := s.Save(r, rw); err != nil {
		return err
	}

	return nil
}

func (cs CookieStore) GetFlashMessages(
	r *http.Request,
	rw http.ResponseWriter,
	key string,
) ([]string, error) {
	s, err := cs.store.Get(r, "flashMsg")
	if err != nil {
		return nil, err
	}

	var msgs []string
	if key != "" {
		for _, f := range s.Flashes(key) {
			msgs = append(msgs, f.(string))
		}
	} else {
		for _, f := range s.Flashes() {
			msgs = append(msgs, f.(string))
		}
	}

	return msgs, nil
}
