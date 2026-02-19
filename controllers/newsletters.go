package controllers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"mortenvistisen/config"
	"mortenvistisen/email"
	"mortenvistisen/internal/storage"
	"mortenvistisen/models"
	"mortenvistisen/queue"
	"mortenvistisen/queue/jobs"
	"mortenvistisen/router/cookies"
	"mortenvistisen/router/routes"
	"mortenvistisen/services"
	"mortenvistisen/views"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"

	"github.com/labstack/echo/v5"
)

type Newsletters struct {
	db         storage.Pool
	insertOnly queue.InsertOnly
	cfg        config.Config
}

const (
	newsletterDailySendCap = 40
	newsletterSendSpacing  = 3 * time.Second
)

func NewNewsletters(db storage.Pool, insertOnly queue.InsertOnly, cfg config.Config) Newsletters {
	return Newsletters{
		db:         db,
		insertOnly: insertOnly,
		cfg:        cfg,
	}
}

func (n Newsletters) Index(etx *echo.Context) error {
	page := int64(1)
	if p := etx.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = int64(parsed)
		}
	}

	perPage := int64(10)
	if pp := etx.QueryParam("per_page"); pp != "" {
		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 &&
			parsed <= 100 {
			perPage = int64(parsed)
		}
	}

	newslettersList, err := models.PaginateNewsletters(
		etx.Request().Context(),
		n.db.Conn(),
		page,
		perPage,
	)
	if err != nil {
		return render(etx, views.InternalError())
	}

	return render(etx, views.NewsletterIndex(newslettersList))
}

func (n Newsletters) Show(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return render(etx, views.BadRequest())
	}
	newsletterID := int32(parsed)

	newsletter, err := models.FindNewsletter(etx.Request().Context(), n.db.Conn(), newsletterID)
	if err != nil {
		return render(etx, views.NotFound())
	}

	return render(etx, views.NewsletterShow(newsletter))
}

func (n Newsletters) New(etx *echo.Context) error {
	return render(etx, views.NewsletterNew())
}

type CreateNewsletterFormPayload struct {
	Title           string `json:"title"`
	MetaTitle       string `json:"metaTitle"`
	MetaDescription string `json:"metaDescription"`
	IsPublished     bool   `json:"isPublished"`
	ReleasedAt      string `json:"releasedAt"`
	Content         string `json:"content"`
}

func (n Newsletters) Create(etx *echo.Context) error {
	ctx := etx.Request().Context()

	var payload CreateNewsletterFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			ctx,
			"could not parse CreateNewsletterFormPayload",
			"error",
			err,
		)

		return render(etx, views.NotFound())
	}

	data := models.CreateNewsletterData{
		Title:           payload.Title,
		MetaTitle:       payload.MetaTitle,
		MetaDescription: payload.MetaDescription,
		IsPublished:     payload.IsPublished,
		ReleasedAt: func() time.Time {
			if payload.ReleasedAt == "" {
				return time.Time{}
			}
			if t, err := time.Parse("2006-01-02", payload.ReleasedAt); err == nil {
				return t
			}
			return time.Time{}
		}(),
		Content: payload.Content,
	}

	tx, err := n.db.BeginTx(ctx)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, "Failed to create newsletter"); flashErr != nil {
			return flashErr
		}

		return etx.Redirect(http.StatusSeeOther, routes.NewsletterNew.URL())
	}
	defer tx.Rollback(ctx)

	newsletter, err := models.CreateNewsletter(
		ctx,
		tx,
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to create newsletter: %v", err)); flashErr != nil {
			return flashErr
		}
		return etx.Redirect(http.StatusSeeOther, routes.NewsletterNew.URL())
	}

	scheduledJobs := 0
	if newsletter.IsPublished {
		scheduledJobs, err = n.scheduleNewsletterReleaseEmails(ctx, tx, newsletter)
		if err != nil {
			if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to schedule newsletter delivery: %v", err)); flashErr != nil {
				return flashErr
			}

			return etx.Redirect(http.StatusSeeOther, routes.NewsletterNew.URL())
		}
	}

	if err := n.db.CommitTx(ctx, tx); err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, "Failed to create newsletter"); flashErr != nil {
			return flashErr
		}

		return etx.Redirect(http.StatusSeeOther, routes.NewsletterNew.URL())
	}

	successMessage := "Newsletter created successfully"
	if scheduledJobs > 0 {
		successMessage = fmt.Sprintf(
			"Newsletter created and %d delivery jobs were scheduled",
			scheduledJobs,
		)
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, successMessage); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.NewsletterShow.URL(newsletter.ID))
}

func (n Newsletters) Edit(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return render(etx, views.BadRequest())
	}
	newsletterID := int32(parsed)

	newsletter, err := models.FindNewsletter(etx.Request().Context(), n.db.Conn(), newsletterID)
	if err != nil {
		return render(etx, views.NotFound())
	}

	return render(etx, views.NewsletterUpdate(newsletter))
}

type UpdateNewsletterFormPayload struct {
	Title           string `json:"title"`
	MetaTitle       string `json:"metaTitle"`
	MetaDescription string `json:"metaDescription"`
	IsPublished     bool   `json:"isPublished"`
	ReleasedAt      string `json:"releasedAt"`
	Content         string `json:"content"`
}

func (n Newsletters) Update(etx *echo.Context) error {
	ctx := etx.Request().Context()

	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return render(etx, views.BadRequest())
	}
	newsletterID := int32(parsed)

	var payload UpdateNewsletterFormPayload
	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			ctx,
			"could not parse UpdateNewsletterFormPayload",
			"error",
			err,
		)

		return render(etx, views.NotFound())
	}

	data := models.UpdateNewsletterData{
		ID:              newsletterID,
		Title:           payload.Title,
		MetaTitle:       payload.MetaTitle,
		MetaDescription: payload.MetaDescription,
		IsPublished:     payload.IsPublished,
		ReleasedAt: func() time.Time {
			if payload.ReleasedAt == "" {
				return time.Time{}
			}
			if t, err := time.Parse("2006-01-02", payload.ReleasedAt); err == nil {
				return t
			}
			return time.Time{}
		}(),
		Content: payload.Content,
	}

	tx, err := n.db.BeginTx(ctx)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, "Failed to update newsletter"); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(
			http.StatusSeeOther,
			routes.NewsletterEdit.URL(newsletterID),
		)
	}
	defer tx.Rollback(ctx)

	currentNewsletter, err := models.FindNewsletter(ctx, tx, newsletterID)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, "Failed to load newsletter"); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(
			http.StatusSeeOther,
			routes.NewsletterEdit.URL(newsletterID),
		)
	}

	newsletter, err := models.UpdateNewsletter(
		ctx,
		tx,
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to update newsletter: %v", err)); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(
			http.StatusSeeOther,
			routes.NewsletterEdit.URL(newsletterID),
		)
	}

	scheduledJobs := 0
	becamePublished := !currentNewsletter.IsPublished && newsletter.IsPublished
	if becamePublished {
		scheduledJobs, err = n.scheduleNewsletterReleaseEmails(ctx, tx, newsletter)
		if err != nil {
			if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to schedule newsletter delivery: %v", err)); flashErr != nil {
				return render(etx, views.InternalError())
			}
			return etx.Redirect(
				http.StatusSeeOther,
				routes.NewsletterEdit.URL(newsletterID),
			)
		}
	}

	if err := n.db.CommitTx(ctx, tx); err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, "Failed to update newsletter"); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(
			http.StatusSeeOther,
			routes.NewsletterEdit.URL(newsletterID),
		)
	}

	successMessage := "Newsletter updated successfully"
	if scheduledJobs > 0 {
		successMessage = fmt.Sprintf(
			"Newsletter updated and %d delivery jobs were scheduled",
			scheduledJobs,
		)
	}
	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, successMessage); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.NewsletterShow.URL(newsletter.ID))
}

func (n Newsletters) Destroy(etx *echo.Context) error {
	parsed, err := strconv.ParseInt(etx.Param("id"), 10, 32)
	if err != nil {
		return render(etx, views.BadRequest())
	}
	newsletterID := int32(parsed)

	err = models.DestroyNewsletter(etx.Request().Context(), n.db.Conn(), newsletterID)
	if err != nil {
		if flashErr := cookies.AddFlash(etx, cookies.FlashError, fmt.Sprintf("Failed to delete newsletter: %v", err)); flashErr != nil {
			return render(etx, views.InternalError())
		}
		return etx.Redirect(http.StatusSeeOther, routes.NewsletterIndex.URL())
	}

	if flashErr := cookies.AddFlash(etx, cookies.FlashSuccess, "Newsletter destroyed successfully"); flashErr != nil {
		return render(etx, views.InternalError())
	}

	return etx.Redirect(http.StatusSeeOther, routes.NewsletterIndex.URL())
}

func (n Newsletters) scheduleNewsletterReleaseEmails(
	ctx context.Context,
	tx pgx.Tx,
	newsletter models.Newsletter,
) (int, error) {
	subscribers, err := models.AllSubscribers(ctx, tx)
	if err != nil {
		return 0, err
	}

	sort.Slice(subscribers, func(i, j int) bool {
		return subscribers[i].ID < subscribers[j].ID
	})

	readURL := newsletterPublicURL(newsletter)

	scheduleBase := time.Now().UTC().Add(10 * time.Second).Truncate(time.Second)

	var insertParams []river.InsertManyParams
	sendIndex := 0
	for _, subscriber := range subscribers {
		if !subscriber.IsVerified {
			continue
		}

		emailAddress := strings.TrimSpace(subscriber.Email)
		if emailAddress == "" {
			continue
		}

		unsubscribeToken, err := createSubscriberUnsubscribeToken(
			ctx,
			tx,
			n.cfg.Auth.Pepper,
			subscriber.ID,
		)
		if err != nil {
			return 0, err
		}
		unsubscribeURL := newsletterUnsubscribeURL(unsubscribeToken)

		newsletterEmail := email.NewsletterRelease{
			NewsletterTitle: newsletter.Title,
			IssueLabel:      newsletterIssueLabel(newsletter),
			Highlights:      newsletter.MetaDescription,
			NewsletterHTML:  services.MarkdownToHTML(newsletter.Content),
			ReadURL:         readURL,
			UnsubscribeURL:  unsubscribeURL,
		}

		htmlBody, err := newsletterEmail.ToHTML()
		if err != nil {
			return 0, err
		}

		textBody, err := newsletterEmail.ToText()
		if err != nil {
			return 0, err
		}

		scheduledAt := scheduledNewsletterSendTime(scheduleBase, sendIndex)
		insertParams = append(insertParams, river.InsertManyParams{
			Args: jobs.SendMarketingEmailArgs{
				Data: email.MarketingData{
					To:             []string{emailAddress},
					From:           "newsletter@mortenvistisen.com",
					Subject:        newsletter.Title,
					HTMLBody:       htmlBody,
					TextBody:       textBody,
					UnsubscribeURL: unsubscribeURL,
					Tags:           []string{"newsletter_release"},
					Metadata: map[string]string{
						"newsletter_id": strconv.Itoa(int(newsletter.ID)),
						"subscriber_id": strconv.Itoa(int(subscriber.ID)),
					},
				},
			},
			InsertOpts: &river.InsertOpts{
				ScheduledAt: scheduledAt,
			},
		})

		sendIndex++
	}

	if len(insertParams) == 0 {
		return 0, nil
	}

	if _, err := n.insertOnly.InsertManyTx(ctx, tx, insertParams); err != nil {
		return 0, err
	}

	return len(insertParams), nil
}

func scheduledNewsletterSendTime(base time.Time, sendIndex int) time.Time {
	dayOffset := sendIndex / newsletterDailySendCap
	positionInDay := sendIndex % newsletterDailySendCap

	return base.
		AddDate(0, 0, dayOffset).
		Add(time.Duration(positionInDay) * newsletterSendSpacing)
}

func newsletterIssueLabel(newsletter models.Newsletter) string {
	if newsletter.ReleasedAt.IsZero() {
		return ""
	}

	return newsletter.ReleasedAt.Format("January 2, 2006")
}

func newsletterPublicURL(newsletter models.Newsletter) string {
	baseURL := strings.TrimRight(config.BaseURL, "/")
	if newsletter.Slug != "" {
		return fmt.Sprintf("%s/newsletters/%s", baseURL, newsletter.Slug)
	}

	return fmt.Sprintf("%s%s", baseURL, routes.NewsletterShow.URL(newsletter.ID))
}

func newsletterUnsubscribeURL(token string) string {
	base := strings.TrimRight(config.BaseURL, "/")
	if token == "" {
		return fmt.Sprintf("%s/unsubscribe", base)
	}

	return fmt.Sprintf("%s/unsubscribe?token=%s", base, url.QueryEscape(token))
}
