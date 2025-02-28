package workers

import (
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

type WorkerDependencies struct {
	Conn    *pgxpool.Pool
	Emailer services.Mail
}

func SetupWorkers(deps WorkerDependencies) (*river.Workers, error) {
	workers := river.NewWorkers()

	if err := river.AddWorkerSafely(workers, &EmailJobWorker{
		emailer: deps.Emailer,
	}); err != nil {
		return nil, err
	}

	if err := river.AddWorkerSafely(workers, &SubscriberCleanupWorker{
		conn: deps.Conn,
	}); err != nil {
		return nil, err
	}

	if err := river.AddWorkerSafely(workers, &ScheduleNewsletterReleaseWorker{
		emailer: deps.Emailer,
		conn:    deps.Conn,
	}); err != nil {
		return nil, err
	}

	return workers, nil
}
