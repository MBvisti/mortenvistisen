package workers

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/url"
	"time"

	"github.com/MBvisti/mortenvistisen/config"
	"github.com/MBvisti/mortenvistisen/emails"
	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/queue/jobs"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/MBvisti/mortenvistisen/views/paths"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

type ScheduleNewsletterReleaseWorker struct {
	emailer services.Mail
	conn    *pgxpool.Pool
	river.WorkerDefaults[jobs.ScheduleNewsletterRelease]
}

func (w *ScheduleNewsletterReleaseWorker) Work(
	ctx context.Context,
	job *river.Job[jobs.ScheduleNewsletterRelease],
) error {
	tx, err := w.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	newsletter, err := models.GetNewsletterByID(
		ctx,
		w.conn,
		job.Args.NewsletterID,
	)
	if err != nil {
		return err
	}

	subscribers, err := models.GetVerifiedSubscribers(ctx, tx)
	if err != nil {
		return err
	}

	const (
		emailsPerDay = 50
	)

	// totalDays := int(
	// 	math.Ceil(float64(len(subscribers)) / float64(emailsPerDay)),
	// )

	minutesBetweenEmails := 1 + rand.Intn(4) // Random number between 2-5

	startTime := time.Now()

	var insertMany []river.InsertManyParams
	for i, subscriber := range subscribers {
		dayOffset := i / emailsPerDay
		emailNumberForDay := i % emailsPerDay

		scheduleTime := startTime.
			Add(time.Duration(dayOffset) * 24 * time.Hour).
			Add(time.Duration(emailNumberForDay*minutesBetweenEmails) * time.Minute)

		unsubTkn, err := models.NewToken(
			ctx,
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
				ctx,
				"failed to create unsubscribe token",
				"error", err,
				"subscriber_id", subscriber.ID,
			)
			return err
		}

		html, txt, err := emails.NewsletterMail{
			Title:   newsletter.Title,
			Content: newsletter.Content,
			UnsubscribeLink: fmt.Sprintf(
				"%s%s?token=%s&email=%s",
				config.Cfg.GetFullDomain(),
				paths.Get(ctx, paths.UnsubscribeEvent),
				url.QueryEscape(unsubTkn.Hash),
				url.QueryEscape(subscriber.Email),
			),
		}.Generate(ctx)
		if err != nil {
			return err
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
	}

	client, err := river.ClientFromContextSafely[pgx.Tx](ctx)
	if err != nil {
		return err
	}

	if _, err := client.InsertManyTx(ctx, tx, insertMany); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (w *ScheduleNewsletterReleaseWorker) Timeout(
	*river.Job[jobs.ScheduleNewsletterRelease],
) time.Duration {
	return 5 * time.Minute
}
