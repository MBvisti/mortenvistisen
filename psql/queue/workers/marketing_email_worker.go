package workers

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/mbvisti/mortenvistisen/clients"
	"github.com/mbvisti/mortenvistisen/config"
	"github.com/mbvisti/mortenvistisen/models"
	"github.com/mbvisti/mortenvistisen/psql"
	"github.com/mbvisti/mortenvistisen/psql/queue/jobs"
	"github.com/riverqueue/river"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type MarketingEmailJobWorker struct {
	river.WorkerDefaults[jobs.MarketingEmailJobArgs]
	emailClient clients.Email
	db          psql.Postgres
}

func NewMarketingEmailWorker(
	emailClient clients.Email,
	db psql.Postgres,
) *MarketingEmailJobWorker {
	return &MarketingEmailJobWorker{
		emailClient: emailClient,
		db:          db,
	}
}

func (w *MarketingEmailJobWorker) Work(
	ctx context.Context,
	job *river.Job[jobs.MarketingEmailJobArgs],
) error {
	tracer := otel.Tracer(config.Cfg.ServiceName)
	ctx, span := tracer.Start(ctx, "marketing_email_job_worker",
		trace.WithAttributes(
			attribute.Int64("job.id", job.ID),
			attribute.String("email.to", job.Args.To),
			attribute.String("email.subject", job.Args.Subject),
			attribute.String("email.category", "marketing"),
		),
	)
	defer span.End()

	start := time.Now()
	slog.InfoContext(ctx, "Processing marketing email job",
		"job_id", job.ID,
		"to", job.Args.To,
		"subject", job.Args.Subject,
		"subscriber_id", job.Args.SubscriberID,
	)

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
		slog.ErrorContext(ctx, "Failed to send marketing email",
			"job_id", job.ID,
			"to", job.Args.To,
			"error", err,
			"duration", time.Since(start),
		)
		span.RecordError(err)
		
		if job.Args.NewsletterID != "" && job.Args.SubscriberID != "" {
			newsletterID, parseErr := uuid.Parse(job.Args.NewsletterID)
			subscriberID, parseErr2 := uuid.Parse(job.Args.SubscriberID)
			if parseErr == nil && parseErr2 == nil {
				_, trackErr := models.UpdateNewsletterEmailSendStatus(
					ctx,
					w.db.Pool,
					newsletterID,
					subscriberID,
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
			}
		}
		
		return err
	}

	if job.Args.NewsletterID != "" && job.Args.SubscriberID != "" {
		newsletterID, parseErr := uuid.Parse(job.Args.NewsletterID)
		subscriberID, parseErr2 := uuid.Parse(job.Args.SubscriberID)
		if parseErr == nil && parseErr2 == nil {
			_, trackErr := models.UpdateNewsletterEmailSendStatus(
				ctx,
				w.db.Pool,
				newsletterID,
				subscriberID,
				"sent",
				"",
			)
			if trackErr != nil {
				slog.ErrorContext(ctx, "Failed to update email send status",
					"error", trackErr,
					"newsletter_id", job.Args.NewsletterID,
					"subscriber_id", job.Args.SubscriberID,
				)
			}
		}
	}

	slog.InfoContext(ctx, "Marketing email sent successfully",
		"job_id", job.ID,
		"to", job.Args.To,
		"duration", time.Since(start),
	)

	return nil
}
