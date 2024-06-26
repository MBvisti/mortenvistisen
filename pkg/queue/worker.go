package queue

import (
	"github.com/MBvisti/mortenvistisen/repository/psql/database"
	"github.com/MBvisti/mortenvistisen/services"
	"github.com/riverqueue/river"
)

type WorkerDependencies struct {
	Db         *database.Queries
	MailClient services.Email
}

func SetupWorkers(deps WorkerDependencies) (*river.Workers, error) {
	workers := river.NewWorkers()

	if err := river.AddWorkerSafely(workers, &EmailJobWorker{
		Sender: &deps.MailClient,
	}); err != nil {
		return nil, err
	}

	return workers, nil
}
