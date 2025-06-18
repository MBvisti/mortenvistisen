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

type EmailJobWorker struct {
	emailClient clients.Email
	db          psql.Postgres
	river.WorkerDefaults[jobs.EmailJobArgs]
}

func (w *EmailJobWorker) Work(
	ctx context.Context,
	job *river.Job[jobs.EmailJobArgs],
) error {
	tracer := otel.Tracer(config.Cfg.ServiceName)
	start := time.Now()

	ctx, span := tracer.Start(ctx, "email_job",
		trace.WithAttributes(
			attribute.Int64("job.id", job.ID),
			attribute.String("job.kind", job.Kind),
			attribute.String("email.type", job.Args.Type),
			attribute.String("email.to", job.Args.To),
			attribute.String("email.subject", job.Args.Subject),
			attribute.Int("job.attempt", job.Attempt),
		),
	)
	defer span.End()

	var err error
	if job.Args.Type == "transaction" {
		span.SetAttributes(attribute.String("email.category", "transaction"))
		err = w.emailClient.SendTransaction(
			ctx,
			clients.EmailPayload{
				To:       job.Args.To,
				From:     job.Args.From,
				Subject:  job.Args.Subject,
				HtmlBody: job.Args.HtmlVersion,
				TextBody: job.Args.TextVersion,
			},
		)
	}

	if job.Args.Type != "transaction" {
		span.SetAttributes(attribute.String("email.category", "marketing"))

		// Generate unsubscribe link if subscriber ID is provided
		var unsubscribe clients.Unsubscribe
		if job.Args.SubscriberID != "" {
			subscriberID, err := uuid.Parse(job.Args.SubscriberID)
			if err != nil {
				slog.ErrorContext(ctx, "Invalid subscriber ID format",
					"job_id", job.ID,
					"subscriber_id", job.Args.SubscriberID,
					"error", err,
				)
			} else {
				unsubscribeLink, err := models.GenerateUnsubscribeLink(ctx, w.db.Pool, subscriberID)
				if err != nil {
					slog.ErrorContext(ctx, "Failed to generate unsubscribe link",
						"job_id", job.ID,
						"subscriber_id", job.Args.SubscriberID,
						"error", err,
					)
				} else {
					unsubscribe = clients.Unsubscribe{
						Email: job.Args.To,
						Link:  unsubscribeLink,
					}
				}
			}
		}

		err = w.emailClient.SendMarketing(
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
	}

	duration := time.Since(start)

	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("job.success", false))
		slog.ErrorContext(ctx, "Email job failed",
			"job_id", job.ID,
			"error", err,
			"duration", duration,
			"attempt", job.Attempt,
		)
		return err
	}

	span.SetAttributes(attribute.Bool("job.success", true))
	slog.InfoContext(ctx, "Email job completed successfully",
		"job_id", job.ID,
		"duration", duration,
		"email_type", job.Args.Type,
	)

	return nil
}
