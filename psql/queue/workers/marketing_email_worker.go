package workers

import (
	"context"
	"log/slog"

	"github.com/mbvisti/mortenvistisen/clients"
	"github.com/mbvisti/mortenvistisen/models"
	"github.com/mbvisti/mortenvistisen/psql"
	"github.com/mbvisti/mortenvistisen/psql/queue/jobs"
	"github.com/riverqueue/river"
)

type MarketingEmailJobWorker struct {
	river.WorkerDefaults[jobs.MarketingEmailJobArgs]
	emailClient clients.Email
	db          psql.Postgres
}

func (w *MarketingEmailJobWorker) Work(
	ctx context.Context,
	job *river.Job[jobs.MarketingEmailJobArgs],
) error {
	unsubscribe := clients.Unsubscribe{
		Email: job.Args.To,
		Link:  job.Args.UnsubscribeLink,
	}

	err := w.emailClient.SendMarketing(
		ctx,
		clients.EmailPayload{
			To:       job.Args.To,
			From:     job.Args.From,
			Subject:  job.Args.Subject,
			HtmlBody: job.Args.HtmlVersion,
			TextBody: job.Args.TextVersion,
		},
		unsubscribe,
	)
	if err != nil {
		_, trackErr := models.UpdateNewsletterEmailSendStatus(
			ctx,
			w.db.Pool,
			job.Args.NewsletterID,
			job.Args.SubscriberID,
			"failed",
			err.Error(),
		)
		if trackErr != nil {
			slog.ErrorContext(ctx, "Failed to update email send status",
				"error", trackErr,
				"newsletter_id", job.Args.NewsletterID,
				"subscriber_id", job.Args.SubscriberID,
			)
		}

		return err
	}

	_, err = models.UpdateNewsletterEmailSendStatus(
		ctx,
		w.db.Pool,
		job.Args.NewsletterID,
		job.Args.SubscriberID,
		"sent",
		"",
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to update email send status",
			"error", err,
			"newsletter_id", job.Args.NewsletterID,
			"subscriber_id", job.Args.SubscriberID,
		)

		return err
	}

	return nil
}
