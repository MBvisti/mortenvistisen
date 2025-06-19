package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/mbvisti/mortenvistisen/config"
	"github.com/mbvisti/mortenvistisen/emails"
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

func ScheduleNewsletterRelease(
	ctx context.Context,
	psql psql.Postgres,
	newsletterID uuid.UUID,
) error {
	newsletter, err := models.GetNewsletterByID(ctx, psql.Pool, newsletterID)
	if err != nil {
		return err
	}

	verifiedSubscribers, err := models.GetVerifiedSubscribers(ctx, psql.Pool)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"failed to get verified subscribers",
			"error",
			err,
			"newsletter_id",
			newsletterID,
		)
		return fmt.Errorf("failed to get verified subscribers: %w", err)
	}

	if len(verifiedSubscribers) > MaxDailyEmails {
		verifiedSubscribers = verifiedSubscribers[:MaxDailyEmails]
	}

	tx, err := psql.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	baseTime := getNextSendTime()
	emailJobs := make([]river.InsertManyParams, 0, len(verifiedSubscribers))

	html, txt, err := emails.SubscriberWelcome{}.Generate(ctx)
	if err != nil {
		return err
	}

	for i, subscriber := range verifiedSubscribers {
		expiration := time.Now().Add(30 * 24 * time.Hour)

		token, err := models.NewHashedToken(ctx, tx, models.NewTokenPayload{
			Expiration: expiration,
			Meta: models.MetaInformation{
				Resource:   models.ResourceSubscriber,
				ResourceID: subscriber.ID,
				Scope:      models.ScopeUnsubscribe,
			},
		})
		if err != nil {
			return err
		}

		unsubscribeURL := fmt.Sprintf("%s/unsubscribe/%s",
			config.Cfg.GetFullDomain(),
			token.Value,
		)

		scheduledAt := baseTime.Add(
			time.Duration(calculateEmailDelay(i)) * time.Minute,
		)

		jobArgs := jobs.MarketingEmailJobArgs{
			To:              subscriber.Email,
			From:            "MBV <noreply@mortenvistisen.com>",
			Subject:         newsletter.Title,
			HtmlVersion:     html.String(),
			TextVersion:     txt.String(),
			SubscriberID:    subscriber.ID.String(),
			UnsubscribeLink: unsubscribeURL,
		}

		emailJobs = append(emailJobs, river.InsertManyParams{
			Args: jobArgs,
			InsertOpts: &river.InsertOpts{
				MaxAttempts: 5,
				ScheduledAt: scheduledAt,
				Tags:        []string{"newsletter"},
			},
		})
	}

	_, err = psql.Queue().InsertManyTx(ctx, tx, emailJobs)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"failed to insert email jobs",
			"error",
			err,
			"newsletter_id",
			newsletter.ID,
		)
		return fmt.Errorf("failed to insert email jobs: %w", err)
	}

	return nil
}

func getNextSendTime() time.Time {
	now := time.Now()
	sendTime := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		SendStartHour,
		SendStartMinute,
		0,
		0,
		now.Location(),
	)

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
