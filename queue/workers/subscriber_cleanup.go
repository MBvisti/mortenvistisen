package workers

import (
	"context"
	"log/slog"
	"time"

	"github.com/MBvisti/mortenvistisen/models"
	"github.com/MBvisti/mortenvistisen/queue/jobs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

type SubscriberCleanupWorker struct {
	conn *pgxpool.Pool
	river.WorkerDefaults[jobs.SubscriberCleanupJobArgs]
}

func (w *SubscriberCleanupWorker) Work(
	ctx context.Context,
	job *river.Job[jobs.SubscriberCleanupJobArgs],
) error {
	slog.Info("STARTING CLEAN")
	if err := models.ClearOldUnverifiedSubs(ctx, w.conn); err != nil {
		slog.ErrorContext(ctx, "SubscriberCleanupWorker", "error", err)
		return err
	}
	slog.Info("ENDING CLEAN")

	return nil
}

func (w *SubscriberCleanupWorker) Timeout(
	*river.Job[jobs.SubscriberCleanupJobArgs],
) time.Duration {
	return 10 * time.Minute
}
