package services

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
	"github.com/mbvisti/mortenvistisen/models"
	"github.com/mbvisti/mortenvistisen/psql"
	"github.com/mbvisti/mortenvistisen/psql/queue/jobs"
	"github.com/riverqueue/river"
)

const (
	MaxDailyEmails  = 35
	MinDelayMinutes = 1
	MaxDelayMinutes = 3
	SendStartHour   = 9
	SendStartMinute = 0
)

type NewsletterSendingService struct {
	db    psql.Postgres
	queue *river.Client[pgx.Tx]
}

func NewNewsletterSendingService(database psql.Postgres, queueClient *river.Client[pgx.Tx]) *NewsletterSendingService {
	return &NewsletterSendingService{
		db:    database,
		queue: queueClient,
	}
}

func (s *NewsletterSendingService) MarkNewsletterReadyToSend(ctx context.Context, newsletterID uuid.UUID) error {
	verifiedSubscribers, err := models.GetVerifiedSubscribers(ctx, s.db.Pool)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get verified subscribers", "error", err, "newsletter_id", newsletterID)
		return fmt.Errorf("failed to get verified subscribers: %w", err)
	}

	totalRecipients := len(verifiedSubscribers)
	if totalRecipients == 0 {
		slog.WarnContext(ctx, "no verified subscribers found", "newsletter_id", newsletterID)
		return fmt.Errorf("no verified subscribers found")
	}

	_, err = models.MarkNewsletterReadyToSend(ctx, s.db.Pool, models.MarkNewsletterReadyToSendPayload{
		ID:              newsletterID,
		Now:             time.Now(),
		TotalRecipients: totalRecipients,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to mark newsletter ready to send", "error", err, "newsletter_id", newsletterID)
		return fmt.Errorf("failed to mark newsletter ready to send: %w", err)
	}

	slog.InfoContext(ctx, "newsletter marked ready to send", "newsletter_id", newsletterID, "total_recipients", totalRecipients)
	return nil
}

func (s *NewsletterSendingService) SendNewsletterToSubscribers(ctx context.Context, newsletter models.Newsletter) error {
	if newsletter.SendStatus != "ready_to_send" {
		return fmt.Errorf("newsletter is not ready to send, current status: %s", newsletter.SendStatus)
	}

	slog.InfoContext(ctx, "starting newsletter send process", "newsletter_id", newsletter.ID, "title", newsletter.Title)

	now := time.Now()
	startedAt := now
	_, err := models.UpdateNewsletterSendStatus(ctx, s.db.Pool, models.UpdateNewsletterSendStatusPayload{
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

	verifiedSubscribers, err := models.GetVerifiedSubscribers(ctx, s.db.Pool)
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
		unsubscribeLink, err := GenerateUnsubscribeLink(ctx, s.db, subscriber.ID)
		if err != nil {
			slog.ErrorContext(ctx, "failed to generate unsubscribe link",
				"error", err,
				"subscriber_id", subscriber.ID,
				"newsletter_id", newsletter.ID)
			continue
		}

		delayMinutes := calculateEmailDelay(i)
		scheduledAt := baseTime.Add(time.Duration(delayMinutes) * time.Minute)

		jobArgs := jobs.MarketingEmailJobArgs{
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

	_, err = s.queue.InsertMany(ctx, emailJobs)
	if err != nil {
		slog.ErrorContext(ctx, "failed to insert email jobs", "error", err, "newsletter_id", newsletter.ID)
		return fmt.Errorf("failed to insert email jobs: %w", err)
	}

	completedAt := now
	_, err = models.UpdateNewsletterSendStatus(ctx, s.db.Pool, models.UpdateNewsletterSendStatusPayload{
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

func (s *NewsletterSendingService) ProcessNewslettersReadyToSend(ctx context.Context) error {
	newsletters, err := models.GetNewslettersReadyToSend(ctx, s.db.Pool)
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

		err := s.SendNewsletterToSubscribers(ctx, newsletter)
		if err != nil {
			slog.ErrorContext(ctx, "failed to send newsletter", "error", err, "newsletter_id", newsletter.ID)
			continue
		}

		totalEmailsToday += newsletter.EmailsSent
	}

	slog.InfoContext(ctx, "newsletter processing completed", "total_emails_sent_today", totalEmailsToday)
	return nil
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

	// Use crypto/rand for security compliance
	randomRange := maxDelay - minDelay + 1
	randomBig, err := rand.Int(rand.Reader, big.NewInt(int64(randomRange)))
	if err != nil {
		// Fallback to minimum delay if crypto/rand fails
		return baseDelay + minDelay
	}
	randomDelay := int(randomBig.Int64()) + minDelay

	return baseDelay + randomDelay
}

func stripHTML(html string) string {
	return strings.ReplaceAll(strings.ReplaceAll(html, "<", ""), ">", "")
}
