package jobs

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mbvisti/mortenvistisen/config"
	"github.com/mbvisti/mortenvistisen/models"
	"github.com/mbvisti/mortenvistisen/psql"
	"github.com/riverqueue/river"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

const (
	MaxDailyEmails  = 35
	MinDelayMinutes = 1
	MaxDelayMinutes = 3
	SendStartHour   = 9
	SendStartMinute = 0
)

type NewsletterProcessingJobArgs struct {
}

func (NewsletterProcessingJobArgs) Kind() string {
	return "newsletter_processing_job"
}

type NewsletterProcessingJobWorker struct {
	river.WorkerDefaults[NewsletterProcessingJobArgs]
	db    psql.Postgres
	queue *river.Client[pgx.Tx]
}

func NewNewsletterProcessingJobWorker(database psql.Postgres, queueClient *river.Client[pgx.Tx]) *NewsletterProcessingJobWorker {
	return &NewsletterProcessingJobWorker{
		db:    database,
		queue: queueClient,
	}
}

func (w *NewsletterProcessingJobWorker) Work(ctx context.Context, job *river.Job[NewsletterProcessingJobArgs]) error {
	tracer := otel.Tracer("newsletter-processing-worker")
	ctx, span := tracer.Start(ctx, "newsletter-processing-job")
	defer span.End()

	span.SetAttributes(
		attribute.String("job.kind", job.Kind),
		attribute.String("job.id", fmt.Sprintf("%d", job.ID)),
	)

	slog.InfoContext(ctx, "starting newsletter processing job",
		"job_id", job.ID,
		"attempt", job.Attempt)

	err := w.processNewslettersReadyToSend(ctx)
	if err != nil {
		span.RecordError(err)
		slog.ErrorContext(ctx, "failed to process newsletters ready to send", "error", err)
		return fmt.Errorf("failed to process newsletters ready to send: %w", err)
	}

	slog.InfoContext(ctx, "newsletter processing job completed successfully",
		"job_id", job.ID,
		"attempt", job.Attempt)

	return nil
}

func (w *NewsletterProcessingJobWorker) processNewslettersReadyToSend(ctx context.Context) error {
	newsletters, err := models.GetNewslettersReadyToSend(ctx, w.db.Pool)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get newsletters ready to send", "error", err)
		return fmt.Errorf("failed to get newsletters ready to send: %w", err)
	}

	if len(newsletters) == 0 {
		slog.InfoContext(ctx, "no newsletters ready to send")
		return nil
	}

	totalEmailsToday := 0

	for _, newsletter := range newsletters {
		if totalEmailsToday >= MaxDailyEmails {
			slog.WarnContext(ctx, "daily email limit reached, skipping remaining newsletters",
				"daily_limit", MaxDailyEmails,
				"emails_sent_today", totalEmailsToday)
			break
		}

		remainingQuota := MaxDailyEmails - totalEmailsToday
		if newsletter.TotalRecipients > remainingQuota {
			slog.WarnContext(ctx, "newsletter would exceed daily limit, skipping",
				"newsletter_id", newsletter.ID,
				"recipients", newsletter.TotalRecipients,
				"remaining_quota", remainingQuota)
			continue
		}

		err := w.sendNewsletterToSubscribers(ctx, newsletter)
		if err != nil {
			slog.ErrorContext(ctx, "failed to send newsletter", "error", err, "newsletter_id", newsletter.ID)
			continue
		}

		totalEmailsToday += newsletter.EmailsSent
	}

	slog.InfoContext(ctx, "newsletter processing completed", "total_emails_sent_today", totalEmailsToday)
	return nil
}

func (w *NewsletterProcessingJobWorker) sendNewsletterToSubscribers(ctx context.Context, newsletter models.Newsletter) error {
	if newsletter.SendStatus != "ready_to_send" {
		return fmt.Errorf("newsletter is not ready to send, current status: %s", newsletter.SendStatus)
	}

	slog.InfoContext(ctx, "starting newsletter send process", "newsletter_id", newsletter.ID, "title", newsletter.Title)

	now := time.Now()
	startedAt := now
	_, err := models.UpdateNewsletterSendStatus(ctx, w.db.Pool, models.UpdateNewsletterSendStatusPayload{
		ID:               newsletter.ID,
		Now:              now,
		SendStatus:       "sending",
		SendingStartedAt: &startedAt,
		EmailsSent:       0,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to update newsletter status to sending", "error", err, "newsletter_id", newsletter.ID)
		return fmt.Errorf("failed to update newsletter status: %w", err)
	}

	verifiedSubscribers, err := models.GetVerifiedSubscribers(ctx, w.db.Pool)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get verified subscribers", "error", err, "newsletter_id", newsletter.ID)
		return fmt.Errorf("failed to get verified subscribers: %w", err)
	}

	if len(verifiedSubscribers) > MaxDailyEmails {
		slog.WarnContext(ctx, "subscriber count exceeds daily limit",
			"newsletter_id", newsletter.ID,
			"subscriber_count", len(verifiedSubscribers),
			"daily_limit", MaxDailyEmails)
		verifiedSubscribers = verifiedSubscribers[:MaxDailyEmails]
	}

	baseTime := getNextSendTime()
	emailJobs := make([]river.InsertManyParams, 0, len(verifiedSubscribers))

	for i, subscriber := range verifiedSubscribers {
		unsubscribeLink, err := w.generateUnsubscribeLink(ctx, subscriber.ID)
		if err != nil {
			slog.ErrorContext(ctx, "failed to generate unsubscribe link",
				"error", err,
				"subscriber_id", subscriber.ID,
				"newsletter_id", newsletter.ID)
			continue
		}

		delayMinutes := calculateEmailDelay(i)
		scheduledAt := baseTime.Add(time.Duration(delayMinutes) * time.Minute)

		jobArgs := MarketingEmailJobArgs{
			To:              subscriber.Email,
			From:            "Newsletter <newsletter@example.com>",
			Subject:         newsletter.Title,
			HtmlVersion:     newsletter.Content,
			TextVersion:     stripHTML(newsletter.Content),
			SubscriberID:    subscriber.ID.String(),
			UnsubscribeLink: unsubscribeLink,
		}

		emailJobs = append(emailJobs, river.InsertManyParams{
			Args: jobArgs,
		})

		slog.DebugContext(ctx, "newsletter email job scheduled",
			"newsletter_id", newsletter.ID,
			"subscriber_email", subscriber.Email,
			"scheduled_at", scheduledAt,
			"delay_minutes", delayMinutes)
	}

	if len(emailJobs) == 0 {
		slog.ErrorContext(ctx, "no email jobs created", "newsletter_id", newsletter.ID)
		return fmt.Errorf("no email jobs could be created")
	}

	_, err = w.queue.InsertMany(ctx, emailJobs)
	if err != nil {
		slog.ErrorContext(ctx, "failed to insert email jobs", "error", err, "newsletter_id", newsletter.ID)
		return fmt.Errorf("failed to insert email jobs: %w", err)
	}

	completedAt := now
	_, err = models.UpdateNewsletterSendStatus(ctx, w.db.Pool, models.UpdateNewsletterSendStatusPayload{
		ID:                 newsletter.ID,
		Now:                time.Now(),
		SendStatus:         "sent",
		SendingStartedAt:   &startedAt,
		SendingCompletedAt: &completedAt,
		EmailsSent:         len(emailJobs),
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to update newsletter status to sent", "error", err, "newsletter_id", newsletter.ID)
		return fmt.Errorf("failed to update newsletter status: %w", err)
	}

	slog.InfoContext(ctx, "newsletter sending completed",
		"newsletter_id", newsletter.ID,
		"emails_scheduled", len(emailJobs),
		"first_email_at", baseTime,
		"last_email_at", baseTime.Add(time.Duration(calculateEmailDelay(len(emailJobs)-1))*time.Minute))

	return nil
}

func (w *NewsletterProcessingJobWorker) generateUnsubscribeLink(ctx context.Context, subscriberID uuid.UUID) (string, error) {
	tx, err := w.db.BeginTx(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	expiration := time.Now().Add(30 * 24 * time.Hour)

	token, err := models.NewHashedToken(ctx, tx, models.NewTokenPayload{
		Expiration: expiration,
		Meta: models.MetaInformation{
			Resource:   models.ResourceSubscriber,
			ResourceID: subscriberID,
			Scope:      models.ScopeUnsubscribe,
		},
	})
	if err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	unsubscribeURL := fmt.Sprintf("%s/unsubscribe/%s",
		config.Cfg.GetFullDomain(),
		token.Value)

	return unsubscribeURL, nil
}

func getNextSendTime() time.Time {
	now := time.Now()
	sendTime := time.Date(now.Year(), now.Month(), now.Day(), SendStartHour, SendStartMinute, 0, 0, now.Location())

	if sendTime.Before(now) {
		sendTime = sendTime.Add(24 * time.Hour)
	}

	return sendTime
}

func calculateEmailDelay(index int) int {
	minDelay := MinDelayMinutes
	maxDelay := MaxDelayMinutes

	baseDelay := index * 2

	randomRange := maxDelay - minDelay + 1
	randomBig, err := rand.Int(rand.Reader, big.NewInt(int64(randomRange)))
	if err != nil {
		return baseDelay + minDelay
	}
	randomDelay := int(randomBig.Int64()) + minDelay

	return baseDelay + randomDelay
}

func stripHTML(html string) string {
	return strings.ReplaceAll(strings.ReplaceAll(html, "<", ""), ">", "")
}
