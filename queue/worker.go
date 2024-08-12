package queue

import (
	"github.com/MBvisti/mortenvistisen/pkg/mail_client"
	"github.com/MBvisti/mortenvistisen/psql/database"
	"github.com/riverqueue/river"
)

type WorkerDependencies struct {
	DB      *database.Queries
	Emailer mail_client.AwsSimpleEmailService
}

func SetupWorkers(deps WorkerDependencies) (*river.Workers, error) {
	workers := river.NewWorkers()

	if err := river.AddWorkerSafely(workers, &EmailJobWorker{
		emailer: &deps.Emailer,
	}); err != nil {
		return nil, err
	}

	return workers, nil
}
